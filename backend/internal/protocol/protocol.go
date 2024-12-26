package protocol

import (
	"context"
	"errors"
	"fmt"
	"math"
	"mime"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/mujhtech/s3ase/config"
	tusHandler "github.com/tus/tusd/v2/pkg/handler"
	tusS3Store "github.com/tus/tusd/v2/pkg/s3store"
)

var (
	reForwardedHost  = regexp.MustCompile(`host="?([^;"]+)`)
	reForwardedProto = regexp.MustCompile(`proto=(https?)`)
	reMimeType       = regexp.MustCompile(`^[a-z]+\/[a-z0-9\-\+\.]+$`)
	// We only allow certain URL-safe characters in upload IDs. URL-safe in this means
	// that their are allowed in a URI's path component according to RFC 3986.
	// See https://datatracker.ietf.org/doc/html/rfc3986#section-3.3
	reValidUploadId = regexp.MustCompile(`^[A-Za-z0-9\-._~%!$'()*+,;=/:@]*$`)
)

type Protocol struct {
	s3Store *tusS3Store.S3Store
	config  *config.Config

	CompleteUploads chan HookEvent

	TerminatedUploads chan HookEvent

	UploadProgress chan HookEvent

	CreatedUploads chan HookEvent
}

func NewProtocol() (*Protocol, error) {
	return &Protocol{}, nil
}

func (p *Protocol) PostFileV2(w http.ResponseWriter, r *http.Request) {
	c := p.getContext(w, r)

	// Parse headers
	contentType := r.Header.Get("Content-Type")
	contentDisposition := r.Header.Get("Content-Disposition")
	willCompleteUpload := isIETFDraftUploadComplete(r)

	info := FileInfo{
		MetaData: make(MetaData),
	}

	size, sizeIsDeferred, err := getIETFDraftUploadLength(r)
	if err != nil {
		p.sendError(c, err)
		return
	}

	if !sizeIsDeferred {
		info.Size = size
	} else {
		// Error out if the storage does not support upload length deferring, but we need it.
		// if !handler.composer.UsesLengthDeferrer {
		// 	p.sendError(c, tusHandler.ErrNotImplemented)
		// 	return
		// }

		info.SizeIsDeferred = true
	}

	// Parse Content-Type and Content-Disposition to get file type or file name
	if contentType != "" {
		fileType, _, err := mime.ParseMediaType(contentType)
		if err != nil {
			p.sendError(c, err)
			return
		}

		info.MetaData["filetype"] = fileType
	}

	if contentDisposition != "" {
		_, values, err := mime.ParseMediaType(contentDisposition)
		if err != nil {
			p.sendError(c, err)
			return
		}

		if values["filename"] != "" {
			info.MetaData["filename"] = values["filename"]
		}
	}

	resp := HTTPResponse{
		StatusCode: http.StatusCreated,
		Header:     HTTPHeader{},
	}

	// 1. Create upload resource
	// if p.config.PreUploadCreateCallback != nil {
	// 	resp2, changes, err := handler.config.PreUploadCreateCallback(newHookEvent(c, info))
	// 	if err != nil {
	// 		p.sendError(c, err)
	// 		return
	// 	}
	// 	resp = resp.MergeWith(resp2)

	// 	// Apply changes returned from the pre-create hook.
	// 	if changes.ID != "" {
	// 		if err := validateUploadId(changes.ID); err != nil {
	// 			p.sendError(c, err)
	// 			return
	// 		}

	// 		info.ID = changes.ID
	// 	}

	// 	if changes.MetaData != nil {
	// 		info.MetaData = changes.MetaData
	// 	}

	// 	if changes.Storage != nil {
	// 		info.Storage = changes.Storage
	// 	}
	// }

	upload, err := p.s3Store.NewUpload(c, toTusFileInfo(info))
	if err != nil {
		p.sendError(c, err)
		return
	}

	info, err = uploadToFileInfo(c, upload)
	if err != nil {
		p.sendError(c, err)
		return
	}

	//id := info.ID
	//url := p.absFileURL(r, id)
	limits := p.getIETFDraftUploadLimits(info)
	//resp.Header["Location"] = url
	resp.Header["Upload-Limit"] = limits

	// Send 104 response
	//w.Header().Set("Location", url)
	//w.Header().Set("Upload-Draft-Interop-Version", string(currentUploadDraftInteropVersion))
	//w.Header().Set("Upload-Limit", limits)
	//w.WriteHeader(104)

	// handler.Metrics.incUploadsCreated()
	// c.log = c.log.With("id", id)
	// c.log.Info("UploadCreated", "size", info.Size, "url", url)

	// if handler.config.NotifyCreatedUploads {
	// 	handler.CreatedUploads <- newHookEvent(c, info)
	// }

	// 2. Lock upload
	// if handler.composer.UsesLocker {
	// 	lock, err := handler.lockUpload(c, id)
	// 	if err != nil {
	// 		p.sendError(c, err)
	// 		return
	// 	}

	// 	defer lock.Unlock()
	// }

	// 3. Write chunk
	resp, err = p.writeChunk(c, resp, upload, info)
	if err != nil {
		p.sendError(c, err)
		return
	}

	// 4. Finish upload, if necessary
	if willCompleteUpload && info.SizeIsDeferred {
		info, err = uploadToFileInfo(c, upload)
		if err != nil {
			p.sendError(c, err)
			return
		}

		uploadLength := info.Offset

		lengthDeclarableUpload := p.s3Store.AsLengthDeclarableUpload(upload)
		if err := lengthDeclarableUpload.DeclareLength(c, uploadLength); err != nil {
			p.sendError(c, err)
			return
		}

		info.Size = uploadLength
		info.SizeIsDeferred = false

		resp, err = p.finishUploadIfComplete(c, resp, upload, info)
		if err != nil {
			p.sendError(c, err)
			return
		}

	}

	p.sendResp(c, resp)
}

