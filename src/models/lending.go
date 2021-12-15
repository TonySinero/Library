package models

type Lending struct {

	ModelBase
	Name         string  `json:"name"`
	Surname      string  `json:"surname"`
	SecondName   string  `json:"second-name"`
	Passport     string  `json:"passport"`
	Books        string  `json:"books"`
	ReturnDate   string  `json:"return-date"`
	Price        float32 `json:"price"`
}

// CreateLendingInput for model book create.

type CreateLendingInput struct {

	Name         string  `json:"name"         binding:"required"`
	Surname      string  `json:"surname"      binding:"required"`
	SecondName   string  `json:"second-name"  binding:"required"`
	Passport     string  `json:"passport"     binding:"required"`
	Books        string  `json:"books"        binding:"required"`
	ReturnDate   string  `json:"return-date"  binding:"required"`
	Price        float32 `json:"price"`
}
