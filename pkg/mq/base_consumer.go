package mq

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"os"
)

type BaseMqConsumer struct {
	BaseMqClient
}

func NewMqConsumer(
	rabbitmqConnection *amqp.Connection,
	exchangeName string,
	exchangeType string,
	routingKey string,
) (BaseMqConsumer, error) {
	channel, err := rabbitmqConnection.Channel()
	if err != nil {
		return BaseMqConsumer{}, err
	}

	consumer := BaseMqConsumer{
		BaseMqClient{
			Channel:      channel,
			ExchangeName: exchangeName,
			ExchangeType: exchangeType,
			RoutingKey:   routingKey,
		},
	}

	err = consumer.declareRoutingTopology()
	if err != nil {
		return BaseMqConsumer{}, err
	}

	return consumer, nil
}

func (consumer BaseMqConsumer) ConsumeMessages(quit chan bool) {
	queueName := consumer.getQueueName()
	err := consumer.Channel.Qos(1, 0, true)

	if err != nil {
		log.Errorf("Error setting consumer QOS settings for %s queue: %s", queueName, err)
		os.Exit(1)
	}

	msgPipe, err := consumer.Channel.Consume(
		queueName,
		"",
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		log.Errorf("Error starting consumer for %s queue: %s", queueName, err)
		os.Exit(1)
	}

	go func() {
		for msg := range msgPipe {
			var rainbowTableMessage RainbowTableMessage
			err := json.Unmarshal(msg.Body, &rainbowTableMessage)

			if err != nil {
				log.Warnf("Failed to unmarshal Rainbow Table message %s")
				continue
			}

			log.Infof("Got RainbowTable message: %+v", msg.Body)
		}
	}()

	log.Infof("Started consumer on queue %s...", queueName)
	<-quit
	log.Infof("Quit signal received, stopping %s consumer", queueName)

	err = consumer.Channel.Close()
	if err != nil {
		log.Errorf("Error closing channel for %s consumer: %s", queueName, err)
	}
}
