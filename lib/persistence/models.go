package persistence

import "time"

const (
	STATUS_TODO = iota
	STATUS_INPROGRESS
	STATUS_BACKLOG
	STATUS_ON_HOLD
	STATUS_COMPLETED
)

var (
	TaskStatus_name = map[int8]string{
		STATUS_TODO:       "todo",
		STATUS_INPROGRESS: "in-progress",
		STATUS_BACKLOG:    "backlog",
		STATUS_ON_HOLD:    "on-hold",
		STATUS_COMPLETED:  "completed",
	}

	TaskStatus_value = map[string]int8{
		"todo":        STATUS_TODO,
		"in-progress": STATUS_INPROGRESS,
		"backlog":     STATUS_BACKLOG,
		"on-hold":     STATUS_ON_HOLD,
		"completed":   STATUS_COMPLETED,
	}
)

type TaskStatus int8

func (x TaskStatus) String() string {
	if val, ok := TaskStatus_name[int8(x)]; ok {
		return val
	}
	return "inactive"
}

type Task struct {
	ID                  int64      `db:"id"`
	Title               string     `db:"title" json:"title"`
	Content             string     `db:"content" json:"content"`
	HTMLStylizedContent string     `db:"stylizedContent" json:"stylizedContent"`
	Status              TaskStatus `db:"status" json:"-"`
	Discarded           bool       `db:"discarded" json:"-"`

	CreatedAt  time.Time `db:"created_at" json:"created_at"`
	CreatedBy  string    `db:"created_by" json:"created_by"`
	ModifiedAt time.Time `db:"modified_at" json:"modified_at"`
	ModifiedBy string    `db:"modified_by" json:"modified_by"`
}
