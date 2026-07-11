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

func TestUpdateFolderService_Run(t *testing.T) {
	type args struct {
		ctx      context.Context
		app      *models.App
		user     *models.User
		folderId string
		body     *dto.CreateFolderRequestDto
	}
	type testCase struct {
		name    string
		args    args
		mockFn  func(s *UpdateFolderService)
		wantErr error
	}
	tests := []testCase{
		{
			name: "successfully update folder",
			args: args{
				ctx:      context.Background(),
				app:      &models.App{ID: "app-id"},
				user:     &models.User{ID: "user-id"},
				folderId: "folder-id",
				body: &dto.CreateFolderRequestDto{
					Name:        "updated-folder",
					Description: "updated description",
				},
			},
			mockFn: func(s *UpdateFolderService) {
				folderRepo, _ := s.FolderRepo.(*mocks.MockFolderRepository)
				folderRepo.EXPECT().
					FindFolderByID(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&models.Folder{ID: "folder-id", AppID: "app-id"}, nil)
				folderRepo.EXPECT().
					UpdateFolder(gomock.Any(), gomock.Any()).
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
				body: &dto.CreateFolderRequestDto{
					Name:        "updated-folder",
					Description: "updated description",
				},
			},
			mockFn: func(s *UpdateFolderService) {
				folderRepo, _ := s.FolderRepo.(*mocks.MockFolderRepository)
				folderRepo.EXPECT().
					FindFolderByID(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, errors.New("folder not found"))
			},
			wantErr: errors.New("folder not found"),
		},
		{
			name: "update folder fails",
			args: args{
				ctx:      context.Background(),
				app:      &models.App{ID: "app-id"},
				user:     &models.User{ID: "user-id"},
				folderId: "folder-id",
				body: &dto.CreateFolderRequestDto{
					Name:        "updated-folder",
					Description: "updated description",
				},
			},
			mockFn: func(s *UpdateFolderService) {
				folderRepo, _ := s.FolderRepo.(*mocks.MockFolderRepository)
				folderRepo.EXPECT().
					FindFolderByID(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&models.Folder{ID: "folder-id", AppID: "app-id"}, nil)
				folderRepo.EXPECT().
					UpdateFolder(gomock.Any(), gomock.Any()).
					Times(1).
					Return(errors.New("failed to update folder"))
			},
			wantErr: errors.New("failed to update folder"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			service := &UpdateFolderService{
				App:        tt.args.app,
				FolderRepo: mocks.NewMockFolderRepository(ctrl),
				User:       tt.args.user,
				FolderId:   tt.args.folderId,
				Body:       tt.args.body,
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
