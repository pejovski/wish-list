package repository

import (
	"testing"
)

func TestMapItemToDomainItem(t *testing.T) {

	i := &Item{
		ProductId: "",
		Name:      "Galaxy",
		Brand:     "Samsung",
		Price:     800,
		Image:     "galaxy.jpg",
		Active:    true,
	}

	di := mapItemToDomainItem(i)

	if i.ProductId != di.ProductId {
		t.Error("ProductId not equal")
	}
}
