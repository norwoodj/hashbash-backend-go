package rabbitmq

import (
	"encoding/json"

	"github.com/norwoodj/hashbash-backend-go/pkg/rainbow"
	"github.com/norwoodj/rabbitmq-client-go/rabbitmq"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

type SearchRainbowTableWorker struct {
	rainbowTableSearchJobService *rainbow.TableSearchJobService
}

func (worker *SearchRainbowTableWorker) HandleMessage(message *amqp.Delivery) error {
	var messageContent RainbowTableSearchRequestMessage
	err := json.Unmarshal(message.Body, &messageContent)

	if err != nil {
		log.Errorf("Failed to deserialize rainbow table search request message: %s", err)
		return err
	}

	log.Infof(
		"SearchRainbowTable consumer got request to search rainbow table %d for hash %s",
		messageContent.RainbowTableId,
		messageContent.Hash,
	)

	return worker.rainbowTableSearchJobService.RunSearchJob(messageContent.SearchId)
}

func NewSearchRainbowTableConsumer(
	connection *rabbitmq.ServerConnection,
	rainbowTableSearchJobService *rainbow.TableSearchJobService,
) (*rabbitmq.Consumer, error) {
	consumerWorker := &SearchRainbowTableWorker{rainbowTableSearchJobService: rainbowTableSearchJobService}
	return rabbitmq.NewConsumer(
		connection,
		consumerWorker,
		taskExchangeName,
		amqp.ExchangeTopic,
		searchRainbowTableRoutingKey,
		true,
	)
}
