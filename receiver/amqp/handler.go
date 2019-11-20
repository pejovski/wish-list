package amqp

import (
	"encoding/json"
	"github.com/pejovski/wish-list/controller"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"time"
)

type Handler interface {
	ProductUpdated(d *amqp.Delivery)
	ProductDeleted(d *amqp.Delivery)
	ProductPriceUpdated(d *amqp.Delivery)
}

type handler struct {
	controller controller.Controller
}

func NewHandler(c controller.Controller) Handler {
	s := handler{
		controller: c,
	}

	return s
}

func (h handler) ProductUpdated(d *amqp.Delivery) {

	msg := struct {
		Id string `json:"id"`
	}{}

	err := json.Unmarshal(d.Body, &msg)
	if err != nil {
		logrus.Errorln("Failed to read body", err)
		h.reject(d)
		return
	}

	err = h.controller.UpdateProduct(msg.Id)
	if err != nil {
		logrus.Errorln("Failed to update product", err)
		h.reject(d)
		return
	}

	logrus.Infof("Product %s successfully updated", msg.Id)
	h.ack(d)
}

func (h handler) ProductDeleted(d *amqp.Delivery) {
	msg := struct {
		Id string `json:"id"`
	}{}

	err := json.Unmarshal(d.Body, &msg)
	if err != nil {
		logrus.Errorln("Failed to read body", err)
		h.reject(d)
		return
	}

	err = h.controller.DeleteProduct(msg.Id)
	if err != nil {
		logrus.Errorln("Failed to delete product", err)
		h.reject(d)
		return
	}

	logrus.Infof("Product %s successfully deleted", msg.Id)
	h.ack(d)
}

func (h handler) ProductPriceUpdated(d *amqp.Delivery) {
	msg := struct {
		Id    string  `json:"id"`
		Price float32 `json:"price"`
	}{}

	err := json.Unmarshal(d.Body, &msg)
	if err != nil {
		logrus.Errorln("Failed to read body", err)
		h.reject(d)
		return
	}

	err = h.controller.UpdateProductPrice(msg.Id, msg.Price)
	if err != nil {
		logrus.Errorln("Failed to update product price", err)
		h.reject(d)
		return
	}

	logrus.Infof("Price of product %s successfully updated", msg.Id)
	h.ack(d)
}

func (h handler) reject(d *amqp.Delivery) {
	time.Sleep(5 * time.Second)
	if err := d.Reject(true); err != nil {
		logrus.Errorln("Failed to reject msg", err)
	}
}

func (h handler) ack(d *amqp.Delivery) {
	if err := d.Ack(false); err != nil {
		logrus.Errorln("Failed to ack msg", err)
	}
}
