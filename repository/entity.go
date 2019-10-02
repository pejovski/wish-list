package repository

type Item struct {
	ProductId string  `bson:"product_id"`
	Name      string  `bson:"name"`
	Brand     string  `bson:"brand"`
	Price     float32 `bson:"price"`
	Image     string  `bson:"image"`
	Active    bool    `bson:"active"`
}
