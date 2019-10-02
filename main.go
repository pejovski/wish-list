package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	_ "github.com/joho/godotenv/autoload"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/pejovski/wish-list/controller"
	"github.com/pejovski/wish-list/factory"
	"github.com/pejovski/wish-list/gateway/catalog"
	"github.com/pejovski/wish-list/pkg/logger"
	amqpReceiver "github.com/pejovski/wish-list/receiver/amqp"
	"github.com/pejovski/wish-list/repository"
	httpServer "github.com/pejovski/wish-list/server/http"
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

	wishRepository := repository.NewWish(mongoClient)
	catalogGateway := catalog.NewGateway(retryablehttp.NewClient(), os.Getenv("CATALOG_API_HOST"))

	wishController := controller.NewWish(wishRepository, catalogGateway)

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

	serverHandler := httpServer.NewHandler(wishController)
	serverRouter := httpServer.NewRouter(serverHandler)

	server := factory.CreateHttpServer(serverRouter, fmt.Sprintf(":%s", os.Getenv("APP_PORT")))
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.Fatalf(err.Error())
		}
	}()
	logrus.Infof("Server started at port: %s", os.Getenv("APP_PORT"))

	gracefulShutdown(server, mongoClient)
}

func gracefulShutdown(server *http.Server, mongoClient *mongo.Client) {
	// Create channel for shutdown signals.
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, os.Kill, syscall.SIGTERM)

	// Receive shutdown signals.
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), serverShutdownTimeout)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logrus.Errorf("Server shutdown failed: %s", err)
	}
	logrus.Println("Server exited properly")

	ctxMongo, cancelMongo := context.WithTimeout(context.Background(), mongoShutdownTimeout)
	defer cancelMongo()

	if err := mongoClient.Disconnect(ctxMongo); err != nil {
		logrus.Errorf("MongoDB shutdown failed: %s", err)
	}
	logrus.Println("MongoDB closed properly")
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
