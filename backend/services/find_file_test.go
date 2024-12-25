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

func TestFindFileService_Run(t *testing.T) {
	type args struct {
		ctx    context.Context
		fileId string
		app    *models.App
		user   *models.User
	}
	type testCase struct {
		name    string
		args    args
		mockFn  func(s *FindFileService)
		want    *models.File
		wantErr error
	}
	tests := []testCase{
		{
			name: "successfully find file",
			args: args{
				ctx:    context.Background(),
				fileId: "file-id",
				app:    &models.App{ID: "app-id"},
				user:   &models.User{ID: "user-id"},
			},
			mockFn: func(s *FindFileService) {
				fileRepo, _ := s.FileRepo.(*mocks.MockFileRepository)
				fileRepo.EXPECT().
					FindFileByID(gomock.Any(), "file-id").
					Return(&models.File{
						ID:       "file-id",
						Name:     "test.jpg",
						Size:     1024,
						MimeType: "image/jpeg",
						AppID:    "app-id",
					}, nil)
			},
			want: &models.File{
				ID:       "file-id",
				Name:     "test.jpg",
				Size:     1024,
				MimeType: "image/jpeg",
				AppID:    "app-id",
			},
			wantErr: nil,
		},
		{
			name: "file not found",
			args: args{
				ctx:    context.Background(),
				fileId: "non-existent-file",
				app:    &models.App{ID: "app-id"},
				user:   &models.User{ID: "user-id"},
			},
			mockFn: func(s *FindFileService) {
				fileRepo, _ := s.FileRepo.(*mocks.MockFileRepository)
				fileRepo.EXPECT().
					FindFileByID(gomock.Any(), "non-existent-file").
					Return(nil, errors.New("file not found"))
			},
			want:    nil,
			wantErr: errors.New("file not found"),
		},
		{
			name: "database error",
			args: args{
				ctx:    context.Background(),
				fileId: "file-id",
				app:    &models.App{ID: "app-id"},
				user:   &models.User{ID: "user-id"},
			},
			mockFn: func(s *FindFileService) {
				fileRepo, _ := s.FileRepo.(*mocks.MockFileRepository)
				fileRepo.EXPECT().
					FindFileByID(gomock.Any(), "file-id").
					Return(nil, errors.New("database error"))
			},
			want:    nil,
			wantErr: errors.New("database error"),
		},
		{
			name: "invalid file id",
			args: args{
				ctx:    context.Background(),
				fileId: "",
				app:    &models.App{ID: "app-id"},
				user:   &models.User{ID: "user-id"},
			},
			mockFn: func(s *FindFileService) {
				fileRepo, _ := s.FileRepo.(*mocks.MockFileRepository)
				fileRepo.EXPECT().
					FindFileByID(gomock.Any(), "").
					Return(nil, errors.New("invalid file id"))
			},
			want:    nil,
			wantErr: errors.New("invalid file id"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			service := &FindFileService{
				FileID:     tt.args.fileId,
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
