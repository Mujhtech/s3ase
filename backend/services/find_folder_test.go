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

func TestFindFolderService_Run(t *testing.T) {
	type args struct {
		ctx      context.Context
		folderId string
		app      *models.App
		user     *models.User
	}
	type testCase struct {
		name    string
		args    args
		mockFn  func(s *FindFolderService)
		want    *models.Folder
		wantErr error
	}
	tests := []testCase{
		{
			name: "successfully find folder",
			args: args{
				ctx:      context.Background(),
				folderId: "folder-id",
				app:      &models.App{ID: "app-id"},
				user:     &models.User{ID: "user-id"},
			},
			mockFn: func(s *FindFolderService) {
				folderRepo, _ := s.FolderRepo.(*mocks.MockFolderRepository)
				folderRepo.EXPECT().
					FindFolderByID(gomock.Any(), "folder-id").
					Return(&models.Folder{
						ID:   "folder-id",
						Name: "Test Folder",
					}, nil)
			},
			want: &models.Folder{
				ID:   "folder-id",
				Name: "Test Folder",
			},
			wantErr: nil,
		},
		{
			name: "folder not found",
			args: args{
				ctx:      context.Background(),
				folderId: "folder-id",
				app:      &models.App{ID: "app-id"},
				user:     &models.User{ID: "user-id"},
			},
			mockFn: func(s *FindFolderService) {
				folderRepo, _ := s.FolderRepo.(*mocks.MockFolderRepository)
				folderRepo.EXPECT().
					FindFolderByID(gomock.Any(), "folder-id").
					Return(nil, errors.New("folder not found"))
			},
			want:    nil,
			wantErr: errors.New("folder not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			service := &FindFolderService{
				FolderId:   tt.args.folderId,
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
