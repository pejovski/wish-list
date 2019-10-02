package factory

import (
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

func CreateAmqpChannel(url string) *amqp.Channel {
	conn, err := amqp.Dial(url)
	if err != nil {
		logrus.Fatalf("%s: %s", "Failed to connect to RabbitMQ", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		logrus.Fatalf("%s: %s", "Failed to open a channel", err)
	}

	chErr := make(chan *amqp.Error)

	go func() {
		x, ok := <-chErr
		if !ok {
			return
		}
		logrus.Fatalf("RabbitMQ channel closed. Reason: %s, Err: %s", x.Reason, x.Error())
	}()

	ch.NotifyClose(chErr)

	return ch

}
