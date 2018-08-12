package twilio

import (
	"encoding/json"
	"github.com/subosito/twilio"
	"net/http"
	"theAmazingCodeExample/app/config"
)

var (
	AccountSid   = config.GetConfig().TWILIO_SID
	AuthToken    = config.GetConfig().TWILIO_AUTH_TOKEN
	AccountPhone = config.GetConfig().TWILIO_ACC_PHONE
)

type CheckPhoneResult struct {
	CountryCode string `json:"country_code"`
	PhoneNumber string `json:"phone_number"`
}

func SendVerificationSMS(verificationCode string, to string) error {

	// Initialize twilio client
	c := twilio.NewClient(AccountSid, AuthToken, nil)

	// Send Message
	params := twilio.MessageParams{
		Body: "Your verification code is: " + verificationCode,
	}
	_, _, err := c.Messages.Send(AccountPhone, to, params)
	if err != nil {
		return err
	}

	return nil

}

func ValidatePhoneNumber(number string) (bool, error, CheckPhoneResult) {

	//Create client
	var result CheckPhoneResult

	client := &http.Client{}

	//Create request
	request, err := http.NewRequest(http.MethodGet, "https://"+AccountSid+":"+AuthToken+"@lookups.twilio.com/v1/PhoneNumbers/"+number, nil)
	if err != nil {
		return false, err, result
	}

	//Fetch Request
	response, err := client.Do(request)
	if err != nil {
		return false, err, result
	}
	defer response.Body.Close()

	json.NewDecoder(response.Body).Decode(&result)

	//Check response
	if response.StatusCode != http.StatusOK {
		return false, err, result
	} else {
		println(result.PhoneNumber)
		println(result.CountryCode)
		return true, nil, result
	}

}
