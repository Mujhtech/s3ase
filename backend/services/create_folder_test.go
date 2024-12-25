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

func TestCreateFolderService_Run(t *testing.T) {
	type args struct {
		ctx  context.Context
		app  *models.App
		user *models.User
		body *dto.CreateFolderRequestDto
	}
	type testCase struct {
		name    string
		args    args
		mockFn  func(s *CreateFolderService)
		want    *models.Folder
		wantErr error
	}
	tests := []testCase{
		{
			name: "successfully create folder",
			args: args{
				ctx:  context.Background(),
				app:  &models.App{ID: "app-id"},
				user: &models.User{ID: "user-id"},
				body: &dto.CreateFolderRequestDto{
					Name:        "test-folder",
					Description: "test description",
				},
			},
			mockFn: func(s *CreateFolderService) {
				folderRepo, _ := s.FolderRepo.(*mocks.MockFolderRepository)
				folderRepo.EXPECT().
					CreateFolder(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil)
				folderRepo.EXPECT().
					FindFolderByID(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&models.Folder{
						ID:          "folder-id",
						Name:        "test-folder",
						Description: null.NewString("test description", true),
					}, nil)
			},
			want: &models.Folder{
				ID:          "folder-id",
				Name:        "test-folder",
				Description: null.NewString("test description", true),
			},
			wantErr: nil,
		},
		{
			name: "fail to create folder",
			args: args{
				ctx:  context.Background(),
				app:  &models.App{ID: "app-id"},
				user: &models.User{ID: "user-id"},
				body: &dto.CreateFolderRequestDto{
					Name:        "test-folder",
					Description: "test description",
				},
			},
			mockFn: func(s *CreateFolderService) {
				folderRepo, _ := s.FolderRepo.(*mocks.MockFolderRepository)
				folderRepo.EXPECT().
					CreateFolder(gomock.Any(), gomock.Any()).
					Times(1).
					Return(errors.New("failed to create folder"))
			},
			want:    nil,
			wantErr: errors.New("failed to create folder"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			service := &CreateFolderService{
				App:        tt.args.app,
				FolderRepo: mocks.NewMockFolderRepository(ctrl),
				User:       tt.args.user,
				Body:       tt.args.body,
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
