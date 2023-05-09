package grpc

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/madappgang/identifo/v2/model"
	pp "github.com/madappgang/identifo/v2/user_payload_provider/grpc/payload_provider"
	ggrpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewTokenPayloadProvider(address string, timeout time.Duration) (model.TokenPayloadProvider, error) {
	conn, err := ggrpc.Dial(
		address,
		ggrpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	if timeout == 0 {
		timeout = time.Second
	}

	client := pp.NewPayloadProviderServiceClient(conn)

	return &GRPCClient{
		Client:   client,
		Timeout:  timeout,
		Closable: conn,
	}, nil
}

type GRPCClient struct {
	Client  pp.PayloadProviderServiceClient
	Timeout time.Duration

	Closable io.Closer
}

func (p *GRPCClient) TokenPayloadForApp(appId, appName, userID string) (map[string]interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), p.Timeout)
	defer cancel()

	resp, err := p.Client.TokenPayload(ctx, &pp.TokenPayloadRequest{
		UserId:  userID,
		AppId:   appId,
		AppName: appName,
	})
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}

	dec := json.NewDecoder(strings.NewReader(resp.PayloadJson))
	if err := dec.Decode(&result); err != nil {
		return nil, fmt.Errorf("getting user token payload from provider, could not parse response body with error: %v", err)
	}

	return result, nil
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
