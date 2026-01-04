package main

import (
	"log"
	"net"

	"rjfield.com/backend/account"
	"rjfield.com/backend/asset"
	"rjfield.com/backend/generated/pb"
	"rjfield.com/backend/price"
	"rjfield.com/backend/user"

	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterAccountServiceServer(s, &account.AccountServer{})
	pb.RegisterPriceServiceServer(s, &price.PriceServer{})
	pb.RegisterUserServiceServer(s, &user.UserServer{})
	pb.RegisterAssetServiceServer(s, &asset.AssetServer{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
