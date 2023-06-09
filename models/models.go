package models

import (
	"context"
	"errors"
	"log"
	"write/database"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ErrUnknownCollection = errors.New("unknown type of collection or document")
)

var databaseName string

func Init(mongoUrl, dbName string) {
	databaseName = dbName
	database.Init(mongoUrl, databaseName)
}

func getCollectionName(doc any) (string, error) {
	return READING_COLLECTION, nil
}

func Save(ctx context.Context, document any) (any, error) {
	collectionName, err := getCollectionName(document)
	if err != nil {
		return nil, err
	}

	res, err := database.RunQuery(func(client *mongo.Client) (interface{}, error) {
		collection := client.Database(databaseName).Collection(collectionName)

		_, err := collection.InsertOne(ctx, document)
		if err != nil {
			logrus.Errorf("error while inserting document for %s\n", collectionName)
			return nil, err
		}
		return nil, nil
	})

	return res, nil
}

func DecodeIntoReading(ctx context.Context, cursor *mongo.Cursor) ([]Reading, error) {
	readings := []Reading{}
	err := cursor.All(ctx, &readings)
	if err != nil {
		log.Printf("error while decoding into reading\n")
		return nil, err
	}
	return readings, nil
}

func Get(ctx context.Context, collectionName string, search any, saveIn any, findOptions *options.FindOptions) error {
	_, err := database.RunQuery(func(client *mongo.Client) (interface{}, error) {
		var err error
		collection := client.Database(databaseName).Collection(collectionName)

		cursor, err := collection.Find(ctx, search, findOptions)
		if err != nil {
			log.Printf("error while finding document for %s\n", collectionName)
			return nil, err
		}

		err = cursor.All(ctx, &saveIn)
		if err != nil {
			log.Printf("error while decoding into reading\n")
			return nil, err
		}
		if err != nil {
			log.Printf("error while decoding document for %s\n", collectionName)
			return nil, err
		}

		return nil, nil
	})

	if err != nil {
		log.Printf("error while running database query for %s\n", collectionName)
		return err
	}
	return nil
}

func Update(ctx context.Context, collectionName string, search, update any) error {
	_, err := database.RunQuery(func(client *mongo.Client) (interface{}, error) {
		var err error
		collection := client.Database(databaseName).Collection(collectionName)

		_, err = collection.UpdateMany(ctx, search, update)
		if err != nil {
			log.Printf("error while deleting document for %s\n", collectionName)
			return nil, err
		}

		return nil, nil
	})
	if err != nil {
		log.Printf("error while running database query for %s\n", collectionName)
		return err
	}
	return err
}

func Delete(ctx context.Context, collectionName string, search any) (int64, error) {
	res, err := database.RunQuery(func(client *mongo.Client) (interface{}, error) {
		var err error
		collection := client.Database(databaseName).Collection(collectionName)

		delRes, err := collection.DeleteMany(ctx, search)
		if err != nil {
			log.Printf("error while deleting document for %s\n", collectionName)
			return nil, err
		}

		return delRes.DeletedCount, nil
	})
	if err != nil {
		log.Printf("error while running database query for %s\n", collectionName)
		return 0, err
	}
	return res.(int64), err
}

func Aggregate(ctx context.Context, collectionName string, search any, decodeInto any, aggregateOptions *options.AggregateOptions) error {
	_, err := database.RunQuery(func(client *mongo.Client) (interface{}, error) {
		var err error
		collection := client.Database(databaseName).Collection(collectionName)

		cursor, err := collection.Aggregate(ctx, search, aggregateOptions)
		if err != nil {
			log.Printf("error while finding document for %s\n", collectionName)
			return nil, err
		}

		err = cursor.All(ctx, &decodeInto)
		if err != nil {
			log.Printf("error while decoding into reading\n")
			return nil, err
		}
		if err != nil {
			log.Printf("error while decoding document for %s\n", collectionName)
			return nil, err
		}

		return nil, nil
	})

	if err != nil {
		log.Printf("error while running database query for %s\n", collectionName)
		return err
	}
	return err
}
