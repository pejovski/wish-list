package controller

import (
	"github.com/pejovski/wish-list/gateway/catalog"
	"github.com/pejovski/wish-list/repository"

	myerr "github.com/pejovski/wish-list/error"
	"github.com/pejovski/wish-list/model"
	"github.com/sirupsen/logrus"
)

type Controller interface {
	UpdateProduct(productId string) error
	DeactivateProduct(productId string) error
	DeleteProduct(productId string) error
	UpdateProductPrice(productId string, price float32) error

	AddItem(userId string, productId string) error
	RemoveItem(userId string, productId string) error

	GetList(userId string) (model.List, error)
}

type controller struct {
	repository     repository.Repository
	productGateway catalog.Gateway
}

func New(r repository.Repository, g catalog.Gateway) Controller {
	return controller{repository: r, productGateway: g}
}

func (c controller) AddItem(userId string, productId string) error {

	// get item from repo
	item, err := c.repository.Item(userId, productId)
	if err != nil {
		logrus.Errorf("GetItem failed for product %s, user %s Error: %s", productId, userId, err)
		return err
	}

	// check if item exist from repo
	if item != nil {
		logrus.Errorf("Item exist failure for product %s, user %s Error: %s", productId, userId, myerr.ErrItemAlreadyExist)
		return myerr.ErrItemAlreadyExist
	}

	// item doesn't exist so create a new one
	err = c.repository.CreateItem(userId, productId)
	if err != nil {
		logrus.Errorf("CreateItem failed for product %s, user %s Error: %s", productId, userId, err)
		return err
	}

	//update item async
	updateFailed := make(chan bool)
	go func(updateFiled chan bool) {

		// get product data from repo
		product, err := c.repository.Product(productId)
		if err != nil {
			logrus.Errorf("Unexpected failure for product %s, user %s Error: %s", productId, userId, err)
			updateFiled <- true
			return
		}

		// check if product exist from repo
		if product != nil {
			logrus.Infof("Product %s exist in some wish-list", productId)

			// update item with product data
			err = c.repository.UpdateItem(userId, product)
			if err != nil {
				logrus.Errorf("Unexpected failure for product %s, user %s Error: %s", productId, userId, err)
				updateFiled <- true
				return
			}

			// item was updated with product data
			// no failure
			updateFiled <- false
			return
		}

		err = c.UpdateProduct(productId)
		if err != nil {
			logrus.Errorf("UpdateProduct async failed for product %s. Error: %s", productId, err)
			updateFiled <- true
			return
		}

		// item was updated with product data
		// no failure
		updateFiled <- false
	}(updateFailed)

	// waiting for updateFailed result in a goroutine
	go func(updateFiled chan bool) {
		// if update failed, remove the item
		if <-updateFiled {
			_ = c.RemoveItem(userId, productId)
		}
	}(updateFailed)

	return nil
}

func (c controller) RemoveItem(userId string, productId string) error {
	go func() {
		err := c.repository.DeleteItem(userId, productId)
		if err != nil {
			logrus.Errorf("DeleteProduct failed for product %s, user %s Error: %s", productId, userId, err)
		}
	}()

	return nil
}

func (c controller) GetList(userId string) (model.List, error) {
	list, err := c.repository.List(userId)
	if err != nil {
		logrus.Errorf("Get List failed for user %s Error: %s", userId, err)
		return nil, err
	}

	return list, nil
}

func (c controller) UpdateProduct(productId string) error {

	// get product data from external domain
	product, err := c.productGateway.Product(productId)
	if err != nil {
		logrus.Errorf("Get Product failed for product %s. Error: %s", productId, err)
		return err
	}

	// no product found
	// product could has been deleted
	if product == nil {
		return c.DeactivateProduct(productId)
	}

	err = c.repository.UpdateProduct(product)
	if err != nil {
		logrus.Errorf("UpdateProduct failed for product %s. Error: %s", productId, err)
		return err
	}

	return nil
}

func (c controller) DeactivateProduct(productId string) error {
	err := c.repository.DeactivateProduct(productId)
	if err != nil {
		logrus.Errorf("DeactivateProduct failed for product %s. Error: %s", productId, err)
		return err
	}

	return nil
}

func (c controller) DeleteProduct(productId string) error {
	err := c.repository.DeleteProduct(productId)
	if err != nil {
		logrus.Errorf("DeleteProduct failed for product %s. Error: %s", productId, err)
		return err
	}

	return nil
}

func (c controller) UpdateProductPrice(productId string, price float32) error {
	err := c.repository.UpdateProductPrice(productId, price)
	if err != nil {
		logrus.Errorf("UpdateProductPrice failed for product %s. Error: %s", productId, err)
		return err
	}

	return nil
}
