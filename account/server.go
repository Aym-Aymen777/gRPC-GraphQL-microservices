package account

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"

	"github.com/Aym-Aymen777/gRPC-GraphQL-microservices/account/pb"
)

type grpcServer struct {
	pb.UnimplementedAccountServiceServer
	Service Service
}

func ListenGRPC(s Service, port string) error  {
	listen, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal(err)
		return err
	}
	serv := grpc.NewServer()
	pb.RegisterAccountServiceServer(serv, &grpcServer{Service: s})
	if err := serv.Serve(listen); err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

func (s *grpcServer) CreateAccount(ctx context.Context, req *pb.CreateAccountRequest) (*pb.CreateAccountResponse, error) {
	account, err := s.Service.CreateAccount(ctx, req.Username, req.Email)
	if err != nil {
		return nil, err
	}
	return &pb.CreateAccountResponse{
		Id: account.ID,
	}, nil
}

func (s *grpcServer) GetAccountDetails(ctx context.Context, req *pb.GetAccountDetailsRequest) (*pb.GetAccountDetailsResponse, error) {
	account, err := s.Service.GetAccountDetails(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.GetAccountDetailsResponse{
		Id:       account.ID,
		Username: account.Username,
		Email:    account.Email,
	}, nil
}

func (s *grpcServer) GetAccounts(ctx context.Context, req *pb.GetAccountsRequest) (*pb.GetAccountsResponse, error) {
	accounts, err := s.Service.GetAccounts(ctx, req.Skip, req.Limit)
	if err != nil {
		return nil, err
	}
	response := &pb.GetAccountsResponse{}
	for _, account := range accounts {
		response.Accounts = append(response.Accounts, &pb.GetAccountDetailsResponse{
			Id:       account.ID,
			Username: account.Username,
			Email:    account.Email,
		})
	}
	return response, nil
}
