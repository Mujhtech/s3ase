package protocol

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"mime"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
	"github.com/mujhtech/s3ase/config"
	"github.com/mujhtech/s3ase/database/models"
	"github.com/mujhtech/s3ase/internal/pkg/sse"
	"github.com/mujhtech/s3ase/job"
	"github.com/mujhtech/s3ase/job/handlers"
	"github.com/mujhtech/s3ase/services"
	"github.com/rs/zerolog"
	tusHandler "github.com/tus/tusd/v2/pkg/handler"
	"github.com/tus/tusd/v2/pkg/memorylocker"
	tusS3Store "github.com/tus/tusd/v2/pkg/s3store"
)

var (
	reMimeType = regexp.MustCompile(`^[a-z]+\/[a-z0-9\-\+\.]+$`)
	// We only allow certain URL-safe characters in upload IDs. URL-safe in this means
	// that their are allowed in a URI's path component according to RFC 3986.
	// See https://datatracker.ietf.org/doc/html/rfc3986#section-3.3
	reValidUploadId = regexp.MustCompile(`^[A-Za-z0-9\-._~%!$'()*+,;=/:@]*$`)
)

var mimeInlineBrowserWhitelist = map[string]struct{}{
	"text/plain": {},

	"image/png":  {},
	"image/jpeg": {},
	"image/gif":  {},
	"image/bmp":  {},
	"image/webp": {},

	"audio/wave":      {},
	"audio/wav":       {},
	"audio/x-wav":     {},
	"audio/x-pn-wav":  {},
	"audio/webm":      {},
	"video/webm":      {},
	"audio/ogg":       {},
	"video/ogg":       {},
	"application/ogg": {},
}

type Protocol struct {
	s3  *s3.Client
	cfg *config.Config

	CompleteUploads chan HookEvent

	TerminatedUploads chan HookEvent

	UploadProgress chan HookEvent

	CreatedUploads chan HookEvent

	job    *job.Job
	sse    sse.Streamer
	locker *memorylocker.MemoryLocker
}

func NewProtocol(cfg *config.Config, s3 *s3.Client, job *job.Job, sse sse.Streamer) (*Protocol, error) {

	return &Protocol{
		s3:     s3,
		sse:    sse,
		cfg:    cfg,
		job:    job,
		locker: memorylocker.New(),
	}, nil
}

