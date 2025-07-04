package repository

import (
	"api_gateway/internal/domain"
	"context"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

var _ domain.AuthRepository = (*authRepositoryImpl)(nil)

type authRepositoryImpl struct {
	db *mongo.Database
}

func NewAuthRepository(db *mongo.Database) domain.AuthRepository {
	return &authRepositoryImpl{
		db: db,
	}
}

func (ar *authRepositoryImpl) Create(ctx context.Context, user *domain.User) error {
	collection := ar.db.Collection("auth")
	_, err := collection.InsertOne(ctx, user)
	if err != nil {
		return err
	}

	return nil
}

func (ar *authRepositoryImpl) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	collection := ar.db.Collection("auth")
	var user domain.User
	err := collection.FindOne(ctx, map[string]interface{}{"email": email}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
