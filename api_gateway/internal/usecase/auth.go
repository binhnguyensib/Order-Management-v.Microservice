package usecase

import (
	"api_gateway/internal/domain"
	"api_gateway/utils"
	"context"
	"errors"

	"github.com/sirupsen/logrus"
)

const (
	CustomerCreationRoutingKey = "customer.create"
)

var (
	_      domain.AuthUsecase = (*authUsecaseImpl)(nil)
	Logger                    = logrus.New()
)

type authUsecaseImpl struct {
	authRepo      domain.AuthRepository
	messageBroker domain.MessageBroker
}

func NewAuthUsecase(authRepo domain.AuthRepository, messageBroker domain.MessageBroker) domain.AuthUsecase {
	return &authUsecaseImpl{
		authRepo:      authRepo,
		messageBroker: messageBroker,
	}
}

func (au *authUsecaseImpl) Register(ctx context.Context, req *domain.RegisterRequest) (*domain.User, error) {
	existingUser, _ := au.authRepo.GetByEmail(ctx, req.Email)
	if existingUser != nil {
		Logger.Error("email already exists")
		return nil, errors.New("email already exist")
	} else {
		user := &domain.User{
			Email:    req.Email,
			Password: req.Password,
			Name:     req.Name,
		}
		if err := user.HashPassword(); err != nil {
			return nil, err
		}

		if err := au.authRepo.Create(ctx, user); err != nil {
			return nil, err
		}
		if au.messageBroker != nil {
			event := domain.Event{
				Type: domain.UserRegistered,
				Payload: domain.UserRegisteredPayload{
					UserID: user.User_id.Hex(),
					Email:  user.Email,
					Name:   user.Name,
				},
			}
			go func() {
				if err := au.messageBroker.Publish(context.Background(), CustomerCreationRoutingKey, event); err != nil {
					Logger.Errorf("Failed to publish customer creation event: %v", err)
				}
			}()
		}
		return user, nil
	}

}

func (au *authUsecaseImpl) RegisterAdmin(ctx context.Context, req *domain.RegisterRequest) (*domain.User, error) {
	existingUser, _ := au.authRepo.GetByEmail(ctx, req.Email)
	if existingUser != nil {
		Logger.Error("email already exists")
		return nil, errors.New("email already exist")
	} else {
		user := &domain.User{
			Email:    req.Email,
			Password: req.Password,
			Role:     "admin",
		}
		if err := user.HashPassword(); err != nil {
			return nil, err
		}

		if err := au.authRepo.Create(ctx, user); err != nil {
			return nil, err
		}
		return user, nil
	}

}

func (au *authUsecaseImpl) Login(ctx context.Context, req *domain.LoginRequest) (*domain.LoginResponse, error) {
	user, err := au.authRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if !user.CheckPassword(req.Password) {
		return nil, errors.New("invalid credentials")
	}
	token, err := utils.GenerateJWT(user.Email, user.User_id.Hex(), user.Role)
	if err != nil {
		return nil, err
	}
	return &domain.LoginResponse{
		AccessToken: token,
		User:        user,
	}, nil
}
