package rainbow

import (
	"strings"
)

// Sometimes, I hate go. I mean, I love you, but also please implement generics already
func abs(x int) int {
	if x < 0 {
		return -x
	}

	return x
}

type reductionFunctionFamily func([]byte, int) string

func getDefaultReductionFunctionFamily(passwordLength int, charset string) reductionFunctionFamily {
	characterSetLength := len(charset)
	return func(digest []byte, chainIndex int) string {
		plaintextBuilder := strings.Builder{}

		for i := 0; i < passwordLength; i++ {
			value := abs((int(digest[i]) ^ chainIndex) % characterSetLength)
			plaintextBuilder.WriteByte(charset[value])
		}

		return plaintextBuilder.String()
	}
}