func (p *Protocol) UploadFile(s services.CreateFileService, w http.ResponseWriter, r *http.Request) {
	c := p.getContext(w, r)

	s3Store := tusS3Store.New(s.App.Bucket, p.s3)

	if s.Body.FolderID != "" {
		//s3Store.ObjectPrefix = s.Body.FolderID
		folder, err := s.GetFolder(c)

		if err != nil {
			p.sendError(c, err)
			return
		}

		if folder != nil {
			s3Store.ObjectPrefix = folder.ID
		}
	}

	// Parse headers
	contentType := r.Header.Get("Content-Type")
	contentDisposition := r.Header.Get("Content-Disposition")
	// The current browser flow sends one complete file per request. The tus S3
	// store is retained underneath for bounded multipart streaming.
	willCompleteUpload := true

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
	file, err := s.Run(c, info.ID, info.Size, info.MetaData)
	if err != nil {
		p.sendError(c, err)
		return
	}
	//if p.config.PreUploadCreateCallback != nil {
	// resp2, changes, err := p.config.PreUploadCreateCallback(newHookEvent(c, info))
	// if err != nil {
	// 	p.sendError(c, err)
	// 	return
	// }
	// resp = resp.MergeWith(resp2)

	// // Apply changes returned from the pre-create hook.
	if file.ID != "" {
		if err := validateUploadId(file.ID); err != nil {
			p.sendError(c, err)
			return
		}

		info.ID = file.ID
	}

	// if changes.MetaData != nil {
	// 	info.MetaData = changes.MetaData
	// }

	// if changes.Storage != nil {
	// 	info.Storage = changes.Storage
	// }
	//}

	upload, err := s3Store.NewUpload(c, toTusFileInfo(info))
	if err != nil {
		p.failUpload(c, s, file, err)
		p.sendError(c, err)
		return
	}

	info, err = uploadToFileInfo(c, upload)
	if err != nil {
		p.failUpload(c, s, file, err)
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
	if err = p.sse.Publish(c, s.App.ID, sse.EventTypeUploadStarted, sse.UploadProgress{
		FileID:   file.ID,
		Name:     file.Name,
		Status:   sse.UploadProgressStatusStarted,
		Progress: 0,
	}); err != nil {
		log.Printf("failed to publish upload started event: %v", err)
	}
	p.enqueueWebhook(s.App.ID, "upload.started", file)
	file.Status = models.FileStatusUploading
	if err = s.FileRepo.UpdateFile(c, file); err != nil {
		_ = s3Store.AsTerminatableUpload(upload).Terminate(c)
		p.failUpload(c, s, file, err)
		p.sendError(c, err)
		return
	}

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
	resp, err = p.writeChunk(c, s, s3Store, resp, upload, file, info)
	if err != nil {
		_ = s3Store.AsTerminatableUpload(upload).Terminate(c)
		p.failUpload(c, s, file, err)
		p.sendError(c, err)
		return
	}

	// 4. Finish upload, if necessary
	if willCompleteUpload && info.SizeIsDeferred {
		info, err = uploadToFileInfo(c, upload)
		if err != nil {
			p.failUpload(c, s, file, err)
			p.sendError(c, err)
			return
		}

		uploadLength := info.Offset

		lengthDeclarableUpload := s3Store.AsLengthDeclarableUpload(upload)
		if err := lengthDeclarableUpload.DeclareLength(c, uploadLength); err != nil {
			p.failUpload(c, s, file, err)
			p.sendError(c, err)
			return
		}

		info.Size = uploadLength
		info.SizeIsDeferred = false

		resp, err = p.finishUploadIfComplete(c, s, resp, upload, file, info)
		if err != nil {
			p.failUpload(c, s, file, err)
			p.sendError(c, err)
			return
		}

	}

	body, err := json.Marshal(struct {
		Data    *models.File `json:"data"`
		Message string       `json:"message"`
	}{
		Data:    file,
		Message: "file uploaded",
	})
	if err != nil {
		p.sendError(c, err)
		return
	}
	resp.Header["Content-Type"] = "application/json; charset=utf-8"
	resp.Body = string(body)

	p.sendResp(c, resp)
}

func (p *Protocol) failUpload(c *httpContext, s services.CreateFileService, file *models.File, uploadErr error) {
	if file == nil {
		return
	}
	file.Status = models.FileStatusFailed
	if err := s.FileRepo.UpdateFile(c, file); err != nil {
		zerolog.Ctx(c).Error().Err(err).Msg("failed to persist upload failure")
	}
	if err := p.sse.Publish(c, file.AppID, sse.EventTypeUploadFailed, sse.UploadProgress{
		FileID: file.ID,
		Name:   file.Name,
		Status: sse.UploadProgressStatusFailed,
	}); err != nil {
		zerolog.Ctx(c).Error().Err(err).Msg("failed to publish upload failure")
	}
	p.enqueueWebhook(file.AppID, "upload.failed", file)
	zerolog.Ctx(c).Error().Err(uploadErr).Str("file_id", file.ID).Msg("upload failed")
}

func (p *Protocol) enqueueWebhook(appID, event string, data interface{}) {
	rawData, err := json.Marshal(data)
	if err != nil {
		return
	}
	payload, err := json.Marshal(handlers.WebhookPayload{
		ID: uuid.NewString(), AppID: appID, Event: event, CreatedAt: time.Now(), Data: rawData,
	})
	if err != nil {
		return
	}
	if err := p.job.Client.Enqueue(job.QueueNameDefault, job.JobNameWebhook, &job.ClientPayload{Data: payload}); err != nil {
		log.Printf("failed to enqueue webhook: %v", err)
	}
}

func (p *Protocol) GetFile(w http.ResponseWriter, r *http.Request) {
	c := p.getContext(w, r)

	id, err := extractIDFromPath(r.URL.Path)
	if err != nil {
		p.sendError(c, err)
		return
	}
	//c.log = c.log.With("id", id)

	// if handler.composer.UsesLocker {
	// 	lock, err := handler.lockUpload(c, id)
	// 	if err != nil {
	// 		p.sendError(c, err)
	// 		return
	// 	}

	// 	defer lock.Unlock()
	// }

	s3Store := tusS3Store.New("", p.s3)

	upload, err := s3Store.GetUpload(c, id)
	if err != nil {
		p.sendError(c, err)
		return
	}

	info, err := upload.GetInfo(c)
	if err != nil {
		p.sendError(c, err)
		return
	}

	contentType, contentDisposition := filterContentType(toFileInfo(info))
	resp := HTTPResponse{
		StatusCode: http.StatusOK,
		Header: HTTPHeader{
			"Content-Length":      strconv.FormatInt(info.Offset, 10),
			"Content-Type":        contentType,
			"Content-Disposition": contentDisposition,
		},
		Body: "", // Body is intentionally left empty, and we copy it manually in later.
	}

	// If no data has been uploaded yet, respond with an empty "204 No Content" status.
	if info.Offset == 0 {
		resp.StatusCode = http.StatusNoContent
		p.sendResp(c, resp)
		return
	}

	src, err := upload.GetReader(c)
	if err != nil {
		p.sendError(c, err)
		return
	}

	p.sendResp(c, resp)
	if _, err = io.Copy(w, src); err != nil {
		_ = src.Close()
		p.sendError(c, err)
		return
	}
	if err := src.Close(); err != nil {
		p.sendError(c, err)
	}
}

func extractIDFromPath(path string) (string, error) {
	return strings.Trim(path, "/"), nil
}

func filterContentType(info FileInfo) (contentType string, contentDisposition string) {
	filetype := info.MetaData["filetype"]

	if reMimeType.MatchString(filetype) {
		// If the filetype from metadata is well formed, we forward use this
		// for the Content-Type header. However, only whitelisted mime types
		// will be allowed to be shown inline in the browser
		contentType = filetype
		if _, isWhitelisted := mimeInlineBrowserWhitelist[filetype]; isWhitelisted {
			contentDisposition = "inline"
		} else {
			contentDisposition = "attachment"
		}
	} else {
		// If the filetype from the metadata is not well formed, we use a
		// default type and force the browser to download the content.
		contentType = "application/octet-stream"
		contentDisposition = "attachment"
	}

	// Add a filename to Content-Disposition if one is available in the metadata
	if filename, ok := info.MetaData["filename"]; ok {
		contentDisposition += ";filename=" + strconv.Quote(filename)
	}

	return contentType, contentDisposition
}

// writeChunk reads the body from the requests r and appends it to the upload
// with the corresponding id. Afterwards, it will set the necessary response
// headers but will not send the response.
func (p *Protocol) writeChunk(c *httpContext, s services.CreateFileService, s3store tusS3Store.S3Store, resp HTTPResponse, upload tusHandler.Upload, file *models.File, info FileInfo) (HTTPResponse, error) {
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
		if p.cfg.Protocol.MaxSize > 0 {
			// Ensure that the upload does not exceed the maximum upload size
			maxSize = p.cfg.Protocol.MaxSize - offset
		} else {
			// If no upload limit is given, we allow arbitrary sizes
			maxSize = math.MaxInt64
		}
	}
	if length > 0 {
		maxSize = length
	}

	zerolog.Ctx(c).Info().Msgf("ChunkWriteStart maxSize %v offset %v", maxSize, offset)

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
			if err := c.resC.SetReadDeadline(time.Now().Add(p.cfg.Protocol.NetworkTimeout)); err != nil {
				zerolog.Ctx(c).Warn().Msgf("NetworkTimeoutError error %v", err)
			}

			// The write deadline is updated accordingly to ensure that we can also write responses.
			if err := c.resC.SetWriteDeadline(time.Now().Add(2 * p.cfg.Protocol.NetworkTimeout)); err != nil {
				zerolog.Ctx(c).Warn().Msgf("NetworkTimeoutError error %v", err)
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

		stopProgress := p.sendProgressMessages(c, s, file, info)

		bytesWritten, err = upload.WriteChunk(c, offset, c.body)
		stopProgress()

		// If we encountered an error while reading the body from the HTTP request, log it, but only include
		// it in the response, if the store did not also return an error.
		bodyErr := c.body.hasError()
		if bodyErr != nil {
			zerolog.Ctx(c).Error().Msgf("BodyReadError error %v", bodyErr.Error())
			if err == nil {
				err = bodyErr
			}
		}

		// Terminate the upload if it was stopped, as indicated by the ErrUploadStoppedByServer error.
		terminateUpload := errors.Is(bodyErr, tusHandler.ErrUploadStoppedByServer)
		if terminateUpload {
			if terminateErr := p.terminateUpload(c, s, s3store, upload, file, info); terminateErr != nil {
				// We only log this error and not show it to the user since this
				// termination error is not relevant to the uploading client
				zerolog.Ctx(c).Error().Msgf("UploadStopTerminateError error %v", terminateErr.Error())
			}
		}
	}

	zerolog.Ctx(c).Info().Msgf("ChunkWriteComplete bytesWritten %v", bytesWritten)

	// Send new offset to client
	newOffset := offset + bytesWritten
	resp.Header["Upload-Offset"] = strconv.FormatInt(newOffset, 10)
	//handler.Metrics.incBytesReceived(uint64(bytesWritten))
	info.Offset = newOffset

	// We try to finish the upload, even if an error occurred. If we have a previous error,
	// we return it and its HTTP response.
	finishResp, finishErr := p.finishUploadIfComplete(c, s, resp, upload, file, info)
	if err != nil {
		return resp, err
	}

	return finishResp, finishErr
}

