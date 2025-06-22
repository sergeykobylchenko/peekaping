package monitor_maintenance

import (
	"context"

	"go.uber.org/zap"
)

type Service interface {
	Create(ctx context.Context, monitorID string, maintenanceID string) (*Model, error)
	FindByID(ctx context.Context, id string) (*Model, error)
	Delete(ctx context.Context, id string) error
	FindByMonitorID(ctx context.Context, monitorID string) ([]*Model, error)
	DeleteByMonitorID(ctx context.Context, monitorID string) error
	DeleteByMaintenanceID(ctx context.Context, maintenanceID string) error
}

type ServiceImpl struct {
	repository Repository
	logger     *zap.SugaredLogger
}

func NewService(
	repository Repository,
	logger *zap.SugaredLogger,
) Service {
	return &ServiceImpl{
		repository,
		logger.Named("[monitor-maintenance-service]"),
	}
}

func (mr *ServiceImpl) Create(ctx context.Context, monitorID string, maintenanceID string) (*Model, error) {
	createModel := &Model{
		MonitorID:     monitorID,
		MaintenanceID: maintenanceID,
	}

	return mr.repository.Create(ctx, createModel)
}

func (mr *ServiceImpl) FindByID(ctx context.Context, id string) (*Model, error) {
	return mr.repository.FindByID(ctx, id)
}

func (mr *ServiceImpl) Delete(ctx context.Context, id string) error {
	return mr.repository.Delete(ctx, id)
}

func (mr *ServiceImpl) FindByMonitorID(ctx context.Context, monitorID string) ([]*Model, error) {
	return mr.repository.FindByMonitorID(ctx, monitorID)
}

func (mr *ServiceImpl) DeleteByMonitorID(ctx context.Context, monitorID string) error {
	return mr.repository.DeleteByMonitorID(ctx, monitorID)
}

func (mr *ServiceImpl) DeleteByMaintenanceID(ctx context.Context, maintenanceID string) error {
	return mr.repository.DeleteByMaintenanceID(ctx, maintenanceID)
}
