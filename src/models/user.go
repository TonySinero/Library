package models

type User struct {

	ID           uint    `json:"id" gorm:"primaryKey"`
	Name         string  `json:"name"`
	Surname      string  `json:"surname"`
	SecondName   string  `json:"second-name"`
	Passport     string  `json:"passport"`
	DateOfBirth  string  `json:"date-of-birth"`
	Email        string  `json:"email"`
	Address      string  `json:"address"`
	Indebtedness string  `json:"indebtedness"`
}

// CreateUserInput for model book create.

type CreateUserInput struct {

	Name         string  `json:"name"          binding:"required"`
	Surname      string  `json:"surname"       binding:"required"`
	SecondName   string  `json:"second-name"   binding:"required"`
	Passport     string  `json:"passport"      binding:"required"`
	DateOfBirth  string  `json:"date-of-birth" binding:"required"`
	Email        string  `json:"email"         binding:"required"`
	Address      string  `json:"address"       binding:"required"`
	Indebtedness string  `json:"indebtedness"`

}
