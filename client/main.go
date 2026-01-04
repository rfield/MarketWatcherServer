package main

import (
	"context"
	"flag"
	"log"
	"time"

	"rjfield.com/backend/generated/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
	// acct    = flag.String("account", "123abc", "Account ID to get")
	// product = flag.String("product", "MSFT", "Price ID to get")
)

func main() {
	flag.Parse()
	conn, err := grpc.NewClient(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	ast := pb.NewAssetServiceClient(conn)
	arsp, err := ast.ListAssetsForUser(ctx, &pb.ListAssetsForUserRequest{
		UserId: "b0159a28-5a03-403f-ba2f-192b41a32d9a",
	})
	if err != nil {
		log.Fatalf("could not list assets: %v", err)
	}
	for _, asset := range arsp.GetAssets() {
		log.Printf("Asset: UserID=%s, AccountName=%s, Ticker=%s, HoldingAmount=%.2f",
			asset.GetUserId(), asset.GetAccountName(), asset.GetTicker(), asset.GetHoldingAmount())
	}

	/*

		ac := pb.NewAccountServiceClient(conn)
		lr, err := ac.Login(ctx, &pb.LoginRequest{
			Username: "user1",
			Password: "pass1",
		})
		if err != nil {
			log.Fatalf("could not login: %v", err)
		}
		log.Printf("Login: %s", lr.GetAccountId())

		ar, err := ac.GetAccount(ctx, &pb.GetAccountRequest{AccountId: *acct})
		if err != nil {
			log.Fatalf("could not get account: %v", err)
		}
		log.Printf("Account: %s", ar.GetAccount())

		pc := pb.NewPriceServiceClient(conn)

		pr, err := pc.GetPrice(ctx, &pb.GetPriceRequest{PriceId: *product})
		if err != nil {
			log.Fatalf("could not get price: %v", err)
		}
		log.Printf("Price: %s", pr.GetPrice())

		uc := pb.NewUserServiceClient(conn)

		lir, err := uc.AuthenticateUser(ctx, &pb.AuthenticateUserRequest{
			Username: "rjfield777",
			Password: "foo",
		})
		if err != nil {
			log.Fatalf("could not login: %v", err)
		}
		log.Printf("Login: %s", lir.GetUserId())

		// ur, err := uc.GetUser(ctx, &pb.GetUserRequest{UserId: "38ce4dbb-73ff-4af2-9ac8-acb0c8bf8c52"})
		ur, err := uc.GetUser(ctx, &pb.GetUserRequest{UserId: lir.GetUserId()})
		if err != nil {
			log.Fatalf("could not get user: %v", err)
		}
		log.Printf("User: %s", ur.GetUser())
	*/

	// stream, err := pc.StreamPrices(ctx, &pb.StreamPricesRequest{
	// 	PriceIds: []string{"MSFT", "AAPL", "GOOG"}})
	// if err != nil {
	// 	log.Fatalf("could not stream prices: %v", err)
	// }
	// for {
	// 	sp, err := stream.Recv()
	// 	if err == io.EOF {
	// 		break
	// 	}
	// 	if err != nil {
	// 		log.Fatalf("error receiving price stream: %v", err)
	// 	}
	// 	log.Printf("Streamed Price: %s", sp.GetPrice())
	// }
}
