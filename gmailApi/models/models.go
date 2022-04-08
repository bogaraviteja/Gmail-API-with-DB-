package models

import "gorm.io/gorm"

type (
	Person struct {
		gorm.Model
		Name    string
		Gender  string
		Email   string
		Address string
		Pincode int64
	}
	SentEmails struct {
		gorm.Model
		Name      string
		Email     string
		Subject   string
		Content   string
		Signature string
	}
)
