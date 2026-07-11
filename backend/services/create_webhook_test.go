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

func TestCreateWebhookService_Run(t *testing.T) {
	type args struct {
		ctx  context.Context
		app  *models.App
		user *models.User
		body *dto.CreateWebhookRequestDto
	}
	type testCase struct {
		name    string
		args    args
		mockFn  func(s *CreateWebhookService)
		want    *models.Webhook
		wantErr error
	}
	tests := []testCase{
		{
			name: "successfully create webhook",
			args: args{
				ctx:  context.Background(),
				app:  &models.App{ID: "app-id"},
				user: &models.User{ID: "user-id"},
				body: &dto.CreateWebhookRequestDto{
					Name:        "test-webhook",
					Description: "test description",
					Url:         "https://test.com/webhook",
					Events:      []string{"file.created", "file.updated"},
				},
			},
			mockFn: func(s *CreateWebhookService) {
				memberRepo, _ := s.AppMemberRepo.(*mocks.MockAppMemberRepository)
				memberRepo.EXPECT().FindAppMemberByAppIDAndUserId(gomock.Any(), "app-id", "user-id").
					Return(&models.AppMember{Role: models.AppMemberRoleOwner}, nil)
				webhookRepo, _ := s.WebhookRepo.(*mocks.MockWebhookRepository)
				webhookRepo.EXPECT().
					CreateWebhook(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil)
				webhookRepo.EXPECT().
					FindWebhookByID(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&models.Webhook{
						ID:          "webhook-id",
						Name:        "test-webhook",
						Description: null.NewString("test description", true),
						URL:         "https://test.com/webhook",
					}, nil)
			},
			want: &models.Webhook{
				ID:          "webhook-id",
				Name:        "test-webhook",
				Description: null.NewString("test description", true),
				URL:         "https://test.com/webhook",
			},
			wantErr: nil,
		},
		{
			name: "fail to create webhook",
			args: args{
				ctx:  context.Background(),
				app:  &models.App{ID: "app-id"},
				user: &models.User{ID: "user-id"},
				body: &dto.CreateWebhookRequestDto{
					Name:        "test-webhook",
					Description: "test description",
					Url:         "https://test.com/webhook",
					Events:      []string{"file.created"},
				},
			},
			mockFn: func(s *CreateWebhookService) {
				memberRepo, _ := s.AppMemberRepo.(*mocks.MockAppMemberRepository)
				memberRepo.EXPECT().FindAppMemberByAppIDAndUserId(gomock.Any(), "app-id", "user-id").
					Return(&models.AppMember{Role: models.AppMemberRoleOwner}, nil)
				webhookRepo, _ := s.WebhookRepo.(*mocks.MockWebhookRepository)
				webhookRepo.EXPECT().
					CreateWebhook(gomock.Any(), gomock.Any()).
					Times(1).
					Return(errors.New("failed to create webhook"))
			},
			want:    nil,
			wantErr: errors.New("failed to create webhook"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			service := &CreateWebhookService{
				App:           tt.args.app,
				AppMemberRepo: mocks.NewMockAppMemberRepository(ctrl),
				WebhookRepo:   mocks.NewMockWebhookRepository(ctrl),
				User:          tt.args.user,
				Body:          tt.args.body,
			}

			if tt.mockFn != nil {
				tt.mockFn(service)
			}

			got, err := service.Run(tt.args.ctx)
			if tt.wantErr != nil {
				require.Error(t, err)
				require.Equal(t, tt.wantErr, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
			}
		})
	}
}
