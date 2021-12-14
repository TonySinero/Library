
package models

import (
	"github.com/library/src/utils"
)

// Admin model.

type Admin struct {
	ModelBase
	Email    string `json:"user"`
	Password string `json:"password"`
}

// UserInput for userModel.

type AdminInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// EncryptPassword method.

func (u *Admin) EncryptPassword() {
	passwordEncrypt := utils.Encrypt(u.Password)
	u.Password = passwordEncrypt
}

// ValidatePassword method.

func (u *Admin) ValidatePassword(password string) bool {
	if err := utils.CompareHash(u.Password, password); err != nil {
		return false
	}
	return true
}
