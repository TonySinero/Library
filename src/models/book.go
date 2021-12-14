package models

// Book model.

type Book struct {
	ModelBase
	Title           string  `json:"title"`
	Price           float32 `json:"price"`
	Author          string  `json:"author"`
	Year            uint16  `json:"year"`
	NumberOfCopies  uint16  `json:"number-of-copies"`
	NumberOfPages   uint16  `json:"number-of-pages"`
	Image           string  `json:"image"`
	PricePerDay     float32 `json:"price-per-day"`
	CategoryID      uint    `json:"category"`
}

// CreateBookInput for model book create.

type CreateBookInput struct {

	Title           string  `json:"title"            binding:"required"`
	Price           float32 `json:"price"            binding:"required"`
	Author          string  `json:"author"           binding:"required"`
	Year            uint16  `json:"year"             binding:"required"`
	NumberOfCopies  uint16  `json:"number-of-copies" binding:"required"`
	NumberOfPages   uint16  `json:"number-of-pages"  binding:"required"`
	Image           string  `json:"image"            binding:"required"`
	PricePerDay     float32 `json:"price-per-day"    binding:"required"`
	CategoryID      uint    `json:"category"`
}
