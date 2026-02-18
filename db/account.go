package db

import (
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"rjfield.com/backend/generated/pb"
)

// listAccounts fetches all accounts for a given user ID.
func listAccounts(user_id string) ([]*pb.Account, error) {
	log.Printf("listAccounts() - listing accounts for user ID: %s", user_id)
	// SQL statement with placeholders ($1, $2, etc. for Postgres)
	sqlStatement := `
		SELECT id, name FROM accounts WHERE user_id = $1`
	rows, err := db.Query(sqlStatement, user_id)
	if err != nil {
		return nil, fmt.Errorf("unable to query rows: %v", err)
	}
	defer rows.Close()

	var accounts []*pb.Account

	for rows.Next() {
		var account pb.Account
		var accountId string
		err := rows.Scan(&accountId, &account.AccountName)
		if err != nil {
			return nil, fmt.Errorf("unable to scan row: %v", err)
		}
		account.Name = "users/" + user_id + "/accounts/" + accountId
		accounts = append(accounts, &account)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %v", err)
	}

	log.Printf("listAccounts() - found %d accounts for user ID: %s", len(accounts), user_id)
	return accounts, nil
}
