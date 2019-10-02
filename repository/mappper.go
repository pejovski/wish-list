package repository

import "github.com/pejovski/wish-list/domain"

func mapItemToDomainItem(item *Item) *domain.Item {
	return &domain.Item{
		Product: &domain.Product{
			ProductId: item.ProductId,
			Name:      item.Name,
			Brand:     item.Brand,
			Price:     item.Price,
			Image:     item.Image,
		},
		Active: item.Active,
	}
}
