package catalog

type Product struct {
	Id    string  `json:"id"`
	Name  string  `json:"name"`
	Brand string  `json:"brand"`
	Price float32 `json:"price"`
	Image string  `json:"image"`
}
