package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	_ "modernc.org/sqlite"
)

type SQLiteTodoRepo struct {
	db *sql.DB
}

var migrations []string = []string{
	`CREATE TABLE IF NOT EXISTS todos (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		completed_at TEXT
	)`,
}

func NewSQLiteTodoRepo(ctx context.Context, databaseName string) (*SQLiteTodoRepo, error) {
	db, err := sql.Open("sqlite", databaseName)
	if err != nil {
		return nil, fmt.Errorf("creating database connection: %v", err)
	}

	for num, migration := range migrations {
		if _, err := db.ExecContext(ctx, migration); err != nil {
			return nil, fmt.Errorf("executing migration %d: %v", num, err)
		}

		log.Printf("Executed migration %d: %s", num, migration)
	}

	return &SQLiteTodoRepo{db}, nil
}

func (r *SQLiteTodoRepo) CreateTodo(ctx context.Context, title string) (Todo, error) {
	query := `INSERT INTO todos (title) VALUES (?) RETURNING id`

	var id int
	if err := r.db.QueryRowContext(ctx, query, title).Scan(&id); err != nil {
		return Todo{}, fmt.Errorf("creating todo: %v", err)
	}

	return Todo{ID: id, Title: title}, nil
}

func (r *SQLiteTodoRepo) GetComplete(ctx context.Context) ([]Todo, error) {
	query := `SELECT id, title, completed_at FROM todos WHERE completed_at IS NOT NULL ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("querying for complete todos: %v", err)
	}

	defer rows.Close()

	todos := make([]Todo, 0)
	for rows.Next() {
		var id int
		var title string
		var rawCompletedAt string

		if err := rows.Scan(&id, &title, &rawCompletedAt); err != nil {
			return nil, fmt.Errorf("reading results: %v", err)
		}

		completedAt, err := time.Parse("2006-01-02 15:04:05", rawCompletedAt)
		if err != nil {
			return nil, fmt.Errorf("parsing malformed completion time %q for todo %d: %v", rawCompletedAt, id, err)
		}

		todos = append(todos, Todo{ID: id, Title: title, CompletedAt: &completedAt})
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("reading results: %v", err)
	}

	return todos, nil
}

func (r *SQLiteTodoRepo) GetIncomplete(ctx context.Context) ([]Todo, error) {
	query := `SELECT id, title FROM todos WHERE completed_at IS NULL ORDER BY id`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("querying for incomplete todos: %v", err)
	}

	defer rows.Close()

	todos := make([]Todo, 0)
	for rows.Next() {
		var id int
		var title string

		if err := rows.Scan(&id, &title); err != nil {
			return nil, fmt.Errorf("reading results: %v", err)
		}

		todos = append(todos, Todo{ID: id, Title: title})
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("reading results: %v", err)
	}

	return todos, nil
}

func (r *SQLiteTodoRepo) SetComplete(ctx context.Context, id int) (Todo, error) {
	query := `UPDATE todos SET completed_at = datetime('now') WHERE id = ? RETURNING title, completed_at`

	var title string
	var rawCompletedAt string

	if err := r.db.QueryRowContext(ctx, query, id).Scan(&title, &rawCompletedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Todo{}, fmt.Errorf("%w: no todo with ID %d", ErrTodoNotFound, id)
		}

		return Todo{}, fmt.Errorf("marking todo %d as complete: %v", id, err)
	}

	completedAt, err := time.Parse("2006-01-02 15:04:05", rawCompletedAt)
	if err != nil {
		return Todo{}, fmt.Errorf("parsing malformed completion time %q for todo %d: %v", rawCompletedAt, id, err)
	}

	return Todo{id, title, &completedAt}, nil
}

func (r *SQLiteTodoRepo) SetIncomplete(ctx context.Context, id int) (Todo, error) {
	query := `UPDATE todos SET completed_at = NULL WHERE id = ? RETURNING title`

	var title string
	if err := r.db.QueryRowContext(ctx, query, id).Scan(&title); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Todo{}, fmt.Errorf("%w: no todo with ID %d", ErrTodoNotFound, id)
		}

		return Todo{}, fmt.Errorf("marking todo %d as incomplete: %v", id, err)
	}

	return Todo{ID: id, Title: title}, nil
}
