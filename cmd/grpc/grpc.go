package main

import (
	"context"
	"fmt"
	"github.com/go-chi/cors"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"homework/config"
	"homework/internal/api"
	"homework/internal/api/middleware"
	"homework/internal/infrastructure/app/oncall"
	"homework/internal/service"
	"homework/pkg/api/order/v1"
	gw "homework/pkg/api/order/v1"
	"log"
	"net"
	"net/http"
	"sync"
)

func startGrpcServer(ctx context.Context, cancelFunc context.CancelFunc, orderService *service.OrderService, producer *oncall.KafkaProducer) *sync.WaitGroup {
	cfg := config.MustNewApiConfig()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GrpcPort))
	log.Println("Start")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	errGw := gw.RegisterOrderHandlerFromEndpoint(ctx, mux, cfg.GrpcENDPOINT, opts)
	if errGw != nil {
		log.Fatalf("failed to RegisterOrderHandlerFromEndpoint: %v", errGw)
	}

	gwServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.HttpPort),
		Handler: cors.AllowAll().Handler(mux),
	}

	go func() {
		promHandler := promhttp.Handler()
		mux.HandlePath(http.MethodGet, "/metrics", func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
			promHandler.ServeHTTP(w, r)
		})

		err := gwServer.ListenAndServe()
		if err != nil {
			cancelFunc()
		}
	}()

	grpcServer := grpc.NewServer(grpc.ChainUnaryInterceptor(middleware.OnCall(producer)))

	order.RegisterOrderServer(grpcServer, api.NewOrderService(orderService))
	go func() {
		err := grpcServer.Serve(lis)
		if err != nil {
			cancelFunc()
		}
	}()

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		<-ctx.Done()
		grpcServer.GracefulStop()
		wg.Done()

		gwServer.Shutdown(ctx)
		wg.Done()
	}()

	return &wg
}
