package status_page

import (
	"peekaping/src/modules/heartbeat"
	"peekaping/src/modules/monitor"

	"go.uber.org/dig"
	"go.uber.org/zap"
)

func RegisterDependencies(container *dig.Container) {
	container.Provide(NewStatusPageRepository)
	container.Provide(NewService)
	container.Provide(func(service Service, monitorService monitor.Service, heartbeatService heartbeat.Service, logger *zap.SugaredLogger) *Controller {
		return NewController(service, monitorService, heartbeatService, logger)
	})
	container.Provide(NewRoute)
}
