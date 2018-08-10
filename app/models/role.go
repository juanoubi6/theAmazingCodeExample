package models

import (
	"theAmazingCodeExample/app/common"
)

type Role struct {
	ID          uint         `gorm:"primary_key"`
	Description string       `gorm:"not null"`
	Permissions []Permission `gorm:"many2many:permission_x_role;"`
}

const (
	ADMIN = 1
	USER  = 2
)

func GetRoleById(id uint) (Role, bool, error) {

	var role Role

	r := common.GetDatabase().Where("id = ?", id).First(&role)
	if r.RecordNotFound() {
		return role, false, nil
	}

	if r.Error != nil {
		return role, true, r.Error
	}

	return role, true, nil
}

func GetRoles(limit int, offset int) ([]Role, int, error) {

	var roles []Role
	var quantity int

	//Get roles
	r := common.GetDatabase().Limit(limit).Offset(offset).Find(&roles)
	if r.Error != nil {
		return roles, 0, r.Error
	}

	//Get role quantity
	r = common.GetDatabase().Table("roles").Count(&quantity)
	if r.Error != nil {
		return roles, 0, r.Error
	}

	return roles, quantity, nil

}
