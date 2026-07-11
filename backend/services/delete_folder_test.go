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

func TestDeleteFolderService_Run(t *testing.T) {
	type args struct {
		ctx      context.Context
		app      *models.App
		user     *models.User
		folderId string
	}
	type testCase struct {
		name    string
		args    args
		mockFn  func(s *DeleteFolderService)
		wantErr error
	}
	tests := []testCase{
		{
			name: "successfully delete folder",
			args: args{
				ctx:      context.Background(),
				app:      &models.App{ID: "app-id"},
				user:     &models.User{ID: "user-id"},
				folderId: "folder-id",
			},
			mockFn: func(s *DeleteFolderService) {
				folderRepo, _ := s.FolderRepo.(*mocks.MockFolderRepository)
				fileRepo, _ := s.FileRepo.(*mocks.MockFileRepository)
				folderRepo.EXPECT().
					FindFolderByID(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&models.Folder{ID: "folder-id", AppID: "app-id"}, nil)
				fileRepo.EXPECT().FindFilesByFolderID(gomock.Any(), "folder-id").Return([]*models.File{}, nil)
				folderRepo.EXPECT().
					DeleteFolder(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil)
			},
			wantErr: nil,
		},
		{
			name: "folder not found",
			args: args{
				ctx:      context.Background(),
				app:      &models.App{ID: "app-id"},
				user:     &models.User{ID: "user-id"},
				folderId: "folder-id",
			},
			mockFn: func(s *DeleteFolderService) {
				folderRepo, _ := s.FolderRepo.(*mocks.MockFolderRepository)
				folderRepo.EXPECT().
					FindFolderByID(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, errors.New("folder not found"))
			},
			wantErr: errors.New("folder not found"),
		},
		{
			name: "delete folder fails",
			args: args{
				ctx:      context.Background(),
				app:      &models.App{ID: "app-id"},
				user:     &models.User{ID: "user-id"},
				folderId: "folder-id",
			},
			mockFn: func(s *DeleteFolderService) {
				folderRepo, _ := s.FolderRepo.(*mocks.MockFolderRepository)
				fileRepo, _ := s.FileRepo.(*mocks.MockFileRepository)
				folderRepo.EXPECT().
					FindFolderByID(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&models.Folder{ID: "folder-id", AppID: "app-id"}, nil)
				fileRepo.EXPECT().FindFilesByFolderID(gomock.Any(), "folder-id").Return([]*models.File{}, nil)
				folderRepo.EXPECT().
					DeleteFolder(gomock.Any(), gomock.Any()).
					Times(1).
					Return(errors.New("failed to delete folder"))
			},
			wantErr: errors.New("failed to delete folder"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			service := &DeleteFolderService{
				App:        tt.args.app,
				FolderRepo: mocks.NewMockFolderRepository(ctrl),
				FileRepo:   mocks.NewMockFileRepository(ctrl),
				User:       tt.args.user,
				FolderId:   tt.args.folderId,
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
