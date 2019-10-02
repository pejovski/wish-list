package domain

import (
	"errors"
	"github.com/streadway/amqp"
	"net/http"
)

var (
	ErrItemAlreadyExist = errors.New("item already exist")
)

type List []*Item

type Item struct {
	*Product
	Active bool `json:"active"`
}

type Product struct {
	ProductId string  `json:"product_id"`
	Name      string  `json:"name"`
	Brand     string  `json:"brand"`
	Price     float32 `json:"price"`
	Image     string  `json:"image"`
}

type Receiver interface {
	Receive()
}

type AmqpHandler interface {
	ProductUpdated(d *amqp.Delivery)
	ProductDeleted(d *amqp.Delivery)
	ProductPriceUpdated(d *amqp.Delivery)
}

type HttpHandler interface {
	GetList() http.HandlerFunc
	AddItem() http.HandlerFunc
	RemoveItem() http.HandlerFunc
}

type WishController interface {
	UpdateProduct(productId string) error
	DeactivateProduct(productId string) error
	DeleteProduct(productId string) error
	UpdateProductPrice(productId string, price float32) error

	AddItem(userId string, productId string) error
	RemoveItem(userId string, productId string) error

	GetList(userId string) (List, error)
}

type WishRepository interface {
	Product(productId string) (*Product, error)
	UpdateProduct(product *Product) error
	DeactivateProduct(productId string) error
	DeleteProduct(productId string) error
	UpdateProductPrice(productId string, price float32) error

	Item(userId string, productId string) (*Item, error)
	CreateItem(userId string, productId string) error
	DeleteItem(userId string, productId string) error
	UpdateItem(userId string, product *Product) error

	List(userId string) (List, error)
}

type CatalogGateway interface {
	Product(id string) (*Product, error)
}
