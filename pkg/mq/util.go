package mq

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/streadway/amqp"
)

const taskExchangeName = "task"
const generateRainbowTableQueueName = "generateRainbowTable"
const deleteRainbowTableQueueName = "deleteRainbowTable"
const searchRainbowTableQueueName = "searchRainbowTable"

func formatRabbitMqDsn(
	hostname string,
	username string,
	password string,
	port int,
) string {
	return fmt.Sprintf(
		"amqp://%s:%s@%s:%d/",
		username,
		password,
		hostname,
		port,
	)
}

func AddRabbitMqFlags(flags *pflag.FlagSet) {
	flags.StringP("rabbitmq-host", "r", "localhost", "The hostname or IP address of the hashbash rabbitmq server")
	flags.StringP("rabbitmq-username", "n", "guest", "The username with which to authenticate to the rabbitmq server")
	flags.StringP("rabbitmq-password", "q", "guest", "The password with which to authenticate to the rabbitmq server")
	flags.IntP("rabbitmq-port", "m", 5672, "The port on which to connect to the hashbash rabbitmq server")
}

func AcquireMqConnection() (*amqp.Connection, error) {
	censoredRabbitDsn := formatRabbitMqDsn(
		viper.GetString("rabbitmq-host"),
		viper.GetString("rabbitmq-username"),
		"********",
		viper.GetInt("rabbitmq-port"),
	)

	log.Infof("Connecting to RabbitMQ server %s", censoredRabbitDsn)

	rabbitDsn := formatRabbitMqDsn(
		viper.GetString("rabbitmq-host"),
		viper.GetString("rabbitmq-username"),
		viper.GetString("rabbitmq-password"),
		viper.GetInt("rabbitmq-port"),
	)

	conn, err := amqp.Dial(rabbitDsn)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
