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

func TestFindFoldersService_Run(t *testing.T) {
	type args struct {
		ctx  context.Context
		app  *models.App
		user *models.User
	}
	type testCase struct {
		name    string
		args    args
		mockFn  func(s *FindFoldersService)
		want    []*models.Folder
		wantErr error
	}
	tests := []testCase{
		{
			name: "successfully find folders",
			args: args{
				ctx:  context.Background(),
				app:  &models.App{ID: "app-id"},
				user: &models.User{ID: "user-id"},
			},
			mockFn: func(s *FindFoldersService) {
				folderRepo, _ := s.FolderRepo.(*mocks.MockFolderRepository)
				folderRepo.EXPECT().
					FindFoldersByAppID(gomock.Any(), "app-id").
					Return([]*models.Folder{
						{ID: "folder-1", Name: "Folder 1"},
						{ID: "folder-2", Name: "Folder 2"},
					}, nil)
			},
			want: []*models.Folder{
				{ID: "folder-1", Name: "Folder 1"},
				{ID: "folder-2", Name: "Folder 2"},
			},
			wantErr: nil,
		},
		{
			name: "error finding folders",
			args: args{
				ctx:  context.Background(),
				app:  &models.App{ID: "app-id"},
				user: &models.User{ID: "user-id"},
			},
			mockFn: func(s *FindFoldersService) {
				folderRepo, _ := s.FolderRepo.(*mocks.MockFolderRepository)
				folderRepo.EXPECT().
					FindFoldersByAppID(gomock.Any(), "app-id").
					Return(nil, errors.New("failed to find folders"))
			},
			want:    nil,
			wantErr: errors.New("failed to find folders"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			service := &FindFoldersService{
				App:        tt.args.app,
				FolderRepo: mocks.NewMockFolderRepository(ctrl),
				FileRepo:   mocks.NewMockFileRepository(ctrl),
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
