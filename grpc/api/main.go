package main

import (
	"context"
	pb "currency/proto/currency"
	"log"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	grpcServerAddress = "localhost:50051"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	conn, err := grpc.NewClient(grpcServerAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}

	defer conn.Close()
	grpcClient := pb.NewCurrencyServiceClient(conn)

	server := NewServer(grpcClient)

	go func() {
		log.Println("Starting HTTP server on :8080")
		if err := server.httpServer.ListenAndServe(); err != nil {
			panic(err)
		}
	}()

	<-ctx.Done()
	log.Println("Closing...")
	server.httpServer.Shutdown(ctx)
}
