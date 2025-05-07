package handlers

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	pb "vk-Service/task2/grpc"
)

func (s *Server) Subscribe(req *pb.SubscribeRequest, stream pb.PubSub_SubscribeServer) error {
	s.Logger.Info("subscribe request [attempt]", "subject", req.GetKey())
	key := req.GetKey()

	if key == "" {
		s.Logger.Warn("subscribe request [empty subject]")
		return status.Error(codes.InvalidArgument, "key is required")
	}

	cb := func(msg interface{}) {
		str, ok := msg.(string)
		if !ok {
			s.Logger.Error("subscribe request [callback]", "err", "message should be string")
			return
		}

		err := stream.Send(&pb.Event{Data: str})
		if err != nil {
			s.Logger.Error("subscribe request [callback]", "err", err)
		}
	}

	sub, err := s.Internal.Subscribe(key, cb)
	if err != nil {
		return status.Errorf(codes.Internal, "subscribe failed: %v", err)
	}
	defer sub.Unsubscribe()
	<-stream.Context().Done()
	s.Logger.Info("subscribe request [done]", "subject", req.GetKey())
	return nil
}
