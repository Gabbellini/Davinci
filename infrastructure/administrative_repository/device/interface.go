package device

import (
	"context"
	"davinci/domain/entities"
)

type Repository interface {
	// Create insert a new device in the database.
	Create(ctx context.Context, device entities.Device, userId int64) error

	// Update a device in the database.
	Update(ctx context.Context, deviceId int64, device entities.Device, userId int64) error

	// Delete remove a device from the database.
	Delete(ctx context.Context, deviceId int64, userId int64) error

	// GetAll return all devices from the database.
	GetAll(ctx context.Context, userId int64) ([]entities.Device, error)

	// GetById return a device by id.
	GetById(ctx context.Context, deviceId int64, userId int64) (*entities.Device, error)

	// GetDeviceByName return the device by name.
	GetDeviceByName(ctx context.Context, deviceName string, userId int64) (*entities.Device, error)
}
