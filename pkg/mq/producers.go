package mq

import (
	"github.com/norwoodj/hashbash-backend-go/pkg/rabbit"
	"github.com/streadway/amqp"
)

type HashbashMqProducers struct {
	DeleteRainbowTableProducer   *rabbit.Producer
	GenerateRainbowTableProducer *rabbit.Producer
	SearchRainbowTableProducer   *rabbit.Producer
}

func newProducer(connection *rabbit.ServerConnection, routingKey string) (*rabbit.Producer, error) {
	serializer := rabbit.JsonMessageSerializer{}
	return rabbit.NewProducer(
		connection,
		serializer,
		taskExchangeName,
		amqp.ExchangeTopic,
		routingKey,
		true,
	)

}

func CreateProducers(connection *rabbit.ServerConnection) (HashbashMqProducers, error) {
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
