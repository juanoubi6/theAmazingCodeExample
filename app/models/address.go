package models

import "theAmazingCodeExample/app/common"

type Address struct {
	ID           uint    `gorm:"primary_key"`
	Address      string  `gorm:"not null"`
	Floor        string  `gorm:"null"`
	Apartment    string  `gorm:"null"`
	MainAddress  bool    `gorm:"not null"`
	PostalCodeID uint    `gorm:"not null" json:"-"`
	Latitude     float64 `gorm:"type:float(10,6);"`
	Longitude    float64 `gorm:"type:float(10,6);"`
	PostalCode   PostalCode
	UserID       uint `gorm:"not null" json:"-"`
}

func (addressData *Address) Save() error {

	err := common.GetDatabase().Create(addressData).Error
	if err != nil {
		return err
	}

	return nil
}

func (addressData *Address) Modify() error {

	r := common.GetDatabase()

	err := r.Save(addressData).Error
	if err != nil {
		return err
	}

	return nil

}

func (addressData *Address) Delete() error {

	err := common.GetDatabase().Delete(addressData).Error
	if err != nil {
		return err
	}

	return nil

}

func GetAddressById(id uint) (Address, bool, error) {

	addressData := Address{}

	r := common.GetDatabase()

	r = r.Preload("PostalCode").Where("id = ?", id).First(&addressData)
	if r.RecordNotFound() {
		return addressData, false, nil
	}
	if r.Error != nil {
		return addressData, true, r.Error
	}

	return addressData, true, nil
}

func GetAddressesForUser(userID uint) ([]Address, error) {

	var addresses []Address

	db := common.GetDatabase().Where("user_id = ?", userID).Preload("PostalCode").Find(&addresses)
	if db.Error != nil {
		return addresses, db.Error
	}

	return addresses, nil

}
