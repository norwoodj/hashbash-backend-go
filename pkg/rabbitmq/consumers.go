package rabbitmq

import (
	"github.com/norwoodj/rabbitmq-client-go/rabbitmq"
)

type HashbashMqConsumerWorkers struct {
	HashbashDeleteRainbowTableConsumer   *rabbitmq.Consumer
	HashbashGenerateRainbowTableConsumer *rabbitmq.Consumer
	HashbashSearchRainbowTableConsumer   *rabbitmq.Consumer
}

func CreateConsumerWorkers(connection *rabbitmq.ServerConnection) (HashbashMqConsumerWorkers, error) {
	deleteRainbowTableConsumer, err0 := NewDeleteRainbowTableConsumer(connection)
	generateRainbowTableConsumer, err1 := NewGenerateRainbowTableConsumer(connection)
	searchRainbowTableConsumer, err2 := NewSearchRainbowTableConsumer(connection)

	for _, e := range []error{err0, err1, err2} {
		if e != nil {
			return HashbashMqConsumerWorkers{}, e
		}
	}

	return HashbashMqConsumerWorkers{
		HashbashDeleteRainbowTableConsumer:   deleteRainbowTableConsumer,
		HashbashGenerateRainbowTableConsumer: generateRainbowTableConsumer,
		HashbashSearchRainbowTableConsumer:   searchRainbowTableConsumer,
	}, nil
}
