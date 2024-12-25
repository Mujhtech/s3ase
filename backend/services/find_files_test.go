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

func TestFindFilesService_Run(t *testing.T) {
	type args struct {
		ctx  context.Context
		app  *models.App
		user *models.User
	}
	type testCase struct {
		name    string
		args    args
		mockFn  func(s *FindFilesService)
		want    []*models.File
		wantErr error
	}
	tests := []testCase{
		{
			name: "successfully find files",
			args: args{
				ctx:  context.Background(),
				app:  &models.App{ID: "app-id"},
				user: &models.User{ID: "user-id"},
			},
			mockFn: func(s *FindFilesService) {
				fileRepo, _ := s.FileRepo.(*mocks.MockFileRepository)
				fileRepo.EXPECT().
					FindFilesByAppIDWithQuery(gomock.Any(), "app-id", gomock.Any()).
					Return([]*models.File{
						{
							ID:       "file-1",
							Name:     "test1.jpg",
							Size:     1024,
							MimeType: "image/jpeg",
							AppID:    "app-id",
						},
						{
							ID:       "file-2",
							Name:     "test2.pdf",
							Size:     2048,
							MimeType: "application/pdf",
							AppID:    "app-id",
						},
					}, nil)
			},
			want: []*models.File{
				{
					ID:       "file-1",
					Name:     "test1.jpg",
					Size:     1024,
					MimeType: "image/jpeg",
					AppID:    "app-id",
				},
				{
					ID:       "file-2",
					Name:     "test2.pdf",
					Size:     2048,
					MimeType: "application/pdf",
					AppID:    "app-id",
				},
			},
			wantErr: nil,
		},
		{
			name: "no files found",
			args: args{
				ctx:  context.Background(),
				app:  &models.App{ID: "app-id"},
				user: &models.User{ID: "user-id"},
			},
			mockFn: func(s *FindFilesService) {
				fileRepo, _ := s.FileRepo.(*mocks.MockFileRepository)
				fileRepo.EXPECT().
					FindFilesByAppIDWithQuery(gomock.Any(), gomock.Any(), gomock.Any()).
					Return([]*models.File{}, nil)
			},
			want:    []*models.File{},
			wantErr: nil,
		},
		{
			name: "error finding files",
			args: args{
				ctx:  context.Background(),
				app:  &models.App{ID: "app-id"},
				user: &models.User{ID: "user-id"},
			},
			mockFn: func(s *FindFilesService) {
				fileRepo, _ := s.FileRepo.(*mocks.MockFileRepository)
				fileRepo.EXPECT().
					FindFilesByAppIDWithQuery(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, errors.New("database error"))
			},
			want:    nil,
			wantErr: errors.New("database error"),
		},
		{
			name: "invalid app id",
			args: args{
				ctx:  context.Background(),
				app:  &models.App{ID: ""},
				user: &models.User{ID: "user-id"},
			},
			mockFn: func(s *FindFilesService) {
				fileRepo, _ := s.FileRepo.(*mocks.MockFileRepository)
				fileRepo.EXPECT().
					FindFilesByAppIDWithQuery(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, errors.New("invalid app id"))
			},
			want:    nil,
			wantErr: errors.New("invalid app id"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			s := &FindFilesService{
				FolderRepo: mocks.NewMockFolderRepository(ctrl),
				FileRepo:   mocks.NewMockFileRepository(ctrl),
				App:        tt.args.app,
				User:       tt.args.user,
				Query:      &dto.FileQueryDto{}, // Add this line to initialize Query
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
