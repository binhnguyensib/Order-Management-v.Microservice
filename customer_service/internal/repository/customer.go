package repository

import (
	"context"
	rdis "customer_service/config"
	"customer_service/internal/domain"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var (
	_      domain.CustomerRepository = (*customerRepositoryImpl)(nil)
	Logger                           = logrus.New()
)

type customerRepositoryImpl struct {
	conn *mongo.Database
}

func NewCustomerRepository(db *mongo.Database) domain.CustomerRepository {
	return &customerRepositoryImpl{
		conn: db,
	}
}

func (cr *customerRepositoryImpl) GetAll(ctx context.Context) ([]*domain.Customer, error) {
	collection := cr.conn.Collection("customers")
	cursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())
	var customers []*domain.Customer
	for cursor.Next(context.TODO()) {
		var customer domain.Customer
		if err := cursor.Decode(&customer); err != nil {
			return nil, err
		}
		customers = append(customers, &customer)
	}

	return customers, nil
}

func (cr *customerRepositoryImpl) GetByID(ctx context.Context, id string) (*domain.Customer, error) {
	redisKey := fmt.Sprintf("customer:%s", id)
	cached, err := rdis.Get(ctx, redisKey)
	if err == nil && cached != "" {
		var c domain.Customer
		if ummarshalErr := json.Unmarshal([]byte(cached), &c); ummarshalErr == nil {
			Logger.Info("Cache hit for", redisKey)
			return &c, nil
		}
	}
	collection := cr.conn.Collection("customers")
	var customer domain.Customer
	ObjectID, ok := bson.ObjectIDFromHex(id)
	if ok != nil {
		log.Printf("Error converting ID to ObjectID: %v", ok)
		return nil, ok
	}
	err = collection.FindOne(ctx, bson.M{"_id": ObjectID}).Decode(&customer)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, err
		}
		return nil, err
	}
	bytes, marshalErr := json.Marshal(customer)
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
	return &customer, nil
}

func (cr *customerRepositoryImpl) Create(ctx context.Context, customer *domain.CustomerRequest) (*domain.Customer, error) {
	collection := cr.conn.Collection("customers")
	result, err := collection.InsertOne(ctx, customer)
	if err != nil {
		return nil, err
	}
	customerID, ok := result.InsertedID.(bson.ObjectID)
	if !ok {
		return nil, fmt.Errorf("failed to convert inserted ID to ObjectID")
	}

	Logger.WithFields(logrus.Fields{
		"method": "Create",
		"id":     customerID,
	}).Info("Add new customer successfully")

	createdCustomer := &domain.Customer{
		UserID: customer.UserID,
		Id:     customerID,
		Name:   customer.Name,
		Email:  customer.Email,
		Phone:  customer.Phone,
	}

	return createdCustomer, nil
}

func (cr *customerRepositoryImpl) Update(ctx context.Context, id string, customerReq *domain.CustomerRequest) (*domain.Customer, error) {
	collection := cr.conn.Collection("customers")
	ObjectID, ok := bson.ObjectIDFromHex(id)
	if ok != nil {
		log.Printf("Error converting ID to ObjectID: %v", ok)
		return nil, ok
	}

	updateFields := bson.M{}
	if customerReq.Name != "" {
		updateFields["name"] = customerReq.Name
	}
	if customerReq.Email != "" {
		updateFields["email"] = customerReq.Email
	}
	if customerReq.Phone != "" {
		updateFields["phone"] = customerReq.Phone
	}

	update := bson.M{"$set": updateFields}

	otps := options.FindOneAndUpdate().SetReturnDocument(options.After)
	result := collection.FindOneAndUpdate(ctx, bson.M{"_id": ObjectID}, update, otps)
	if result.Err() != nil {

		if result.Err() == mongo.ErrNoDocuments {
			return nil, result.Err()
		}
		return nil, result.Err()
	}

	Logger.WithFields(logrus.Fields{
		"method": "Update",
		"id":     id,
	}).Info("Update customer successfully")

	var updatedCustomer domain.Customer
	if err := result.Decode(&updatedCustomer); err != nil {
		return nil, err
	}

	redisKey := fmt.Sprintf("customer:%s", id)
	go func() {
		ctxBg := context.Background()
		rdis.Del(ctxBg, redisKey)
	}()

	return &updatedCustomer, nil
}

func (cr *customerRepositoryImpl) Delete(ctx context.Context, id string) (*domain.Customer, error) {
	collection := cr.conn.Collection("customers")
	ObjectID, ok := bson.ObjectIDFromHex(id)
	if ok != nil {
		log.Printf("Error converting ID to ObjectID: %v", ok)
		return nil, ok
	}

	var deletedCustomer domain.Customer

	err := collection.FindOneAndDelete(ctx, bson.M{"_id": ObjectID}).Decode(&deletedCustomer)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, err
		}
		return nil, err
	}
	Logger.WithFields(logrus.Fields{
		"method": "Delete",
		"id":     id,
	}).Info("Delete customer successfully")

	redisKey := fmt.Sprintf("customer:%s", id)
	go func() {
		ctxBg := context.Background()
		rdis.Del(ctxBg, redisKey)
	}()

	return &deletedCustomer, nil
}
