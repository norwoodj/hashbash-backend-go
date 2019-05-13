package mq

import "github.com/streadway/amqp"

type HashbashMqConsumers struct {
	GenerateRainbowTableConsumer *BaseMqConsumer
	DeleteRainbowTableConsumer   *BaseMqConsumer
	SearchRainbowTableConsumer   *BaseMqConsumer
}

func CreateConsumers(connection *amqp.Connection) (HashbashMqConsumers, error) {
	generateRainbowTableConsumer, err0 := NewMqConsumer(connection, taskExchangeName, "topic", generateRainbowTableQueueName)
	deleteRainbowTableConsumer, err1 := NewMqConsumer(connection, taskExchangeName, "topic", deleteRainbowTableQueueName)
	searchRainbowTableConsumer, err2 := NewMqConsumer(connection, taskExchangeName, "topic", searchRainbowTableQueueName)

	for _, e := range []error{err0, err1, err2} {
		if e != nil {
			return HashbashMqConsumers{}, e
		}
	}

	return HashbashMqConsumers{
		GenerateRainbowTableConsumer: &generateRainbowTableConsumer,
		DeleteRainbowTableConsumer:   &deleteRainbowTableConsumer,
		SearchRainbowTableConsumer:   &searchRainbowTableConsumer,
	}, nil
}
