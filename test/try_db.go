package main

import (
	// "database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"rjfield.com/backend/generated/pb"
	db "rjfield.com/backend/user"
)

// const (
// 	host     = "localhost"
// 	port     = 5432
// 	user     = "postgres"
// 	password = "admin"
// 	dbname   = "postgres"

func main() {
	// Connection string
	// psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
	// 	host, port, user, password, dbname)

	// // Open a database connection
	// db, err := sql.Open("postgres", psqlInfo)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer db.Close() // Ensure the connection pool is closed when main exits

	// // Verify the connection is active
	// err = db.Ping()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println("Successfully connected to the database!")

	user := &pb.User{
		Credentials: &pb.Credentials{
			Username:     "bob",
			PasswordHash: "securepassword",
		},
		GivenName:  "Bob",
		FamilyName: "Wonderland",
		Email:      "bob.wonderland@example.com",
	}

	createdUser, err := db.CreateUser(user)
	if err != nil {
		log.Fatalf("Error creating user: %v", err)
	}
	fmt.Printf("Created user: %+v\n", createdUser)

	// // --- Insert a row ---
	// insertRow(db, "joe1", "bar", "Joe", "Schmoe", "jschmoe@example.com")

	// // --- Fetch a user by ID ---
	// user, err := getUserByID(db, "b0159a28-5a03-403f-ba2f-192b41a32d9a")
	// if err != nil {
	// 	log.Fatalf("Error fetching user: %v", err)
	// }
	// fmt.Printf("Fetched user: %+v\n", user)

	// // --- List all users ---
	// listUsers(db)

	// --- Update a user ---
	// updateUser(db, "a55c4356-66b7-41dc-963b-2cb0c93a8f8b", "Jonathan")

	// --- Delete a user ---
	// deleteUser(db, "a55c4356-66b7-41dc-963b-2cb0c93a8f8b")
}

/*
// insertRow inserts a new record into the 'users' table and prints the new ID.
func insertRow(db *sql.DB, un, pw, gn, fn, em string) {
	// SQL statement with placeholders ($1, $2, etc. for Postgres)
	sqlStatement := `
		INSERT INTO users (username, password, given_name, family_name, email)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id` // RETURNING id is used to get the auto-generated ID back

	// Use QueryRow for an INSERT statement that returns a single row (the ID in this case)
	var newID string
	err := db.QueryRow(sqlStatement, un, pw, gn, fn, em).Scan(&newID)

	if err != nil {
		log.Fatalf("Unable to insert row: %v", err)
	}

	fmt.Printf("New record ID is: %s\n", newID)
}

// getUserByID fetches a single user record by their ID.
func getUserByID(db *sql.DB, id string) (*pb.User, error) {
	sqlStatement := `SELECT id, username, password, email FROM users WHERE id=$1;`

	// QueryRow retrieves at most a single database row.
	row := db.QueryRow(sqlStatement, id)

	var user pb.User
	var creds pb.Credentials
	user.Credentials = &creds

	// Scan copies the column values from the matched row into the variables
	// pointed to by the arguments. Errors (including sql.ErrNoRows) are handled here.
	switch err := row.Scan(&user.UserId, &user.Credentials.Username, &user.Credentials.PasswordHash, &user.Email); err {
	case sql.ErrNoRows:
		return nil, fmt.Errorf("user with ID %s not found", id)
	case nil:
		return &user, nil
	default:
		return nil, fmt.Errorf("query error: %v", err)
	}
}

// READ Operation: Fetches all users from the database and prints them
func listUsers(db *sql.DB) {
	rows, err := db.Query(`SELECT id, given_name, family_name, email FROM users ORDER BY id`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var u pb.User
		var creds pb.Credentials
		u.Credentials = &creds

		err := rows.Scan(&u.UserId, &u.Credentials.Username, &u.Credentials.PasswordHash, &u.Email)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("ID: %s, Name: %s %s, Email: %s\n", u.UserId, u.Credentials.Username, u.Credentials.PasswordHash, u.Email)
	}

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
}

// UPDATE Operation: Updates the first name of a user based on their ID
func updateUser(db *sql.DB, id string, newFirstName string) {
	sqlStatement := `
		UPDATE users
		SET given_name = $2
		WHERE id = $1`
	res, err := db.Exec(sqlStatement, id, newFirstName)
	if err != nil {
		log.Fatal(err)
	}
	count, err := res.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Updated %d record(s)\n", count)
}

// DELETE Operation: Deletes a user based on their ID
func deleteUser(db *sql.DB, id string) {
	sqlStatement := `
		DELETE FROM users
		WHERE id = $1`
	res, err := db.Exec(sqlStatement, id)
	if err != nil {
		log.Fatal(err)
	}
	count, err := res.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Deleted %d record(s)\n", count)
}
*/
