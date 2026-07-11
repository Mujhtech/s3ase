package services

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/guregu/null"
	"github.com/mujhtech/s3ase/api/dto"
	"github.com/mujhtech/s3ase/database/models"
	"github.com/mujhtech/s3ase/mocks"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestCreateAppService_Run(t *testing.T) {

	type args struct {
		ctx  context.Context
		User *models.User
		Body *dto.CreateAppRequestDto
	}
	tests := []struct {
		name    string
		args    args
		mockFn  func(s *CreateAppService)
		want    *models.App
		wantErr bool
	}{
		{
			name: "should_create_app_success",
			args: args{
				ctx: context.Background(),
				User: &models.User{
					ID:                   "1",
					Name:                 "test",
					Email:                "test@test.com",
					EmailVerified:        true,
					AuthenticationMethod: "github",
				},
				Body: &dto.CreateAppRequestDto{
					Name:        "test",
					Description: "test",
					Region:      "us-east-1",
					Slug:        "test",
				},
			},
			want: &models.App{
				Name:        "test",
				OwnerID:     "1",
				Slug:        "test",
				Region:      null.NewString("us-east-1", true),
				Description: null.NewString("test", true),
				Bucket:      "test",
				Metadata:    models.Metadata{},
			},
			wantErr: false,
			mockFn: func(s *CreateAppService) {
				a, _ := s.AppRepo.(*mocks.MockAppRepository)
				a.EXPECT().CreateApp(gomock.Any(), gomock.Any()).Times(1).Return(nil)

				am, _ := s.AppMemberRepo.(*mocks.MockAppMemberRepository)
				am.EXPECT().CreateAppMember(gomock.Any(), gomock.Any()).Times(1).Return(nil)
			},
		},
		{
			name: "should_return_error_on_app_create_failure",
			args: args{
				ctx: context.Background(),
				User: &models.User{
					ID:                   "1",
					Name:                 "error",
					Email:                "error@test.com",
					EmailVerified:        false,
					AuthenticationMethod: "github",
				},
				Body: &dto.CreateAppRequestDto{
					Name:        "fail_app",
					Description: "this should fail",
					Region:      "us-east-1",
					Slug:        "fail-app",
				},
			},
			want:    nil,
			wantErr: true,
			mockFn: func(s *CreateAppService) {
				a, _ := s.AppRepo.(*mocks.MockAppRepository)
				a.EXPECT().
					CreateApp(gomock.Any(), gomock.Any()).
					Times(1).
					Return(fmt.Errorf("mock create app error"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			s := CreateAppService{
				AppRepo:       mocks.NewMockAppRepository(ctrl),
				AppMemberRepo: mocks.NewMockAppMemberRepository(ctrl),
				User:          tt.args.User,
				DefaultRegion: "us-east-1",
				Body:          tt.args.Body,
				S3:            successfulBucketStore{},
			}

			if tt.mockFn != nil {
				tt.mockFn(&s)
			}

			got, err := s.Run(tt.args.ctx)

			if tt.wantErr {
				require.NotNil(t, err)
				return
			}

			require.Nil(t, err)
			require.NotEmpty(t, got.Name)
			require.NotEmpty(t, got.Slug)
			require.NotEmpty(t, got.OwnerID)

			// set the time to zero value
			got.CreatedAt = time.Time{}
			got.UpdatedAt = time.Time{}
			got.ID = ""

			require.Equal(t, tt.want, got)
		})
	}

}

func TestSlugifyProducesS3SafeName(t *testing.T) {
	require.Equal(t, "my-production-app", slugify("  My Production App!  "))
	require.Equal(t, "a-b-c", slugify("A---B___C"))
}

func TestCreateAppServiceForcesConfiguredStorageRegion(t *testing.T) {
	ctrl := gomock.NewController(t)
	appRepo := mocks.NewMockAppRepository(ctrl)
	memberRepo := mocks.NewMockAppMemberRepository(ctrl)
	buckets := &recordingBucketStore{}
	appRepo.EXPECT().CreateApp(gomock.Any(), gomock.Any()).Return(nil)
	memberRepo.EXPECT().CreateAppMember(gomock.Any(), gomock.Any()).Return(nil)

	app, err := (&CreateAppService{
		Body:               &dto.CreateAppRequestDto{Name: "R2 App", Region: "eu-west-2"},
		DefaultRegion:      "auto",
		ForceDefaultRegion: true,
		AppRepo:            appRepo,
		AppMemberRepo:      memberRepo,
		User:               &models.User{ID: "owner-id"},
		S3:                 buckets,
	}).Run(context.Background())

	require.NoError(t, err)
	require.Equal(t, "auto", buckets.region)
	require.Equal(t, "auto", app.Region.String)
}

type successfulBucketStore struct{}

func (successfulBucketStore) CheckOrCreateNewBucket(context.Context, string, string) (string, error) {
	return "", nil
}

type recordingBucketStore struct {
	region string
}

func (s *recordingBucketStore) CheckOrCreateNewBucket(_ context.Context, _ string, region string) (string, error) {
	s.region = region
	return region, nil
}
