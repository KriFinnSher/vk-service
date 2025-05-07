package handlers

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"log/slog"
	"vk-Service/task1/subpub"
	pb "vk-Service/task2/grpc"
)

type Server struct {
	pb.UnimplementedPubSubServer
	Internal subpub.SubPub
	Logger   *slog.Logger
}

func (s *Server) Publish(_ context.Context, req *pb.PublishRequest) (*emptypb.Empty, error) {
	s.Logger.Info("publish request [attempt]", "subject", req.GetKey(), "message", req.GetData())
	if req.GetKey() == "" {
		s.Logger.Warn("publish request [empty subject]", "message", req.GetData())
		return nil, status.Error(codes.InvalidArgument, "key is required")
	}

	if err := s.Internal.Publish(req.GetKey(), req.GetData()); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to publish message: %v", err)
	}
	s.Logger.Info("publish request [done]", "subject", req.GetKey(), "message", req.GetData())

	return &emptypb.Empty{}, nil
}
