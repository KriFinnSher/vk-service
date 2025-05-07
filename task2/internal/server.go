package internal

import (
	"google.golang.org/grpc"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"vk-Service/task1/subpub"
	pb "vk-Service/task2/grpc"
	"vk-Service/task2/internal/config"
	"vk-Service/task2/internal/handlers"
)

func RunServer(logger *slog.Logger) {
	lis, err := net.Listen("tcp", ":"+config.AppConfig.Server.Port)
	if err != nil {
		logger.Error("failed to listen", "port", config.AppConfig.Server.Port)
		return
	}

	server := grpc.NewServer()
	pb.RegisterPubSubServer(server, &handlers.Server{
		Internal: subpub.NewSubPub(),
		Logger:   logger,
	})

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		logger.Info("server started", "addr", lis.Addr())
		if err = server.Serve(lis); err != nil {
			logger.Error("failed to serve", "err", err)
		}
	}()

	<-stop
	logger.Info("shutdown signal received")

	server.GracefulStop()
	logger.Info("server gracefully stopped")

	wg.Wait()
}
