//go:build integration

package kafka

import (
	"context"
	"github.com/IBM/sarama"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"homework/internal/dto"
	"homework/tests/kafka/kafka"
	"testing"
	"time"
)

type OnCallTestSuite struct {
	suite.Suite
	ctx   context.Context
	kafka *kafka.Kafka
}

func TestOnCall(t *testing.T) {
	suite.Run(t, new(OnCallTestSuite))
}

func (s *OnCallTestSuite) SetupSuite() {
	s.T().Parallel()
	s.ctx = context.Background()
	s.kafka = kafka.NewOnCallFromEnv(s.ctx, "_call")
}

func (s *OnCallTestSuite) TearDownSuite() {
	s.kafka.Close()
}

func (s *OnCallTestSuite) SetupTest() {
	s.kafka.SetUp()
}

func (s *OnCallTestSuite) TearDownTest() {
	s.kafka.TearDown()
}

func (s *OnCallTestSuite) TestAsyncSend() {
	onCallMessage := NewOnCallMessage()
	err := s.kafka.OnCallSender.SendAsyncMessage(onCallMessage)
	require.NoError(s.T(), err)
}

func (s *OnCallTestSuite) TestSend() {
	onCallMessage := NewOnCallMessage()
	err := s.kafka.OnCallSender.SendMessage(&onCallMessage)
	require.NoError(s.T(), err)
}

func (s *OnCallTestSuite) TestSendMessages() {
	onCallMessages := []dto.CallMessage{NewOnCallMessage(), NewOnCallMessage()}
	err := s.kafka.OnCallSender.SendMessages(onCallMessages)
	require.NoError(s.T(), err)
}

func (s *OnCallTestSuite) TestGetMessages() {
	onCallMessage := NewOnCallMessage()

	received := false
	err := s.kafka.OnCallReceiver.Subscribe(s.kafka.Topic, func(message *sarama.ConsumerMessage) {
		received = true

		var m dto.CallMessage
		err := m.Unmarshal(message.Value)
		require.NoError(s.T(), err)
		require.EqualExportedValues(s.T(), onCallMessage, m)
	})
	require.NoError(s.T(), err)

	err = s.kafka.OnCallSender.SendMessage(&onCallMessage)
	require.NoError(s.T(), err)

	time.Sleep(time.Second * 3)
	require.True(s.T(), received)
}
