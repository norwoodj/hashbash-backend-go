package rabbitmq

import (
	"github.com/norwoodj/rabbitmq-client-go/rabbitmq"
	"github.com/streadway/amqp"
)

type HashbashMqProducers struct {
	DeleteRainbowTableProducer   *rabbitmq.Producer
	GenerateRainbowTableProducer *rabbitmq.Producer
	SearchRainbowTableProducer   *rabbitmq.Producer
}

func newProducer(connection *rabbitmq.ServerConnection, routingKey string) (*rabbitmq.Producer, error) {
	serializer := rabbitmq.JsonMessageSerializer{}
	return rabbitmq.NewProducer(
		connection,
		serializer,
		taskExchangeName,
		amqp.ExchangeTopic,
		routingKey,
		true,
	)

}

func CreateProducers(connection *rabbitmq.ServerConnection) (HashbashMqProducers, error) {
	deleteRainbowTableProducer, err0 := newProducer(connection, deleteRainbowTableRoutingKey)
	generateRainbowTableProducer, err1 := newProducer(connection, generateRainbowTableRoutingKey)
	searchRainbowTableProducer, err2 := newProducer(connection, searchRainbowTableRoutingKey)

	for _, e := range []error{err0, err1, err2} {
		if e != nil {
			return HashbashMqProducers{}, e
		}
	}
	return HashbashMqProducers{
		DeleteRainbowTableProducer:   deleteRainbowTableProducer,
		GenerateRainbowTableProducer: generateRainbowTableProducer,
		SearchRainbowTableProducer:   searchRainbowTableProducer,
	}, nil
}
