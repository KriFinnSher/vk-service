package internal

import (
	"google.golang.org/grpc"
	"log/slog"
	"net"
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

	logger.Info("server started", "addr", lis.Addr())

	server := grpc.NewServer()
	pb.RegisterPubSubServer(server, &handlers.Server{
		Internal: subpub.NewSubPub(),
		Logger:   logger,
	})
	if err = server.Serve(lis); err != nil {
		logger.Error("failed to serve", "err", err, "addr", lis.Addr())
		return
	}
}
