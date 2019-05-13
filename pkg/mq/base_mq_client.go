package mq

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

const deadLetterExchangeName = "DLX"
const deadLetterExchangeType = "direct"

type BaseMqClient struct {
	Channel      *amqp.Channel
	ExchangeName string
	ExchangeType string
	RoutingKey   string
}

func (client BaseMqClient) getDeadLetterQueueName() string {
	return fmt.Sprintf("%s.%s.dlq", client.ExchangeName, client.RoutingKey)
}

func (client BaseMqClient) getDeadLetterRoutingKey() string {
	return fmt.Sprintf("%s.%s", client.ExchangeName, client.RoutingKey)
}

func (client BaseMqClient) getQueueName() string {
	return fmt.Sprintf("%s.%s", client.ExchangeName, client.RoutingKey)
}

func (client BaseMqClient) declareDeadLetterRouting() error {
	log.Infof("Declaring rabbitmq dead letter exchange %s", deadLetterExchangeName)
	err := client.Channel.ExchangeDeclare(
		deadLetterExchangeName,
		deadLetterExchangeType,
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return err
	}

	deadLetterQueueName := client.getDeadLetterQueueName()
	log.Infof("Declaring rabbitmq dead letter queue %s", deadLetterQueueName)
	_, err = client.Channel.QueueDeclare(
		deadLetterQueueName,
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return err
	}

	deadLetterRoutingKey := client.getDeadLetterRoutingKey()
	log.Infof("Binding rabbitmq dead letter queue %s to exchange %s on routing key %s", deadLetterQueueName, deadLetterExchangeName, deadLetterRoutingKey)
	return client.Channel.QueueBind(
		deadLetterQueueName,
		deadLetterRoutingKey,
		deadLetterExchangeName,
		false,
		nil,
	)
}

func (client BaseMqClient) declareQueueRouting() error {
	log.Infof("Declaring %s rabbitmq exchange %s", client.ExchangeType, client.ExchangeName)
	err := client.Channel.ExchangeDeclare(
		client.ExchangeName,
		client.ExchangeType,
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return err
	}

	queueName := client.getQueueName()
	log.Infof("Declaring rabbitmq queue %s", queueName)
	dlqOptions := amqp.Table{
		"x-dead-letter-exchange":    deadLetterExchangeName,
		"x-dead-letter-routing-key": queueName,
	}

	_, err = client.Channel.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		dlqOptions,
	)

	if err != nil {
		return err
	}

	log.Infof("Binding rabbitmq queue %s to exchange %s on routing key %s", queueName, client.ExchangeName, client.RoutingKey)
	return client.Channel.QueueBind(
		queueName,
		client.RoutingKey,
		client.ExchangeName,
		false,
		nil,
	)
}

func (client BaseMqClient) declareRoutingTopology() error {
	err0 := client.declareQueueRouting()
	err1 := client.declareDeadLetterRouting()

	for _, e := range []error{err0, err1} {
		if e != nil {
			return e
		}
	}

	return nil
}
