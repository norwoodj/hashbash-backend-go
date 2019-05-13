package mq

import (
	"github.com/streadway/amqp"
)

type SearchRainbowTableConsumer struct {
	BaseMqConsumerWorker
}

type HashbashMqConsumerWorkers struct {
	HashbashDeleteRainbowTableConsumer   *BaseMqConsumerWorker
	HashbashGenerateRainbowTableConsumer *BaseMqConsumerWorker
	HashbashSearchRainbowTableConsumer   *BaseMqConsumerWorker
}

func newBaseMqConsumer(
	worker ConsumerWorker,
	connection *amqp.Connection,
	exchangeName string,
	exchangeType string,
	routingKey string,
) (*BaseMqConsumerWorker, error) {
	baseConsumer := BaseMqConsumerWorker{
		BaseMqClient: BaseMqClient{
			exchangeName: exchangeName,
			exchangeType: exchangeType,
			routingKey:   routingKey,
		},
		worker: worker,
	}

	err := baseConsumer.createConsumer(connection)
	if err != nil {
		return nil, err
	}

	err = baseConsumer.declareRoutingTopology()
	return &baseConsumer, err
}

func CreateConsumerWorkers(connection *amqp.Connection) (HashbashMqConsumerWorkers, error) {
	deleteRainbowTableConsumer, err0 := NewDeleteRainbowTableConsumer(connection)
	generateRainbowTableConsumer, err1 := NewGenerateRainbowTableConsumer(connection)
	searchRainbowTableConsumer, err2 := NewSearchRainbowTableConsumer(connection)

	for _, e := range []error{err0, err1, err2} {
		if e != nil {
			return HashbashMqConsumerWorkers{}, e
		}
	}

	return HashbashMqConsumerWorkers{
		HashbashDeleteRainbowTableConsumer:   deleteRainbowTableConsumer,
		HashbashGenerateRainbowTableConsumer: generateRainbowTableConsumer,
		HashbashSearchRainbowTableConsumer:   searchRainbowTableConsumer,
	}, nil
}
