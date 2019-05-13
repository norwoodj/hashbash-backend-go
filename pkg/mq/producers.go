package mq

import "github.com/streadway/amqp"

type HashbashMqProducers struct {
	CreateRainbowTableProducer *BaseMqProducer
	DeleteRainbowTableProducer *BaseMqProducer
	SearchRainbowTableProducer *BaseMqProducer
}

func CreateProducers(connection *amqp.Connection) (HashbashMqProducers, error) {
	createRainbowTableProducer, err0 := NewMqProducer(connection, taskExchangeName, "topic", generateRainbowTableQueueName)
	deleteRainbowTableProducer, err1 := NewMqProducer(connection, taskExchangeName, "topic", deleteRainbowTableQueueName)
	searchRainbowTableProducer, err2 := NewMqProducer(connection, taskExchangeName, "topic", searchRainbowTableQueueName)

	for _, e := range []error{err0, err1, err2} {
		if e != nil {
			return HashbashMqProducers{}, e
		}
	}
	return HashbashMqProducers{
		CreateRainbowTableProducer: &createRainbowTableProducer,
		DeleteRainbowTableProducer: &deleteRainbowTableProducer,
		SearchRainbowTableProducer: &searchRainbowTableProducer,
	}, nil
}
