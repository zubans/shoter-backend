package services

import (
	"context"
	"errors"
	"shuter-go/internal/dto"
)

type UserService struct {
}

func New() *UserService {
	return &UserService{}
}

func (u *UserService) create(ctx context.Context, req dto.CredentialsRequest) error {
	return errors.New("hello")
}
