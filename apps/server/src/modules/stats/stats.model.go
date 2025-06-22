package stats

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type StatPeriod string

const (
	StatMinutely StatPeriod = "minutely"
	StatHourly   StatPeriod = "hourly"
	StatDaily    StatPeriod = "daily"
)

type Stat struct {
	ID          primitive.ObjectID `json:"id"`
	MonitorID   primitive.ObjectID `json:"monitor_id"`
	Timestamp   time.Time          `json:"timestamp"`
	Ping        float64            `json:"ping"`
	PingMin     float64            `json:"ping_min"`
	PingMax     float64            `json:"ping_max"`
	Up          int                `json:"up"`
	Down        int                `json:"down"`
	Maintenance int                `json:"maintenance"`
}
