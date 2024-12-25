package services

import (
	"context"
	"errors"
	"testing"

	"github.com/guregu/null"
	"github.com/mujhtech/s3ase/api/dto"
	"github.com/mujhtech/s3ase/database/models"
	"github.com/mujhtech/s3ase/mocks"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestUpdateAppService_Run(t *testing.T) {
	type args struct {
		ctx   context.Context
		appId string
		user  *models.User
		body  *dto.CreateAppRequestDto
	}
	type testCase struct {
		name    string
		args    args
		mockFn  func(s *UpdateAppService)
		wantErr error
	}
	tests := []testCase{
		{
			name: "successfully update app",
			args: args{
				ctx:   context.Background(),
				appId: "app-id",
				user:  &models.User{ID: "user-id"},
				body: &dto.CreateAppRequestDto{
					Name:        "updated-app",
					Description: "updated description",
				},
			},
			mockFn: func(s *UpdateAppService) {
				appRepo, _ := s.AppRepo.(*mocks.MockAppRepository)
				appRepo.EXPECT().
					FindAppByID(gomock.Any(), "app-id").
					Return(&models.App{
						ID:      "app-id",
						OwnerID: "user-id",
					}, nil)
				appRepo.EXPECT().
					UpdateApp(gomock.Any(), &models.App{
						ID:          "app-id",
						OwnerID:     "user-id",
						Name:        "updated-app",
						Description: null.NewString("updated description", true),
					}).
					Return(nil)
			},
			wantErr: nil,
		},
		{
			name: "app not found",
			args: args{
				ctx:   context.Background(),
				appId: "app-id",
				user:  &models.User{ID: "user-id"},
				body:  &dto.CreateAppRequestDto{},
			},
			mockFn: func(s *UpdateAppService) {
				appRepo, _ := s.AppRepo.(*mocks.MockAppRepository)
				appRepo.EXPECT().
					FindAppByID(gomock.Any(), "app-id").
					Return(nil, errors.New("app not found"))
			},
			wantErr: errors.New("app not found"),
		},
		{
			name: "user not owner",
			args: args{
				ctx:   context.Background(),
				appId: "app-id",
				user:  &models.User{ID: "user-id"},
				body:  &dto.CreateAppRequestDto{},
			},
			mockFn: func(s *UpdateAppService) {
				appRepo, _ := s.AppRepo.(*mocks.MockAppRepository)
				appRepo.EXPECT().
					FindAppByID(gomock.Any(), "app-id").
					Return(&models.App{
						ID:      "app-id",
						OwnerID: "different-user-id",
					}, nil)
			},
			wantErr: errors.New("user is not owner of the app"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			service := &UpdateAppService{
				AppId:   tt.args.appId,
				AppRepo: mocks.NewMockAppRepository(ctrl),
				User:    tt.args.user,
				Body:    tt.args.body,
			}

			if tt.mockFn != nil {
				tt.mockFn(service)
			}

			err := service.Run(tt.args.ctx)
			if tt.wantErr != nil {
				require.Error(t, err)
				require.Equal(t, tt.wantErr.Error(), err.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}
