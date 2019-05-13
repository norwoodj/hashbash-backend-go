package rabbit

import (
	"fmt"

	"github.com/streadway/amqp"
)

type Producer struct {
	BaseClient
	messageSerializer MessageSerializer
}

func NewProducer(
	connection *ServerConnection,
	serializer MessageSerializer,
	exchangeName string,
	exchangeType string,
	routingKey string,
	declareRoutingTopology bool,
) (*Producer, error) {
	baseClient, err := newBaseRabbitmqClient(
		connection,
		exchangeName,
		exchangeType,
		routingKey,
		declareRoutingTopology,
	)

	if err != nil {
		return nil, err
	}

	return &Producer{
		BaseClient:        baseClient,
		messageSerializer: serializer,
	}, nil
}

func (producer *Producer) PublishMessage(msg interface{}) error {
	serializedMsg, err := producer.messageSerializer.SerializeMessage(msg)
	if err != nil {
		return fmt.Errorf("failed to serialize message for publishing: %s", err)
	}

	return producer.channel.Publish(
		producer.exchangeName,
		producer.routingKey,
		false,
		false,
		amqp.Publishing{ContentType: producer.messageSerializer.GetContentType(), Body: serializedMsg},
	)
}
