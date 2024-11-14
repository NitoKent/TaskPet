package user

import (
	"database/sql"
	"fmt"
	"time"
)

type UserStore struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *UserStore {
	return &UserStore{db: db}
}

type UserStatusResponse struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Balance    int    `json:"balance"`
	Status     string `json:"status"`
	ReferrerID *int   `json:"referrer_id"`
}

type Task struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int    `json:"price"`
	IsActive    bool   `json:"is_active"`
}

type UserTask struct {
	ID             int       `json:"id"`
	UserID         string    `json:"user_id"`
	TaskID         int       `json:"task_id"`
	Completed      bool      `json:"completed"`
	CompletionDate time.Time `json:"completion_date"`
}

func (store *UserStore) GetUserStatus(userID string) (UserStatusResponse, error) {
	var user UserStatusResponse
	err := store.db.QueryRow("SELECT id, name, balance, status FROM users WHERE id = $1", userID).
		Scan(&user.ID, &user.Name, &user.Balance, &user.Status)
	if err != nil {
		return UserStatusResponse{}, err
	}
	return user, nil
}

func (store *UserStore) GetTopUsersByBalance(limit int) ([]UserStatusResponse, error) {
	var users []UserStatusResponse

	rows, err := store.db.Query("SELECT id, name, balance, status FROM users ORDER BY balance DESC LIMIT $1", limit)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch top users: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var user UserStatusResponse
		if err := rows.Scan(&user.ID, &user.Name, &user.Balance, &user.Status); err != nil {
			return nil, fmt.Errorf("failed to scan user row: %w", err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return users, nil
}

func (s *UserStore) GetTaskByID(taskID int) (*Task, error) {
	var task Task
	query := "SELECT id, name, description, price, is_active FROM tasks WHERE id = $1"
	err := s.db.QueryRow(query, taskID).Scan(&task.ID, &task.Name, &task.Description, &task.Price, &task.IsActive)
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (s *UserStore) GetUserTask(userID string, taskID int) (*UserTask, error) {
	var userTask UserTask
	query := "SELECT id, user_id, task_id, completed, completion_date FROM user_tasks WHERE user_id = $1 AND task_id = $2"
	err := s.db.QueryRow(query, userID, taskID).Scan(&userTask.ID, &userTask.UserID, &userTask.TaskID, &userTask.Completed, &userTask.CompletionDate)
	if err != nil {
		return nil, err
	}
	return &userTask, nil
}

func (s *UserStore) UpdateUserBalance(userID string, price int) error {
	query := `UPDATE users SET balance = balance + $1 WHERE id = $2`
	_, err := s.db.Exec(query, price, userID)
	return err
}

func (s *UserStore) CreateUserTask(userID string, taskID int) error {
	query := `INSERT INTO user_tasks (user_id, task_id) VALUES ($1, $2)`
	_, err := s.db.Exec(query, userID, taskID)
	return err
}

func (s *UserStore) CompleteUserTask(userID string, taskID int) error {
	query := `UPDATE user_tasks SET completed = TRUE, completion_date = NOW() WHERE user_id = $1 AND task_id = $2`
	_, err := s.db.Exec(query, userID, taskID)
	return err
}

func (s *UserStore) GetUserByID(userID string) (*UserStatusResponse, error) {
	var user UserStatusResponse
	query := "SELECT id, name, balance, status, referrer_id FROM users WHERE id = $1"
	err := s.db.QueryRow(query, userID).Scan(&user.ID, &user.Name, &user.Balance, &user.Status, &user.ReferrerID)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserStore) SetReferrerID(userID string, referrerID int) error {
	query := "UPDATE users SET referrer_id = $1 WHERE id = $2"
	_, err := s.db.Exec(query, referrerID, userID)
	return err
}
