package main

import (
	"context"
	"fmt"
	"github.com/go-chi/cors"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"homework/config"
	"homework/internal/api"
	"homework/internal/service"
	"homework/pkg/api/order/v1"
	gw "homework/pkg/api/order/v1"
	"log"
	"net"
	"net/http"
	"sync"
)

type GRPC struct {
}

func startGrpcServer(ctx context.Context, cancelFunc context.CancelFunc, orderService *service.Order) *sync.WaitGroup {
	var wg sync.WaitGroup

	cfg, err := config.NewApiConfig()
	if err != nil {
		log.Fatalln(err)
	}

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

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := gwServer.ListenAndServe()
		if err != nil {
			cancelFunc()
		}
	}()

	grpcServer := grpc.NewServer()
	order.RegisterOrderServer(grpcServer, api.NewOrderService(orderService))
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := grpcServer.Serve(lis)
		if err != nil {
			cancelFunc()
		}
	}()

	go func() {
		<-ctx.Done()
		grpcServer.Stop()
		gwServer.Close()
	}()

	return &wg
}
