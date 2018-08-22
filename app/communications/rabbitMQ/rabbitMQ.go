package rabbitMQ

import (
	"github.com/streadway/amqp"
	"theAmazingCodeExample/app/common"
	"math/rand"
)

type RabbitMQTask interface{
	GetMessageBytes() ([]byte,error)
	GetQueue() (string)
}

func PublishMessageOnQueue(newTask RabbitMQTask) (error) {

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

func RPCcall(newTask RabbitMQTask) ([]byte,error){

	var response []byte

	ch := common.GetRabbitMQChannel()
	defer ch.Close()

	checkPhoneQueue, err := ch.QueueDeclare(
		newTask.GetQueue(),
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return response,err
	}

	anonQueue, err := ch.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil,
	)
	if err != nil {
		return response,err
	}

	anonQueueMessages, err := ch.Consume(
		anonQueue.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return response,err
	}

	corrId := randomString(32)

	messageBody, err := newTask.GetMessageBytes()
	if err != nil {
		return response,err
	}

	err = ch.Publish(
		"",
		checkPhoneQueue.Name,
		false,
		false,
		amqp.Publishing{
			ContentType:   "text/plain",
			CorrelationId: corrId,
			ReplyTo:       anonQueue.Name,
			Body:          messageBody,
		})
	if err != nil {
		return response,err
	}

	for d := range anonQueueMessages {
		if corrId == d.CorrelationId {
			response = d.Body
			break
		}
	}

	return response,nil

}

func randomString(l int) string {
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		bytes[i] = byte(randInt(65, 90))
	}
	return string(bytes)
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}