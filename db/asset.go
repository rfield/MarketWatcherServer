package db

import (
	"log"

	_ "github.com/lib/pq"
	"rjfield.com/backend/generated/pb"
)

// ListAssets fetches all assets for a given user ID and account
func listAssets(userId, accountId string) ([]*pb.Asset, error) {
	log.Printf("listAssets() - fetching assets for user ID: %s, account ID: %s", userId, accountId)
	sqlStatement :=
		`select id, user_id, account_id, ticker, holding_amount
	 from assets
	where
		user_id = $1 and
		account_id = $2
	order by
		user_id, account_id`

	rows, err := db.Query(sqlStatement, userId, accountId)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var assets []*pb.Asset

	for rows.Next() {
		var a pb.Asset
		var userId string
		var accountId string
		var assetId string
		err := rows.Scan(&assetId, &userId, &accountId, &a.Ticker, &a.HoldingAmount)
		if err != nil {
			log.Fatal(err)
		}
		a.Name = "users/" + userId + "/accounts/" + accountId + "/assets/" + assetId
		assets = append(assets, &a)
	}

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("listAssetsByUser() - returning %d assets", len(assets))
	return assets, nil
}
