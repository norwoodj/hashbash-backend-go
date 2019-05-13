package rabbitmq

import (
	"github.com/norwoodj/rabbitmq-client-go/rabbitmq"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

type DeleteRainbowTableWorker struct{}

func (worker *DeleteRainbowTableWorker) HandleMessage(message *amqp.Delivery) error {
	log.Infof("DeleteRainbowTableConsumer got message: %+v", message)
	return nil
}

func NewDeleteRainbowTableConsumer(connection *rabbitmq.ServerConnection) (*rabbitmq.Consumer, error) {
	consumerWorker := &DeleteRainbowTableWorker{}
	return rabbitmq.NewConsumer(
		connection,
		consumerWorker,
		taskExchangeName,
		amqp.ExchangeTopic,
		deleteRainbowTableRoutingKey,
		true,
	)
}
