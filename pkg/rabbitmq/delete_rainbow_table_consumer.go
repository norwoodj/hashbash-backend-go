package rabbitmq

import (
	"encoding/json"
	"github.com/norwoodj/hashbash-backend-go/pkg/dao"
	"github.com/norwoodj/rabbitmq-client-go/rabbitmq"
	"github.com/rs/zerolog/log"
	amqp "github.com/rabbitmq/amqp091-go"
)

type DeleteRainbowTableWorker struct {
	rainbowTableService dao.RainbowTableService
}

func (worker *DeleteRainbowTableWorker) HandleMessage(message *amqp.Delivery) error {
	var messageContent RainbowTableIdMessage
	err := json.Unmarshal(message.Body, &messageContent)
	if err != nil {
		log.Error().Err(err).Msg("Failed to deserialize rainbow table ID message for deletion")
		return err
	}

	log.Info().Msgf("DeleteRainbowTable consumer got request for deletion of rainbow table %d", messageContent.RainbowTableId)
	err = worker.rainbowTableService.DeleteRainbowTableById(messageContent.RainbowTableId)

	if err != nil {
		if !dao.IsRainbowTableNotExistsError(err) {
			log.Error().Err(err).Msgf("Unknown error occurred deleting rainbow table with id %d", messageContent.RainbowTableId)
		}

		return err
	}

	return nil
}

func NewDeleteRainbowTableConsumer(
	connection *rabbitmq.ServerConnection,
	rainbowTableService dao.RainbowTableService,
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
