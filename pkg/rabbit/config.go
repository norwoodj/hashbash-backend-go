package rabbit

type Config struct {
	// Connection fields
	Hostname string
	Username string
	Password string

	// Routing Configuration
	DeadLetterExchangeName string
	DeadLetterQueueSuffix  string

	// Naming Configuration
	QueueNamingStrategy QueueNamingStrategy
}

func NewConfig(
	hostname string,
	username string,
	password string,
) *Config {
	return &Config{
		Hostname:               hostname,
		Username:               username,
		Password:               password,
		DeadLetterExchangeName: DefaultDeadLetterExchangeName,
		DeadLetterQueueSuffix:  DefaultDeadLetterQueueSuffix,
		QueueNamingStrategy:    DefaultQueueNamingStrategy{},
	}
}
