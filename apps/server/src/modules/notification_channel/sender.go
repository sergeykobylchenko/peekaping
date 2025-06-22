package notification_channel

import (
	"context"
	"peekaping/src/modules/heartbeat"
	"peekaping/src/modules/monitor"
)

type NotificationChannelProvider interface {
	Send(ctx context.Context, configJSON, message string, monitor *monitor.Model, heartbeat *heartbeat.Model) error
	Validate(configJSON string) error
	Unmarshal(configJSON string) (any, error)
}