// finishUploadIfComplete checks whether an upload is completed (i.e. upload offset
// matches upload size) and if so, it will call the data store's FinishUpload
// function and emit the necessary events for the hooks.
func (p *Protocol) finishUploadIfComplete(c *httpContext, s services.CreateFileService, resp HTTPResponse, upload tusHandler.Upload, file *models.File, info FileInfo) (HTTPResponse, error) {
	// If the upload is completed, ...

	zerolog.Ctx(c).Info().Msgf("ChunkWriteComplete offset %v size %v sizeIsDeferred %v", info.Offset, info.Size, info.SizeIsDeferred)

	// if !info.SizeIsDeferred && info.Offset == info.Size {
	if info.Offset == info.Size {
		var err error
		// ... allow the data storage to finish and cleanup the upload
		if err = upload.FinishUpload(c); err != nil {
			return resp, err
		}

		// ... and call pre-finish callback and send post-finish notification.
		resp, err = p.emitFinishEvents(c, s, resp, file, info)
		if err != nil {
			return resp, err
		}
	}

	return resp, nil
}

// emitFinishEvents calls the PreFinishResponseCallback function and sends
// the necessary message on the CompleteUpload channel.
func (p *Protocol) emitFinishEvents(c *httpContext, s services.CreateFileService, resp HTTPResponse, file *models.File, info FileInfo) (HTTPResponse, error) {
	// if handler.config.PreFinishResponseCallback != nil {
	// 	resp2, err := handler.config.PreFinishResponseCallback(newHookEvent(c, info))
	// 	if err != nil {
	// 		return resp, err
	// 	}
	// 	resp = resp.MergeWith(resp2)
	// }
	// update file size
	file.Size = info.Size
	file.Status = models.FileStatusCompleted

	zerolog.Ctx(c).Info().Msgf("File Info%v", info)
	zerolog.Ctx(c).Info().Msgf("File%v", file)

	if err := s.FileRepo.UpdateFile(c, file); err != nil {
		return resp, err
	}

	zerolog.Ctx(c).Info().Msgf("UploadFinished size %v", info.Size)
	// handler.Metrics.incUploadsFinished()

	// if handler.config.NotifyCompleteUploads {
	// 	handler.CompleteUploads <- newHookEvent(c, info)
	// }

	if err := p.sse.Publish(c, file.AppID, sse.EventTypeUploadCompleted, sse.UploadProgress{
		FileID:   file.ID,
		Name:     file.Name,
		Status:   sse.UploadProgressStatusCompleted,
		Progress: 100,
	}); err != nil {
		log.Printf("failed to publish upload completed event: %v", err)
	}
	p.enqueueWebhook(file.AppID, "upload.completed", file)

	// handle update file status

	return resp, nil
}

