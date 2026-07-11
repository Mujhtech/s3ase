package services

import (
	"context"
	"testing"

	"github.com/guregu/null"
	"github.com/mujhtech/s3ase/database/models"
	"github.com/mujhtech/s3ase/mocks"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

type recordingFileStore struct {
	bucket string
	key    string
	region string
}

func (s *recordingFileStore) DeleteFile(_ context.Context, bucket, key, region string) error {
	s.bucket, s.key, s.region = bucket, key, region
	return nil
}

func TestDeleteFileServiceRun(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := mocks.NewMockFileRepository(ctrl)
	file := &models.File{
		ID: "file-id", AppID: "app-id", PublicID: "object-id",
		FolderID: null.NewString("folder-id", true),
	}
	repo.EXPECT().FindFileByID(gomock.Any(), "file-id").Return(file, nil)
	repo.EXPECT().DeleteFile(gomock.Any(), "file-id").Return(nil)
	objects := &recordingFileStore{}

	deleted, err := (&DeleteFileService{
		App: &models.App{
			ID: "app-id", Bucket: "bucket", Region: null.NewString("eu-west-2", true),
		},
		FileID: "file-id", FileRepo: repo, S3: objects,
	}).Run(context.Background())

	require.NoError(t, err)
	require.Equal(t, file, deleted)
	require.Equal(t, "bucket", objects.bucket)
	require.Equal(t, "folder-id/object-id", objects.key)
	require.Equal(t, "eu-west-2", objects.region)
}
