package service

import (
	"context"
	"time"
	"write/database"
	"write/logger"
	"write/models"

	"go.mongodb.org/mongo-driver/mongo"
	"go.opentelemetry.io/otel"
)

var (
	serviceTracer = otel.Tracer("write-service")
)

// WriteReading - writes reading to database
func WriteReading(ctx context.Context, reading models.Reading) error {
	ctx, span := serviceTracer.Start(ctx, "WriteReading")
	defer span.End()

	logEntry := logger.
		WithTracingFields(ctx).
		WithField("uid", reading.Uid)

	if reading.DateTime == 0 {
		reading.DateTime = time.Now().Unix()
	}
	_, err := database.RunQuery(func(client *mongo.Client) (interface{}, error) {
		collection := client.
			Database(database.DatabaseName).
			Collection(models.READING_COLLECTION)

		_, err := collection.InsertOne(ctx, reading)
		if err != nil {
			logEntry.Info(err)
			return nil, err
		}
		return nil, nil
	})
	if err != nil {
		logEntry.Error(err)
		return err
	}
	return nil
}
