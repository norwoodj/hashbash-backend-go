package rainbow

import (
	"math/rand"
	"strings"
	"time"
)

type RandomStringGenerator struct {
	rng *rand.Rand
}

func NewRandomStringGenerator(seed int) *RandomStringGenerator {
	fullSeed := int64(seed * time.Now().Nanosecond())
	source := rand.NewSource(fullSeed)
	return &RandomStringGenerator{rng: rand.New(source)}
}

func (rsg *RandomStringGenerator) NewRandomString(characterSet string, stringLength int64) string {
	stringBuilder := strings.Builder{}
	characterSetSize := len(characterSet)
	var i int64
	for i = 0; i < stringLength; i++ {
		index := rsg.rng.Intn(characterSetSize)
		stringBuilder.WriteByte((characterSet)[index])
	}

	return stringBuilder.String()
}
