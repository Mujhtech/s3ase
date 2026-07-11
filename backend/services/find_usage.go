package services

import (
	"context"
	"time"

	"github.com/mujhtech/s3ase/database/models"
	"github.com/mujhtech/s3ase/database/store"
)

type DailyUsage struct {
	Date         string `json:"date"`
	Uploads      int64  `json:"uploads"`
	StorageBytes int64  `json:"storage_bytes"`
}

type Usage struct {
	FileCount    int64        `json:"file_count"`
	StorageBytes int64        `json:"storage_bytes"`
	Daily        []DailyUsage `json:"daily"`
}

type FindUsageService struct {
	App      *models.App
	FileRepo store.FileRepository
	Now      time.Time
}

func (s *FindUsageService) Run(ctx context.Context) (*Usage, error) {
	files, err := s.FileRepo.FindFilesByAppID(ctx, s.App.ID)
	if err != nil {
		return nil, err
	}
	now := s.Now.UTC()
	if now.IsZero() {
		now = time.Now().UTC()
	}
	start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC).AddDate(0, 0, -29)
	dailyByDate := make(map[string]*DailyUsage, 30)
	daily := make([]DailyUsage, 30)
	for index := range daily {
		date := start.AddDate(0, 0, index).Format("2006-01-02")
		daily[index] = DailyUsage{Date: date}
		dailyByDate[date] = &daily[index]
	}
	usage := &Usage{Daily: daily}
	for _, file := range files {
		if file.Status != models.FileStatusCompleted {
			continue
		}
		usage.FileCount++
		usage.StorageBytes += file.Size
		if day := dailyByDate[file.CreatedAt.UTC().Format("2006-01-02")]; day != nil {
			day.Uploads++
			day.StorageBytes += file.Size
		}
	}
	return usage, nil
}
