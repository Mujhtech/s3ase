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

func TestDeleteWebhookService_Run(t *testing.T) {
	type args struct {
		ctx       context.Context
		app       *models.App
		user      *models.User
		webhookId string
	}
	type testCase struct {
		name    string
		args    args
		mockFn  func(s *DeleteWebhookService)
		wantErr error
	}
	tests := []testCase{
		{
			name: "successfully delete webhook",
			args: args{
				ctx:       context.Background(),
				app:       &models.App{ID: "app-id"},
				user:      &models.User{ID: "user-id"},
				webhookId: "webhook-id",
			},
			mockFn: func(s *DeleteWebhookService) {
				webhookRepo, _ := s.WebhookRepo.(*mocks.MockWebhookRepository)
				webhookRepo.EXPECT().
					FindWebhookByID(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&models.Webhook{ID: "webhook-id"}, nil)
				webhookRepo.EXPECT().
					DeleteWebhook(gomock.Any(), gomock.Any()).
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
			},
			mockFn: func(s *DeleteWebhookService) {
				webhookRepo, _ := s.WebhookRepo.(*mocks.MockWebhookRepository)
				webhookRepo.EXPECT().
					FindWebhookByID(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, errors.New("webhook not found"))
			},
			wantErr: errors.New("webhook not found"),
		},
		{
			name: "delete webhook fails",
			args: args{
				ctx:       context.Background(),
				app:       &models.App{ID: "app-id"},
				user:      &models.User{ID: "user-id"},
				webhookId: "webhook-id",
			},
			mockFn: func(s *DeleteWebhookService) {
				webhookRepo, _ := s.WebhookRepo.(*mocks.MockWebhookRepository)
				webhookRepo.EXPECT().
					FindWebhookByID(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&models.Webhook{ID: "webhook-id"}, nil)
				webhookRepo.EXPECT().
					DeleteWebhook(gomock.Any(), gomock.Any()).
					Times(1).
					Return(errors.New("failed to delete webhook"))
			},
			wantErr: errors.New("failed to delete webhook"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			service := &DeleteWebhookService{
				App:         tt.args.app,
				WebhookRepo: mocks.NewMockWebhookRepository(ctrl),
				User:        tt.args.user,
				WebhookId:   tt.args.webhookId,
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
