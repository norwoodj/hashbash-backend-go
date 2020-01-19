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
	flags.String("rabbitmq-host", "localhost", "The hostname or IP address of the hashbash rabbitmq metrics")
	flags.String("rabbitmq-username", "guest", "The username with which to authenticate to the rabbitmq metrics")
	flags.String("rabbitmq-password", "guest", "The password with which to authenticate to the rabbitmq metrics")
	flags.Int("rabbitmq-port", 5672, "The port on which to connect to the hashbash rabbitmq metrics")
	flags.String("rabbitmq-vhost", "", "The virtual host to connect to on the rabbitmq server")
}

func AcquireMqConnectionOrDie() *rabbitmq.ServerConnection {
	rabbitConfig := rabbitmq.NewConfig(
		viper.GetString("rabbitmq-host"),
		viper.GetInt("rabbitmq-port"),
		viper.GetString("rabbitmq-username"),
		viper.GetString("rabbitmq-password"),
		viper.GetString("rabbitmq-vhost"),
	)

	connection, err := rabbitmq.NewServerConnection(rabbitConfig)

	if err != nil {
		log.Errorf("Failed to create rabbitmq connection: %s", err)
		os.Exit(1)
	}

	return connection
}
