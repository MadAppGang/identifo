package mock

import "context"

type Storage struct{}

func (u *Storage) Ready(ctx context.Context) error {
	return nil
}

func (u *Storage) Connect(ctx context.Context) error {
	return nil
}

func (u *Storage) Close(ctx context.Context) error {
	return nil
}
