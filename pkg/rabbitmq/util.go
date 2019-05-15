package rabbitmq

import (
	"github.com/norwoodj/rabbitmq-client-go/rabbitmq"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const taskExchangeName = "task"
const generateRainbowTableRoutingKey = "generateRainbowTable"
const deleteRainbowTableRoutingKey = "deleteRainbowTable"
const searchRainbowTableRoutingKey = "searchRainbowTable"

func AddRabbitMqFlags(flags *pflag.FlagSet) {
	flags.StringP("rabbitmq-host", "r", "localhost", "The hostname or IP address of the hashbash rabbitmq server")
	flags.StringP("rabbitmq-username", "n", "guest", "The username with which to authenticate to the rabbitmq server")
	flags.StringP("rabbitmq-password", "q", "guest", "The password with which to authenticate to the rabbitmq server")
	flags.IntP("rabbitmq-port", "m", 5672, "The port on which to connect to the hashbash rabbitmq server")
}

func AcquireMqConnectionOrDie() *rabbitmq.ServerConnection {
	rabbitConfig := rabbitmq.NewConfig(
		viper.GetString("rabbitmq-hostname"),
		viper.GetString("rabbitmq-username"),
		viper.GetString("rabbitmq-password"),
	)

	connection, err := rabbitmq.NewServerConnection(rabbitConfig)

	if err != nil {
		log.Errorf("Failed to create rabbitmq connection: %s", err)
		os.Exit(1)
	}

	return connection
}
