package models

import "time"

type User struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
}

type Board struct {
	ID        int       `json:"id"`
	OwnerID   int       `json:"owner_id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
}

type Column struct {
	ID       int     `json:"id"`
	BoardID  int     `json:"board_id"`
	Title    string  `json:"title"`
	Position float64 `json:"position"`
}

type Task struct {
	ID          int        `json:"id"`
	ColumnID    int        `json:"column_id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Position    float64    `json:"position"`
	AssigneeID  *int       `json:"assignee_id"`
	DueDate     *time.Time `json:"due_date"`
	CreatedAt   time.Time  `json:"created_at"`
}

type Subtask struct {
	ID        int    `json:"id"`
	TaskID    int    `json:"task_id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

type Comment struct {
	ID        int       `json:"id"`
	TaskID    int       `json:"task_id"`
	UserID    int       `json:"user_id"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
}

type Attachment struct {
	ID        int       `json:"id"`
	TaskID    int       `json:"task_id"`
	UserID    int       `json:"user_id"`
	Filename  string    `json:"filename"`
	DiskPath  string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
}

type ActivityLog struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	TaskID    *int      `json:"task_id"`
	BoardID   *int      `json:"board_id"`
	Event     string    `json:"event"`
	CreatedAt time.Time `json:"created_at"`
}