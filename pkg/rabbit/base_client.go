package rabbit

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

type BaseClient struct {
	channel      *amqp.Channel
	config       *Config
	exchangeName string
	exchangeType string
	routingKey   string
}

func newBaseRabbitmqClient(
	connection *ServerConnection,
	exchangeName string,
	exchangeType string,
	routingKey string,
	declareRoutingTopology bool,
) (BaseClient, error) {
	channel, err := connection.connection.Channel()

	if err != nil {
		return BaseClient{}, err
	}

	baseClient := BaseClient{
		channel: channel,
		config: connection.config,
		exchangeName: exchangeName,
		exchangeType: exchangeType,
		routingKey: routingKey,
	}

	if declareRoutingTopology {
		err := baseClient.declareRoutingTopology()

		if err != nil {
			return baseClient, fmt.Errorf(
				"failed to declare routing topology for exchange/routing-key %s/%s: %s",
				exchangeName,
				routingKey,
				err,
			)
		}
	}

	return baseClient, nil
}

func (client BaseClient) declareRoutingTopology() error {
	err0 := client.declareQueueRouting()
	err1 := client.declareDeadLetterRouting()

	for _, e := range []error{err0, err1} {
		if e != nil {
			return e
		}
	}

	return nil
}

func (client BaseClient) getQueueName() string {
	return client.config.QueueNamingStrategy.GetQueueName(client.exchangeName, client.routingKey)
}

func (client BaseClient) getDeadLetterQueueName() string {
	return fmt.Sprintf("%s.%s.%s", client.exchangeName, client.routingKey, client.config.DeadLetterQueueSuffix)
}

func (client BaseClient) getDeadLetterRoutingKey() string {
	return fmt.Sprintf("%s.%s", client.exchangeName, client.routingKey)
}

func (client BaseClient) declareDeadLetterRouting() error {
	log.Infof("Declaring rabbitmq dead letter exchange %s", client.config.DeadLetterExchangeName)
	err := client.channel.ExchangeDeclare(
		client.config.DeadLetterExchangeName,
		amqp.ExchangeDirect,
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
	_, err = client.channel.QueueDeclare(
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
	log.Infof("Binding rabbitmq dead letter queue %s to exchange %s on routing key %s", deadLetterQueueName, client.config.DeadLetterExchangeName, deadLetterRoutingKey)
	return client.channel.QueueBind(
		deadLetterQueueName,
		deadLetterRoutingKey,
		client.config.DeadLetterExchangeName,
		false,
		nil,
	)
}

func (client BaseClient) declareQueueRouting() error {
	log.Infof("Declaring %s rabbitmq exchange %s", client.exchangeType, client.exchangeName)
	err := client.channel.ExchangeDeclare(
		client.exchangeName,
		client.exchangeType,
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return err
	}

	deadLetterQueueName := client.getQueueName()
	log.Infof("Declaring rabbitmq queue %s", deadLetterQueueName)
	dlqOptions := amqp.Table{
		"x-dead-letter-exchange":    client.config.DeadLetterExchangeName,
		"x-dead-letter-routing-key": deadLetterQueueName,
	}

	_, err = client.channel.QueueDeclare(
		deadLetterQueueName,
		true,
		false,
		false,
		false,
		dlqOptions,
	)

	if err != nil {
		return err
	}

	log.Infof("Binding rabbitmq queue %s to exchange %s on routing key %s", deadLetterQueueName, client.exchangeName, client.routingKey)
	return client.channel.QueueBind(
		deadLetterQueueName,
		client.routingKey,
		client.exchangeName,
		false,
		nil,
	)
}
