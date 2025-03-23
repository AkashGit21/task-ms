package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"sync"

	"github.com/AkashGit21/task-ms/lib/persistence"
	"github.com/AkashGit21/task-ms/utils"
	_ "github.com/go-sql-driver/mysql"
)

type TaskPersistenceLayer struct {
	db *sql.DB
	sync.Mutex
}

type TaskOps interface {
	Exists(id int64) (bool, error)
	SaveRecord(record persistence.Task) (int64, error)
	UpdateRecord(id int64, record persistence.Task) error
	FetchRecords() ([]persistence.Task, error)
	GetRecord(id int64) (*persistence.Task, error)
	DeactivateRecord(id int64) error
}

func NewTaskPersistenceLayer() (TaskOps, error) {
	database := utils.GetEnvValue("TASK_DB_NAME", "task_db")
	username := utils.GetEnvValue("TASK_DB_USER", "task_user")
	password := utils.GetEnvValue("TASK_DB_PASSWORD", "task_password")
	host := utils.GetEnvValue("TASK_DB_HOST", "localhost")
	port := utils.GetEnvValue("TASK_DB_PORT", "3306")

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

	return &TaskPersistenceLayer{
		db: db,
	}, nil
}

/** TODO: Implement functionality to check if record exists in DB **/
func (tpl *TaskPersistenceLayer) Exists(id int64) (bool, error) {
	return false, nil
}

/** TODO: Inserts a new TASK record **/
func (tpl *TaskPersistenceLayer) SaveRecord(record persistence.Task) (int64, error) {
	return -1, nil
}

/** TODO: Update an existing TASK record **/
func (tpl *TaskPersistenceLayer) UpdateRecord(id int64, record persistence.Task) error {
	return nil
}

/** TODO: Returns all the active TASK records with pagination and filters in descending order of modifications **/
func (tpl *TaskPersistenceLayer) FetchRecords() ([]persistence.Task, error) {
	return nil, nil
}

/** TODO: Query to fetch the TASK with given ID **/
func (pdb *TaskPersistenceLayer) GetRecord(id int64) (*persistence.Task, error) {
	return nil, nil
}

/**: TODO: Soft delete TASK record with given ID**/
func (tpl *TaskPersistenceLayer) DeactivateRecord(id int64) error {
	return nil
}
