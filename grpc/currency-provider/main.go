package main

import (
	"context"
	pb "currency/proto/currency"
	"fmt"
	"log"
	"net"
	"net/http"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	httpClient := &http.Client{}

	grpcServer := grpc.NewServer()
	pb.RegisterCurrencyServiceServer(grpcServer, &Server{httpClient: httpClient})

	go func() {
		fmt.Println("gRPC server listening on :50051")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("Closing...")
	grpcServer.GracefulStop()
}
