package rabbitMQ

import (
	"bytes"
	"encoding/gob"
)

type RabbitMQTask interface{
	GetMessageBytes() ([]byte,error)
	GetConfigInfo() (string,string,string)
}

type SmsTask struct{
	Exchange string
	Queue	 string
	ExchangeType	 string
	PhoneNumber	string
	Message interface{}
}

func NewSmsTask (phoneNumber string, message string) SmsTask{
	return SmsTask{
		Exchange:"sms_exchange",
		Queue:"sms_queue",
		ExchangeType:"direct",
		PhoneNumber:phoneNumber,
		Message:message,
	}
}

func (t *SmsTask) GetMessageBytes () ([]byte,error){
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(t.Message)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (t *SmsTask) GetConfigInfo () (exchangeName string,queueName string,exchangeType string){
	return t.Exchange,t.Queue,t.ExchangeType
}
