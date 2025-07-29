package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"product_service/internal/domain"
	"time"

	rdis "product_service/config"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var _ domain.ProductRepository = (*productRepositoryImpl)(nil)
var Logger = logrus.New()

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
	redisKey := fmt.Sprintf("product:%s", id)
	cached, err := rdis.Get(ctx, redisKey)
	if err == nil && cached != "" {
		var p domain.Product
		if ummarshalErr := json.Unmarshal([]byte(cached), &p); ummarshalErr == nil {
			Logger.Info("Cache hit for", redisKey)
			return &p, nil
		}
	}
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

	bytes, marshalErr := json.Marshal(product)
	if marshalErr == nil {
		go func() {
			ctxBg := context.Background()
			rdis.Set(ctxBg, redisKey, bytes, time.Hour)
		}()
	}
	Logger.WithFields(
		logrus.Fields{
			"method": "GetByID",
			"id":     id}).Info("Fetch product from DB")
	return &product, nil
}

func (pr *productRepositoryImpl) GetAndUpdateByName(ctx context.Context, name string, newPrice float64) (*domain.PriceUpdateResult, error) {
	collection := pr.conn.Collection("products")
	filter := bson.M{"name": name}
	update := bson.M{"$set": bson.M{"price": newPrice}}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.Before)
	result := collection.FindOneAndUpdate(ctx, filter, update, opts)

	if result.Err() != nil {
		if errors.Is(result.Err(), mongo.ErrNoDocuments) {
			return &domain.PriceUpdateResult{
				ProductName: name,
				NewPrice:    newPrice,
				Status:      "not_found",
				Message:     fmt.Sprintf("product with name %s not found", name),
			}, nil
		}
		return nil, result.Err()
	}

	var productBeforeUpdate domain.Product
	err := result.Decode(&productBeforeUpdate)
	if err != nil {
		return nil, err
	}

	Logger.WithFields(logrus.Fields{
		"method":   "GetAndUpdateByName",
		"name":     name,
		"oldPrice": productBeforeUpdate.Price,
		"newPrice": newPrice,
	}).Info("Product price updated successfully")

	return &domain.PriceUpdateResult{
		ProductName: productBeforeUpdate.Name,
		OldPrice:    productBeforeUpdate.Price,
		NewPrice:    newPrice,
		Status:      "success",
		Message:     "Product price updated successfully",
	}, nil
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

	Logger.WithFields(logrus.Fields{
		"method": "Create",
		"id":     productID,
	}).Info("Add new product successfully")

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
			return nil, result.Err()
		}
		return nil, result.Err()
	}

	Logger.WithFields(logrus.Fields{
		"method": "Update",
		"id":     id,
	}).Info("Update product successfully")

	var updatedProduct domain.Product
	if err := result.Decode(&updatedProduct); err != nil {
		return nil, err
	}

	redisKey := fmt.Sprintf("product:%s", id)
	go func() {
		ctxBg := context.Background()
		rdis.Del(ctxBg, redisKey)
	}()

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

	Logger.WithFields(logrus.Fields{
		"method": "Delete",
		"id":     id,
	}).Info("Delete product successfully")

	var deletedProduct domain.Product
	if err := result.Decode(&deletedProduct); err != nil {
		return nil, err
	}

	redisKey := fmt.Sprintf("product:%s", id)
	go func() {
		ctxBg := context.Background()
		rdis.Del(ctxBg, redisKey)
	}()

	return &deletedProduct, nil
}
