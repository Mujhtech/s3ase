package handler

import (
	"net/http"
	"net/url"

	"github.com/mujhtech/s3ase/api/dto"
	"github.com/mujhtech/s3ase/api/middleware"
	"github.com/mujhtech/s3ase/internal/pkg/request"
	"github.com/mujhtech/s3ase/internal/pkg/response"
	"github.com/mujhtech/s3ase/services"
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

	page := ParsePage(r)

	perPage := ParsePerPage(r)

	return &dto.FileQueryDto{
		FolderID: folderId,
		Page:     page,
		PerPage:  perPage,
	}
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
		_ = response.InternalServerError(w, r, err)
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
		_ = response.InternalServerError(w, r, err)
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

	dst := new(dto.CreateFileRequestDto)

	if err := request.ReadBody(r, dst); err != nil {
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

func (h *Handler) DeleteFile(w http.ResponseWriter, r *http.Request) {}
