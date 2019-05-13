package mq

import (
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

type DeleteRainbowTableWorker struct{}

func (worker *DeleteRainbowTableWorker) handleMessage(message *amqp.Delivery) error {
	log.Infof("DeleteRainbowTableConsumer got message: %+v", message)
	return nil
}

func NewDeleteRainbowTableConsumer(connection *amqp.Connection) (*BaseMqConsumerWorker, error) {
	consumerWorker := &DeleteRainbowTableWorker{}
	return newBaseMqConsumer(consumerWorker, connection, taskExchangeName, "topic", deleteRainbowTableRoutingKey)
}
