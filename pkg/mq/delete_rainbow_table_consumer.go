package mq

import (
	"github.com/norwoodj/hashbash-backend-go/pkg/rabbit"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

type DeleteRainbowTableWorker struct{}

func (worker *DeleteRainbowTableWorker) HandleMessage(message *amqp.Delivery) error {
	log.Infof("DeleteRainbowTableConsumer got message: %+v", message)
	return nil
}

func NewDeleteRainbowTableConsumer(connection *rabbit.ServerConnection) (*rabbit.Consumer, error) {
	consumerWorker := &DeleteRainbowTableWorker{}
	return rabbit.NewConsumer(
		connection,
		consumerWorker,
		taskExchangeName,
		amqp.ExchangeTopic,
		deleteRainbowTableRoutingKey,
		true,
	)
}
