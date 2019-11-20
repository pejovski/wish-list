package main

import (
	"fmt"
	"github.com/pejovski/wish-list/pkg/signals"
	mongo2 "github.com/pejovski/wish-list/repository/mongo"
	"github.com/pejovski/wish-list/server/api"
	"os"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	_ "github.com/joho/godotenv/autoload"
	"github.com/pejovski/wish-list/controller"
	"github.com/pejovski/wish-list/factory"
	"github.com/pejovski/wish-list/gateway/catalog"
	"github.com/pejovski/wish-list/pkg/logger"
	amqpReceiver "github.com/pejovski/wish-list/receiver/amqp"
	"github.com/sirupsen/logrus"
)

const (
	serverShutdownTimeout = 3 * time.Second
	mongoShutdownTimeout  = 2 * time.Second
)

func init() {
	initLogger()
}

func main() {
	mongoClient := factory.CreateMongoClient(fmt.Sprintf(
		"mongodb://%s:%s",
		os.Getenv("MONGO_HOST"),
		os.Getenv("MONGO_PORT"),
	))

	wishRepository := mongo2.NewRepository(mongoClient)
	catalogGateway := catalog.NewGateway(retryablehttp.NewClient(), os.Getenv("CATALOG_API_HOST"))

	wishController := controller.New(wishRepository, catalogGateway)

	amqpCh := factory.CreateAmqpChannel(fmt.Sprintf(
		"amqp://%s:%s@%s:%s/%s",
		os.Getenv("RABBITMQ_USER"),
		os.Getenv("RABBITMQ_PASSWORD"),
		os.Getenv("RABBITMQ_HOST"),
		os.Getenv("RABBITMQ_PORT"),
		os.Getenv("RABBITMQ_VHOST"),
	))

	amqpHandler := amqpReceiver.NewHandler(wishController)
	receiver := amqpReceiver.NewReceiver(amqpCh, amqpHandler)
	// Receive events in goroutines
	receiver.Receive()

	// ToDo mongo shutdown
	ctx := signals.Context()

	serverAPI := api.NewServer(wishController)
	serverAPI.Run(ctx)

	logrus.Infof("allowing %s for graceful shutdown to complete", serverShutdownTimeout)
	<-time.After(serverShutdownTimeout)
}

func initLogger() {
	file, err := os.OpenFile("logstash.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		logrus.Panicln("Failed to log to file, using default stderr.", err.Error())
	}

	logrus.SetReportCaller(true)
	logrus.AddHook(logger.New(file, logger.DefaultFormatter(logrus.Fields{"type": os.Getenv("APP_NAME"), "env": os.Getenv("APP_ENV")})))

	logrus.Infoln("Logger is initialized.")
}
