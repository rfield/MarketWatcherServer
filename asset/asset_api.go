package asset

import (
	"context"
	"log"

	"rjfield.com/backend/db"
	"rjfield.com/backend/generated/pb"
)

type AssetServer struct {
	pb.UnimplementedAssetServiceServer
}

func (s *AssetServer) ListAssets(ctx context.Context, req *pb.ListAssetsRequest) (*pb.ListAssetsReply, error) {

	user_id := db.UserIDFromResourceName(req.GetParent())
	account_id := db.AccountIDFromResourceName(req.GetParent())

	log.Printf("ListAssets() - received: %v, %v", user_id, account_id)

	a, err := db.ListAssets(user_id, account_id)
	if err != nil {
		log.Printf("ListAssets() - error fetching assets for user %s, account %s: %v", user_id, account_id, err)
		return nil, err
	}
	log.Printf("ListAssets() - fetched %d assets for user %s", len(a), user_id)
	return &pb.ListAssetsReply{
		Assets: a,
	}, nil
}

func (s *AssetServer) ListAssetsForUser(ctx context.Context, req *pb.ListAssetsForUserRequest) (*pb.ListAssetsForUserReply, error) {
	log.Printf("ListAssetsForUser() - received: %v", req.GetUserId())

	a, err := db.ListAssetsByUser(req.GetUserId())
	if err != nil {
		log.Printf("ListAssetsForUser() - error fetching assets for user %s: %v", req.GetUserId(), err)
		return nil, err
	}
	log.Printf("ListAssetsForUser() - fetched %d assets for user %s", len(a), req.GetUserId())
	return &pb.ListAssetsForUserReply{
		Assets: a,
	}, nil
}
