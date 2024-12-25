package handler

import (
	"net/http"
	"net/url"

	"github.com/mujhtech/s3ase/api/dto"
	"github.com/mujhtech/s3ase/api/middleware"
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
