package mq

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

type MqConsumerWorker interface {
	ConsumeMessages(chan bool, chan error)

	createConsumer(*amqp.Connection) error
	declareRoutingTopology() error
	handleMessage(*amqp.Delivery) error
}

type BaseMqConsumer struct {
	BaseMqClient
}

func (consumer *BaseMqConsumer) createConsumer(rabbitmqConnection *amqp.Connection) error {
	var err error
	consumer.channel, err = rabbitmqConnection.Channel()
	return err
}

func (consumer *BaseMqConsumer) ConsumeMessages(quit chan bool, startErrorChannel chan error) {
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
			err := consumer.handleMessage(&msg)
			if err != nil {
				log.Warnf("Failed to handle message: %s", err)
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

func (consumer *BaseMqConsumer) handleMessage(message *amqp.Delivery) error {
	log.Infof("BaseConsumer got message: %+v", message)
	return nil
}
