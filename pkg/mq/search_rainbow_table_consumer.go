package mq

import (
	"github.com/norwoodj/hashbash-backend-go/pkg/rabbit"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

type SearchRainbowTableWorker struct{}

func (worker *SearchRainbowTableWorker) HandleMessage(message *amqp.Delivery) error {
	log.Infof("SearchRainbowTableConsumer got message: %+v", message)
	return nil
}

func NewSearchRainbowTableConsumer(connection *rabbit.ServerConnection) (*rabbit.Consumer, error) {
	consumerWorker := &SearchRainbowTableWorker{}
	return rabbit.NewConsumer(
		connection,
		consumerWorker,
		taskExchangeName,
		amqp.ExchangeTopic,
		searchRainbowTableRoutingKey,
		true,
	)
}
