// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0
// source: queries.sql

package db

import (
	"context"
	"time"
)

const createTask = `-- name: CreateTask :exec
INSERT INTO tasks (type, value, state, creation_time, last_update_time)
VALUES ($1, $2, $3, $4, $5)
`

type CreateTaskParams struct {
	Type           int32
	Value          int32
	State          string
	CreationTime   time.Time
	LastUpdateTime time.Time
}

func (q *Queries) CreateTask(ctx context.Context, arg CreateTaskParams) error {
	_, err := q.db.ExecContext(ctx, createTask,
		arg.Type,
		arg.Value,
		arg.State,
		arg.CreationTime,
		arg.LastUpdateTime,
	)
	return err
}

const getTaskById = `-- name: GetTaskById :one
SELECT id, type, value, state, creation_time, last_update_time
FROM tasks
WHERE id = $1
`

func (q *Queries) GetTaskById(ctx context.Context, id int32) (Task, error) {
	row := q.db.QueryRowContext(ctx, getTaskById, id)
	var i Task
	err := row.Scan(
		&i.ID,
		&i.Type,
		&i.Value,
		&i.State,
		&i.CreationTime,
		&i.LastUpdateTime,
	)
	return i, err
}

const updateTaskState = `-- name: UpdateTaskState :exec
UPDATE tasks
SET state = $1, last_update_time = $2
WHERE id = $3
`

type UpdateTaskStateParams struct {
	State          string
	LastUpdateTime time.Time
	ID             int32
}

func (q *Queries) UpdateTaskState(ctx context.Context, arg UpdateTaskStateParams) error {
	_, err := q.db.ExecContext(ctx, updateTaskState, arg.State, arg.LastUpdateTime, arg.ID)
	return err
}
