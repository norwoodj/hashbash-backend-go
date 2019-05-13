package rabbitmq

import (
	"github.com/norwoodj/rabbitmq-client-go/rabbitmq"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

type SearchRainbowTableWorker struct{}

func (worker *SearchRainbowTableWorker) HandleMessage(message *amqp.Delivery) error {
	log.Infof("SearchRainbowTableConsumer got message: %+v", message)
	return nil
}

func NewSearchRainbowTableConsumer(connection *rabbitmq.ServerConnection) (*rabbitmq.Consumer, error) {
	consumerWorker := &SearchRainbowTableWorker{}
	return rabbitmq.NewConsumer(
		connection,
		consumerWorker,
		taskExchangeName,
		amqp.ExchangeTopic,
		searchRainbowTableRoutingKey,
		true,
	)
}
