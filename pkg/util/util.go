package util

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
