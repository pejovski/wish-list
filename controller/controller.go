package controller

import (
	"github.com/pejovski/wish-list/domain"
	"github.com/sirupsen/logrus"
)

type Wish struct {
	wishRepository domain.WishRepository
	productGateway domain.CatalogGateway
}

func NewWish(r domain.WishRepository, g domain.CatalogGateway) Wish {
	return Wish{wishRepository: r, productGateway: g}
}

func (c Wish) AddItem(userId string, productId string) error {

	// get item from repo
	item, err := c.wishRepository.Item(userId, productId)
	if err != nil {
		logrus.Errorf("GetItem failed for product %s, user %s Error: %s", productId, userId, err)
		return err
	}

	// check if item exist from repo
	if item != nil {
		logrus.Errorf("Item exist failure for product %s, user %s Error: %s", productId, userId, domain.ErrItemAlreadyExist)
		return domain.ErrItemAlreadyExist
	}

	// item doesn't exist so create a new one
	err = c.wishRepository.CreateItem(userId, productId)
	if err != nil {
		logrus.Errorf("CreateItem failed for product %s, user %s Error: %s", productId, userId, err)
		return err
	}

	//update item async
	updateFailed := make(chan bool)
	go func(updateFiled chan bool) {

		// get product data from repo
		product, err := c.wishRepository.Product(productId)
		if err != nil {
			logrus.Errorf("Unexpected failure for product %s, user %s Error: %s", productId, userId, err)
			updateFiled <- true
			return
		}

		// check if product exist from repo
		if product != nil {
			logrus.Infof("Product %s exist in some wish-list", productId)

			// update item with product data
			err = c.wishRepository.UpdateItem(userId, product)
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

func (c Wish) RemoveItem(userId string, productId string) error {
	go func() {
		err := c.wishRepository.DeleteItem(userId, productId)
		if err != nil {
			logrus.Errorf("DeleteProduct failed for product %s, user %s Error: %s", productId, userId, err)
		}
	}()

	return nil
}

func (c Wish) GetList(userId string) (domain.List, error) {
	list, err := c.wishRepository.List(userId)
	if err != nil {
		logrus.Errorf("Get List failed for user %s Error: %s", userId, err)
		return nil, err
	}

	return list, nil
}

func (c Wish) UpdateProduct(productId string) error {

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

	err = c.wishRepository.UpdateProduct(product)
	if err != nil {
		logrus.Errorf("UpdateProduct failed for product %s. Error: %s", productId, err)
		return err
	}

	return nil
}

func (c Wish) DeactivateProduct(productId string) error {
	err := c.wishRepository.DeactivateProduct(productId)
	if err != nil {
		logrus.Errorf("DeactivateProduct failed for product %s. Error: %s", productId, err)
		return err
	}

	return nil
}

func (c Wish) DeleteProduct(productId string) error {
	err := c.wishRepository.DeleteProduct(productId)
	if err != nil {
		logrus.Errorf("DeleteProduct failed for product %s. Error: %s", productId, err)
		return err
	}

	return nil
}

func (c Wish) UpdateProductPrice(productId string, price float32) error {
	err := c.wishRepository.UpdateProductPrice(productId, price)
	if err != nil {
		logrus.Errorf("UpdateProductPrice failed for product %s. Error: %s", productId, err)
		return err
	}

	return nil
}
