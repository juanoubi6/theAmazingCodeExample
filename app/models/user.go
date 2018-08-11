package models

import (
	"github.com/jinzhu/gorm"
	"github.com/rs/xid"
	"strings"
	"theAmazingCodeExample/app/common"
)

type User struct {
	gorm.Model
	GUID                 string `gorm:"type:char(20);unique_index:idx_unique_guid_object" json:"ID"`
	Name                 string `gorm:"null"`
	LastName             string `gorm:"null"`
	Password             string `gorm:"null" json:"-"`
	Email                string
	Phone                string         `gorm:"null"`
	PasswordRecoveryCode string         `gorm:"null" json:"-"`
	RoleID               uint           `gorm:"not null" json:"-"`
	Role                 Role           `gorm:"ForeignKey:RoleID"`
	Addresses            []Address      `gorm:"ForeignKey:UserID"`
	ProfilePicture       ProfilePicture `gorm:"ForeignKey:UserID"`
	Enabled              bool           `gorm:"default:true"`
}

type EmailConfirmation struct {
	ID     uint `gorm:"primary_key" json:"-"`
	UserID uint
	Code   string
}

type PhoneConfirmation struct {
	ID     uint `gorm:"primary_key" json:"-"`
	UserID uint
	Code   string
}

type ProfilePicture struct {
	ID     uint   `gorm:"primary_key"`
	Url    string `gorm:"not null"`
	S3Key  string `json:"-"`
	UserID uint   `json:"-"`
}

func (userData *User) Save() error {

	//Add GUID
	userData.GUID = xid.New().String()

	err := common.GetDatabase().Create(userData).Error
	if err != nil {
		return err
	}

	return nil
}

func (userData *User) Modify() error {

	r := common.GetDatabase()

	err := r.Save(&userData).Error
	if err != nil {
		return err
	}

	return nil

}

func (userData *User) Delete() error {

	err := common.GetDatabase().Delete(&userData).Error
	if err != nil {
		return err
	}

	return nil

}

func (userData *User) GetMainAddress() (Address, bool, error) {

	addressData := Address{}

	r := common.GetDatabase()

	r = r.Where("user_id = ? and main_address = ?", userData.ID, true).First(&addressData)
	if r.RecordNotFound() {
		return addressData, false, nil
	}
	if r.Error != nil {
		return addressData, true, r.Error
	}

	return addressData, true, nil

}

func (emailConfirmationData *EmailConfirmation) Save() error {

	err := common.GetDatabase().Create(emailConfirmationData).Error
	if err != nil {
		return err
	}

	return nil
}

func (profilePictureData *ProfilePicture) Save() error {

	err := common.GetDatabase().Create(profilePictureData).Error
	if err != nil {
		return err
	}

	return nil
}

func GetUserById(id uint) (user User, found bool, err error) {

	user = User{}

	r := common.GetDatabase()

	r = r.Unscoped().Preload("ProfilePicture").Preload("Role").Preload("Addresses").Preload("Addresses.PostalCode").Where("id = ?", id).First(&user)
	if r.RecordNotFound() {
		return user, false, nil
	}

	if r.Error != nil {
		return user, true, r.Error
	}

	return user, true, nil
}

func GetUserByGuid(guid string) (User, bool, error) {

	userData := User{}

	r := common.GetDatabase()

	r = r.Unscoped().Where("guid = ?", guid).Preload("Role").Preload("Addresses").First(&userData)
	if r.RecordNotFound() {
		return userData, false, nil
	}
	if r.Error != nil {
		return userData, true, r.Error
	}

	return userData, true, nil
}

func GetUserByEmail(email string) (User, bool, error) {

	user := User{}

	r := common.GetDatabase()

	r = r.Where("email = ?", email).First(&user)
	if r.RecordNotFound() {
		return user, false, nil
	}
	if r.Error != nil {
		return user, true, r.Error
	}

	return user, true, nil
}

