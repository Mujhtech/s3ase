package services

import (
	"context"
	"fmt"
	"testing"

	"github.com/mujhtech/s3ase/api/dto"
	"github.com/mujhtech/s3ase/database/models"
	"github.com/mujhtech/s3ase/errors"
	"github.com/mujhtech/s3ase/mocks"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestUpdateApiKeyService_Run(t *testing.T) {
	type args struct {
		ctx      context.Context
		app      *models.App
		user     *models.User
		req      *dto.CreateApiKeyRequestDto
		role     models.AppMemberRole
		apiKeyId string
	}
	type testCase struct {
		name    string
		args    args
		mockFn  func(s *UpdateApiKeyService)
		wantErr error
	}
	tests := []testCase{
		{
			name: "owner can update api key",
			args: args{
				ctx:      context.Background(),
				app:      &models.App{ID: "app-id"},
				user:     &models.User{ID: "user-id"},
				req:      &dto.CreateApiKeyRequestDto{Name: "updated-name"},
				role:     models.AppMemberRoleOwner,
				apiKeyId: "key-id",
			},
			mockFn: func(s *UpdateApiKeyService) {
				apiKey, _ := s.ApiKeyRepo.(*mocks.MockApiKeyRepository)
				apiKey.EXPECT().FindApiKeyByID(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&models.ApiKey{ID: "key-id", Name: "original-name"}, nil)
				apiKey.EXPECT().
					UpdateApiKey(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil)

			},
			wantErr: nil,
		},
		{
			name: "member cannot update api key",
			args: args{
				ctx:      context.Background(),
				app:      &models.App{ID: "app-id"},
				user:     &models.User{ID: "user-id"},
				req:      &dto.CreateApiKeyRequestDto{Name: "updated-name"},
				role:     models.AppMemberRoleMember,
				apiKeyId: "key-id",
			},
			mockFn: func(s *UpdateApiKeyService) {
				apiKey, _ := s.ApiKeyRepo.(*mocks.MockApiKeyRepository)
				apiKey.EXPECT().FindApiKeyByID(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&models.ApiKey{ID: "key-id", Name: "original-name"}, nil)
				apiKey.EXPECT().
					UpdateApiKey(gomock.Any(), gomock.Any()).
					Times(1).
					Return(fmt.Errorf("not authorized"))
			},
			wantErr: errors.ErrNotAuthorized,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			service := &UpdateApiKeyService{
				App:        tt.args.app,
				ApiKeyRepo: mocks.NewMockApiKeyRepository(ctrl),
				User:       tt.args.user,
				Body:       tt.args.req,
				ApiKeyId:   tt.args.apiKeyId,
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
