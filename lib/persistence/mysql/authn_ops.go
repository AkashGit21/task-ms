package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"sync"

	"github.com/AkashGit21/task-ms/lib/persistence"
	"github.com/AkashGit21/task-ms/utils"
)

type UserPersistenceLayer struct {
	db *sql.DB
	sync.Mutex
}

type AuthnOps interface {
	FetchActiveRecord(string) (*persistence.User, error)
}

func NewUserPersistenceLayer() (AuthnOps, error) {
	database := utils.GetEnvValue("AUTH_DB_NAME", "auth_db")
	username := utils.GetEnvValue("AUTH_DB_USER", "auth_user")
	password := utils.GetEnvValue("AUTH_DB_PASSWORD", "auth_password")
	host := utils.GetEnvValue("AUTH_DB_HOST", "localhost")
	port := utils.GetEnvValue("AUTH_DB_PORT", "3306")

	// Create a DSN (Data Source Name) for the MySQL connection.
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", username, password, host, port, database)

	// Open a connection to the MySQL database.
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		utils.ErrorLog("Could not open database connection: ", err)
		return nil, err
	}

	// Verify the connection by pinging the database.
	if err := db.PingContext(context.TODO()); err != nil {
		utils.ErrorLog("Could not ping database: ", err)
		return nil, err
	}

	return &UserPersistenceLayer{
		db: db,
	}, nil
}

/** Fetch active user entries with given username **/
func (upl *UserPersistenceLayer) FetchActiveRecord(username string) (*persistence.User, error) {
	query := "SELECT id, username, enc_password, created_at FROM users WHERE username = ? AND discarded = false"

	// Execute the query with the primary key value
	var user persistence.User
	err := upl.db.QueryRow(query, username).Scan(
		&user.ID, &user.Username, &user.EncryptedPassword, &user.CreatedAt,
	)
	utils.ErrorLog("FEtchActiveRecord error: ", err)
	utils.InfoLog("user", user.ID, user.EncryptedPassword, user.CreatedAt)

	// Check for errors
	if err == sql.ErrNoRows {
		utils.InfoLog("ErrNoRows found")
		// No record found with the given identifier, not an error.
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &user, nil
}
