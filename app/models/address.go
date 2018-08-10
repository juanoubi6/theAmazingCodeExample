package models

type Address struct {
	ID           uint   `gorm:"primary_key"`
	Address      string `gorm:"not null"`
	Floor        string `gorm:"null"`
	Apartment    string `gorm:"null"`
	MainAddress  bool   `gorm:"not null"`
	PostalCodeID uint   `gorm:"not null" json:"-"`
	PostalCode   PostalCode
	UserID       uint `gorm:"not null" json:"-"`
}
