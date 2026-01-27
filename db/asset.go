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
		`select u.username, a.name, ast.ticker, ast.holding_amount
	 from accounts a, users u, assets ast
	where
		a.user_id = u.id and
		ast.user_id = u.id and
		ast.account_id = a.id and
		u.id = $1 and 
		a.id = $2
	order by a.name`

	rows, err := db.Query(sqlStatement, userId, accountId)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var assets []*pb.Asset

	for rows.Next() {
		var a pb.Asset
		err := rows.Scan(&a.UserId, &a.AccountName, &a.Ticker, &a.HoldingAmount)
		if err != nil {
			log.Fatal(err)
		}
		a.Name = "users/" + userId + "/accounts/" + accountId + "/assets/" + a.Ticker
		assets = append(assets, &a)
		log.Printf("listAssetsByUser() - UserID: %s, AccountName: %s, Ticker: %s, HoldingAmount: %f", a.UserId, a.AccountName, a.Ticker, a.HoldingAmount)
	}

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("listAssetsByUser() - returning %d assets", len(assets))
	return assets, nil
}
