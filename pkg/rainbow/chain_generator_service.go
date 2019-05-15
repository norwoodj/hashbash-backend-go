package rainbow

import (
	"fmt"
	"github.com/norwoodj/hashbash-backend-go/pkg/model"
)

type chainGeneratorService struct {
	hashFunction            hashFunction
	reductionFunctionFamily reductionFunctionFamily
}

func newChainGeneratorService(
	hashFunction hashFunction,
	reductionFunctionFamily reductionFunctionFamily,
) *chainGeneratorService {
	return &chainGeneratorService{
		hashFunction:            hashFunction,
		reductionFunctionFamily: reductionFunctionFamily,
	}
}

func (service *chainGeneratorService) generateRainbowChainLinkFromHash(
	digest []byte,
	nextChainIndex int,
	numLinks int,
) rainbowChainLink {
	if numLinks > 0 {
		plaintext := service.reductionFunctionFamily(digest, nextChainIndex)
		return service.generateRainbowChainLinkFromPlaintext(plaintext, nextChainIndex+1, numLinks)
	}

	return rainbowChainLink{hashedPlaintext: digest}
}

func (service *chainGeneratorService) generateRainbowChainLinkFromPlaintext(
	plaintext string,
	nextChainIndex int,
	numLinks int,
) rainbowChainLink {
	// Hash the plaintext, generating the first link
	chainLink := rainbowChainLink{
		plaintext:       plaintext,
		hashedPlaintext: service.hashFunction.apply(plaintext),
	}

	// From this link to the end of the chain
	for i := 0; i < numLinks-1; i++ {
		// Hash the current key, then reduce it to the next key
		reducedPlaintext := service.reductionFunctionFamily(chainLink.hashedPlaintext, nextChainIndex+i)
		hashedDigest := service.hashFunction.apply(reducedPlaintext)

		chainLink.plaintext = reducedPlaintext
		chainLink.hashedPlaintext = hashedDigest
	}

	return chainLink
}

func (service *chainGeneratorService) generateRainbowChain(startPlaintext string, chainLength int) model.RainbowChain {
	endingLink := service.generateRainbowChainLinkFromPlaintext(startPlaintext, 0, chainLength)
	return model.RainbowChain{
		StartPlaintext: startPlaintext,
		EndHash:        fmt.Sprintf("%x", endingLink.hashedPlaintext),
	}
}
