package mq

import (
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

type GenerateRainbowTableConsumer struct {
	BaseMqConsumer
}

func NewGenerateRainbowTableConsumer(connection *amqp.Connection) (MqConsumerWorker, error) {
	baseConsumer, err := newBaseMqConsumer(connection, taskExchangeName, "topic", generateRainbowTableRoutingKey)
	return &GenerateRainbowTableConsumer{baseConsumer}, err
}

func (consumer *GenerateRainbowTableConsumer) handleMessage(message *amqp.Delivery) error {
	log.Infof("GenerateRainbowTableConsumer got message: %+v", message)
	return nil
}
