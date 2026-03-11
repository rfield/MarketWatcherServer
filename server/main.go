package main

import (
	"log"
	"net"
	"strconv"

	"rjfield.com/backend/account"
	"rjfield.com/backend/asset"
	"rjfield.com/backend/config"
	"rjfield.com/backend/generated/pb"
	"rjfield.com/backend/notification"
	"rjfield.com/backend/price"
	"rjfield.com/backend/user"

	"google.golang.org/grpc"
)

func main() {
	config := config.GetConfig()
	grpcPort := config.Grpc.Port

	log.Printf("main() - starting server on port %d", grpcPort)

	lis, err := net.Listen("tcp", ":"+strconv.Itoa(grpcPort))
	if err != nil {
		log.Fatalf("main() - failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterAccountServiceServer(s, &account.AccountServer{})
	pb.RegisterPriceServiceServer(s, &price.PriceServer{})
	pb.RegisterUserServiceServer(s, &user.UserServer{})
	pb.RegisterAssetServiceServer(s, &asset.AssetServer{})
	pb.RegisterNotificationServiceServer(s, &notification.NotificationServer{})
	log.Printf("main() - server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("main() - failed to serve: %v", err)
	}
}
