package rainbow

import (
	"github.com/norwoodj/hashbash-backend-go/pkg/model"
)

func initializeErrorChannels(numThreads int) []chan error {
	errorChannels := make([]chan error, numThreads)

	for i := range errorChannels {
		errorChannels[i] = make(chan error)
	}

	return errorChannels
}

func initializeFoundChannels(numBatches int) []chan model.RainbowChain {
	foundChannels := make([]chan model.RainbowChain, numBatches)

	for i := range foundChannels {
		foundChannels[i] = make(chan model.RainbowChain)
	}

	return foundChannels
}
