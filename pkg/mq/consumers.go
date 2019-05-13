package mq

import (
	"github.com/streadway/amqp"
)

type SearchRainbowTableConsumer struct {
	BaseMqConsumer
}

type HashbashMqConsumers struct {
	HashbashDeleteRainbowTableConsumer   MqConsumerWorker
	HashbashGenerateRainbowTableConsumer MqConsumerWorker
	HashbashSearchRainbowTableConsumer   MqConsumerWorker
}

func newBaseMqConsumer(
	connection *amqp.Connection,
	exchangeName string,
	exchangeType string,
	routingKey string,
) (BaseMqConsumer, error) {
	baseConsumer := BaseMqConsumer{
		BaseMqClient{
			exchangeName: exchangeName,
			exchangeType: exchangeType,
			routingKey:   routingKey,
		},
	}

	err := baseConsumer.createConsumer(connection)
	if err != nil {
		return BaseMqConsumer{}, err
	}

	err = baseConsumer.declareRoutingTopology()
	return baseConsumer, err
}

func NewSearchRainbowTableConsumer(connection *amqp.Connection) (MqConsumerWorker, error) {
	baseConsumer, err := newBaseMqConsumer(connection, taskExchangeName, "topic", searchRainbowTableRoutingKey)
	return &SearchRainbowTableConsumer{baseConsumer}, err
}

func CreateConsumers(connection *amqp.Connection) (HashbashMqConsumers, error) {
	deleteRainbowTableConsumer, err0 := NewDeleteRainbowTableConsumer(connection)
	generateRainbowTableConsumer, err1 := NewGenerateRainbowTableConsumer(connection)
	searchRainbowTableConsumer, err2 := NewSearchRainbowTableConsumer(connection)

	for _, e := range []error{err0, err1, err2} {
		if e != nil {
			return HashbashMqConsumers{}, e
		}
	}

	return HashbashMqConsumers{
		HashbashDeleteRainbowTableConsumer:   deleteRainbowTableConsumer,
		HashbashGenerateRainbowTableConsumer: generateRainbowTableConsumer,
		HashbashSearchRainbowTableConsumer:   searchRainbowTableConsumer,
	}, nil
}
