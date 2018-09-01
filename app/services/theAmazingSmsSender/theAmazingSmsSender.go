package theAmazingSmsSender

import (
	"encoding/json"
	"errors"
	"theAmazingCodeExample/app/communications/rabbitMQ"
	"theAmazingCodeExample/app/communications/rabbitMQ/tasks"
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

func SendSms(phoneNumber string, code string) error {

	if err := rabbitMQ.PublishMessageOnQueue(tasks.NewSmsTask(phoneNumber, code)); err != nil {
		return err
	}

	return nil

}
