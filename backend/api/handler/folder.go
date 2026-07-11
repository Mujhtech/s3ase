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
	FolderParamID = "folder_param_id"
)

func getFolderIdFromPath(r *http.Request) (string, error) {
	rawRef, err := pathParamOrError(r, FolderParamID)
	if err != nil {
		return "", err
	}

	return url.PathUnescape(rawRef)
}

func (h *Handler) GetFolders(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	session, app, err := middleware.AuthSessionAndAppFrom(ctx)

	if err != nil {
		_ = response.Unauthorized(w, r, err)
		return
	}

	findFoldersService := services.FindFoldersService{
		App:        app,
		FolderRepo: h.store.FolderRepo,
		FileRepo:   h.store.FileRepo,
		User:       session.User,
	}

	folders, err := findFoldersService.Run(ctx)

	if err != nil {
		_ = response.Error(w, r, err)
		return
	}

	_ = response.Ok(w, r, "folders retrieved", folders)
}

func (h *Handler) GetFolder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	session, app, err := middleware.AuthSessionAndAppFrom(ctx)

	if err != nil {
		_ = response.Unauthorized(w, r, err)
		return
	}

	folderId, err := getFolderIdFromPath(r)

	if err != nil {
		_ = response.BadRequest(w, r, err)
		return
	}

	findFolderService := services.FindFolderService{
		FolderId:   folderId,
		App:        app,
		FolderRepo: h.store.FolderRepo,
		FileRepo:   h.store.FileRepo,
		User:       session.User,
	}

	folder, err := findFolderService.Run(ctx)

	if err != nil {
		_ = response.Error(w, r, err)
		return
	}

	_ = response.Ok(w, r, "folder retrieved", folder)
}

func (h *Handler) CreateFolder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	session, app, err := middleware.AuthSessionAndAppFrom(ctx)

	if err != nil {
		_ = response.Unauthorized(w, r, err)
		return
	}

	dst := new(dto.CreateFolderRequestDto)

	if err := request.ReadBody(r, dst); err != nil {
		_ = response.BadRequest(w, r, err)
		return
	}

	createFolderService := services.CreateFolderService{
		App:        app,
		FolderRepo: h.store.FolderRepo,
		FileRepo:   h.store.FileRepo,
		User:       session.User,
		Body:       dst,
	}

	folder, err := createFolderService.Run(ctx)

	if err != nil {
		_ = response.Error(w, r, err)
		return
	}

	_ = response.Created(w, r, "folder created", folder)
}

func (h *Handler) UpdateFolder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	session, app, err := middleware.AuthSessionAndAppFrom(ctx)

	if err != nil {
		_ = response.Unauthorized(w, r, err)
		return
	}

	folderId, err := getFolderIdFromPath(r)

	if err != nil {
		_ = response.BadRequest(w, r, err)
		return
	}

	dst := new(dto.CreateFolderRequestDto)

	if err := request.ReadBody(r, dst); err != nil {
		_ = response.BadRequest(w, r, err)
		return
	}

	updateFolderService := services.UpdateFolderService{
		App:        app,
		FolderId:   folderId,
		FolderRepo: h.store.FolderRepo,
		FileRepo:   h.store.FileRepo,
		User:       session.User,
		Body:       dst,
	}

	if err = updateFolderService.Run(ctx); err != nil {
		_ = response.Error(w, r, err)
		return
	}

	_ = response.Ok(w, r, "folder updated", nil)
}

func (h *Handler) DeleteFolder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	session, app, err := middleware.AuthSessionAndAppFrom(ctx)

	if err != nil {
		_ = response.Unauthorized(w, r, err)
		return
	}

	folderId, err := getFolderIdFromPath(r)

	if err != nil {
		_ = response.BadRequest(w, r, err)
		return
	}

	deleteFolderService := services.DeleteFolderService{
		App:        app,
		FolderId:   folderId,
		FolderRepo: h.store.FolderRepo,
		FileRepo:   h.store.FileRepo,
		User:       session.User,
	}

	if err = deleteFolderService.Run(ctx); err != nil {
		_ = response.Error(w, r, err)
		return
	}

	_ = response.Ok(w, r, "folder deleted", nil)
}
