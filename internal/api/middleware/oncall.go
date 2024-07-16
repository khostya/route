package middleware

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"homework/internal/dto"
	"log"
	"time"
)

type onCallProducer interface {
	SendAsyncMessage(message dto.OnCallMessage) error
}

func OnCall(producer onCallProducer) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		raw, _ := protojson.Marshal((req).(proto.Message))

		err = producer.SendAsyncMessage(dto.OnCallMessage{
			CalledAt: time.Now(),
			Method:   info.FullMethod,
			Args:     string(raw),
		})
		if err != nil {
			log.Printf("[interceptor.OnCall] error:%v", err.Error())
		}

		resp, err = handler(ctx, req)
		return
	}
}
