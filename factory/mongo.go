package factory

import (
	"context"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

const mongoTimeout = 2 * time.Second

func CreateMongoClient(uri string) *mongo.Client {
	ctx, _ := context.WithTimeout(context.Background(), mongoTimeout)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		logrus.Fatalln("Failed to connect to MongoDB", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		logrus.Fatalln("Failed to ping MongoDB", err)
	}

	logrus.Infoln("Connected to MongoDB!")
	return client
}
