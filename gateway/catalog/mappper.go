package catalog

import (
	"github.com/pejovski/wish-list/model"
)

func (g gateway) mapProductToDomainProduct(p *Product) *model.Product {
	return &model.Product{
		ProductId: p.Id,
		Name:      p.Name,
		Brand:     p.Brand,
		Price:     p.Price,
		Image:     p.Image,
	}
}
