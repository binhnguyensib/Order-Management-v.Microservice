package domain

import "context"

type AuthRepository interface {
	Create(ctx context.Context, user *User) error
	GetByEmail(ctx context.Context, email string) (*User, error)
}

type AuthUsecase interface {
	Register(ctx context.Context, req *RegisterRequest) (*User, error)
	RegisterAdmin(ctx context.Context, req *RegisterRequest) (*User, error)
	Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error)
}
