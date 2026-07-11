package protocol

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/mujhtech/s3ase/database/models"
	"github.com/mujhtech/s3ase/internal/pkg/sse"
	"github.com/mujhtech/s3ase/services"
	tusHandler "github.com/tus/tusd/v2/pkg/handler"
	tusS3Store "github.com/tus/tusd/v2/pkg/s3store"
)

// TusUpload serves the standard tus 1.0 creation, offset, append, and
// termination requests. The direct POST /files endpoint remains available for
// API clients which upload a file in one request.
func (p *Protocol) TusUpload(s services.CreateFileService, basePath string, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var file *models.File
	var err error

	if r.Method == http.MethodPost {
		folder, folderErr := s.GetFolder(ctx)
		if folderErr != nil {
			http.Error(w, folderErr.Error(), http.StatusBadRequest)
			return
		}
		if s.Body.FolderID != "" && folder == nil {
			http.Error(w, "folder not found", http.StatusNotFound)
			return
		}
	} else {
		fileID, idErr := tusFileIDFromPath(r.URL.Path, basePath)
		if idErr != nil {
			http.Error(w, idErr.Error(), http.StatusBadRequest)
			return
		}
		file, err = s.FileRepo.FindFileByID(ctx, fileID)
		if err != nil {
			http.Error(w, "upload not found", http.StatusNotFound)
			return
		}
		if file.AppID != s.App.ID {
			http.Error(w, "upload not found", http.StatusNotFound)
			return
		}
		if file.FolderID.Valid {
			s.Body.FolderID = file.FolderID.String
		}
	}

	s3Store := tusS3Store.New(s.App.Bucket, p.s3)
	if s.Body.FolderID != "" {
		s3Store.ObjectPrefix = s.Body.FolderID
	}
	composer := tusHandler.NewStoreComposer()
	s3Store.UseIn(composer)
	p.locker.UseIn(composer)

	config := tusHandler.Config{
		BasePath:           strings.TrimSuffix(basePath, "/") + "/",
		StoreComposer:      composer,
		MaxSize:            p.cfg.Protocol.MaxSize,
		DisableDownload:    true,
		DisableTermination: false,
		Cors:               &tusHandler.CorsConfig{Disable: true},
		NetworkTimeout:     p.cfg.Protocol.NetworkTimeout,
	}

	config.PreUploadCreateCallback = func(event tusHandler.HookEvent) (tusHandler.HTTPResponse, tusHandler.FileInfoChanges, error) {
		created, createErr := s.Run(event.Context, "", event.Upload.Size, event.Upload.MetaData)
		if createErr != nil {
			return tusHandler.HTTPResponse{}, tusHandler.FileInfoChanges{}, createErr
		}
		file = created
		file.Status = models.FileStatusUploading
		if updateErr := s.FileRepo.UpdateFile(event.Context, file); updateErr != nil {
			return tusHandler.HTTPResponse{}, tusHandler.FileInfoChanges{}, updateErr
		}

		p.publishUploadEvent(event.Context, file, sse.EventTypeUploadStarted, sse.UploadProgressStatusStarted, 0)
		p.enqueueWebhook(file.AppID, "upload.started", file)

		return tusHandler.HTTPResponse{
			Header: tusHandler.HTTPHeader{"X-File-ID": file.ID},
		}, tusHandler.FileInfoChanges{ID: file.ID}, nil
	}

	config.PreFinishResponseCallback = func(event tusHandler.HookEvent) (tusHandler.HTTPResponse, error) {
		completed, findErr := p.findTusFile(event.Context, s, event.Upload.ID)
		if findErr != nil {
			return tusHandler.HTTPResponse{}, findErr
		}
		completed.Size = event.Upload.Size
		completed.Status = models.FileStatusCompleted
		if updateErr := s.FileRepo.UpdateFile(event.Context, completed); updateErr != nil {
			return tusHandler.HTTPResponse{}, updateErr
		}
		p.publishUploadEvent(event.Context, completed, sse.EventTypeUploadCompleted, sse.UploadProgressStatusCompleted, 100)
		p.enqueueWebhook(completed.AppID, "upload.completed", completed)
		return tusHandler.HTTPResponse{Header: tusHandler.HTTPHeader{"X-File-ID": completed.ID}}, nil
	}

	h, err := tusHandler.NewHandler(config)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	recorder := &statusResponseWriter{ResponseWriter: w}
	trimmedBase := strings.TrimSuffix(basePath, "/")
	http.StripPrefix(trimmedBase, h).ServeHTTP(recorder, r)

	if r.Method == http.MethodPost && file != nil && recorder.status >= http.StatusBadRequest {
		p.failTusUpload(ctx, s, file, fmt.Errorf("tus upload creation failed with status %d", recorder.status))
		return
	}

	if recorder.status < http.StatusOK || recorder.status >= http.StatusMultipleChoices || file == nil {
		return
	}

	switch r.Method {
	case http.MethodPatch:
		offset, parseErr := strconv.ParseInt(recorder.Header().Get("Upload-Offset"), 10, 64)
		if parseErr == nil {
			progress := int64(0)
			if file.Size > 0 {
				progress = min(100, offset*100/file.Size)
			}
			if progress < 100 {
				p.publishUploadEvent(ctx, file, sse.EventTypeUploadProgress, sse.UploadProgressStatusUploading, progress)
			}
		}
	case http.MethodDelete:
		file.Status = models.FileStatusCancelled
		if updateErr := s.FileRepo.UpdateFile(ctx, file); updateErr != nil {
			log.Printf("failed to persist cancelled upload: %v", updateErr)
		}
		p.publishUploadEvent(ctx, file, sse.EventTypeUploadCancelled, sse.UploadProgressStatusCancelled, 0)
		p.enqueueWebhook(file.AppID, "upload.cancelled", file)
	}
}

