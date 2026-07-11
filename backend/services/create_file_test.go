package services

import (
	"context"
	"testing"

	"github.com/mujhtech/s3ase/api/dto"
	"github.com/mujhtech/s3ase/database/models"
	"github.com/mujhtech/s3ase/mocks"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestCreateFileServiceRun(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := mocks.NewMockFileRepository(ctrl)
	repo.EXPECT().CreateFile(gomock.Any(), gomock.Any()).DoAndReturn(
		func(_ context.Context, file *models.File) error {
			require.Equal(t, "app-id", file.AppID)
			require.Equal(t, "user-id", file.UploadedBy)
			require.Equal(t, "report.pdf", file.Name)
			require.Equal(t, "pdf", file.Extension)
			require.Equal(t, "application/pdf", file.MimeType)
			require.False(t, file.FolderID.Valid)
			require.Equal(t, int64(42), file.Size)
			require.Equal(t, models.FileStatusStarted, file.Status)
			return nil
		},
	)

	service := CreateFileService{
		App:      &models.App{ID: "app-id"},
		User:     &models.User{ID: "user-id"},
		FileRepo: repo,
		Body:     &dto.CreateFileRequestDto{Name: " report.pdf "},
	}
	file, err := service.Run(context.Background(), "", 42, map[string]string{"filetype": "application/pdf"})
	require.NoError(t, err)
	require.NotEmpty(t, file.ID)
	require.Equal(t, file.ID, file.PublicID)
}

func TestCreateFileServiceRejectsUnsafeName(t *testing.T) {
	service := CreateFileService{
		App:  &models.App{ID: "app-id"},
		User: &models.User{ID: "user-id"},
		Body: &dto.CreateFileRequestDto{Name: "../secret.txt"},
	}
	file, err := service.Run(context.Background(), "", 1, nil)
	require.Error(t, err)
	require.Nil(t, file)
}

func TestCreateFileServiceUsesValidClientFileID(t *testing.T) {
	const fileID = "33333333-3333-4333-8333-333333333333"
	ctrl := gomock.NewController(t)
	repo := mocks.NewMockFileRepository(ctrl)
	repo.EXPECT().CreateFile(gomock.Any(), gomock.Any()).DoAndReturn(
		func(_ context.Context, file *models.File) error {
			require.Equal(t, fileID, file.ID)
			require.Equal(t, fileID, file.PublicID)
			return nil
		},
	)

	service := CreateFileService{
		App: &models.App{ID: "app-id"}, User: &models.User{ID: "user-id"},
		FileRepo: repo, Body: &dto.CreateFileRequestDto{Name: "resumable.txt"},
	}
	file, err := service.Run(context.Background(), "", 10, map[string]string{"file_id": fileID})
	require.NoError(t, err)
	require.Equal(t, fileID, file.ID)
}

func TestCreateFileServiceRejectsInvalidClientFileID(t *testing.T) {
	service := CreateFileService{
		App: &models.App{ID: "app-id"}, User: &models.User{ID: "user-id"},
		Body: &dto.CreateFileRequestDto{Name: "resumable.txt"},
	}
	file, err := service.Run(context.Background(), "", 10, map[string]string{"file_id": "not-a-uuid"})
	require.Error(t, err)
	require.Nil(t, file)
}
