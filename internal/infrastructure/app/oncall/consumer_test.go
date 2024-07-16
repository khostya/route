package oncall

import (
	"context"
	"github.com/IBM/sarama"
	mock_kafka "github.com/IBM/sarama/mocks"
	"github.com/stretchr/testify/require"
	"homework/internal/infrastructure/kafka"
	"testing"
)

func TestKafkaConsumer_Subscribe(t *testing.T) {
	topic := "call"
	consumer := mock_kafka.NewConsumer(t, nil)
	consumer.SetTopicMetadata(map[string][]int32{topic: {1}})
	receiver := NewKafkaReceiver(&kafka.Consumer{
		SingleConsumer: consumer,
	})

	consumer.ExpectConsumePartition(topic, 1, -1).YieldMessage(&sarama.ConsumerMessage{})
	ctx, cancel := context.WithCancel(context.Background())

	err := receiver.Subscribe("call", func(message *sarama.ConsumerMessage) {
		cancel()
	})

	require.NoError(t, err)
	<-ctx.Done()
}
