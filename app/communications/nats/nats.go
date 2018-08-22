package nats

import (
	"theAmazingCodeExample/app/common"
	"time"
)

type NatsMessage interface{
	GetTopic() (string)
	GetMessageBytes() ([]byte,error)
}

func SendNatsMessage(newMessage NatsMessage) ([]byte,error){

	conn := common.GetNatsConnection()

	messageBody,err := newMessage.GetMessageBytes()
	if err != nil{
		return nil, err
	}

	response, err := conn.Request(newMessage.GetTopic(),messageBody,20*time.Second)
	if err != nil{
		println(err.Error())
		return nil,err
	}

	return response.Data,nil

}