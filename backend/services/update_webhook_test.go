package services

import (
	"context"
	"errors"
	"testing"

	"github.com/mujhtech/s3ase/api/dto"
	"github.com/mujhtech/s3ase/database/models"
	"github.com/mujhtech/s3ase/mocks"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestUpdateWebhookService_Run(t *testing.T) {
	type args struct {
		ctx       context.Context
		app       *models.App
		user      *models.User
		webhookId string
		body      *dto.CreateWebhookRequestDto
	}
	type testCase struct {
		name    string
		args    args
		mockFn  func(s *UpdateWebhookService)
		wantErr error
	}
	tests := []testCase{
		{
			name: "successfully update webhook",
			args: args{
				ctx:       context.Background(),
				app:       &models.App{ID: "app-id"},
				user:      &models.User{ID: "user-id"},
				webhookId: "webhook-id",
				body: &dto.CreateWebhookRequestDto{
					Name:        "updated-webhook",
					Description: "updated description",
					Url:         "https://updated.com/webhook",
					Events:      []string{"file.deleted"},
				},
			},
			mockFn: func(s *UpdateWebhookService) {
				memberRepo, _ := s.AppMemberRepo.(*mocks.MockAppMemberRepository)
				memberRepo.EXPECT().FindAppMemberByAppIDAndUserId(gomock.Any(), "app-id", "user-id").
					Return(&models.AppMember{Role: models.AppMemberRoleOwner}, nil)
				webhookRepo, _ := s.WebhookRepo.(*mocks.MockWebhookRepository)
				webhookRepo.EXPECT().
					FindWebhookByID(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&models.Webhook{ID: "webhook-id", AppID: "app-id"}, nil)
				webhookRepo.EXPECT().
					UpdateWebhook(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil)
			},
			wantErr: nil,
		},
		{
			name: "webhook not found",
			args: args{
				ctx:       context.Background(),
				app:       &models.App{ID: "app-id"},
				user:      &models.User{ID: "user-id"},
				webhookId: "webhook-id",
				body: &dto.CreateWebhookRequestDto{
					Name: "updated-webhook",
				},
			},
			mockFn: func(s *UpdateWebhookService) {
				memberRepo, _ := s.AppMemberRepo.(*mocks.MockAppMemberRepository)
				memberRepo.EXPECT().FindAppMemberByAppIDAndUserId(gomock.Any(), "app-id", "user-id").
					Return(&models.AppMember{Role: models.AppMemberRoleOwner}, nil)
				webhookRepo, _ := s.WebhookRepo.(*mocks.MockWebhookRepository)
				webhookRepo.EXPECT().
					FindWebhookByID(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, errors.New("webhook not found"))
			},
			wantErr: errors.New("webhook not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			service := &UpdateWebhookService{
				App:           tt.args.app,
				AppMemberRepo: mocks.NewMockAppMemberRepository(ctrl),
				WebhookRepo:   mocks.NewMockWebhookRepository(ctrl),
				User:          tt.args.user,
				WebhookId:     tt.args.webhookId,
				Body:          tt.args.body,
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
