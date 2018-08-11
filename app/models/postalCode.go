package models

import "theAmazingCodeExample/app/common"

type PostalCode struct {
	ID         uint   `gorm:"primary_key" json:"-"`
	PostalCode string `gorm:"not null"`
}

func GetPostalCodeByCode(postalCode string) (PostalCode, bool, error) {

	existingPostalCode := PostalCode{}

	db := common.GetDatabase()

	r := db.Where("postal_code = ?", postalCode).First(&existingPostalCode)
	if r.RecordNotFound() {
		return existingPostalCode, false, nil
	}
	if r.Error != nil {
		return existingPostalCode, true, r.Error
	}

	return existingPostalCode, true, nil

}