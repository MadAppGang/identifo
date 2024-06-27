package grpc

import (
	"context"
	"encoding/json"
	"io"

	"github.com/madappgang/identifo/v2/model"
	pp "github.com/madappgang/identifo/v2/user_payload_provider/grpc/payload_provider"
)

// Here is the gRPC server that GRPCClient talks to.
type GRPCServer struct {
	// This is the real implementation
	Impl model.TokenPayloadProvider
	pp.UnimplementedPayloadProviderServiceServer
}

func (s *GRPCServer) TokenPayload(ctx context.Context, request *pp.TokenPayloadRequest) (*pp.TokenPayloadResponse, error) {
	payload, err := s.Impl.TokenPayloadForApp(request.AppId, request.AppName, request.UserId)
	if err != nil {
		return nil, err
	}

	paylodBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	response := &pp.TokenPayloadResponse{
		PayloadJson: string(paylodBytes),
	}

	return response, nil
}

func (s *GRPCServer) Close(ctx context.Context, _ *pp.CloseRequest) (*pp.CloseResponse, error) {
	c, ok := s.Impl.(io.Closer)
	if ok {
		c.Close()
	}
	return &pp.CloseResponse{}, nil
}
