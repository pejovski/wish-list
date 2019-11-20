package model

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
