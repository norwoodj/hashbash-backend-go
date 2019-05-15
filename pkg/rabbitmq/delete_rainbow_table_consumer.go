package rabbitmq

import (
	"encoding/json"
	"github.com/norwoodj/hashbash-backend-go/pkg/service"
	"github.com/norwoodj/rabbitmq-client-go/rabbitmq"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

type DeleteRainbowTableWorker struct {
	rainbowTableService service.RainbowTableService
}

func (worker *DeleteRainbowTableWorker) HandleMessage(message *amqp.Delivery) error {
	var messageContent RainbowTableIdMessage
	err := json.Unmarshal(message.Body, &messageContent)
	if err != nil {
		log.Errorf("Failed to deserialize rainbow table ID message for deletion: %s", err)
		return err
	}

	log.Infof("DeleteRainbowTable consumer got request for deletion of rainbow table %d", messageContent.RainbowTableId)
	err = worker.rainbowTableService.DeleteRainbowTableById(messageContent.RainbowTableId)

	if err != nil {
		if !service.IsRainbowTableNotExistsError(err) {
			log.Errorf("Unknown error occurred deleting rainbow table with id %d: %s", messageContent.RainbowTableId, err)
		}

		return err
	}

	return nil
}

func NewDeleteRainbowTableConsumer(
	connection *rabbitmq.ServerConnection,
	rainbowTableService service.RainbowTableService,
) (*rabbitmq.Consumer, error) {
	consumerWorker := &DeleteRainbowTableWorker{rainbowTableService: rainbowTableService}
	return rabbitmq.NewConsumer(
		connection,
		consumerWorker,
		taskExchangeName,
		amqp.ExchangeTopic,
		deleteRainbowTableRoutingKey,
		true,
	)
}
