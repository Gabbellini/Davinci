package device

import (
	"davinci/domain/entities"
	"context"
)

type Repository interface {
	GetByName(ctx context.Context, deviceName string, userId int64) (*entities.Device, error)
}
