package mysql

import (
	"context"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/AkashGit21/task-ms/lib/persistence"
	"github.com/AkashGit21/task-ms/utils"
	_ "github.com/go-sql-driver/mysql"
)

const DEFAULT_MAX_PAGE_SIZE = 50

type TaskPersistenceLayer struct {
	db *sql.DB
	sync.Mutex
}

type TaskOps interface {
	SaveRecord(persistence.Task) (int64, error)
	UpdateRecord(id string, record persistence.Task) (int64, error)
	FetchRecords(persistence.TaskFilterParams) ([]persistence.Task, error)
	GetRecord(string) (*persistence.Task, error)
	DeactivateRecord(string) error
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

/** Inserts a new TASK record **/
func (tpl *TaskPersistenceLayer) SaveRecord(record persistence.Task) (int64, error) {
	_, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	stmt, err := tpl.db.Prepare("INSERT INTO tasks (id, title, content, stylized_content, status, created_at, modified_at, created_by, modified_by) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)")

	if err != nil {
		return int64(-1), err
	}
	defer stmt.Close()

	tpl.Lock()
	defer tpl.Unlock()
	// Execute the SQL statement to insert the new row
	res, err := stmt.Exec(record.ID, record.Title, record.Content, record.HTMLStylizedContent, record.Status, record.CreatedAt, record.ModifiedAt, record.CreatedBy, record.ModifiedBy)
	if err != nil {
		return int64(-1), err
	}

	return res.RowsAffected()
}

/** Update an existing TASK record **/
func (tpl *TaskPersistenceLayer) UpdateRecord(id string, record persistence.Task) (int64, error) {
	query := "UPDATE tasks SET title = ?, content = ?, stylized_content = ?, status = ?, modified_at = ?, modified_by = ? WHERE id = ? AND discarded = false"

	tpl.Lock()
	defer tpl.Unlock()

	// Execute the query with the given ID
	res, err := tpl.db.Exec(query, record.Title, record.Content, record.HTMLStylizedContent, record.Status, record.ModifiedAt, record.ModifiedBy, id)

	// Check for errors
	if err == sql.ErrNoRows {
		// No record found with the given identifier, not an error.
		return int64(-1), nil
	} else if err != nil {
		return int64(-1), err
	}

	return res.RowsAffected()
}

/** Returns all the active TASK records with cursor-based pagination and status filter in descending order of modifications **/
func (tpl *TaskPersistenceLayer) FetchRecords(params persistence.TaskFilterParams) ([]persistence.Task, error) {
	query := `SELECT id, title, content, stylized_content, status, created_at, created_by, modified_at, modified_by FROM tasks WHERE discarded = FALSE`
	args := []interface{}{}
	var cursorData persistence.CursorData

	if params.Status != nil {
		query += ` AND status = ?`
		args = append(args, *params.Status)
	}

	if params.Cursor != nil && *params.Cursor != "" {
		decodedCursor, err := base64.StdEncoding.DecodeString(*params.Cursor)
		if err != nil {
			return nil, fmt.Errorf("invalid cursor: %w", err)
		}

		if err := json.Unmarshal(decodedCursor, &cursorData); err != nil {
			return nil, fmt.Errorf("invalid cursor data: %w", err)
		}
		query += ` AND (modified_at, id) < (?, ?)`
		args = append(args, cursorData.ModifiedAt, cursorData.ID)
	}

	query += ` ORDER BY modified_at DESC, id DESC`

	if params.Limit > 0 {
		query += ` LIMIT ?`
		args = append(args, params.Limit)
	} else {
		query += fmt.Sprintf(` LIMIT %d`, DEFAULT_MAX_PAGE_SIZE)
	}

	rows, err := tpl.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch tasks: %w", err)
	}
	defer rows.Close()

	tasks := []persistence.Task{}
	for rows.Next() {
		var task persistence.Task
		if err := rows.Scan(
			&task.ID, &task.Title, &task.Content, &task.HTMLStylizedContent,
			&task.Status, &task.CreatedAt, &task.CreatedBy, &task.ModifiedAt,
			&task.ModifiedBy); err != nil {
			return nil, fmt.Errorf("failed to scan task: %w", err)
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

/** Query to fetch the TASK with given ID **/
func (pdb *TaskPersistenceLayer) GetRecord(id string) (*persistence.Task, error) {
	query := "SELECT id, title, content, stylized_content, status, created_at, modified_at, created_by, modified_by FROM tasks WHERE id = ? AND discarded = false"

	// Execute the query with the primary key value
	var task persistence.Task
	err := pdb.db.QueryRow(query, id).Scan(
		&task.ID, &task.Title, &task.Content, &task.HTMLStylizedContent,
		&task.Status, &task.CreatedAt, &task.ModifiedAt, &task.CreatedBy, &task.ModifiedBy,
	)

	// Check for errors
	if err == sql.ErrNoRows {
		// No record found with the given identifier, not an error.
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &task, nil
}

/**: Soft delete TASK record with given ID**/
func (tpl *TaskPersistenceLayer) DeactivateRecord(id string) error {
	query := "UPDATE tasks SET discarded = 1 WHERE id =?"
	tpl.Lock()
	defer tpl.Unlock()

	res, err := tpl.db.Exec(query, id)
	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return errors.New("no rows affected")
	}
	return err
}
