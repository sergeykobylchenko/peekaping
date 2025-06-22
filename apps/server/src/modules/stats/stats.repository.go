package stats

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Repository interface {
	GetOrCreateStat(ctx context.Context, monitorID primitive.ObjectID, timestamp time.Time, period StatPeriod) (*Stat, error)
	UpsertStat(ctx context.Context, stat *Stat, period StatPeriod) error
	FindStatsByMonitorIDAndTimeRange(ctx context.Context, monitorID primitive.ObjectID, since, until time.Time, period StatPeriod) ([]*Stat, error)
}
