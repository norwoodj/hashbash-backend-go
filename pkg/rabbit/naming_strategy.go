package rabbit

import "fmt"

type QueueNamingStrategy interface {
	GetQueueName(exchangeName string, routingKey string) string
}

type DefaultQueueNamingStrategy struct {}

func (DefaultQueueNamingStrategy) GetQueueName(exchangeName string, routingKey string) string {
	return fmt.Sprintf("%s.%s", exchangeName, routingKey)
}
