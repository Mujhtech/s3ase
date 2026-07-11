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

type fakeAppObjectStore struct {
	bucket string
	region string
}

func (f *fakeAppObjectStore) DeleteBucket(_ context.Context, bucket, region string) error {
	f.bucket, f.region = bucket, region
	return nil
}

func TestDeleteAppServiceRun(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := mocks.NewMockAppRepository(ctrl)
	repo.EXPECT().FindAppByID(gomock.Any(), "app-id").Return(&models.App{
		ID: "app-id", OwnerID: "owner-id", Name: "My App", Bucket: "my-app",
		Region: null.StringFrom("eu-west-2"),
	}, nil)
	repo.EXPECT().DeleteApp(gomock.Any(), "app-id").Return(nil)
	objects := new(fakeAppObjectStore)

	err := (&DeleteAppService{
		AppID: "app-id", ConfirmName: "My App", AppRepo: repo,
		User: &models.User{ID: "owner-id"}, S3: objects,
	}).Run(context.Background())

	require.NoError(t, err)
	require.Equal(t, "my-app", objects.bucket)
	require.Equal(t, "eu-west-2", objects.region)
}
