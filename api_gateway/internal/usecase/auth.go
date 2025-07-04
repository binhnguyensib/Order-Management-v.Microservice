package usecase

import (
	"api_gateway/internal/domain"
	"api_gateway/utils"
	"context"
)

var _ domain.AuthUsecase = (*authUsecaseImpl)(nil)

type authUsecaseImpl struct {
	authRepo domain.AuthRepository
}

func NewAuthUsecase(authRepo domain.AuthRepository) domain.AuthUsecase {
	return &authUsecaseImpl{
		authRepo: authRepo,
	}
}

func (au *authUsecaseImpl) Register(ctx context.Context, req *domain.RegisterRequest) (*domain.User, error) {
	user := &domain.User{
		Email:    req.Email,
		Password: req.Password,
	}
	if err := au.authRepo.Create(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (au *authUsecaseImpl) Login(ctx context.Context, req *domain.LoginRequest) (*domain.LoginResponse, error) {
	user, err := au.authRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	token, err := utils.GenerateJWT(user.Email, user.User_id.Hex())
	if err != nil {
		return nil, err
	}
	return &domain.LoginResponse{
		AccessToken: token,
		User:        user,
	}, nil
}
