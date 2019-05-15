package rabbitmq

import (
	"encoding/json"
	"fmt"
	"github.com/norwoodj/hashbash-backend-go/pkg/rainbow"
	"github.com/norwoodj/hashbash-backend-go/pkg/service"
	"github.com/norwoodj/rabbitmq-client-go/rabbitmq"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

type GenerateRainbowTableWorker struct {
	rainbowTableService            service.RainbowTableService
	rainbowTableGenerateJobService *rainbow.TableGeneratorJobService
}

func (worker *GenerateRainbowTableWorker) HandleMessage(message *amqp.Delivery) error {
	var rainbowTableGenerateMessage RainbowTableIdMessage
	err := json.Unmarshal(message.Body, &rainbowTableGenerateMessage)
	if err != nil {
		return fmt.Errorf("failed to unmarshal rainbow table generate message")
	}

	log.Infof("GenerateRainbowTable consumer got generate request for rainbow table %d", rainbowTableGenerateMessage.RainbowTableId)
	rainbowTable := worker.rainbowTableService.FindRainbowTableById(rainbowTableGenerateMessage.RainbowTableId)

	if rainbowTable.Name == "" {
		return fmt.Errorf("rainbow table with ID %d not found, cannot generate", rainbowTableGenerateMessage.RainbowTableId)
	}

	return worker.rainbowTableGenerateJobService.RunGenerateJobForTable(rainbowTable)
}

func NewGenerateRainbowTableConsumer(
	connection *rabbitmq.ServerConnection,
	rainbowTableService service.RainbowTableService,
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
