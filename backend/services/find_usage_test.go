package services

import (
	"context"
	"testing"
	"time"

	"github.com/mujhtech/s3ase/database/models"
	"github.com/mujhtech/s3ase/mocks"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestFindUsageServiceRun(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := mocks.NewMockFileRepository(ctrl)
	now := time.Date(2026, time.July, 11, 12, 0, 0, 0, time.UTC)
	repo.EXPECT().FindFilesByAppID(gomock.Any(), "app-id").Return([]*models.File{
		{Status: models.FileStatusCompleted, Size: 100, CreatedAt: now},
		{Status: models.FileStatusCompleted, Size: 50, CreatedAt: now.AddDate(0, 0, -1)},
		{Status: models.FileStatusFailed, Size: 999, CreatedAt: now},
	}, nil)

	usage, err := (&FindUsageService{
		App: &models.App{ID: "app-id"}, FileRepo: repo, Now: now,
	}).Run(context.Background())

	require.NoError(t, err)
	require.Equal(t, int64(2), usage.FileCount)
	require.Equal(t, int64(150), usage.StorageBytes)
	require.Len(t, usage.Daily, 30)
	require.Equal(t, int64(1), usage.Daily[29].Uploads)
	require.Equal(t, int64(100), usage.Daily[29].StorageBytes)
}
