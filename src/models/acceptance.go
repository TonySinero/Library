package models

type Acceptance struct {

	ModelBase
	Name         string  `json:"name"`
	Surname      string  `json:"surname"`
	SecondName   string  `json:"second-name"`
	Passport     string  `json:"passport"`
	Books        string  `json:"books"`
	Condition    string  `json:"condition"`
	Rating       uint8   `json:"rating"`
	FinalPrice   float32 `json:"final-price"`
}

// CreateAcceptanceInput for model book create.

type CreateAcceptanceInput struct {

	Name         string  `json:"name"         binding:"required"`
	Surname      string  `json:"surname"      binding:"required"`
	SecondName   string  `json:"second-name"  binding:"required"`
	Passport     string  `json:"passport"     binding:"required"`
	Books        string  `json:"books"        binding:"required"`
	Condition    string  `json:"condition"    binding:"required"`
	Rating       uint8   `json:"rating"       binding:"required"`
	FinalPrice   float32 `json:"final-price"`

}
