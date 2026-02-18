package account

import (
	"context"
	"log"

	"rjfield.com/backend/db"
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
		},
	}, nil
}

// ListAccounts retrieves all accounts for the given user ID.
func (s *AccountServer) ListAccounts(ctx context.Context, req *pb.ListAccountsRequest) (*pb.ListAccountsReply, error) {
	log.Printf("ListAccounts() - received: %v", req.GetParent())

	user_id := db.UserIDFromResourceName(req.GetParent())

	accounts, err := db.ListAccounts(user_id)
	if err != nil {
		log.Printf("ListAccounts() - error fetching accounts for user %s: %v", req.GetParent(), err)
		return nil, err
	}
	log.Printf("ListAccounts() - fetched %d accounts for user %s", len(accounts), req.GetParent())
	return &pb.ListAccountsReply{
		Accounts: accounts,
	}, nil
}
