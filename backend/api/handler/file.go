package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/mujhtech/s3ase/api/dto"
	"github.com/mujhtech/s3ase/api/middleware"
	"github.com/mujhtech/s3ase/database/models"
	"github.com/mujhtech/s3ase/internal/pkg/response"
	"github.com/mujhtech/s3ase/internal/pkg/sse"
	"github.com/mujhtech/s3ase/job"
	jobHandlers "github.com/mujhtech/s3ase/job/handlers"
	"github.com/mujhtech/s3ase/services"
	"github.com/rs/zerolog"
)

const (
	FileParamID = "file_param_id"
)

func getFileIdFromPath(r *http.Request) (string, error) {
	rawRef, err := pathParamOrError(r, FileParamID)
	if err != nil {
		return "", err
	}

	return url.PathUnescape(rawRef)
}

func getFilesQueryParams(r *http.Request) *dto.FileQueryDto {
	folderId, _ := queryParam(r, "folder_id")
	search, _ := queryParam(r, "search")

	page := ParsePage(r)

	perPage := ParsePerPage(r)

	return &dto.FileQueryDto{
		FolderID: folderId,
		Search:   search,
		Page:     page,
		PerPage:  perPage,
	}
}

func getCreateFileQuery(r *http.Request) (*dto.CreateFileRequestDto, error) {
	folderId, _ := queryParam(r, "folder_id")
	name, err := queryParamOrError(r, "name")

	if err != nil {
		return nil, err
	}

	return &dto.CreateFileRequestDto{
		FolderID: folderId,
		Name:     name,
	}, nil
}

func (h *Handler) GetFiles(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	session, app, err := middleware.AuthSessionAndAppFrom(ctx)

	if err != nil {
		_ = response.Unauthorized(w, r, err)
		return
	}

	query := getFilesQueryParams(r)

	findFilesService := services.FindFilesService{
		App:        app,
		FolderRepo: h.store.FolderRepo,
		FileRepo:   h.store.FileRepo,
		User:       session.User,
		Query:      query,
	}

	files, err := findFilesService.Run(ctx)

	if err != nil {
		_ = response.Error(w, r, err)
		return
	}

	_ = response.Ok(w, r, "files retrieved", files)
}

func (h *Handler) GetFile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	session, app, err := middleware.AuthSessionAndAppFrom(ctx)

	if err != nil {
		_ = response.Unauthorized(w, r, err)
		return
	}

	fileId, err := getFileIdFromPath(r)

	if err != nil {
		_ = response.BadRequest(w, r, err)
		return
	}

	findFileService := services.FindFileService{
		FileID:     fileId,
		App:        app,
		FolderRepo: h.store.FolderRepo,
		FileRepo:   h.store.FileRepo,
		User:       session.User,
	}

	file, err := findFileService.Run(ctx)

	if err != nil {
		_ = response.Error(w, r, err)
		return
	}

	_ = response.Ok(w, r, "file retrieved", file)
}

func (h *Handler) UploadFile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	session, app, err := middleware.AuthSessionAndAppFrom(ctx)

	if err != nil {
		_ = response.Unauthorized(w, r, err)
		return
	}

	dst, err := getCreateFileQuery(r)

	if err != nil {
		_ = response.BadRequest(w, r, err)
		return
	}

	createFileService := services.CreateFileService{
		App:        app,
		FolderRepo: h.store.FolderRepo,
		FileRepo:   h.store.FileRepo,
		User:       session.User,
		Body:       dst,
	}

	h.protocol.UploadFile(createFileService, w, r)

	// id := uuid.New().String()
	// name := "test.txt"
	// progress := 0

	// // Create background context for the goroutine
	// bgCtx := context.Background()

	// if err = h.sse.Publish(ctx, app.ID, sse.EventTypeUploadStarted, UploadProgress{
	// 	FileID:   id,
	// 	Name:     name,
	// 	Status:   UploadProgressStatusStarted,
	// 	Progress: progress,
	// }); err != nil {
	// 	log.Printf("failed to publish upload started event: %v", err)
	// }

	// Use background context in goroutine
	// go func(ctx context.Context, appID string) {
	// 	for progress < 100 {
	// 		progress += 5
	// 		if err = h.sse.Publish(ctx, appID, sse.EventTypeUploadProgress, UploadProgress{
	// 			FileID:   id,
	// 			Name:     name,
	// 			Status:   UploadProgressStatusUploading,
	// 			Progress: progress,
	// 		}); err != nil {
	// 			log.Printf("failed to publish upload progress event: %v", err)
	// 		}
	// 		time.Sleep(1 * time.Second)
	// 		if progress == 100 {
	// 			if err = h.sse.Publish(ctx, appID, sse.EventTypeUploadCompleted, UploadProgress{
	// 				FileID:   id,
	// 				Name:     name,
	// 				Status:   UploadProgressStatusCompleted,
	// 				Progress: progress,
	// 			}); err != nil {
	// 				log.Printf("failed to publish upload completed event: %v", err)
	// 			}
	// 		}
	// 	}
	// }(bgCtx, app.ID)

	//_ = response.Ok(w, r, "file uploaded", nil)
}

