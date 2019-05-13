package mq

import (
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

type GenerateRainbowTableWorker struct{}

func (worker *GenerateRainbowTableWorker) handleMessage(message *amqp.Delivery) error {
	log.Infof("GenerateRainbowTableConsumer got message: %+v", message)
	return nil
}

func NewGenerateRainbowTableConsumer(connection *amqp.Connection) (Consumer, error) {
	consumerWorker := &GenerateRainbowTableWorker{}
	return newBaseMqConsumer(consumerWorker, connection, taskExchangeName, "topic", generateRainbowTableRoutingKey)
}
