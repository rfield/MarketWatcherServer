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

// GetPrice retrieves the current price for the given price ID.
func (s *PriceServer) GetPrice(ctx context.Context, req *pb.GetPriceRequest) (*pb.GetPriceReply, error) {
	log.Printf("GetPrice() - received: %v", req.GetPriceId())
	quotes, err := api.StockQuote().Symbol(req.GetPriceId()).Get()
	if err != nil {
		log.Printf("GetPrice() - Error fetching stock quote for %s: %v", req.GetPriceId(), err)
		return nil, err
	}
	log.Printf("GetPrice() - Fetched stock quote for %s: %v", req.GetPriceId(), quotes)
	return &pb.GetPriceReply{
		Price: &pb.Price{
			PriceId: req.GetPriceId(),
			Price:   quotes[0].Last,
		},
	}, nil
}

// GetPrices retrieves the current price for the given price ID.
func (s *PriceServer) GetPrices(ctx context.Context, req *pb.GetPricesRequest) (*pb.GetPricesReply, error) {
	log.Printf("GetPrices() - received: %v", req.GetPriceIds())
	quotes, err := api.BulkStockQuotes().Symbols(req.GetPriceIds()).Get()
	if err != nil {
		log.Printf("GetPrice() - Error fetching stock quotes for %s: %v", req.GetPriceIds(), err)
		return nil, err
	}
	log.Printf("GetPrice() - Fetched stock quote for %s: %v", req.GetPriceIds(), quotes)
	p := make([]*pb.Price, 0, len(quotes))
	for _, q := range quotes {
		p = append(p, &pb.Price{
			PriceId: q.Symbol,
			Price:   q.Last,
		})
	}
	return &pb.GetPricesReply{
		Prices: p,
	}, nil
}

// StreamPrices streams real-time price updates for the requested price IDs.
// For now, it streams prices for each requested ticker symbol a small number of
// iterations, sleeping between each batch
func (s *PriceServer) StreamPrices(req *pb.StreamPricesRequest, stream pb.PriceService_StreamPricesServer) error {
	log.Printf("StreamPrices() - streaming prices for: %v", req.GetPriceIds())

	for range 2 {

		// No easy way to set a timeout here
		quotes, err := api.BulkStockQuotes().Symbols(req.GetPriceIds()).Get()
		if err != nil {
			log.Printf("StreamPrices() - Error fetching bulk stock quotes for %v: %v", req.GetPriceIds(), err)
			continue
		}
		log.Printf("StreamPrices() - Fetched bulk stock quotes for %v: %v", req.GetPriceIds(), quotes)

		// Map quotes by symbol for easy lookup
		for _, q := range quotes {
			price := &pb.Price{
				PriceId: q.Symbol,
				Price:   q.Last,
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
