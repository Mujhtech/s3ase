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

func TestFindMembersService_Run(t *testing.T) {
	type args struct {
		ctx  context.Context
		app  *models.App
		user *models.User
	}
	type testCase struct {
		name    string
		args    args
		mockFn  func(s *FindMembersService)
		want    []*models.AppMember
		wantErr error
	}
	tests := []testCase{
		{
			name: "successfully find members",
			args: args{
				ctx:  context.Background(),
				app:  &models.App{ID: "app-id"},
				user: &models.User{ID: "user-id"},
			},
			mockFn: func(s *FindMembersService) {
				appMemberRepo, _ := s.AppMemberRepo.(*mocks.MockAppMemberRepository)
				userRepo, _ := s.UserRepo.(*mocks.MockUserRepository)

				members := []*models.AppMember{
					{ID: "member-1", UserId: "user-1"},
					{ID: "member-2", UserId: "user-2"},
				}

				appMemberRepo.EXPECT().
					FindAppMembersByAppID(gomock.Any(), "app-id").
					Return(members, nil)

				userRepo.EXPECT().
					FindUserByID(gomock.Any(), "user-1").
					Return(&models.User{ID: "user-1", Name: "User 1"}, nil)
				userRepo.EXPECT().
					FindUserByID(gomock.Any(), "user-2").
					Return(&models.User{ID: "user-2", Name: "User 2"}, nil)

				members[0].User = &models.User{ID: "user-1", Name: "User 1"}
				members[1].User = &models.User{ID: "user-2", Name: "User 2"}
			},
			want: []*models.AppMember{
				{ID: "member-1", UserId: "user-1", User: &models.User{ID: "user-1", Name: "User 1"}},
				{ID: "member-2", UserId: "user-2", User: &models.User{ID: "user-2", Name: "User 2"}},
			},
			wantErr: nil,
		},
		{
			name: "error finding members",
			args: args{
				ctx:  context.Background(),
				app:  &models.App{ID: "app-id"},
				user: &models.User{ID: "user-id"},
			},
			mockFn: func(s *FindMembersService) {
				appMemberRepo, _ := s.AppMemberRepo.(*mocks.MockAppMemberRepository)
				appMemberRepo.EXPECT().
					FindAppMembersByAppID(gomock.Any(), "app-id").
					Return(nil, errors.New("failed to find members"))
			},
			want:    nil,
			wantErr: errors.New("failed to find members"),
		},
		{
			name: "error finding user details",
			args: args{
				ctx:  context.Background(),
				app:  &models.App{ID: "app-id"},
				user: &models.User{ID: "user-id"},
			},
			mockFn: func(s *FindMembersService) {
				appMemberRepo, _ := s.AppMemberRepo.(*mocks.MockAppMemberRepository)
				userRepo, _ := s.UserRepo.(*mocks.MockUserRepository)

				appMemberRepo.EXPECT().
					FindAppMembersByAppID(gomock.Any(), "app-id").
					Return([]*models.AppMember{{ID: "member-1", UserId: "user-1"}}, nil)

				userRepo.EXPECT().
					FindUserByID(gomock.Any(), "user-1").
					Return(nil, errors.New("failed to find user"))
			},
			want:    nil,
			wantErr: errors.New("failed to find user"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			service := &FindMembersService{
				App:           tt.args.app,
				AppRepo:       mocks.NewMockAppRepository(ctrl),
				AppMemberRepo: mocks.NewMockAppMemberRepository(ctrl),
				UserRepo:      mocks.NewMockUserRepository(ctrl),
				User:          tt.args.user,
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
