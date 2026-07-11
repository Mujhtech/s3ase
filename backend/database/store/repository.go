package store

import (
	"context"

	"github.com/mujhtech/s3ase/api/dto"
	"github.com/mujhtech/s3ase/database/models"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *models.User) error
	UpdateUser(ctx context.Context, user *models.User) error
	FindUserByEmail(ctx context.Context, email string) (*models.User, error)
	FindUserByID(ctx context.Context, id string) (*models.User, error)
}

type TokenRepository interface {
	CreateToken(ctx context.Context, token *models.Token) error
	FindTokenByID(ctx context.Context, id string) (*models.Token, error)
	DeleteToken(ctx context.Context, id string) error
}

type AppRepository interface {
	CreateApp(ctx context.Context, app *models.App) error
	UpdateApp(ctx context.Context, app *models.App) error
	FindAppByID(ctx context.Context, id string) (*models.App, error)
	FindAppsByUserID(ctx context.Context, userID string) ([]*models.App, error)
	DeleteApp(ctx context.Context, id string) error
}

type AppMemberRepository interface {
	CreateAppMember(ctx context.Context, appMember *models.AppMember) error
	FindAppMemberByID(ctx context.Context, id string) (*models.AppMember, error)
	FindAppMemberByAppIDAndUserId(ctx context.Context, appID string, userId string) (*models.AppMember, error)
	FindAppMembersByAppID(ctx context.Context, appID string) ([]*models.AppMember, error)
	FindAppMembersByUserID(ctx context.Context, userID string) ([]*models.AppMember, error)
	DeleteAppMember(ctx context.Context, id string) error
}

type AppSubscriptionRepository interface {
	CreateAppSubscription(ctx context.Context, subscription *models.AppSubscription) error
	FindAppSubscriptionByAppID(ctx context.Context, appID string) (*models.AppSubscription, error)
}

type FolderRepository interface {
	CreateFolder(ctx context.Context, folder *models.Folder) error
	UpdateFolder(ctx context.Context, folder *models.Folder) error
	FindFolderByID(ctx context.Context, id string) (*models.Folder, error)
	FindFoldersByAppID(ctx context.Context, appID string) ([]*models.Folder, error)
	DeleteFolder(ctx context.Context, id string) error
}

type FileRepository interface {
	CreateFile(ctx context.Context, file *models.File) error
	UpdateFile(ctx context.Context, file *models.File) error
	FindFileByID(ctx context.Context, id string) (*models.File, error)
	FindFilesByUserID(ctx context.Context, userID string) ([]*models.File, error)
	FindFilesByAppID(ctx context.Context, appID string) ([]*models.File, error)
	FindFilesByAppIDWithQuery(ctx context.Context, appID string, query *dto.FileQueryDto) ([]*models.File, error)
	FindFilesByFolderID(ctx context.Context, folderID string) ([]*models.File, error)
	DeleteFile(ctx context.Context, id string) error
}

type ApiKeyRepository interface {
	CreateApiKey(ctx context.Context, apiKey *models.ApiKey) error
	UpdateApiKey(ctx context.Context, apiKey *models.ApiKey) error
	FindApiKeyByID(ctx context.Context, id string) (*models.ApiKey, error)
	FindApiKeyByHash(ctx context.Context, hash string) (*models.ApiKey, error)
	FindApiKeysByAppID(ctx context.Context, appId string) ([]*models.ApiKey, error)
	TouchApiKey(ctx context.Context, id string) error
	DeleteApiKey(ctx context.Context, id string) error
}

type WebhookRepository interface {
	CreateWebhook(ctx context.Context, apiKey *models.Webhook) error
	UpdateWebhook(ctx context.Context, apiKey *models.Webhook) error
	FindWebhookByID(ctx context.Context, id string) (*models.Webhook, error)
	FindWebhooksByAppID(ctx context.Context, appId string) ([]*models.Webhook, error)
	DeleteWebhook(ctx context.Context, id string) error
}

type DomainRepository interface {
	CreateDomain(ctx context.Context, apiKey *models.Domain) error
	UpdateDomain(ctx context.Context, apiKey *models.Domain) error
	FindDomainByID(ctx context.Context, id string) (*models.Domain, error)
	FindDomainByDomain(ctx context.Context, domain string) (*models.Domain, error)
	FindDomainByAppID(ctx context.Context, appId string) (*models.Domain, error)
	DeleteDomain(ctx context.Context, id string) error
}
