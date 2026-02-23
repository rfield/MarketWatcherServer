package price

import (
	"log"
	"math/rand"
	"strings"

	"github.com/MarketDataApp/sdk-go/models"
)

// Stub implementations for testing without actual API calls
// Use these calls in place of the real API calls in price_api.go for testing purposes
// especially if you're going to be doing a lot of API calls and want to avoid hitting
// rate limits or incurring costs. You can also use these stubs to test the functionality
// of your gRPC server without relying on the external API.
//
// TODO - eventually, you may want to implement a more sophisticated mock that simulates
// different scenarios (e.g., errors, varying price changes, etc.) for more comprehensive testing.
// And you may want to use a command line flag or environment variable to switch between
// the real API and the stub for easier testing and development.
func stubGetStockQuotes(id string) ([]models.StockQuote, error) {
	change := 0.0
	return []models.StockQuote{
		{
			Symbol: id,
			Last:   100.0 * rand.Float64(),
			Change: &change,
		},
	}, nil
}

func stubGetBulkStockQuotes(ids []string) ([]models.StockQuote, error) {
	change := 0.0
	log.Printf("stubGetBulkStockQuotes() - received: %v", ids)
	//TO DO - break up the ids seperated by comma
	s := ids[0]
	log.Printf("stubGetBulkStockQuotes() - received: %v", s)
	// l := s[1 : len(s)-1]
	l := s
	log.Printf("stubGetBulkStockQuotes() - received: %v", l)
	new_ids := strings.Split(l, ",")
	log.Printf("stubGetBulkStockQuotes() - received: %v", new_ids)

	p := make([]models.StockQuote, 0, len(new_ids))
	for _, id := range new_ids {
		p = append(p, models.StockQuote{
			Symbol: id,
			Last:   100.0 * rand.Float64(),
			Change: &change,
		})
	}
	return p, nil
}
