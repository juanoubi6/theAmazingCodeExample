package sendgrid

import (
	"theAmazingCodeExample/app/models"
	"theAmazingCodeExample/app/communications/nats/messages"
	"theAmazingCodeExample/app/communications/nats"
	"encoding/json"
	"errors"
)

func SendGenericIndividualEmail(subjectValue string, messageValue string, userData models.User) error {

	var result messages.IndividualEmailSendResponse

	natsMessage := messages.IndividualEmailSendRequest{
		Subject:subjectValue,
		Message:messageValue,
		UserEmail:userData.Email,
		UserName:userData.Name,
	}

	response,err := nats.SendNatsMessage(natsMessage)
	if err != nil{
		return err
	}

	if err = json.Unmarshal(response,&result); err != nil{
		return err
	}

	if result.Error == ""{
		return nil
	}else{
		return errors.New(result.Error)
	}

}
