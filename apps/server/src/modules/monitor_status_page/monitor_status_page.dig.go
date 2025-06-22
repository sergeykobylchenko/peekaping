package monitor_status_page

import "go.uber.org/dig"

func RegisterDependencies(container *dig.Container) {
	container.Provide(NewMongoRepository)
	container.Provide(NewService)
}
