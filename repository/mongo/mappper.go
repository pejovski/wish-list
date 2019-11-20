package mongo

import (
	"github.com/pejovski/wish-list/model"
)

func mapItemToDomainItem(item *Item) *model.Item {
	return &model.Item{
		Product: &model.Product{
			ProductId: item.ProductId,
			Name:      item.Name,
			Brand:     item.Brand,
			Price:     item.Price,
			Image:     item.Image,
		},
		Active: item.Active,
	}
}
