package mq

import (
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

type DeleteRainbowTableConsumer struct {
	BaseMqConsumer
}

func NewDeleteRainbowTableConsumer(connection *amqp.Connection) (MqConsumerWorker, error) {
	baseConsumer, err := newBaseMqConsumer(connection, taskExchangeName, "topic", deleteRainbowTableRoutingKey)
	return &DeleteRainbowTableConsumer{baseConsumer}, err
}

func (consumer *DeleteRainbowTableConsumer) handleMessage(message *amqp.Delivery) error {
	log.Infof("DeleteRainbowTableConsumer got message: %+v", message)
	return nil
}
