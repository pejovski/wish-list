package repository

import "github.com/pejovski/wish-list/model"

type Repository interface {
	Product(productId string) (*model.Product, error)
	UpdateProduct(product *model.Product) error
	DeactivateProduct(productId string) error
	DeleteProduct(productId string) error
	UpdateProductPrice(productId string, price float32) error

	Item(userId string, productId string) (*model.Item, error)
	CreateItem(userId string, productId string) error
	DeleteItem(userId string, productId string) error
	UpdateItem(userId string, product *model.Product) error

	List(userId string) (model.List, error)
}
