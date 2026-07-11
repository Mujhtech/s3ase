package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/lib/pq"
	"github.com/mujhtech/s3ase/database"
	errs "github.com/mujhtech/s3ase/errors"
	"github.com/rs/zerolog/log"
)

var (
	Builder = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	// errors
	ErrNotFound = errors.New("not found")
)

const (
	excludeDeleted = "deleted_at IS NULL"
)

type Store struct {
	UserRepo            UserRepository
	AppRepo             AppRepository
	ApiKeyRepo          ApiKeyRepository
	AppMemberRepo       AppMemberRepository
	AppSubscriptionRepo AppSubscriptionRepository
	FolderRepo          FolderRepository
	FileRepo            FileRepository
	TokenRepo           TokenRepository
	WebhookRepo         WebhookRepository
	DomainRepo          DomainRepository
}

func NewStore(db *database.Database) *Store {
	return &Store{
		UserRepo:            NewUserRepository(db),
		AppRepo:             NewAppRepository(db),
		ApiKeyRepo:          NewApiKeyRepository(db),
		AppMemberRepo:       NewAppMemberRepository(db),
		AppSubscriptionRepo: NewAppSubscriptionRepository(db),
		FolderRepo:          NewFolderRepository(db),
		FileRepo:            NewFileRepository(db),
		TokenRepo:           NewTokenRepository(db),
		WebhookRepo:         NewWebhookRepository(db),
		DomainRepo:          NewDomainRepository(db),
	}
}

func ProcessSQLErrorfWithCtx(ctx context.Context, query string, err error, format string, args ...interface{}) error {
	// create fallback error returned if we can't map it
	fallbackErr := fmt.Errorf(format, args...)

	// always log internal error together with message.
	log.Info().Msgf("Query: %v", query)
	log.Error().Err(err).Msgf("%v: [SQL] %v", fallbackErr, err)

	switch {
	case errors.Is(err, sql.ErrNoRows):
		return ErrNotFound
	case isPostgresConstraint(err, "23505"):
		return errs.ErrConflict
	default:
		return fallbackErr
	}
}

func isPostgresConstraint(err error, code pq.ErrorCode) bool {
	var postgresError *pq.Error
	return errors.As(err, &postgresError) && postgresError.Code == code
}

// func toDatabaseValue(value interface{}) interface{} {
// 	if value == nil {
// 		return nil
// 	}

// 	switch v := value.(type) {
// 	case string:
// 		return v
// 	case int:
// 		return v
// 	case int64:
// 		return v
// 	case bool:
// 		return v
// 	default:
// 		return v
// 	}
// }
