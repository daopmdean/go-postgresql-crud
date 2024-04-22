package model

import "gorm.io/gorm"

type Book struct {
	gorm.Model
	Name      string
	Author    string
	Publisher string
}
