package twilio

import (
	"github.com/subosito/twilio"
	"theAmazingCodeExample/app/config"
)

var (
	AccountSid   = config.GetConfig().TWILIO_SID
	AuthToken    = config.GetConfig().TWILIO_AUTH_TOKEN
	AccountPhone = config.GetConfig().TWILIO_ACC_PHONE
)

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
