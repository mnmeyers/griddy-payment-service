package repositories

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"griddy-payment-service/models/payment"
)

type PaymentRepository interface {
	GetPaymentsByCustomerId(ctx context.Context, customerId string) ([]payment.DAO, error)
	CreatePayment(ctx context.Context, dao payment.DAO) (*payment.DAO, error)
	UpdatePayments(ctx context.Context, customerId string, daos []payment.DAO) error
}

func (db *MongoDB) GetPaymentsByCustomerId(ctx context.Context, customerId string) ([]payment.DAO, error) {
	client, err := db.getClient()

	if err != nil {
		return nil, err
	}
	defer client.Disconnect(nil)

	filter := bson.M{"customerId": customerId}
	cursor, err := client.Database(DB_NAME).Collection(PAYMENTS_COLLECTION).Find(nil, filter)

	if err != nil {
		return nil, err
	}

	var payments []payment.DAO
	err = cursor.All(ctx, &payments)

	if err != nil {
		return nil, err
	}

	return payments, nil
}

func (db *MongoDB) CreatePayment(ctx context.Context, dao payment.DAO) (*payment.DAO, error) {

	client, err := db.getClient()

	if err != nil {
		return nil, err
	}

	defer client.Disconnect(nil)

	// First value is the id of the inserted record which is not part of the DTO
	id, err := client.Database(DB_NAME, nil).Collection(PAYMENTS_COLLECTION).InsertOne(nil, dao)
	if err != nil {
		fmt.Printf("Error occurred on saving payment to db: %v", err)
		return nil, err
	}

	fmt.Printf("Payment successfully inserted into db with id: %v", id)

	return &dao, nil
}

func (db *MongoDB) UpdatePayments(ctx context.Context, customerId string, daos []payment.DAO) error {
	client, err := db.getClient()

	if err != nil {
		return err
	}

	defer client.Disconnect(nil)

	var writes []mongo.WriteModel

	for _, dao:=  range daos {
		filter := bson.D{
			{"customerId", customerId},
			{"paymentId", dao.PaymentId},
		}
		update := bson.D{
			{
				"$set", bson.D{
					{"status", dao.Status},
					{"customerId", customerId},
					{"amount", dao.Amount},
					{"paymentId", dao.PaymentId},
				},
			},
		}
		writes = append(writes, mongo.NewUpdateOneModel().SetFilter(filter).SetUpdate(update).SetUpsert(true))
	}

	opts := options.BulkWrite().SetOrdered(false)
	collection := client.Database(DB_NAME).Collection(PAYMENTS_COLLECTION)
	res, err := collection.BulkWrite(nil, writes, opts)
	if err != nil {
		fmt.Printf("Error on bulk update for customer %s with error: %v", customerId, err)
		return err
	}

	fmt.Printf("Upsert results for customer %s: %v", customerId, res)
	return nil
}