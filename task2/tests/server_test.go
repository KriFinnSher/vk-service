package tests

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"testing"
	pb "vk-Service/task2/grpc"
)

const (
	subject = "weather"
	message = "very hot"
	port    = ":8080"
)

func TestConn(t *testing.T) {
	conn, err := grpc.NewClient(
		port,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatalf("failed to create client conn: %v", err)
	}
	defer conn.Close()

	client := pb.NewPubSubClient(conn)

	subReq := &pb.SubscribeRequest{Key: subject}
	stream, err := client.Subscribe(context.Background(), subReq)
	if err != nil {
		t.Fatalf("failed to subscribe: %v", err)
	}

	pubReq := &pb.PublishRequest{Key: subject, Data: message}
	if _, err := client.Publish(context.Background(), pubReq); err != nil {
		t.Fatalf("failed to publish: %v", err)
	}

	msg, err := stream.Recv()
	if err != nil {
		t.Fatalf("failed to receive from stream: %v", err)
	}

	if msg.GetData() != message {
		t.Errorf("Expected message %q, got %q", message, msg.GetData())
	}
}
