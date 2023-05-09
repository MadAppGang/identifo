package grpc

import (
	"github.com/madappgang/identifo/v2/model"
	"github.com/madappgang/identifo/v2/storage/grpc/proto"
	"github.com/madappgang/identifo/v2/storage/grpc/shared"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// NewUserStorage creates and inits MongoDB user storage.
func NewUserStorage(settings model.GRPCSettings) (model.UserStorage, error) {
	// Set up a connection to the server.
	conn, err := grpc.Dial(settings.Address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	uc := proto.NewUserStorageClient(&grpc.ClientConn{})

	user := shared.GRPCClient{Client: uc, Closable: conn}

	return user, nil
}