func (h *Handler) TusUploadFile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	session, app, err := middleware.AuthSessionAndAppFrom(ctx)
	if err != nil {
		_ = response.Unauthorized(w, r, err)
		return
	}

	body := &dto.CreateFileRequestDto{}
	if r.Method == http.MethodPost {
		body, err = getCreateFileQuery(r)
		if err != nil {
			_ = response.BadRequest(w, r, err)
			return
		}
	}

	createFileService := services.CreateFileService{
		App: app, FolderRepo: h.store.FolderRepo, FileRepo: h.store.FileRepo,
		User: session.User, Body: body,
	}
	basePath := "/api/ui/files/uploads"
	if strings.Contains(r.URL.Path, "/api/v1/") {
		basePath = "/api/v1/files/uploads"
	}
	h.protocol.TusUpload(createFileService, basePath, w, r)
}

func (h *Handler) DownloadFile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	session, app, err := middleware.AuthSessionAndAppFrom(ctx)
	if err != nil {
		_ = response.Unauthorized(w, r, err)
		return
	}
	fileID, err := getFileIdFromPath(r)
	if err != nil {
		_ = response.BadRequest(w, r, err)
		return
	}
	file, err := (&services.FindFileService{
		FileID: fileID, App: app, FileRepo: h.store.FileRepo, User: session.User,
	}).Run(ctx)
	if err != nil {
		_ = response.Error(w, r, err)
		return
	}
	if file.Status != models.FileStatusCompleted {
		_ = response.BadRequest(w, r, fmt.Errorf("file is not available for download"))
		return
	}
	object, err := h.s3.OpenFile(ctx, app.Bucket, file.ObjectKey(), app.Region.String)
	if err != nil {
		_ = response.Error(w, r, err)
		return
	}
	defer func() {
		if closeErr := object.Body.Close(); closeErr != nil {
			zerolog.Ctx(ctx).Warn().Err(closeErr).Str("file_id", file.ID).
				Msg("failed to close downloaded object body")
		}
	}()
	contentType := file.MimeType
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Disposition", mime.FormatMediaType("attachment", map[string]string{"filename": file.Name}))
	if object.ContentLength != nil {
		w.Header().Set("Content-Length", fmt.Sprintf("%d", *object.ContentLength))
	}
	w.WriteHeader(http.StatusOK)
	if _, err := io.Copy(w, object.Body); err != nil {
		return
	}
}

func (h *Handler) DeleteFile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	_, app, err := middleware.AuthSessionAndAppFrom(ctx)
	if err != nil {
		_ = response.Unauthorized(w, r, err)
		return
	}
	fileID, err := getFileIdFromPath(r)
	if err != nil {
		_ = response.BadRequest(w, r, err)
		return
	}
	file, err := (&services.DeleteFileService{
		App: app, FileID: fileID, FileRepo: h.store.FileRepo, S3: h.s3,
	}).Run(ctx)
	if err != nil {
		_ = response.Error(w, r, err)
		return
	}
	if err := h.sse.Publish(ctx, app.ID, sse.EventTypeUploadDeleted, sse.UploadProgress{
		FileID: file.ID, Name: file.Name, Status: sse.UploadProgressStatusCancelled,
	}); err != nil {
		zerolog.Ctx(ctx).Warn().Err(err).Str("file_id", file.ID).
			Msg("file deleted but upload-deleted event could not be published")
	}
	if data, marshalErr := json.Marshal(file); marshalErr == nil {
		if payload, marshalErr := json.Marshal(jobHandlers.WebhookPayload{
			ID: uuid.NewString(), AppID: app.ID, Event: "upload.deleted", CreatedAt: time.Now(), Data: data,
		}); marshalErr == nil {
			_ = h.job.Client.Enqueue(job.QueueNameDefault, job.JobNameWebhook, &job.ClientPayload{Data: payload})
		}
	}
	_ = response.Ok(w, r, "file deleted", nil)
}
