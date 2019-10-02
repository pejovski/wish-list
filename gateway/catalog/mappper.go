package catalog

import "github.com/pejovski/wish-list/domain"

func (g Gateway) mapProductToDomainProduct(p *Product) *domain.Product {
	return &domain.Product{
		ProductId: p.Id,
		Name:      p.Name,
		Brand:     p.Brand,
		Price:     p.Price,
		Image:     p.Image,
	}
}
