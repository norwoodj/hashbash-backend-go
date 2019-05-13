package mq

import (
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
)

type BaseMqProducer struct {
	BaseMqClient
}

func NewMqProducer(
	rabbitmqConnection *amqp.Connection,
	exchangeName string,
	exchangeType string,
	routingKey string,
) (BaseMqProducer, error) {
	channel, err := rabbitmqConnection.Channel()
	if err != nil {
		return BaseMqProducer{}, err
	}

	producer := BaseMqProducer{
		BaseMqClient{
			channel:      channel,
			exchangeName: exchangeName,
			exchangeType: exchangeType,
			routingKey:   routingKey,
		},
	}

	err = producer.declareRoutingTopology()
	if err != nil {
		return BaseMqProducer{}, err
	}

	return producer, nil
}

func (producer BaseMqProducer) PublishMessage(msg RainbowTableMessage) error {
	serializedMsg, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to serialize message for publishing: %s", err)
	}

	return producer.channel.Publish(
		producer.exchangeName,
		producer.routingKey,
		false,
		false,
		amqp.Publishing{ContentType: "application/json", Body: []byte(serializedMsg)},
	)
}
