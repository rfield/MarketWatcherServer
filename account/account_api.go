package account

import (
	"context"
	"log"

	"rjfield.com/backend/generated/pb"
)

type AccountServer struct {
	pb.UnimplementedAccountServiceServer
}

// GetAccount retrieves account details for the given account ID.
func (s *AccountServer) GetAccount(ctx context.Context, req *pb.GetAccountRequest) (*pb.GetAccountReply, error) {
	log.Printf("GetAccount() - received: %v", req.GetAccountId())
	return &pb.GetAccountReply{
		Account: &pb.Account{
			AccountName: "Sample Account",
			Balance:     1000.50,
		},
	}, nil
}

// Deprecated: Use user.proto's AuthenticateUser rpc instead.
func (s *AccountServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginReply, error) {
	log.Printf("Login() - received: %s / %s", req.GetUsername(), req.GetPassword())
	return &pb.LoginReply{
		Token:     "sample_token",
		AccountId: "12345",
	}, nil
}
