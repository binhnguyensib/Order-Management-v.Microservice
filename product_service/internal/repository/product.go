package repository

import (
	"context"
	"fmt"
	"product_service/internal/domain"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var _ domain.ProductRepository = (*productRepositoryImpl)(nil)

type productRepositoryImpl struct {
	conn *mongo.Database
}

func NewProductRepository(db *mongo.Database) domain.ProductRepository {
	return &productRepositoryImpl{
		conn: db,
	}
}

func (pr *productRepositoryImpl) GetAll(ctx context.Context) ([]*domain.Product, error) {
	collection := pr.conn.Collection("products")
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var products []*domain.Product
	for cursor.Next(ctx) {
		var product domain.Product
		if err := cursor.Decode(&product); err != nil {
			return nil, err
		}
		products = append(products, &product)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

func (pr *productRepositoryImpl) GetByID(ctx context.Context, id string) (*domain.Product, error) {
	collection := pr.conn.Collection("products")
	var product domain.Product
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	err = collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&product)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, err
		}
		return nil, err
	}

	return &product, nil
}

func (pr *productRepositoryImpl) Create(ctx context.Context, product *domain.ProductRequest) (*domain.Product, error) {
	collection := pr.conn.Collection("products")
	result, err := collection.InsertOne(ctx, product)
	if err != nil {
		return nil, err
	}
	productID, ok := result.InsertedID.(bson.ObjectID)
	if !ok {
		return nil, fmt.Errorf("failed to convert inserted ID to ObjectID: %v", result.InsertedID)
	}

	createdProduct := &domain.Product{
		Id:    productID,
		Name:  product.Name,
		Price: product.Price,
		Stock: product.Stock,
	}
	return createdProduct, nil
}

func (pr *productRepositoryImpl) Update(ctx context.Context, id string, productReq *domain.ProductRequest) (*domain.Product, error) {
	collection := pr.conn.Collection("products")
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	updateFields := bson.M{}
	if productReq.Name != "" {
		updateFields["name"] = productReq.Name
	}
	if productReq.Price > 0 {
		updateFields["price"] = productReq.Price
	}
	if productReq.Stock >= 0 {
		updateFields["stock"] = productReq.Stock
	}
	update := bson.M{"$set": updateFields}

	otps := options.FindOneAndUpdate().SetReturnDocument(options.After)
	result := collection.FindOneAndUpdate(ctx, bson.M{"_id": objectID}, update, otps)
	if result.Err() != nil {

		if result.Err() == mongo.ErrNoDocuments {
			return nil, err
		}
		return nil, result.Err()
	}

	var updatedProduct domain.Product
	if err := result.Decode(&updatedProduct); err != nil {
		return nil, err
	}

	return &updatedProduct, nil
}

func (pr *productRepositoryImpl) Delete(ctx context.Context, id string) (*domain.Product, error) {
	collection := pr.conn.Collection("products")
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	result := collection.FindOneAndDelete(ctx, bson.M{"_id": objectID})
	if result.Err() != nil {
		if result.Err() == mongo.ErrNoDocuments {
			return nil, err
		}
		return nil, result.Err()
	}

	var deletedProduct domain.Product
	if err := result.Decode(&deletedProduct); err != nil {
		return nil, err
	}

	return &deletedProduct, nil
}