func CheckEmailExistence(emailValue string) (bool, error) {

	r := common.GetDatabase().Where("email = ?", emailValue).First(&User{})
	if r.RecordNotFound() {
		return false, nil
	} else {
		return true, nil
	}

}

func GetUserPermissions(userID uint) ([]string, error) {

	userData := User{}
	var permissionList []string

	r := common.GetDatabase().Preload("Role").Preload("Role.Permissions").Where("id = ?", userID).First(&userData)
	if r.RecordNotFound() {
		return []string{}, nil
	}
	if r.Error != nil {
		return []string{}, r.Error
	}

	//For each permission, get it's description
	for _, permissionFound := range userData.Role.Permissions {
		permissionList = append(permissionList, permissionFound.Description)
	}

	return permissionList, nil

}

func CreateEmailConfirmation(newEmailConfirmation EmailConfirmation) (EmailConfirmation, error) {

	err := common.GetDatabase().Create(&newEmailConfirmation).Error
	if err != nil {
		return newEmailConfirmation, err
	}

	return newEmailConfirmation, nil
}

func CreateProfilePicture(newProfilePicture ProfilePicture) (ProfilePicture, error) {

	err := common.GetDatabase().Create(&newProfilePicture).Error
	if err != nil {
		return newProfilePicture, err
	}

	return newProfilePicture, nil
}

func GetUsers(limit int, offset int, idParam string, phoneParam string, emailParam string, nameParam string, lastNameParam string, roleIdsParam []string, columnToOrder string, order string) ([]User, int, error) {

	var users []User
	var quantity int

	//Get user ids
	var userIDs []uint
	db := common.GetDatabase()
	db = db.Table("users").Unscoped()

	if idParam != "" {
		db = db.Where("users.id = ?", idParam)
	}
	if phoneParam != "" {
		db = db.Where("lower(phone) like ?", "%"+strings.ToLower(phoneParam)+"%")
	}
	if emailParam != "" {
		db = db.Where("lower(email) like ?", "%"+strings.ToLower(emailParam)+"%")
	}
	if nameParam != "" {
		db = db.Where("lower(name) like ?", "%"+strings.ToLower(nameParam)+"%")
	}
	if lastNameParam != "" {
		db = db.Where("lower(last_name) like ?", "%"+strings.ToLower(lastNameParam)+"%")
	}
	if len(roleIdsParam) != 0 {
		db = db.Where("role_id in (?)", roleIdsParam)
	}
	db = db.Offset(offset).Limit(limit).Select("users.id").Pluck("users.id", &userIDs)
	if db.Error != nil {
		return users, quantity, db.Error
	}

	//Get users
	db = common.GetDatabase()
	if columnToOrder != "" && order != "" {
		db = db.Order(columnToOrder + " " + order)
	}
	db = db.Unscoped().Preload("Role").Preload("Addresses").Preload("Addresses.PostalCode").Where("id in (?)", userIDs).Find(&users)

	//Count user amount
	db = common.GetDatabase()
	db = db.Table("users").Unscoped()
	if idParam != "" {
		db = db.Where("users.id = ?", idParam)
	}
	if phoneParam != "" {
		db = db.Where("lower(phone) like ?", "%"+strings.ToLower(phoneParam)+"%")
	}
	if emailParam != "" {
		db = db.Where("lower(email) like ?", "%"+strings.ToLower(emailParam)+"%")
	}
	if nameParam != "" {
		db = db.Where("lower(name) like ?", "%"+strings.ToLower(nameParam)+"%")
	}
	if lastNameParam != "" {
		db = db.Where("lower(last_name) like ?", "%"+strings.ToLower(lastNameParam)+"%")
	}
	if len(roleIdsParam) != 0 {
		db = db.Where("role_id in (?)", roleIdsParam)
	}
	db = db.Count(&quantity)
	if db.Error != nil {
		return users, quantity, db.Error
	}

	return users, quantity, nil

}
