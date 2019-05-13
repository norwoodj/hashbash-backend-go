package mq

import (
	"github.com/norwoodj/hashbash-backend-go/pkg/rabbit"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

type GenerateRainbowTableWorker struct{}

func (worker *GenerateRainbowTableWorker) HandleMessage(message *amqp.Delivery) error {
	log.Infof("GenerateRainbowTableConsumer got message: %+v", message)
	return nil
}

func NewGenerateRainbowTableConsumer(connection *rabbit.ServerConnection) (*rabbit.Consumer, error) {
	consumerWorker := &GenerateRainbowTableWorker{}
	return rabbit.NewConsumer(
		connection,
		consumerWorker,
		taskExchangeName,
		amqp.ExchangeTopic,
		generateRainbowTableRoutingKey,
		true,
	)
}