// writeChunk reads the body from the requests r and appends it to the upload
// with the corresponding id. Afterwards, it will set the necessary response
// headers but will not send the response.
func (p *Protocol) writeChunk(c *httpContext, resp HTTPResponse, upload tusHandler.Upload, info FileInfo) (HTTPResponse, error) {
	// Get Content-Length if possible
	r := c.req
	length := r.ContentLength
	offset := info.Offset

	// Test if this upload fits into the file's size
	if !info.SizeIsDeferred && offset+length > info.Size {
		return resp, tusHandler.ErrSizeExceeded
	}

	maxSize := info.Size - offset
	// If the upload's length is deferred and the PATCH request does not contain the Content-Length
	// header (which is allowed if 'Transfer-Encoding: chunked' is used), we still need to set limits for
	// the body size.
	if info.SizeIsDeferred {
		if p.config.Protocol.MaxSize > 0 {
			// Ensure that the upload does not exceed the maximum upload size
			maxSize = p.config.Protocol.MaxSize - offset
		} else {
			// If no upload limit is given, we allow arbitrary sizes
			maxSize = math.MaxInt64
		}
	}
	if length > 0 {
		maxSize = length
	}

	c.log.Info("ChunkWriteStart", "maxSize", maxSize, "offset", offset)

	var bytesWritten int64
	var err error
	// Prevent a nil pointer dereference when accessing the body which may not be
	// available in the case of a malicious request.
	if r.Body != nil {
		// Limit the data read from the request's body to the allowed maximum. We use
		// http.MaxBytesReader instead of io.LimitedReader because it returns an error
		// if too much data is provided (handled in bodyReader) and also stops the server
		// from reading the remaining request body.
		c.body = newBodyReader(c, maxSize)
		c.body.onReadDone = func() {
			// Update the read deadline for every successful read operation. This ensures that the request handler
			// keeps going while data is transmitted but that dead connections can also time out and be cleaned up.
			if err := c.resC.SetReadDeadline(time.Now().Add(p.config.Protocol.NetworkTimeout)); err != nil {
				c.log.Warn("NetworkTimeoutError", "error", err)
			}

			// The write deadline is updated accordingly to ensure that we can also write responses.
			if err := c.resC.SetWriteDeadline(time.Now().Add(2 * p.config.Protocol.NetworkTimeout)); err != nil {
				c.log.Warn("NetworkTimeoutError", "error", err)
			}
		}

		// We use a callback to allow the hook system to cancel an upload. The callback
		// cancels the request context causing the request body to be closed with the
		// provided error.
		info.stopUpload = func(res HTTPResponse) {
			cause := tusHandler.ErrUploadStoppedByServer
			cause.HTTPResponse = cause.HTTPResponse.MergeWith(toTusdHTTPResponse(res))
			c.cancel(cause)
		}

		// if handler.config.NotifyUploadProgress {
		// 	handler.sendProgressMessages(c, info)
		// }

		bytesWritten, err = upload.WriteChunk(c, offset, c.body)

		// If we encountered an error while reading the body from the HTTP request, log it, but only include
		// it in the response, if the store did not also return an error.
		bodyErr := c.body.hasError()
		if bodyErr != nil {
			c.log.Error("BodyReadError", "error", bodyErr.Error())
			if err == nil {
				err = bodyErr
			}
		}

		// Terminate the upload if it was stopped, as indicated by the ErrUploadStoppedByServer error.
		// terminateUpload := errors.Is(bodyErr, tusHandler.ErrUploadStoppedByServer)
		// if terminateUpload && p.composer.UsesTerminater {
		// 	if terminateErr := p.terminateUpload(c, upload, info); terminateErr != nil {
		// 		// We only log this error and not show it to the user since this
		// 		// termination error is not relevant to the uploading client
		// 		c.log.Error("UploadStopTerminateError", "error", terminateErr.Error())
		// 	}
		// }
	}

	c.log.Info("ChunkWriteComplete", "bytesWritten", bytesWritten)

	// Send new offset to client
	newOffset := offset + bytesWritten
	resp.Header["Upload-Offset"] = strconv.FormatInt(newOffset, 10)
	//handler.Metrics.incBytesReceived(uint64(bytesWritten))
	info.Offset = newOffset

	// We try to finish the upload, even if an error occurred. If we have a previous error,
	// we return it and its HTTP response.
	finishResp, finishErr := p.finishUploadIfComplete(c, resp, upload, info)
	if err != nil {
		return resp, err
	}

	return finishResp, finishErr
}

