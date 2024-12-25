package services

import (
	"context"
	"errors"
	"testing"

	"github.com/mujhtech/s3ase/database/models"
	"github.com/mujhtech/s3ase/mocks"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestFindAppsService_Run(t *testing.T) {
	type args struct {
		ctx  context.Context
		user *models.User
	}
	type testCase struct {
		name    string
		args    args
		mockFn  func(s *FindAppsService)
		want    []*models.App
		wantErr error
	}
	tests := []testCase{
		{
			name: "successfully find apps",
			args: args{
				ctx:  context.Background(),
				user: &models.User{ID: "user-id"},
			},
			mockFn: func(s *FindAppsService) {
				appRepo, _ := s.AppRepo.(*mocks.MockAppRepository)
				appRepo.EXPECT().
					FindAppsByUserID(gomock.Any(), "user-id").
					Return([]*models.App{
						{ID: "app-1", Name: "App 1"},
						{ID: "app-2", Name: "App 2"},
					}, nil)
			},
			want: []*models.App{
				{ID: "app-1", Name: "App 1"},
				{ID: "app-2", Name: "App 2"},
			},
			wantErr: nil,
		},
		{
			name: "error finding apps",
			args: args{
				ctx:  context.Background(),
				user: &models.User{ID: "user-id"},
			},
			mockFn: func(s *FindAppsService) {
				appRepo, _ := s.AppRepo.(*mocks.MockAppRepository)
				appRepo.EXPECT().
					FindAppsByUserID(gomock.Any(), "user-id").
					Return(nil, errors.New("failed to find apps"))
			},
			want:    nil,
			wantErr: errors.New("failed to find apps"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			service := &FindAppsService{
				AppRepo: mocks.NewMockAppRepository(ctrl),
				User:    tt.args.user,
			}

			if tt.mockFn != nil {
				tt.mockFn(service)
			}

			got, err := service.Run(tt.args.ctx)
			if tt.wantErr != nil {
				require.Error(t, err)
				require.Equal(t, tt.wantErr.Error(), err.Error())
				require.Nil(t, got)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
			}
		})
	}
}
