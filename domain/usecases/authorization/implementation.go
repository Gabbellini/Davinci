package authorization

import (
	"base/domain/entities"
	"base/infrastructure/repositories/authorization"
	"base/view/http_error"
	"context"
	"strings"
)

type useCases struct {
	repository authorization.Repository
}

func NewUseCases(
	repository authorization.Repository,
) UseCases {
	return &useCases{
		repository: repository,
	}
}

func (u useCases) Login(ctx context.Context, credential entities.Credential) error {
	if strings.TrimSpace(credential.Email) == "" {
		return http_error.NewForbiddenError("Credenciais inválidas")
	}

	if strings.TrimSpace(credential.Password) == "" {
		return http_error.NewForbiddenError("Credenciais inválidas")
	}

	return u.repository.Login(ctx, credential)
}
