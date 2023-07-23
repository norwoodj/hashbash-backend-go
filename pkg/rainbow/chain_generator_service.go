package rainbow

import (
	"github.com/norwoodj/hashbash-backend-go/pkg/model"
)

type chainGeneratorService struct {
	RandomStringGenerator
	hashFunction            HashFunction
	reductionFunctionFamily reductionFunctionFamily
}

func newChainGeneratorService(
	hashFunction HashFunction,
	reductionFunctionFamily reductionFunctionFamily,
	randomStringSeed int,
) *chainGeneratorService {
	return &chainGeneratorService{
		RandomStringGenerator:   *NewRandomStringGenerator(randomStringSeed),
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
		hashedPlaintext: service.hashFunction.Apply(plaintext),
	}

	// From this link to the end of the chain
	for i := 0; i < numLinks-1; i++ {
		// Hash the current key, then reduce it to the next key
		reducedPlaintext := service.reductionFunctionFamily(chainLink.hashedPlaintext, nextChainIndex+i)
		hashedDigest := service.hashFunction.Apply(reducedPlaintext)

		chainLink.plaintext = reducedPlaintext
		chainLink.hashedPlaintext = hashedDigest
	}

	return chainLink
}

func (service *chainGeneratorService) generateRainbowChain(startPlaintext string, chainLength int) model.RainbowChain {
	endingLink := service.generateRainbowChainLinkFromPlaintext(startPlaintext, 0, chainLength)
	return model.RainbowChain{
		StartPlaintext: startPlaintext,
		EndHash:        endingLink.hashedPlaintext,
	}
}
