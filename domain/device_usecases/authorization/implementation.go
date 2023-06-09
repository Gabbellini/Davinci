package authorization

import (
	"context"
	"davinci/domain/entities"
	"davinci/infrastructure/device_repository/authorization"
	device_repository "davinci/infrastructure/device_repository/device"
	user_repository "davinci/infrastructure/device_repository/user"
	"davinci/view/http_error"
	"golang.org/x/crypto/bcrypt"
	"log"
	"strings"
)

type useCases struct {
	repository       authorization.Repository
	userRepository   user_repository.Repository
	deviceRepository device_repository.Repository
}

func NewUseCases(
	repository authorization.Repository,
	userRepository user_repository.Repository,
	deviceRepository device_repository.Repository,
) UseCases {
	return &useCases{
		repository:       repository,
		userRepository:   userRepository,
		deviceRepository: deviceRepository,
	}
}

func (u useCases) Login(ctx context.Context, credential entities.Credential) (*entities.Device, error) {
	if credential.Email = strings.TrimSpace(credential.Email); credential.Email == "" {
		return nil, http_error.NewForbiddenError(http_error.ForbiddenMessage)
	}

	if credential.Password = strings.TrimSpace(credential.Password); credential.Password == "" {
		return nil, http_error.NewForbiddenError(http_error.ForbiddenMessage)
	}

	user, err := u.userRepository.GetByEmail(ctx, credential.Email)
	if err != nil {
		log.Println("[Login] Error GetByEmail", err)
		return nil, http_error.NewForbiddenError(http_error.ForbiddenMessage)
	}

	if user.StatusCode == entities.StatusDeleted {
		return nil, http_error.NewForbiddenError(http_error.ForbiddenMessage)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Credential.Password), []byte(credential.Password))
	if err != nil {
		log.Println("[Login] Error CompareHashAndPassword", err)
		return nil, http_error.NewForbiddenError(http_error.ForbiddenMessage)
	}

	device, err := u.deviceRepository.GetByName(ctx, credential.DeviceName, user.Id)
	if err != nil {
		log.Println("[Login] Error GetByName", err)
		return nil, http_error.NewForbiddenError(http_error.ForbiddenMessage)
	}

	return device, nil
}
