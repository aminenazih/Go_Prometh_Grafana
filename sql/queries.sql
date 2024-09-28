-- name: CreateTask :exec
INSERT INTO tasks (type, value, state, creation_time, last_update_time)
VALUES ($1, $2, $3, $4, $5);

-- name: UpdateTaskState :exec
UPDATE tasks
SET state = $1, last_update_time = $2
WHERE id = $3;

-- name: GetTaskById :one
SELECT id, type, value, state, creation_time, last_update_time
FROM tasks
WHERE id = $1;
