package user

import (
	"github.com/go-errors/errors"
	"theAmazingCodeExample/app/models"
	"theAmazingCodeExample/app/validators"
)

func Validate(password string, email string) error {
	if err := validatePassword(password); err != nil {
		return err
	}

	if err := ValidateEmail(email); err != nil {
		return err
	}

	return nil
}

func validatePassword(password string) error {
	if !validators.IsLongerOrEqualThan(password, 8) {
		return errors.New("The password must be at least 8 characters long")
	}
	if !validators.IsShorterOrEqualThan(password, 20) {
		return errors.New("The password can't be more than 20 characters long")
	}

	return nil
}

func ValidateEmail(emailValue string) error {
	if validators.IsEmpty(emailValue) {
		return errors.New("No email informed")
	}
	if !validators.IsEmail(emailValue) {
		return errors.New("Invalid email address")
	}
	exists, _ := models.CheckEmailExistence(emailValue)
	if exists {
		return errors.New("Email already in use")
	}

	return nil
}

func validateRecoveryData(email string, recoveryCode string) (models.User, error) {

	userData, found, err := models.GetUserByEmail(email)
	if found == false {
		return models.User{}, errors.New("No email informed")
	}
	if err != nil {
		return models.User{}, err
	}
	if userData.PasswordRecoveryCode != recoveryCode {
		return models.User{}, errors.New("Invalid password recovery code")
	}

	return userData, nil
}

func validateLiteralEmail(email string) error {
	if validators.IsEmpty(email) {
		return errors.New("No email informed")
	}
	if !validators.IsEmail(email) {
		return errors.New("Invalid email")
	}

	return nil
}
