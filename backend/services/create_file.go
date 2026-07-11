package services

import (
	"context"
	"errors"
	"fmt"
	"path"
	"regexp"
	"time"

	"strings"

	"github.com/google/uuid"
	"github.com/guregu/null"
	"github.com/mujhtech/s3ase/api/dto"
	"github.com/mujhtech/s3ase/database/models"
	"github.com/mujhtech/s3ase/database/store"
)

type CreateFileService struct {
	App        *models.App
	FolderRepo store.FolderRepository
	FileRepo   store.FileRepository
	User       *models.User
	Body       *dto.CreateFileRequestDto
}

func (c *CreateFileService) Run(ctx context.Context, objectId string, size int64, metadata map[string]string) (*models.File, error) {
	if c.App == nil || c.User == nil || c.Body == nil {
		return nil, fmt.Errorf("missing file upload context")
	}
	c.Body.Name = strings.TrimSpace(c.Body.Name)
	if c.Body.Name == "" || c.Body.Name != path.Base(c.Body.Name) || len(c.Body.Name) > 255 {
		return nil, fmt.Errorf("invalid file name")
	}
	if size < 0 {
		return nil, fmt.Errorf("invalid file size")
	}

	id := uuid.New().String()
	if requestedID := metadata["file_id"]; requestedID != "" {
		parsedID, err := uuid.Parse(requestedID)
		if err != nil {
			return nil, fmt.Errorf("invalid file id")
		}
		id = parsedID.String()
	}

	extension := strings.TrimPrefix(path.Ext(c.Body.Name), ".")
	filename := c.Body.Name
	mimeType := ""
	publicID := objectId

	if metadata["extension"] != "" {
		extension = metadata["extension"]
	}

	if metadata["filename"] != "" {
		filename = metadata["filename"]
	}

	if metadata["filetype"] != "" {
		mimeType = metadata["filetype"]
	}

	if publicID == "" {
		publicID = sanitizeID(id)
	}

	file := &models.File{
		ID:         id,
		AppID:      c.App.ID,
		Name:       filename,
		MimeType:   mimeType,
		FolderID:   null.NewString(c.Body.FolderID, c.Body.FolderID != ""),
		Size:       size,
		PublicID:   publicID,
		UploadedBy: c.User.ID,
		Extension:  extension,
		IsPublic:   true,
		Status:     models.FileStatusStarted,
		Metadata:   map[string]interface{}{},
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if err := c.FileRepo.CreateFile(ctx, file); err != nil {
		return nil, err
	}

	return file, nil
}

func (c *CreateFileService) GetFolder(ctx context.Context) (*models.Folder, error) {
	if c.Body.FolderID == "" {
		return nil, nil
	}
	folder, err := c.FolderRepo.FindFolderByID(ctx, c.Body.FolderID)
	if err != nil && !errors.Is(err, store.ErrNotFound) {
		return nil, err
	}
	if folder != nil {
		if err := requireSameApp(c.App.ID, folder.AppID); err != nil {
			return nil, err
		}
	}
	return folder, nil
}

func sanitizeID(id string) string {
	// Convert to lowercase for consistency
	id = strings.ToLower(id)

	// Replace spaces with hyphens
	id = strings.ReplaceAll(id, " ", "-")

	// Remove special characters except alphanumeric, hyphens and underscores
	reg := regexp.MustCompile(`[^a-z0-9\-_]`)
	id = reg.ReplaceAllString(id, "")

	// Replace multiple hyphens with single hyphen
	reg = regexp.MustCompile(`-+`)
	id = reg.ReplaceAllString(id, "-")

	// Trim hyphens from start and end
	id = strings.Trim(id, "-")

	// Limit length to 64 characters
	if len(id) > 64 {
		id = id[:64]
	}

	return id
}

// Replace existing line with:
