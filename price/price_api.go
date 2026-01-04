package price

import (
	"context"
	"log"
	"time"

	api "github.com/MarketDataApp/sdk-go"
	"rjfield.com/backend/generated/pb"
)

type PriceServer struct {
	pb.UnimplementedPriceServiceServer
}

func (s *PriceServer) GetPrice(ctx context.Context, req *pb.GetPriceRequest) (*pb.GetPriceReply, error) {
	log.Printf("GetPrice() - received: %v", req.GetPriceId())
	return &pb.GetPriceReply{
		Price: &pb.Price{
			Price: 100.50,
		},
	}, nil
}

// TODO use https://www.marketdata.app/sdk/go/golang-stock-api/real-time-stock-api/ for real data

func (s *PriceServer) StreamPrices(req *pb.StreamPricesRequest, stream pb.PriceService_StreamPricesServer) error {
	log.Printf("StreamPrices() - streaming prices for: %v", req.GetPriceIds())

	// For now, the stream is 3 prices in 30 seconds prices for each requested ID
	// In the real implementation, this would be a continuous stream of updates
	// Use a Context to set a timeout on the API calls
	// ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	// defer cancel()

	// client, err := api.GetClient()
	// if err != nil {
	// 	log.Printf("Error getting client: %v", err)
	// 	return err
	// }
	// client.Timeout(1)
	// client.StockQuote().WithContext(context.Background())

	for range 3 {
		for _, id := range req.GetPriceIds() {
			quotes, err := api.StockQuote().Symbol(id).Get()
			if err != nil {
				log.Printf("StreamPrices() - Error fetching stock quote for %s: %v", id, err)
				continue
			}
			log.Printf("StreamPrices() - Fetched stock quote for %s: %v", id, quotes)
			price := &pb.Price{
				PriceId: id,
				// Price:   rand.Float64() * 200,
				Price: quotes[0].Last,
			}
			if err := stream.Send(&pb.StreamPricesReply{Price: price}); err != nil {
				return err
			}
		}
		// Wait before sending the next batch of prices
		// This avoids flooding the client, overwhelming the provider
		// and allows the prices time to change, based on real market movements
		time.Sleep(10 * time.Second)
	}

	log.Printf("StreamPrices() - completed streaming prices for: %v", req.GetPriceIds())
	return nil
}
