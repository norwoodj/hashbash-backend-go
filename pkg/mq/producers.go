package mq

import "github.com/streadway/amqp"

type HashbashMqProducers struct {
	DeleteRainbowTableProducer   *BaseMqProducer
	GenerateRainbowTableProducer *BaseMqProducer
	SearchRainbowTableProducer   *BaseMqProducer
}

func CreateProducers(connection *amqp.Connection) (HashbashMqProducers, error) {
	deleteRainbowTableProducer, err1 := NewMqProducer(connection, taskExchangeName, "topic", deleteRainbowTableRoutingKey)
	generateRainbowTableProducer, err0 := NewMqProducer(connection, taskExchangeName, "topic", generateRainbowTableRoutingKey)
	searchRainbowTableProducer, err2 := NewMqProducer(connection, taskExchangeName, "topic", searchRainbowTableRoutingKey)

	for _, e := range []error{err0, err1, err2} {
		if e != nil {
			return HashbashMqProducers{}, e
		}
	}
	return HashbashMqProducers{
		DeleteRainbowTableProducer:   &deleteRainbowTableProducer,
		GenerateRainbowTableProducer: &generateRainbowTableProducer,
		SearchRainbowTableProducer:   &searchRainbowTableProducer,
	}, nil
}
