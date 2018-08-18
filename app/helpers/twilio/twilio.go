package twilio

import (
	"encoding/json"
	"theAmazingCodeExample/app/config"
	"theAmazingCodeExample/app/helpers/rabbitMQ"
	"theAmazingCodeExample/app/helpers/rabbitMQ/tasks"
	"errors"
)

var (
	AccountSid   = config.GetConfig().TWILIO_SID
	AuthToken    = config.GetConfig().TWILIO_AUTH_TOKEN
	AccountPhone = config.GetConfig().TWILIO_ACC_PHONE
)

type CheckPhoneResult struct {
	CountryCode string `json:"country_code"`
	PhoneNumber string `json:"phone_number"`
	Error 		string `json:"error"`
}

func ValidatePhoneNumber(number string) (bool, error, CheckPhoneResult) {

	var result CheckPhoneResult

	resp,err := rabbitMQ.RPCcall(tasks.NewPhoneCheckTask(number))
	if err != nil{
		return false,err,result
	}

	if err = json.Unmarshal(resp,&result); err != nil{
		return false,err,result
	}

	if result.Error == ""{
		return true,nil,result
	}else{
		return false,errors.New(result.Error),result
	}

}
