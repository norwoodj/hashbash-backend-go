package mq

import (
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

type SearchRainbowTableWorker struct{}

func (worker *SearchRainbowTableWorker) handleMessage(message *amqp.Delivery) error {
	log.Infof("SearchRainbowTableConsumer got message: %+v", message)
	return nil
}

func NewSearchRainbowTableConsumer(connection *amqp.Connection) (*BaseMqConsumerWorker, error) {
	consumerWorker := &SearchRainbowTableWorker{}
	return newBaseMqConsumer(consumerWorker, connection, taskExchangeName, "topic", searchRainbowTableRoutingKey)
}
