package models

type PostalCode struct {
	ID         uint   `gorm:"primary_key" json:"-"`
	PostalCode string `gorm:"not null"`
}
