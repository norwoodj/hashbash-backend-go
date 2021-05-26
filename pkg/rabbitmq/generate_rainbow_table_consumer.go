package rabbitmq

import (
	"encoding/json"
	"fmt"
	"github.com/norwoodj/hashbash-backend-go/pkg/dao"
	"github.com/norwoodj/hashbash-backend-go/pkg/rainbow"
	"github.com/norwoodj/rabbitmq-client-go/rabbitmq"
	"github.com/rs/zerolog/log"
	"github.com/streadway/amqp"
)

type GenerateRainbowTableWorker struct {
	rainbowTableService            dao.RainbowTableService
	rainbowTableGenerateJobService *rainbow.TableGeneratorJobService
}

func (worker *GenerateRainbowTableWorker) HandleMessage(message *amqp.Delivery) error {
	var rainbowTableGenerateMessage RainbowTableIdMessage
	err := json.Unmarshal(message.Body, &rainbowTableGenerateMessage)
	if err != nil {
		return fmt.Errorf("failed to unmarshal rainbow table generate message")
	}

	log.Info().Msgf("GenerateRainbowTable consumer got generate request for rainbow table %d", rainbowTableGenerateMessage.RainbowTableId)
	return worker.rainbowTableGenerateJobService.RunGenerateJobForTable(rainbowTableGenerateMessage.RainbowTableId)
}

func NewGenerateRainbowTableConsumer(
	connection *rabbitmq.ServerConnection,
	rainbowTableService dao.RainbowTableService,
	rainbowTableGenerateJobService *rainbow.TableGeneratorJobService,
) (*rabbitmq.Consumer, error) {
	consumerWorker := &GenerateRainbowTableWorker{
		rainbowTableService:            rainbowTableService,
		rainbowTableGenerateJobService: rainbowTableGenerateJobService,
	}

	return rabbitmq.NewConsumer(
		connection,
		consumerWorker,
		taskExchangeName,
		amqp.ExchangeTopic,
		generateRainbowTableRoutingKey,
		true,
	)
}
