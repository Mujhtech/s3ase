package services

import (
	"context"
	"testing"

	"github.com/mujhtech/s3ase/api/dto"
	"github.com/mujhtech/s3ase/database/models"
	"github.com/mujhtech/s3ase/errors"
	"github.com/mujhtech/s3ase/mocks"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestCreateApiKeyService_Run(t *testing.T) {
	type args struct {
		ctx    context.Context
		app    *models.App
		user   *models.User
		apiKey *dto.CreateApiKeyRequestDto
		role   models.AppMemberRole
	}

	type testCase struct {
		name    string
		args    args
		mockFn  func(s *CreateApiKeyService)
		wantErr error
	}

	tests := []testCase{
		{
			name: "user is owner, should create api key",
			args: args{
				ctx:  context.Background(),
				app:  &models.App{ID: "app-id"},
				user: &models.User{ID: "user-id"},
				apiKey: &dto.CreateApiKeyRequestDto{
					Name:        "test-api-key",
					Description: "test description",
					Access:      "full",
					ExpiredAt:   0,
				},
				role: models.AppMemberRoleOwner,
			},
			mockFn: func(s *CreateApiKeyService) {

				am, _ := s.AppMemberRepo.(*mocks.MockAppMemberRepository)
				am.EXPECT().
					FindAppMemberByAppIDAndUserId(gomock.Any(), "app-id", "user-id").
					Return(&models.AppMember{Role: models.AppMemberRoleOwner}, nil)

				apiKey, _ := s.ApiKeyRepo.(*mocks.MockApiKeyRepository)
				apiKey.EXPECT().
					CreateApiKey(gomock.Any(), gomock.Any()).
					Return(nil)
			},
			wantErr: nil,
		},
		{
			name: "user is not owner, should return error",
			args: args{
				ctx:  context.Background(),
				app:  &models.App{ID: "app-id"},
				user: &models.User{ID: "user-id"},
				apiKey: &dto.CreateApiKeyRequestDto{
					Name:        "test-api-key",
					Description: "test description",
					Access:      "full",
					ExpiredAt:   0,
				},
				role: models.AppMemberRoleMember,
			},
			mockFn: func(s *CreateApiKeyService) {
				am, _ := s.AppMemberRepo.(*mocks.MockAppMemberRepository)
				am.EXPECT().
					FindAppMemberByAppIDAndUserId(gomock.Any(), "app-id", "user-id").
					Return(&models.AppMember{Role: models.AppMemberRoleMember}, nil)
			},
			wantErr: errors.ErrNotAuthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			service := &CreateApiKeyService{
				App:           tt.args.app,
				AppMemberRepo: mocks.NewMockAppMemberRepository(ctrl),
				ApiKeyRepo:    mocks.NewMockApiKeyRepository(ctrl),
				User:          tt.args.user,
				Body:          tt.args.apiKey,
			}

			if tt.mockFn != nil {
				tt.mockFn(service)
			}

			apiKey, err := service.Run(tt.args.ctx)

			if tt.wantErr != nil {
				require.Error(t, err)
				require.Nil(t, apiKey)
				require.Equal(t, tt.wantErr, err)
			} else {
				require.NoError(t, err)
				require.NotNil(t, apiKey)
				require.Equal(t, tt.args.apiKey.Name, apiKey.Name)
			}
		})
	}
}
