package notification_channel

import (
	"go.uber.org/dig"
)

func RegisterDependencies(container *dig.Container) {
	container.Provide(NewRepository)
	container.Provide(NewService)
	container.Provide(NewController)
	container.Provide(NewRoute)
	container.Provide(NewNotificationEventListener)
}
