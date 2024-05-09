package grpc

import (
	"context"
	"io"

	pp "github.com/madappgang/identifo/v2/impersonation/grpc/impersonation_provider"
	"github.com/madappgang/identifo/v2/model"
)

// Here is the gRPC server that GRPCClient talks to.
type GRPCServer struct {
	// This is the real implementation
	Impl model.ImpersonationProvider
	pp.UnimplementedImpersonationProviderServer
}

func (s *GRPCServer) CanImpersonate(ctx context.Context, req *pp.CanImpersonateRequest) (*pp.CanImpersonateResponse, error) {
	ok, err := s.Impl.CanImpersonate(
		ctx,
		req.AppId,
		userFromGRPC(req.AdminUser),
		userFromGRPC(req.ImpersonatedUser),
	)
	if err != nil {
		return nil, err
	}

	return &pp.CanImpersonateResponse{
		Ok: ok,
	}, nil
}

func userFromGRPC(u *pp.User) model.User {
	return model.User{
		ID:         u.Id,
		Email:      u.Email,
		Active:     u.Active,
		AccessRole: u.AccessRole,
		Anonymous:  u.Anonymous,
		Scopes:     u.Scopes,
	}
}

func (s *GRPCServer) Close(ctx context.Context, _ *pp.CloseRequest) (*pp.CloseResponse, error) {
	c, ok := s.Impl.(io.Closer)
	if ok {
		c.Close()
	}
	return &pp.CloseResponse{}, nil
}
