package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"rjfield.com/backend/config"
	"rjfield.com/backend/generated/pb"
)

var db *sql.DB

func init() {

	var err error

	var configData *config.Config

	configData = config.GetConfig()

	// Connection string
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		configData.Database.Host, configData.Database.Port, configData.Database.User, configData.Database.Password, configData.Database.DBName)

	log.Printf("init() - connecting to database at %s:%d with username %s, password %s and database %s",
		configData.Database.Host, configData.Database.Port, configData.Database.User, "****", configData.Database.DBName)

	// Open a database connection
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}
	// defer db.Close() // Ensure the connection pool is closed when main exits

	// Verify the connection is active
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("init() - Successfully connected to the database!")
}

// insertRow inserts a new record into the 'users' table and prints the new ID.
func CreateUser(u *pb.User) (*pb.User, error) {
	log.Printf("CreateUser() - creating user: %v", u)
	// SQL statement with placeholders ($1, $2, etc. for Postgres)
	sqlStatement := `
		INSERT INTO users (username, password, given_name, family_name, email)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id` // RETURNING id is used to get the auto-generated ID back

	// Use QueryRow for an INSERT statement that returns a single row (the ID in this case)
	var newID string
	err := db.QueryRow(sqlStatement, u.Credentials.Username, u.Credentials.PasswordHash, u.GivenName, u.FamilyName, u.Email).Scan(&newID)

	if err != nil {
		return nil, fmt.Errorf("unable to insert row: %v", err)
	}

	log.Printf("CreateUser() - new record ID is: %s", newID)
	u.UserId = newID
	return u, nil
}

// ReadUser fetches a single user record by their ID.
func ReadUser(id string) (*pb.User, error) {
	log.Printf("ReadUser() - fetching user with ID: %s", id)
	sqlStatement := `SELECT id, username, password, given_name, family_name, email FROM users WHERE id=$1;`

	// QueryRow retrieves at most a single database row.
	row := db.QueryRow(sqlStatement, id)

	var user pb.User
	var creds pb.Credentials
	user.Credentials = &creds

	// Scan copies the column values from the matched row into the variables
	// pointed to by the arguments. Errors (including sql.ErrNoRows) are handled here.
	switch err := row.Scan(&user.UserId, &user.Credentials.Username, &user.Credentials.PasswordHash, &user.GivenName, &user.FamilyName, &user.Email); err {
	case sql.ErrNoRows:
		log.Printf("ReadUser() - no user found with ID: %s", id)
		return nil, fmt.Errorf("user with ID %s not found", id)
	case nil:
		log.Printf("ReadUser() - found user: %s", user.GetUserId())
		return &user, nil
	default:
		log.Printf("ReadUser() - query error: %v", err)
		return nil, fmt.Errorf("query error: %v", err)
	}
}

// ListUsers fetches all users from the database
func ListUsers() ([]*pb.User, error) {
	log.Printf("ListUsers() - fetching all users")
	rows, err := db.Query(`SELECT id, given_name, family_name, email FROM users ORDER BY id`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var users []*pb.User

	for rows.Next() {
		var u pb.User
		var creds pb.Credentials
		u.Credentials = &creds

		err := rows.Scan(&u.UserId, &u.Credentials.Username, &u.Credentials.PasswordHash, &u.Email)
		if err != nil {
			log.Fatal(err)
		}
		users = append(users, &u)
		log.Printf("ListUsers() - ID: %s, Name: %s %s, Email: %s", u.UserId, u.Credentials.Username, u.Credentials.PasswordHash, u.Email)
	}

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("ListUsers() - returning %d users", len(users))
	return users, nil
}

// UpdateUser: Updates the first name of a user based on their ID
func UpdateUser(u *pb.User) (*pb.User, error) {
	log.Printf("UpdateUser() - updating user: %v", u)
	sqlStatement := `
		UPDATE users
		SET given_name = $2
		WHERE id = $1`
	res, err := db.Exec(sqlStatement, u.UserId, u.GivenName)
	if err != nil {
		return nil, err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}
	log.Printf("UpdateUser() - Updated %d record(s)", count)
	return u, nil
}

// DeleteUser: Deletes a user based on their ID
func DeleteUser(id string) (string, error) {
	log.Printf("DeleteUser() - deleting user with ID: %s", id)
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
	log.Printf("DeleteUser() - deleted %d record(s)", count)

	return id, nil
}

// AuthenticateUser fetches a single user record by their ID.
func AuthenticateUser(username string, password string) (string, error) {
	log.Printf("AuthenticateUser() - authenticating user: %s", username)
	sqlStatement := `SELECT id FROM users WHERE username = $1 and password = $2;`

	// QueryRow retrieves at most a single database row.
	row := db.QueryRow(sqlStatement, username, password)

	var userId string

	// Scan copies the column values from the matched row into the variables
	// pointed to by the arguments. Errors (including sql.ErrNoRows) are handled here.
	switch err := row.Scan(&userId); err {
	case sql.ErrNoRows:
		log.Printf("AuthenticateUser() - username %s or password incorrect", username)
		return "", fmt.Errorf("username %s or password incorrect", username)
	case nil:
		log.Printf("AuthenticateUser() - user %s authenticated successfully", username)
		return userId, nil
	default:
		log.Printf("AuthenticateUser() - query error: %v", err)
		return "", fmt.Errorf("query error: %v", err)
	}
}

func ListAssetsByUser(userId string) ([]*pb.Asset, error) {
	log.Printf("ListAssetsByUser() - fetching assets for user ID: %s", userId)
	sqlStatement :=
		`select u.username, a.name, ast.ticker, ast.holding_amount
	 from accounts a, users u, assets ast
	where
		a.user_id = u.id and
		ast.user_id = u.id and
		ast.account_id = a.id and
		u.id = $1
	order by a.name`

	rows, err := db.Query(sqlStatement, userId)
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
		assets = append(assets, &a)
		log.Printf("ListAssetsByUser() - UserID: %s, AccountName: %s, Ticker: %s, HoldingAmount: %f", a.UserId, a.AccountName, a.Ticker, a.HoldingAmount)
	}

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("ListAssetsByUser() - returning %d assets", len(assets))
	return assets, nil
}
