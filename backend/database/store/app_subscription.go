package store

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/mujhtech/s3ase/database"
	"github.com/mujhtech/s3ase/database/models"
)

type appSubscriptionRepo struct{ db *database.Database }

func NewAppSubscriptionRepository(db *database.Database) AppSubscriptionRepository {
	return &appSubscriptionRepo{db: db}
}

func (r *appSubscriptionRepo) CreateAppSubscription(ctx context.Context, subscription *models.AppSubscription) error {
	stmt := Builder.Insert("app_subscriptions").Columns("id", "app_id", "plan", "status", "metadata").
		Values(subscription.ID, subscription.AppID, subscription.Plan, subscription.Status, subscription.Metadata)
	sql, args, err := stmt.ToSql()
	if err != nil {
		return err
	}
	if _, err := r.db.GetDB().ExecContext(ctx, sql, args...); err != nil {
		return ProcessSQLErrorfWithCtx(ctx, sql, err, "failed to create app subscription")
	}
	return nil
}

func (r *appSubscriptionRepo) FindAppSubscriptionByAppID(ctx context.Context, appID string) (*models.AppSubscription, error) {
	stmt := Builder.Select("id, app_id, plan, status, metadata, created_at, updated_at, deleted_at").
		From("app_subscriptions").Where(squirrel.Eq{"app_id": appID}).Where(excludeDeleted)
	sql, args, err := stmt.ToSql()
	if err != nil {
		return nil, err
	}
	subscription := new(models.AppSubscription)
	if err := r.db.GetDB().GetContext(ctx, subscription, sql, args...); err != nil {
		return nil, ProcessSQLErrorfWithCtx(ctx, sql, err, "failed to find app subscription")
	}
	return subscription, nil
}
