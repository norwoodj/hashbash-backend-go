package rabbit

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

type ConsumerWorker interface {
	HandleMessage(*amqp.Delivery) error
}

type Consumer struct {
	BaseClient
	worker ConsumerWorker
}

func NewConsumer(
	connection *ServerConnection,
	worker ConsumerWorker,
	exchangeName string,
	exchangeType string,
	routingKey string,
	declareRoutingTopology bool,
) (*Consumer, error) {
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

	return &Consumer{
		BaseClient: baseClient,
		worker:     worker,
	}, nil
}

func (consumer *Consumer) ConsumeMessages(quit chan bool, startErrorChannel chan error) {
	queueName := consumer.getQueueName()
	err := consumer.channel.Qos(1, 0, false)

	if err != nil {
		startErrorChannel <- fmt.Errorf("error setting consumer QOS settings for %s queue: %s", queueName, err)
	}

	msgPipe, err := consumer.channel.Consume(
		queueName,
		"",
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		startErrorChannel <- fmt.Errorf("error starting consumer for %s queue: %s", queueName, err)
	}

	close(startErrorChannel)

	go func() {
		for msg := range msgPipe {
			err := consumer.worker.HandleMessage(&msg)
			if err != nil {
				log.Warnf("Failed to handle message: %s", err)
				err = msg.Nack(false, false)
				if err != nil {
					log.Warnf("Failed to nack message, dropping: %s", err)
				}

				continue
			}

			err = msg.Ack(false)
			if err != nil {
				log.Warnf("Failed to ack message: %s", err)
				continue
			}
		}
	}()

	log.Infof("Started consumer on queue %s...", queueName)
	<-quit
	log.Infof("Quit signal received, stopping %s consumer", queueName)

	err = consumer.channel.Close()
	if err != nil {
		log.Errorf("Error closing channel for %s consumer: %s", queueName, err)
	}
}