// finishUploadIfComplete checks whether an upload is completed (i.e. upload offset
// matches upload size) and if so, it will call the data store's FinishUpload
// function and emit the necessary events for the hooks.
func (p *Protocol) finishUploadIfComplete(c *httpContext, resp HTTPResponse, upload tusHandler.Upload, info FileInfo) (HTTPResponse, error) {
	// If the upload is completed, ...
	if !info.SizeIsDeferred && info.Offset == info.Size {
		var err error
		// ... allow the data storage to finish and cleanup the upload
		if err = upload.FinishUpload(c); err != nil {
			return resp, err
		}

		// ... and call pre-finish callback and send post-finish notification.
		resp, err = p.emitFinishEvents(c, resp, info)
		if err != nil {
			return resp, err
		}
	}

	return resp, nil
}

// emitFinishEvents calls the PreFinishResponseCallback function and sends
// the necessary message on the CompleteUpload channel.
func (p *Protocol) emitFinishEvents(c *httpContext, resp HTTPResponse, info FileInfo) (HTTPResponse, error) {
	// if handler.config.PreFinishResponseCallback != nil {
	// 	resp2, err := handler.config.PreFinishResponseCallback(newHookEvent(c, info))
	// 	if err != nil {
	// 		return resp, err
	// 	}
	// 	resp = resp.MergeWith(resp2)
	// }

	c.log.Info("UploadFinished", "size", info.Size)
	// handler.Metrics.incUploadsFinished()

	// if handler.config.NotifyCompleteUploads {
	// 	handler.CompleteUploads <- newHookEvent(c, info)
	// }

	return resp, nil
}

func (p *Protocol) terminateUpload(c *httpContext, upload tusHandler.Upload, info FileInfo) error {
	terminatableUpload := p.s3Store.AsTerminatableUpload(upload)

	err := terminatableUpload.Terminate(c)
	if err != nil {
		return err
	}

	// if handler.config.NotifyTerminatedUploads {
	// 	handler.TerminatedUploads <- newHookEvent(c, info)
	// }

	c.log.Info("UploadTerminated")
	//handler.Metrics.incUploadsTerminated()

	return nil
}

