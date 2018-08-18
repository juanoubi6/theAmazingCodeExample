package rabbitMQ

import (
	"github.com/streadway/amqp"
	"theAmazingCodeExample/app/common"
)

type RabbitMQTask interface{
	GetMessageBytes() ([]byte,error)
	GetQueue() (string)
}

func PublishMessageOnQueue(newTask RabbitMQTask) error {

	ch := common.GetRabbitMQChannel()
	defer ch.Close()

	queueName := newTask.GetQueue()

	//Queue declared but not needed if created previously
	queue, err := ch.QueueDeclare(
		queueName,
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
		"",
		queue.Name,
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
