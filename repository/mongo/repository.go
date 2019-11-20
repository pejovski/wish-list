package mongo

import (
	"context"
	"github.com/pejovski/wish-list/model"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	repo "github.com/pejovski/wish-list/repository"
)

const (
	database   = "wish"
	collection = "items"
)

type repository struct {
	collection *mongo.Collection
}

func NewRepository(c *mongo.Client) repo.Repository {
	return repository{collection: c.Database(database).Collection(collection)}
}

// get product with full data
func (r repository) Product(productId string) (*model.Product, error) {
	filter := bson.M{
		"product_id": productId,
		"price": bson.M{
			"$exists": true,
		},
	}

	ctx, _ := context.WithTimeout(context.Background(), 1*time.Second)

	result := r.collection.FindOne(ctx, filter)
	if result.Err() != nil {

		if result.Err() == mongo.ErrNoDocuments {
			logrus.Infof("No document for product %s", productId)
			return nil, nil
		}

		logrus.Errorf("FindOne failed for product %s Error: %s", productId, result.Err())
		return nil, result.Err()
	}

	var product *model.Product
	err := result.Decode(&product)
	if err != nil {
		logrus.Errorf("FindOne failed for product %s. Error: %s", productId, err)
		return nil, err
	}

	// ToDo - check why productId was not set
	product.ProductId = productId

	return product, nil
}

func (r repository) UpdateProduct(product *model.Product) error {
	filter := bson.M{
		"product_id": bson.M{
			"$eq": product.ProductId,
		},
	}

	update := bson.M{"$set": bson.M{
		"name":  product.Name,
		"brand": product.Brand,
		"price": product.Price,
		"image": product.Image,
	}}

	ctx, _ := context.WithTimeout(context.Background(), 1*time.Second)
	_, err := r.collection.UpdateMany(
		ctx,
		filter,
		update,
	)

	if err != nil {
		logrus.Errorf("UpdateMany failed for product %s; Error: %s", product.ProductId, err)
		return err
	}

	return nil
}
func (r repository) DeactivateProduct(productId string) error {
	// ToDo
	return nil
}

func (r repository) DeleteProduct(productId string) error {
	ctx, _ := context.WithTimeout(context.Background(), 1*time.Second)
	_, err := r.collection.DeleteMany(ctx, bson.M{"product_id": productId})
	if err != nil {
		logrus.Errorf("DeleteMany failed for product %s; Error: %s", productId, err)
		return err
	}

	return nil
}

func (r repository) UpdateProductPrice(productId string, price float32) error {
	filter := bson.M{
		"product_id": bson.M{
			"$eq": productId,
		},
	}

	update := bson.M{"$set": bson.M{
		"price": price,
	}}

	ctx, _ := context.WithTimeout(context.Background(), 1*time.Second)
	_, err := r.collection.UpdateMany(
		ctx,
		filter,
		update,
	)

	if err != nil {
		logrus.Errorf("UpdateMany failed for product %s; Error: %s", productId, err)
		return err
	}

	return nil
}

func (r repository) Item(userId string, productId string) (*model.Item, error) {

	filter := bson.M{"user_id": userId, "product_id": productId}
	ctx, _ := context.WithTimeout(context.Background(), 1*time.Second)

	result := r.collection.FindOne(ctx, filter)
	if result.Err() != nil {

		if result.Err() == mongo.ErrNoDocuments {
			logrus.Infof("No document for product %s, user %s", productId, userId)
			return nil, nil
		}

		logrus.Errorf("FindOne failed for product %s, user %s Error: %s", productId, userId, result.Err())
		return nil, result.Err()
	}

	var item *model.Item
	err := result.Decode(&item)
	if err != nil {
		logrus.Errorf("FindOne failed for product %s, user %s Error: %s", productId, userId, err)
		return nil, err
	}

	return item, nil
}

func (r repository) CreateItem(userId string, productId string) error {

	ctx, _ := context.WithTimeout(context.Background(), 1*time.Second)
	_, err := r.collection.InsertOne(ctx, bson.M{"user_id": userId, "product_id": productId, "active": true})
	if err != nil {
		logrus.Errorf("InsertOne failed for product %s, user %s Error: %s", productId, userId, err)
		return err
	}

	return nil
}

func (r repository) DeleteItem(userId string, productId string) error {
	ctx, _ := context.WithTimeout(context.Background(), 1*time.Second)
	_, err := r.collection.DeleteOne(ctx, bson.M{"user_id": userId, "product_id": productId})
	if err != nil {
		logrus.Errorf("DeleteOne failed for product %s, user %s Error: %s", productId, userId, err)
		return err
	}

	return nil
}

func (r repository) UpdateItem(userId string, product *model.Product) error {

	filter := bson.M{
		"product_id": bson.M{
			"$eq": product.ProductId,
		},
		"user_id": bson.M{
			"$eq": userId,
		},
	}

	update := bson.M{"$set": bson.M{
		"name":  product.Name,
		"brand": product.Brand,
		"price": product.Price,
		"image": product.Image,
	}}

	ctx, _ := context.WithTimeout(context.Background(), 1*time.Second)
	_, err := r.collection.UpdateOne(
		ctx,
		filter,
		update,
	)

	if err != nil {
		logrus.Errorf("UpdateMany failed for product %s; Error: %s", product.ProductId, err)
		return err
	}

	return nil
}

func (r repository) List(userId string) (model.List, error) {

	list := model.List{}

	filter := bson.M{"user_id": userId}
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)

	cur, err := r.collection.Find(ctx, filter)
	if err != nil {
		logrus.Errorf("Find all failed for user %s Error: %s", userId, err)
		return nil, err
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {

		var item *Item
		err := cur.Decode(&item)
		if err != nil {
			logrus.Errorf("Find all decode failed for user %s Error: %s", userId, err)
			return nil, err
		}

		// ToDo create filter and move this out
		if item.Name == "" {
			continue
		}

		list = append(list, mapItemToDomainItem(item))
	}

	if err := cur.Err(); err != nil {
		logrus.Errorf("Find all failed for user %s Error: %s", userId, err)
		return nil, err
	}

	return list, nil
}
