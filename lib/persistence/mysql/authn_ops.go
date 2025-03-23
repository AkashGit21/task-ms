package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"sync"

	"github.com/AkashGit21/task-ms/utils"
)

type UserPersistenceLayer struct {
	db *sql.DB
	sync.Mutex
}

type AuthnOps interface {
	FetchActiveRecord(username, password string) (bool, error)
}

func NewUserPersistenceLayer() (AuthnOps, error) {
	database := utils.GetEnvValue("AUTH_DB_NAME", "task_db")
	username := utils.GetEnvValue("AUTH_DB_USER", "task_user")
	password := utils.GetEnvValue("AUTH_DB_PASSWORD", "task_password")
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

/** TODO: Fetch active user entries with given username **/
func (upl *UserPersistenceLayer) FetchActiveRecord(username, password string) (bool, error) {
	return true, nil
}
