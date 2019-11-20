package amqp

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"log"
)

const (
	exProductUpdated      = "product_updated"
	exProductDeleted      = "product_deleted"
	exProductPriceUpdated = "product_price_updated"

	queueName = "wish-list"

	exKind        = "fanout"
	prefetchCount = 5
)

type Receiver interface {
	Receive()
}

type receiver struct {
	ch      *amqp.Channel
	handler Handler
}

func NewReceiver(ch *amqp.Channel, h Handler) Receiver {
	s := receiver{
		ch:      ch,
		handler: h,
	}

	return s
}

func (r receiver) Receive() {
	if err := r.ch.Qos(
		prefetchCount,
		0,
		false,
	); err != nil {
		logrus.Fatalln("Failed to set Qos", err)
	}

	exchanges := []string{exProductUpdated, exProductDeleted, exProductPriceUpdated}

	for _, ex := range exchanges {

		dCh := r.deliveryCh(ex)

		switch ex {
		case exProductUpdated:
			go func() {
				for d := range dCh {
					r.handler.ProductUpdated(&d)
				}
			}()
		case exProductDeleted:
			go func() {
				for d := range dCh {
					r.handler.ProductDeleted(&d)
				}
			}()

		case exProductPriceUpdated:
			go func() {
				for d := range dCh {
					r.handler.ProductPriceUpdated(&d)
				}
			}()
		default:
			return
		}
	}

}

func (r *receiver) deliveryCh(ex string) <-chan amqp.Delivery {
	queue := fmt.Sprintf("%s:%s", ex, queueName)

	_, err := r.ch.QueueDeclare(
		queue,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		logrus.Fatalln("Failed to declare a queue", err)
	}

	err = r.ch.ExchangeDeclare(
		ex,
		exKind,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("%s %s: %s", "Failed to declare an exchange", ex, err)
	}

	err = r.ch.QueueBind(
		queue,
		"",
		ex,
		false,
		nil,
	)
	if err != nil {
		logrus.Fatalln("Failed to bind a queue", err)
	}

	logrus.Infof("RabbitMQ queue %s declared\n", queue)

	msgs, err := r.ch.Consume(
		queue,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		logrus.Fatalln("Failed to register a consumer", err)
	}

	return msgs
}
