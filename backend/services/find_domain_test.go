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

func TestFindDomainService_Run(t *testing.T) {
	type args struct {
		ctx  context.Context
		app  *models.App
		user *models.User
	}
	type testCase struct {
		name    string
		args    args
		mockFn  func(s *FindDomainService)
		want    *models.Domain
		wantErr error
	}
	tests := []testCase{
		{
			name: "successfully find domain",
			args: args{
				ctx:  context.Background(),
				app:  &models.App{ID: "app-id"},
				user: &models.User{ID: "user-id"},
			},
			mockFn: func(s *FindDomainService) {
				domainRepo, _ := s.DomainRepo.(*mocks.MockDomainRepository)
				domainRepo.EXPECT().
					FindDomainByAppID(gomock.Any(), "app-id").
					Return(&models.Domain{
						ID:     "domain-id",
						AppID:  "app-id",
						Domain: "test.com",
					}, nil)
			},
			want: &models.Domain{
				ID:     "domain-id",
				AppID:  "app-id",
				Domain: "test.com",
			},
			wantErr: nil,
		},
		{
			name: "domain not found",
			args: args{
				ctx:  context.Background(),
				app:  &models.App{ID: "app-id"},
				user: &models.User{ID: "user-id"},
			},
			mockFn: func(s *FindDomainService) {
				domainRepo, _ := s.DomainRepo.(*mocks.MockDomainRepository)
				domainRepo.EXPECT().
					FindDomainByAppID(gomock.Any(), "app-id").
					Return(nil, errors.New("domain not found"))
			},
			want:    nil,
			wantErr: errors.New("domain not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			service := &FindDomainService{
				App:        tt.args.app,
				DomainRepo: mocks.NewMockDomainRepository(ctrl),
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
