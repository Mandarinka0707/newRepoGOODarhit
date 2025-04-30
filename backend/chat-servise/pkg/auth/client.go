// chat-servise/pkg/auth/client.go
package auth

import (
	"context"

	"errors"

	"github.com/Mandarinka0707/newRepoGOODarhit/proto"
	"google.golang.org/grpc"
)

type Client struct {
	conn proto.AuthServiceClient
}

func NewClient(grpcAddr string) (*Client, error) {
	conn, err := grpc.Dial(grpcAddr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return &Client{
		conn: proto.NewAuthServiceClient(conn),
	}, nil
}

func (c *Client) ValidateToken(ctx context.Context, token string) (*User, error) {
	resp, err := c.conn.ValidateToken(ctx, &proto.ValidateTokenRequest{Token: token})
	if err != nil {
		return nil, err
	}

	if resp.Username == "" {
		return nil, errors.New("username not provided by auth service")
	}

	return &User{
		ID:       resp.UserId,
		Username: resp.Username,
	}, nil
}

type User struct {
	ID       int64
	Username string
}
