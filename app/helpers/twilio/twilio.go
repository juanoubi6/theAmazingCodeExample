package twilio

import (
	"encoding/json"
	"errors"
	"theAmazingCodeExample/app/communications/rabbitMQ"
	"theAmazingCodeExample/app/communications/rabbitMQ/tasks"
	"theAmazingCodeExample/app/config"
)

var (
	AccountSid   = config.GetConfig().TWILIO_SID
	AuthToken    = config.GetConfig().TWILIO_AUTH_TOKEN
	AccountPhone = config.GetConfig().TWILIO_ACC_PHONE
)

func ValidatePhoneNumber(number string) (bool, error, tasks.PhoneCheckTaskResponse) {

	var result tasks.PhoneCheckTaskResponse

	resp, err := rabbitMQ.RPCcall(tasks.NewPhoneCheckTask(number))
	if err != nil {
		return false, err, result
	}

	if err = json.Unmarshal(resp, &result); err != nil {
		return false, err, result
	}

	if result.Error == "" {
		return true, nil, result
	} else {
		return false, errors.New(result.Error), result
	}

}
