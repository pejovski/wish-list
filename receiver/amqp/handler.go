package amqp

import (
	"encoding/json"
	"github.com/pejovski/wish-list/domain"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"time"
)

type Handler struct {
	controller domain.WishController
}

func NewHandler(c domain.WishController) *Handler {
	s := &Handler{
		controller: c,
	}

	return s
}

func (h Handler) ProductUpdated(d *amqp.Delivery) {

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

func (h Handler) ProductDeleted(d *amqp.Delivery) {
	msg := struct {
		Id string `json:"id"`
	}{}

	err := json.Unmarshal(d.Body, &msg)
	if err != nil {
		logrus.Errorln("Failed to read body", err)
		h.reject(d)
		return
	}

	err = h.controller.DeactivateProduct(msg.Id)
	if err != nil {
		logrus.Errorln("Failed to update product", err)
		h.reject(d)
		return
	}

	logrus.Infof("Product %s successfully deleted", msg.Id)
	h.ack(d)
}

func (h Handler) ProductPriceUpdated(d *amqp.Delivery) {
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

func (h Handler) reject(d *amqp.Delivery) {
	time.Sleep(5 * time.Second)
	if err := d.Reject(true); err != nil {
		logrus.Errorln("Failed to reject msg", err)
	}
}

func (h Handler) ack(d *amqp.Delivery) {
	if err := d.Ack(false); err != nil {
		logrus.Errorln("Failed to ack msg", err)
	}
}
