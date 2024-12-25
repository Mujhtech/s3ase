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

func TestFindApiKeysService_Run(t *testing.T) {
	type args struct {
		ctx  context.Context
		app  *models.App
		user *models.User
	}
	type testCase struct {
		name    string
		args    args
		mockFn  func(s *FindApiKeysService)
		want    []*models.ApiKey
		wantErr error
	}
	tests := []testCase{
		{
			name: "successfully find api keys",
			args: args{
				ctx:  context.Background(),
				app:  &models.App{ID: "app-id"},
				user: &models.User{ID: "user-id"},
			},
			mockFn: func(s *FindApiKeysService) {
				apiKeyRepo, _ := s.ApiKeyRepo.(*mocks.MockApiKeyRepository)
				apiKeyRepo.EXPECT().
					FindApiKeysByAppID(gomock.Any(), "app-id").
					Return([]*models.ApiKey{
						{ID: "key-1", Name: "API Key 1"},
						{ID: "key-2", Name: "API Key 2"},
					}, nil)
			},
			want: []*models.ApiKey{
				{ID: "key-1", Name: "API Key 1"},
				{ID: "key-2", Name: "API Key 2"},
			},
			wantErr: nil,
		},
		{
			name: "error finding api keys",
			args: args{
				ctx:  context.Background(),
				app:  &models.App{ID: "app-id"},
				user: &models.User{ID: "user-id"},
			},
			mockFn: func(s *FindApiKeysService) {
				apiKeyRepo, _ := s.ApiKeyRepo.(*mocks.MockApiKeyRepository)
				apiKeyRepo.EXPECT().
					FindApiKeysByAppID(gomock.Any(), "app-id").
					Return(nil, errors.New("failed to find api keys"))
			},
			want:    nil,
			wantErr: errors.New("failed to find api keys"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			service := &FindApiKeysService{
				App:        tt.args.app,
				ApiKeyRepo: mocks.NewMockApiKeyRepository(ctrl),
				AppRepo:    mocks.NewMockAppRepository(ctrl),
				User:       tt.args.user,
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
