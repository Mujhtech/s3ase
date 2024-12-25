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

func TestFindWebhooksService_Run(t *testing.T) {
	type args struct {
		ctx  context.Context
		app  *models.App
		user *models.User
	}
	type testCase struct {
		name    string
		args    args
		mockFn  func(s *FindWebhooksService)
		want    []*models.Webhook
		wantErr error
	}
	tests := []testCase{
		{
			name: "successfully find webhooks",
			args: args{
				ctx:  context.Background(),
				app:  &models.App{ID: "app-id"},
				user: &models.User{ID: "user-id"},
			},
			mockFn: func(s *FindWebhooksService) {
				webhookRepo, _ := s.WebhookRepo.(*mocks.MockWebhookRepository)
				webhookRepo.EXPECT().
					FindWebhooksByAppID(gomock.Any(), "app-id").
					Return([]*models.Webhook{
						{
							ID:    "webhook-1",
							Name:  "Webhook 1",
							URL:   "https://example.com/webhook1",
							AppID: "app-id",
						},
						{
							ID:    "webhook-2",
							Name:  "Webhook 2",
							URL:   "https://example.com/webhook2",
							AppID: "app-id",
						},
					}, nil)
			},
			want: []*models.Webhook{
				{
					ID:    "webhook-1",
					Name:  "Webhook 1",
					URL:   "https://example.com/webhook1",
					AppID: "app-id",
				},
				{
					ID:    "webhook-2",
					Name:  "Webhook 2",
					URL:   "https://example.com/webhook2",
					AppID: "app-id",
				},
			},
			wantErr: nil,
		},
		{
			name: "no webhooks found",
			args: args{
				ctx:  context.Background(),
				app:  &models.App{ID: "app-id"},
				user: &models.User{ID: "user-id"},
			},
			mockFn: func(s *FindWebhooksService) {
				webhookRepo, _ := s.WebhookRepo.(*mocks.MockWebhookRepository)
				webhookRepo.EXPECT().
					FindWebhooksByAppID(gomock.Any(), "app-id").
					Return([]*models.Webhook{}, nil)
			},
			want:    []*models.Webhook{},
			wantErr: nil,
		},
		{
			name: "error finding webhooks",
			args: args{
				ctx:  context.Background(),
				app:  &models.App{ID: "app-id"},
				user: &models.User{ID: "user-id"},
			},
			mockFn: func(s *FindWebhooksService) {
				webhookRepo, _ := s.WebhookRepo.(*mocks.MockWebhookRepository)
				webhookRepo.EXPECT().
					FindWebhooksByAppID(gomock.Any(), "app-id").
					Return(nil, errors.New("some error"))
			},
			want:    nil,
			wantErr: errors.New("some error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			webhookRepo := mocks.NewMockWebhookRepository(ctrl)
			s := &FindWebhooksService{
				WebhookRepo: webhookRepo,
				App:         tt.args.app,
				User:        tt.args.user,
			}

			if tt.mockFn != nil {
				tt.mockFn(s)
			}

			got, err := s.Run(tt.args.ctx)
			if tt.wantErr != nil {
				require.EqualError(t, err, tt.wantErr.Error())
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, tt.want, got)
		})
	}
}
