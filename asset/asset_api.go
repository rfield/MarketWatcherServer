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

func (s *AssetServer) ListAssetsForUser(ctx context.Context, req *pb.ListAssetsForUserRequest) (*pb.ListAssetsForUserReply, error) {
	log.Printf("ListAssetsForUser() - received: %v", req.GetUserId())
	// return &pb.ListAssetsForUserReply{
	// 	Assets: []*pb.Asset{
	// 		{
	// 			UserId:        req.GetUserId(),
	// 			AccountName:   "Sample Account",
	// 			Ticker:        "AAPL",
	// 			HoldingAmount: 100.50,
	// 		},
	// 	},
	// }, nil
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