func (p *Protocol) findTusFile(ctx context.Context, s services.CreateFileService, uploadID string) (*models.File, error) {
	fileID := tusObjectID(uploadID)
	file, err := s.FileRepo.FindFileByID(ctx, fileID)
	if err != nil {
		return nil, err
	}
	if file.AppID != s.App.ID {
		return nil, fmt.Errorf("upload does not belong to this app")
	}
	return file, nil
}

func (p *Protocol) failTusUpload(ctx context.Context, s services.CreateFileService, file *models.File, uploadErr error) {
	file.Status = models.FileStatusFailed
	if err := s.FileRepo.UpdateFile(ctx, file); err != nil {
		log.Printf("failed to persist tus upload failure: %v", err)
	}
	p.publishUploadEvent(ctx, file, sse.EventTypeUploadFailed, sse.UploadProgressStatusFailed, 0)
	p.enqueueWebhook(file.AppID, "upload.failed", file)
	log.Printf("tus upload failed for file %s: %v", file.ID, uploadErr)
}

func (p *Protocol) publishUploadEvent(ctx context.Context, file *models.File, eventType sse.EventType, status sse.UploadProgressStatus, progress int64) {
	folderID := ""
	if file.FolderID.Valid {
		folderID = file.FolderID.String
	}
	if err := p.sse.Publish(ctx, file.AppID, eventType, sse.UploadProgress{
		FileID: file.ID, Name: file.Name, FolderID: folderID, Progress: progress, Status: status,
	}); err != nil {
		log.Printf("failed to publish %s event: %v", eventType, err)
	}
}

func tusFileIDFromPath(requestPath, basePath string) (string, error) {
	escapedID := strings.TrimPrefix(requestPath, strings.TrimSuffix(basePath, "/")+"/")
	if escapedID == requestPath || escapedID == "" {
		return "", fmt.Errorf("missing tus upload id")
	}
	id, err := url.PathUnescape(escapedID)
	if err != nil {
		return "", fmt.Errorf("invalid tus upload id: %w", err)
	}
	fileID := tusObjectID(id)
	if fileID == "" {
		return "", fmt.Errorf("invalid tus upload id")
	}
	return fileID, nil
}

func tusObjectID(uploadID string) string {
	if index := strings.IndexByte(uploadID, '+'); index >= 0 {
		return uploadID[:index]
	}
	return uploadID
}

type statusResponseWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusResponseWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func (w *statusResponseWriter) Write(body []byte) (int, error) {
	if w.status == 0 {
		w.WriteHeader(http.StatusOK)
	}
	return w.ResponseWriter.Write(body)
}
