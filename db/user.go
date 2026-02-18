package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"rjfield.com/backend/generated/pb"
)

// createUser inserts a new record into the 'users' table and prints the new ID.
func createUser(u *pb.User) (*pb.User, error) {
	log.Printf("createUser() - creating user: %v", u)
	// SQL statement with placeholders ($1, $2, etc. for Postgres)
	sqlStatement := `
		INSERT INTO users (username, password, given_name, family_name, email)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id` // RETURNING id is used to get the auto-generated ID back

	// Use QueryRow for an INSERT statement that returns a single row (the ID in this case)
	var newID string
	err := db.QueryRow(sqlStatement, u.Username, u.PasswordHash, u.GivenName, u.FamilyName, u.Email).Scan(&newID)

	if err != nil {
		return nil, fmt.Errorf("unable to insert row: %v", err)
	}

	log.Printf("createUser() - new record ID is: %s", newID)
	u.Name = "users/" + newID

	return u, nil
}

// ReadUser fetches a single user record by their ID.
func readUser(id string) (*pb.User, error) {
	log.Printf("readUser() - fetching user with ID: %s", id)
	sqlStatement := `SELECT id, username, password, given_name, family_name, email FROM users WHERE id=$1;`

	// QueryRow retrieves at most a single database row.
	row := db.QueryRow(sqlStatement, id)

	var user pb.User
	var userId string

	// Scan copies the column values from the matched row into the variables
	// pointed to by the arguments. Errors (including sql.ErrNoRows) are handled here.
	switch err := row.Scan(&userId, &user.Username, &user.PasswordHash, &user.GivenName, &user.FamilyName, &user.Email); err {
	case sql.ErrNoRows:
		log.Printf("readUser() - no user found with ID: %s", id)
		return nil, fmt.Errorf("user with ID %s not found", id)
	case nil:
		log.Printf("readUser() - found user: %s", userId)
		user.Name = "users/" + userId
		return &user, nil
	default:
		log.Printf("readUser() - query error: %v", err)
		return nil, fmt.Errorf("query error: %v", err)
	}
}

// updateUser: Updates the first name of a user based on their ID
func updateUser(u *pb.User) (*pb.User, error) {
	log.Printf("updateUser() - updating user: %v", u)
	userId := UserIDFromResourceName(u.Name)
	sqlStatement := `
		UPDATE users
		SET given_name = $2
		WHERE id = $1`
	res, err := db.Exec(sqlStatement, userId, u.GivenName)
	if err != nil {
		return nil, err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}
	log.Printf("updateUser() - Updated %d record(s)", count)
	return u, nil
}

// deleteUser: Deletes a user based on their ID
func deleteUser(id string) (string, error) {
	log.Printf("deleteUser() - deleting user with ID: %s", id)
	sqlStatement := `
		DELETE FROM users
		WHERE id = $1`
	res, err := db.Exec(sqlStatement, id)
	if err != nil {
		return "", err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return "", err
	}
	log.Printf("deleteUser() - deleted %d record(s)", count)

	return id, nil
}

// listUsers fetches all users from the database
func listUsers(pageSize int32, pageToken string) ([]*pb.User, error) {
	_ = pageSize // TODO: implement pagination
	_ = pageToken

	log.Printf("listUsers() - fetching all users")
	rows, err := db.Query(`SELECT id, given_name, family_name, email FROM users ORDER BY id`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var users []*pb.User

	for rows.Next() {
		var u pb.User
		var userId string

		err := rows.Scan(&userId, &u.GivenName, &u.FamilyName, &u.Email)
		if err != nil {
			log.Fatal(err)
		}
		u.Name = "users/" + userId
		users = append(users, &u)
		log.Printf("listUsers() - ID: %s, Name: %s %s, Email: %s", userId, u.Username, u.PasswordHash, u.Email)
	}

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("listUsers() - returning %d users", len(users))
	return users, nil
}

// AuthenticateUser fetches a single user record by their ID.
func authenticateUser(username string, password string) (string, error) {
	log.Printf("authenticateUser() - authenticating user: %s", username)
	sqlStatement := `SELECT id FROM users WHERE username = $1 and password = $2;`

	// QueryRow retrieves at most a single database row.
	row := db.QueryRow(sqlStatement, username, password)

	var userId string

	// Scan copies the column values from the matched row into the variables
	// pointed to by the arguments. Errors (including sql.ErrNoRows) are handled here.
	switch err := row.Scan(&userId); err {
	case sql.ErrNoRows:
		log.Printf("authenticateUser() - username %s or password incorrect", username)
		return "", fmt.Errorf("username %s or password incorrect", username)
	case nil:
		log.Printf("authenticateUser() - user %s authenticated successfully", username)
		return userId, nil
	default:
		log.Printf("authenticateUser() - query error: %v", err)
		return "", fmt.Errorf("query error: %v", err)
	}
}

// ListAssetsByUser fetches all assets for a given user ID across their accounts.
// func listAssetsByUser(userId string) ([]*pb.Asset, error) {
// 	log.Printf("listAssetsByUser() - fetching assets for user ID: %s", userId)
// 	sqlStatement :=
// 		`select u.username, a.name, ast.ticker, ast.holding_amount
// 	 from accounts a, users u, assets ast
// 	where
// 		a.user_id = u.id and
// 		ast.user_id = u.id and
// 		ast.account_id = a.id and
// 		u.id = $1
// 	order by a.name`

// 	rows, err := db.Query(sqlStatement, userId)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer rows.Close()

// 	var assets []*pb.Asset

// 	for rows.Next() {
// 		var a pb.Asset
// 		err := rows.Scan(&a.UserId, &a.AccountName, &a.Ticker, &a.HoldingAmount)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		assets = append(assets, &a)
// 		log.Printf("listAssetsByUser() - UserID: %s, AccountName: %s, Ticker: %s, HoldingAmount: %f", a.UserId, a.AccountName, a.Ticker, a.HoldingAmount)
// 	}

// 	err = rows.Err()
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	log.Printf("listAssetsByUser() - returning %d assets", len(assets))
// 	return assets, nil
// }
