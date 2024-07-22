package cartstore

import "gorm.io/gorm"

type Item struct {
	gorm.Model
	Id        string `gorm:"uniqueIndex"`
	UserID    string
	ProductID string
	Quantity  int32
}
