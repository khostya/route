package oncall

import (
	"github.com/IBM/sarama"
	"homework/internal/dto"
	"homework/internal/infrastructure/kafka"
)

type KafkaSender struct {
	producer *kafka.Producer
	topic    kafka.Topic
}

func NewKafkaSender(producer *kafka.Producer, topic kafka.Topic) *KafkaSender {
	return &KafkaSender{
		producer,
		topic,
	}
}

func (s *KafkaSender) SendAsyncMessage(message dto.CallMessage) error {
	kafkaMsg, err := s.buildMessage(message)
	if err != nil {
		return err
	}

	s.producer.SendAsyncMessage(kafkaMsg)
	return nil
}

func (s *KafkaSender) SendMessage(message *dto.CallMessage) error {
	kafkaMsg, err := s.buildMessage(*message)
	if err != nil {
		return err
	}

	_, _, err = s.producer.SendSyncMessage(kafkaMsg)
	return err
}

func (s *KafkaSender) SendMessages(messages []dto.CallMessage) error {
	var kafkaMsg []*sarama.ProducerMessage

	for _, m := range messages {
		message, err := s.buildMessage(m)
		kafkaMsg = append(kafkaMsg, message)

		if err != nil {
			return err
		}
	}

	return s.producer.SendSyncMessages(kafkaMsg)
}

func (s *KafkaSender) buildMessage(message dto.CallMessage) (*sarama.ProducerMessage, error) {
	msg, err := message.Marshal()

	if err != nil {
		return nil, err
	}

	return &sarama.ProducerMessage{
		Topic:     string(s.topic),
		Value:     sarama.ByteEncoder(msg),
		Partition: -1,
		Key:       sarama.StringEncoder(message.Method),
	}, nil
}

func (s *KafkaSender) Close() error {
	return s.producer.Close()
}
