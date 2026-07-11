package services

import (
	"context"
	"testing"

	"github.com/mujhtech/s3ase/database/models"
	"github.com/mujhtech/s3ase/errors"
	"github.com/mujhtech/s3ase/mocks"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestDeleteApiKeyService_Run(t *testing.T) {
	type args struct {
		ctx      context.Context
		app      *models.App
		user     *models.User
		apiKeyId string
		role     models.AppMemberRole
	}
	type testCase struct {
		name    string
		args    args
		mockFn  func(s *DeleteApiKeyService)
		wantErr error
	}
	tests := []testCase{
		{
			name: "owner can delete api key",
			args: args{
				ctx:      context.Background(),
				app:      &models.App{ID: "app-id"},
				user:     &models.User{ID: "user-id"},
				apiKeyId: "key-id",
				role:     models.AppMemberRoleOwner,
			},
			mockFn: func(s *DeleteApiKeyService) {
				memberRepo, _ := s.AppMemberRepo.(*mocks.MockAppMemberRepository)
				memberRepo.EXPECT().FindAppMemberByAppIDAndUserId(gomock.Any(), "app-id", "user-id").
					Return(&models.AppMember{Role: models.AppMemberRoleOwner}, nil)
				apiKey, _ := s.ApiKeyRepo.(*mocks.MockApiKeyRepository)
				apiKey.EXPECT().FindApiKeyByID(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&models.ApiKey{ID: "key-id", AppID: "app-id", Name: "original-name"}, nil)
				apiKey.EXPECT().
					DeleteApiKey(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil)
			},
			wantErr: nil,
		},
		{
			name: "member cannot delete api key",
			args: args{
				ctx:      context.Background(),
				app:      &models.App{ID: "app-id"},
				user:     &models.User{ID: "user-id"},
				apiKeyId: "key-id",
				role:     models.AppMemberRoleMember,
			},
			mockFn: func(s *DeleteApiKeyService) {
				memberRepo, _ := s.AppMemberRepo.(*mocks.MockAppMemberRepository)
				memberRepo.EXPECT().FindAppMemberByAppIDAndUserId(gomock.Any(), "app-id", "user-id").
					Return(&models.AppMember{Role: models.AppMemberRoleMember}, nil)
			},
			wantErr: errors.ErrNotAuthorized,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			service := &DeleteApiKeyService{
				App:           tt.args.app,
				AppMemberRepo: mocks.NewMockAppMemberRepository(ctrl),
				ApiKeyRepo:    mocks.NewMockApiKeyRepository(ctrl),
				User:          tt.args.user,
				ApiKeyId:      tt.args.apiKeyId,
			}

			if tt.mockFn != nil {
				tt.mockFn(service)
			}
			err := service.Run(tt.args.ctx)
			if tt.wantErr != nil {
				require.Error(t, err)
				require.Equal(t, tt.wantErr, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
