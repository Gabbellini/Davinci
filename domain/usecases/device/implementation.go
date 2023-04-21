package device

import (
	"base/domain/entities"
	"base/infrastructure/repositories/device"
	"context"
)

type useCases struct {
	deviceRepo device.Repository
}

func NewUseCases(deviceRepo device.Repository) UseCases {
	return &useCases{
		deviceRepo: deviceRepo,
	}
}

func (u useCases) Create(ctx context.Context, device entities.Device) error {
	return u.deviceRepo.Create(ctx, device)
}

func (u useCases) Update(ctx context.Context, device entities.Device) error {
	return u.deviceRepo.Update(ctx, device)
}

func (u useCases) Delete(ctx context.Context, device entities.Device) error {
	return u.deviceRepo.Delete(ctx, device)
}

func (u useCases) GetAll(ctx context.Context, idUser int64) ([]entities.Device, error) {
	return u.deviceRepo.GetAll(ctx, idUser)
}

func (u useCases) GetById(
	ctx context.Context,
	id int64,
	idUser int64,
) (entities.Device, error) {
	return u.deviceRepo.GetById(
		ctx,
		id,
		idUser,
	)
}
