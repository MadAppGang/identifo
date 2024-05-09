package grpc

import (
	"context"
	"io"
	"time"

	pp "github.com/madappgang/identifo/v2/impersonation/grpc/impersonation_provider"
	"github.com/madappgang/identifo/v2/model"
	ggrpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewImpersonationService(address string, timeout time.Duration) (model.ImpersonationProvider, error) {
	conn, err := ggrpc.Dial(
		address,
		ggrpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	if timeout == 0 {
		timeout = time.Second
	}

	client := pp.NewImpersonationProviderClient(conn)

	return &GRPCClient{
		Client:   client,
		Timeout:  timeout,
		Closable: conn,
	}, nil
}

type GRPCClient struct {
	Client  pp.ImpersonationProviderClient
	Timeout time.Duration

	Closable io.Closer
}

func (p *GRPCClient) CanImpersonate(ctx context.Context, appID string, adminUser, impUser model.User) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, p.Timeout)
	defer cancel()

	resp, err := p.Client.CanImpersonate(ctx, &pp.CanImpersonateRequest{
		AppId:            appID,
		AdminUser:        grpcUser(adminUser),
		ImpersonatedUser: grpcUser(impUser),
	})
	if err != nil {
		return false, err
	}

	return resp.Ok, nil
}

func grpcUser(u model.User) *pp.User {
	return &pp.User{
		Id:         u.ID,
		Email:      u.Email,
		Active:     u.Active,
		AccessRole: u.AccessRole,
		Anonymous:  u.Anonymous,
		Scopes:     u.Scopes,
	}
}

func (p *GRPCClient) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), p.Timeout)
	defer cancel()

	p.Client.Close(ctx, &pp.CloseRequest{})
	if p.Closable != nil {
		return p.Closable.Close()
	}

	return nil
}
