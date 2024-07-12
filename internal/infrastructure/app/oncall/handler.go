package oncall

import (
	"fmt"
	"github.com/IBM/sarama"
	"homework/internal/dto"
)

func NewTopicHandler() (<-chan string, HandleFunc) {
	out := make(chan string)
	return out, func(message *sarama.ConsumerMessage) {
		var callMessage dto.CallMessage
		err := callMessage.Unmarshal(message.Value)
		if err != nil {
			out <- fmt.Sprintf("Consumer error: %v \n", err)
			return
		}
		out <- fmt.Sprintf("kafka: %s", callMessage.String())
	}
}
