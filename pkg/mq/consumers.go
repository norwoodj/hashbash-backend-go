package mq

import (
	"github.com/norwoodj/hashbash-backend-go/pkg/rabbit"
)

type HashbashMqConsumerWorkers struct {
	HashbashDeleteRainbowTableConsumer   *rabbit.Consumer
	HashbashGenerateRainbowTableConsumer *rabbit.Consumer
	HashbashSearchRainbowTableConsumer   *rabbit.Consumer
}

func CreateConsumerWorkers(connection *rabbit.ServerConnection) (HashbashMqConsumerWorkers, error) {
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
