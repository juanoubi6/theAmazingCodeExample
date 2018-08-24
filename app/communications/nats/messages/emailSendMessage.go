package messages

import "encoding/json"

type IndividualEmailSendRequest struct {
	Subject   string `json:"subject"`
	Message   string `json:"message"`
	UserName  string `json:"user_name"`
	UserEmail string `json:"user_email"`
}

type IndividualEmailSendResponse struct {
	Message string `json:"message"`
	Error   string `json:"error"`
}

func (e IndividualEmailSendRequest) GetTopic() string {
	return "emailSender:sendIndividualEmail"
}

func (e IndividualEmailSendRequest) GetMessageBytes() ([]byte, error) {
	data, err := json.Marshal(e)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (e IndividualEmailSendResponse) GetMessageBytes() ([]byte, error) {
	data, err := json.Marshal(e)
	if err != nil {
		return nil, err
	}

	return data, nil
}
