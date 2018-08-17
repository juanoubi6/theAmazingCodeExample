package rabbitMQ

import (
	"github.com/streadway/amqp"
	"theAmazingCodeExample/app/common"
)

var exchangeMap = map[string]string{
	"sms": "sms_exchange",
}

func PublishMessageOnExchange(newTask RabbitMQTask) error {

	ch := common.GetRabbitMQChannel()
	defer ch.Close()

	exchangeName,exchangeType,_ := newTask.GetConfigInfo()

	err := ch.ExchangeDeclare(
		exchangeName,
		exchangeType,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	messageBody, err := newTask.GetMessageBytes()
	if err != nil {
		return err
	}

	err = ch.Publish(
		exchangeName,
		"",
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         messageBody,
		})
	if err != nil {
		return err
	}

	return nil

}