// Send the error in the response body. The status code will be looked up in
// ErrStatusCodes. If none is found 500 Internal Error will be used.
func (p *Protocol) sendError(c *httpContext, err error) {
	r := c.req

	var detailedErr Error

	if !errors.As(err, &detailedErr) {
		c.log.Error("InternalServerError", "message", err.Error())
		detailedErr = NewError("ERR_INTERNAL_SERVER_ERROR", err.Error(), http.StatusInternalServerError)
	}

	// If we are sending the response for a HEAD request, ensure that we are not including
	// any response body.
	if r.Method == "HEAD" {
		detailedErr.HTTPResponse.Body = ""
	}

	p.sendResp(c, detailedErr.HTTPResponse)
	//handler.Metrics.incErrorsTotal(detailedErr)
}

// sendResp writes the header to w with the specified status code.
func (p *Protocol) sendResp(c *httpContext, resp HTTPResponse) {
	resp.writeTo(c.res)

	c.log.Info("ResponseOutgoing", "status", resp.StatusCode, "body", resp.Body)
}

// sendProgressMessage will send a notification over the UploadProgress channel
// indicating how much data has been transfered to the server.
// It will stop sending these instances once the provided context is done.
func (p *Protocol) sendProgressMessages(c *httpContext, info tusHandler.FileInfo) {
	hook := newHookEvent(c, info)

	previousOffset := int64(0)
	originalOffset := hook.Upload.Offset

	emitProgress := func() {
		hook.Upload.Offset = originalOffset + c.body.bytesRead()
		if hook.Upload.Offset != previousOffset {
			p.UploadProgress <- hook
			previousOffset = hook.Upload.Offset
		}
	}

	go func() {
		for {
			select {
			case <-c.Done():
				emitProgress()
				return
			case <-time.After(p.config.Protocol.UploadProgressInterval):
				emitProgress()
			}
		}
	}()
}

func validateUploadId(newId string) error {
	if newId == "" {
		// An empty ID from FileInfoChanges is allowed. The store will then
		// just pick an ID.
		return nil
	}

	if strings.HasPrefix(newId, "/") || strings.HasSuffix(newId, "/") {
		// Disallow leading and trailing slashes, as these would be
		// stripped away by extractIDFromPath, which can cause problems and confusion.
		return fmt.Errorf("validation error in FileInfoChanges: ID must not begin or end with a forward slash (got: %s)", newId)
	}

	if !reValidUploadId.MatchString(newId) {
		// Disallow some non-URL-safe characters in the upload ID to
		// prevent issues with URL parsing, which are though to debug for users.
		return fmt.Errorf("validation error in FileInfoChanges: ID must contain only URL-safe character: %s (got: %s)", reValidUploadId.String(), newId)
	}

	return nil
}

func toTusFileInfo(info FileInfo) tusHandler.FileInfo {
	// Convert MetaData to tusHandler.MetaData type
	handlerMetaData := make(tusHandler.MetaData)
	for k, v := range info.MetaData {
		handlerMetaData[k] = v
	}

	return tusHandler.FileInfo{
		ID:             info.ID,
		Size:           info.Size,
		MetaData:       handlerMetaData,
		Storage:        info.Storage,
		SizeIsDeferred: info.SizeIsDeferred,
		IsPartial:      info.IsPartial,
		IsFinal:        info.IsFinal,
		Offset:         info.Offset,
		PartialUploads: info.PartialUploads,
	}
}

func uploadToFileInfo(ctx context.Context, upload tusHandler.Upload) (FileInfo, error) {

	info, err := upload.GetInfo(ctx)

	if err != nil {
		return FileInfo{}, err
	}

	return toFileInfo(info), nil
}

func toFileInfo(info tusHandler.FileInfo) FileInfo {
	// Convert tusHandler.MetaData to MetaData type
	metadata := make(MetaData)
	for k, v := range info.MetaData {
		metadata[k] = v
	}
	return FileInfo{
		ID:             info.ID,
		Size:           info.Size,
		MetaData:       metadata,
		Storage:        info.Storage,
		SizeIsDeferred: info.SizeIsDeferred,
		IsPartial:      info.IsPartial,
		IsFinal:        info.IsFinal,
		Offset:         info.Offset,
		PartialUploads: info.PartialUploads,
	}
}
