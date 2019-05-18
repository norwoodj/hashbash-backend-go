package util

import (
	"math/rand"
	"strings"
)

func IntOrDefault(value int64, defaultValue int64) int64 {
	if value == 0 {
		return defaultValue
	}

	return value
}

func StringOrDefault(value string, defaultValue string) string {
	if value == "" {
		return defaultValue
	}

	return value
}

func RandomString(characterSet *string, stringLength int64) string {
	stringBuilder := strings.Builder{}
	characterSetSize := len(*characterSet)
	var i int64
	for i = 0; i < stringLength; i++ {
		index := rand.Intn(characterSetSize)
		stringBuilder.WriteByte((*characterSet)[index])
	}

	return stringBuilder.String()
}
