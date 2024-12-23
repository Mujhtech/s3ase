package services

import (
	"context"
	"testing"

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
				ID:      "1",
				Name:    "test",
				OwnerID: "1",
				Slug:    "test",
				Region:  null.NewString("us-east-1", true),
				Bucket:  "test",
				Metadata: null.NewString(
					"{}",
					true,
				),
			},
			wantErr: false,
			mockFn: func(s *CreateAppService) {
				a, _ := s.AppRepo.(*mocks.MockAppRepository)
				a.EXPECT().CreateApp(gomock.Any(), gomock.Any()).Times(1).Return(nil)
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
				ApiKeyRepo:    mocks.NewMockApiKeyRepository(ctrl),
				User:          tt.args.User,
				DefaultRegion: "us-east-1",
				Body:          tt.args.Body,
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
			require.NotEmpty(t, got.ID)
			require.NotEmpty(t, got.Name)
			require.NotEmpty(t, got.Slug)
			require.NotEmpty(t, got.OwnerID)

			require.Equal(t, tt.want, got)
		})
	}

}
