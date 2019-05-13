package rabbitmq

import (
	"github.com/norwoodj/rabbitmq-client-go/rabbitmq"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

type GenerateRainbowTableWorker struct{}

func (worker *GenerateRainbowTableWorker) HandleMessage(message *amqp.Delivery) error {
	log.Infof("GenerateRainbowTableConsumer got message: %+v", message)
	return nil
}

func NewGenerateRainbowTableConsumer(connection *rabbitmq.ServerConnection) (*rabbitmq.Consumer, error) {
	consumerWorker := &GenerateRainbowTableWorker{}
	return rabbitmq.NewConsumer(
		connection,
		consumerWorker,
		taskExchangeName,
		amqp.ExchangeTopic,
		generateRainbowTableRoutingKey,
		true,
	)
}
