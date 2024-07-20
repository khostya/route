package oncall

import (
	"github.com/IBM/sarama"
	"homework/internal/dto"
	"homework/internal/infrastructure/kafka"
)

type KafkaProducer struct {
	producer *kafka.Producer
	topic    kafka.Topic
}

func NewKafkaProducer(producer *kafka.Producer, topic kafka.Topic) *KafkaProducer {
	return &KafkaProducer{
		producer,
		topic,
	}
}

func (p *KafkaProducer) SendAsyncMessage(message dto.CallMessage) error {
	kafkaMsg, err := p.buildMessage(message)
	if err != nil {
		return err
	}

	p.producer.SendAsyncMessage(kafkaMsg)
	return nil
}

func (p *KafkaProducer) buildMessage(message dto.CallMessage) (*sarama.ProducerMessage, error) {
	msg, err := message.Marshal()
	if err != nil {
		return nil, err
	}

	return &sarama.ProducerMessage{
		Topic:     string(p.topic),
		Value:     sarama.ByteEncoder(msg),
		Partition: -1,
		Key:       sarama.StringEncoder(message.Method),
	}, nil
}

func (p *KafkaProducer) Close() error {
	return p.producer.Close()
}
