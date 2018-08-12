package sendgrid

import (
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"theAmazingCodeExample/app/config"
	"theAmazingCodeExample/app/models"
)

var OwnerEmail string = "contact@theAmazingCodeExaple.com"
var OwnerName string = "The Amazing Code Example"

func SendGenericIndividualEmail(subjectValue string, messageValue string, userData models.User) error {

	from := mail.NewEmail(OwnerName, OwnerEmail)
	subject := subjectValue
	to := mail.NewEmail(userData.Name, userData.Email)
	plainTextContent := "The Amazing Code Example"
	htmlContent := messageValue

	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	client := sendgrid.NewSendClient(config.GetConfig().SENDGRID_KEY_ID)

	response, err := client.Send(message)
	if err != nil {
		println(err.Error())
		return err
	} else {
		println(response.Body)
		return nil
	}

}
