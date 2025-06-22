package stats

import (
	"peekaping/src/modules/events"

	"go.uber.org/dig"
)

func RegisterDependencies(container *dig.Container) {
	container.Provide(NewMongoRepository)
	container.Provide(NewService)
	container.Invoke(func(s Service, bus *events.EventBus) {
		if impl, ok := s.(*ServiceImpl); ok {
			impl.RegisterEventHandlers(bus)
		}
	})
}