func (p *Protocol) terminateUpload(c *httpContext, s services.CreateFileService, s3Store tusS3Store.S3Store, upload tusHandler.Upload, file *models.File, info FileInfo) error {
	terminatableUpload := s3Store.AsTerminatableUpload(upload)

	err := terminatableUpload.Terminate(c)
	if err != nil {
		return err
	}

	// if handler.config.NotifyTerminatedUploads {
	// 	handler.TerminatedUploads <- newHookEvent(c, info)
	// }

	if err = s.FileRepo.DeleteFile(c, file.ID); err != nil {
		log.Printf("failed to delete file: %v", err)
	}

	if err := p.sse.Publish(c, file.AppID, sse.EventTypeUploadCancelled, sse.UploadProgress{
		FileID:   file.ID,
		Name:     file.Name,
		Status:   sse.UploadProgressStatusCancelled,
		Progress: info.Offset,
	}); err != nil {
		log.Printf("failed to publish upload cancelled event: %v", err)
	}

	zerolog.Ctx(c).Info().Msgf("UploadTerminated")
	//handler.Metrics.incUploadsTerminated()

	return nil
}

// Send the error in the response body. The status code will be looked up in
// ErrStatusCodes. If none is found 500 Internal Error will be used.
func (p *Protocol) sendError(c *httpContext, err error) {
	r := c.req

	var detailedErr Error

	if !errors.As(err, &detailedErr) {
		zerolog.Ctx(c).Error().Msgf("InternalServerError message %v", err.Error())
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

	zerolog.Ctx(c).Info().Msgf("ResponseOutgoing status %v body %v", resp.StatusCode, resp.Body)
}

// sendProgressMessage will send a notification over the UploadProgress channel
// indicating how much data has been transfered to the server.
// It will stop sending these instances once the provided context is done.
func (p *Protocol) sendProgressMessages(c *httpContext, s services.CreateFileService, file *models.File, info FileInfo) func() {
	//hook := newHookEvent(c, info)

	previousOffset := int64(0)
	//originalOffset := hook.Upload.Offset
	originalOffset := int64(0)

	emitProgress := func() {
		//hook.Upload.Offset = originalOffset + c.body.bytesRead()
		//if hook.Upload.Offset != previousOffset {
		//p.UploadProgress <- hook
		//previousOffset = hook.Upload.Offset
		previousOffset = originalOffset + c.body.bytesRead()

		progress := int64(0)

		if info.Size > 0 {
			progress = int64(float64(previousOffset) / float64(info.Size) * 100)
		}

		if err := p.sse.Publish(c, s.App.ID, sse.EventTypeUploadProgress, sse.UploadProgress{
			FileID:   file.ID,
			Name:     file.Name,
			Status:   sse.UploadProgressStatusUploading,
			Progress: progress,
		}); err != nil {
			log.Printf("failed to publish upload progress event: %v", err)
		}
		//}
	}

	done := make(chan struct{})
	stopped := make(chan struct{})
	go func() {
		defer close(stopped)
		ticker := time.NewTicker(p.cfg.Protocol.UploadProgressInterval)
		defer ticker.Stop()
		for {
			select {
			case <-done:
				emitProgress()
				return
			case <-ticker.C:
				emitProgress()
			}
		}
	}()
	var once sync.Once
	return func() {
		once.Do(func() {
			close(done)
			<-stopped
		})
	}
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
