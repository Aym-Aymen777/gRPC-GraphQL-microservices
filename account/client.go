package account

import (
	"context"

	"github.com/Aym-Aymen777/gRPC-GraphQL-microservices/account/pb"
	"google.golang.org/grpc"
)

type Client struct {
	conn    *grpc.ClientConn
	Service pb.AccountServiceClient
}

func NewGRPCClient(url string) (*Client, error) {
	conn, err := grpc.Dial(url, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	Service := pb.NewAccountServiceClient(conn)
	return &Client{conn: conn, Service: Service}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) CreateAccount(username string, email string) (*pb.CreateAccountResponse, error) {
	return c.Service.CreateAccount(context.Background(), &pb.CreateAccountRequest{
		Username: username,
		Email:    email,
	})
}

func (c *Client) GetAccountDetails(id string) (*pb.GetAccountDetailsResponse, error) {
	return c.Service.GetAccountDetails(context.Background(), &pb.GetAccountDetailsRequest{
		Id: id,
	})
}

func (c *Client) GetAccounts(skip, limit uint64) (*pb.GetAccountsResponse, error) {
	return c.Service.GetAccounts(context.Background(), &pb.GetAccountsRequest{
		Skip:  skip,
		Limit: limit,
	})
}
