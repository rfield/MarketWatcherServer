package db

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

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

// User CRUD operations
func CreateUser(u *pb.User) (*pb.User, error) {
	return createUser(u)
}
func ReadUser(id string) (*pb.User, error) {
	return readUser(id)
}
func UpdateUser(u *pb.User) (*pb.User, error) {
	return updateUser(u)
}
func DeleteUser(id string) (string, error) {
	return deleteUser(id)
}
func ListUsers(pageSize int32, pageToken string) ([]*pb.User, error) {
	return listUsers(pageSize, pageToken)
}

// Other User functions
func AuthenticateUser(username string, password string) (string, error) {
	return authenticateUser(username, password)
}
func ListAssetsByUser(userId string) ([]*pb.Asset, error) {
	return listAssetsByUser(userId)
}

// Account CRUD operations
func ListAccounts(userId string) ([]*pb.Account, error) {
	return listAccounts(userId)
}

// Asset CRUD operations
func ListAssets(userId, accountId string) ([]*pb.Asset, error) {
	return listAssets(userId, accountId)
}

// Utility functions

func UserIDFromResourceName(resourceName string) string {
	resourceMap := resourceNameToMap(resourceName)
	return resourceMap["users"]
}

func AccountIDFromResourceName(resourceName string) string {
	resourceMap := resourceNameToMap(resourceName)
	return resourceMap["accounts"]
}

func resourceNameToMap(resourceName string) map[string]string {

	// 1. Split the resource name by "/"; i.e. "publishers/123/books/345"
	parts := strings.Split(resourceName, "/")

	// 2. Create a map to hold the key-value pairs
	resourceMap := make(map[string]string)

	// 3. Iterate through parts and fill the map
	// Assuming an even number of parts (key/value pairs)
	for i := 0; i < len(parts); i += 2 {
		if i+1 < len(parts) {
			resourceMap[parts[i]] = parts[i+1]
		}
	}

	// Output: map[books:345 publishers:123]
	fmt.Println(resourceMap)
	return resourceMap
}
