-- name: ListTodo :many
SELECT * FROM "todo" ORDER BY id DESC LIMIT $1 OFFSET $2;

-- name: CountTodo :one
SELECT count(1) FROM "todo";

-- name: CreateTodo :one
INSERT INTO "todo" (title) VALUES ($1) RETURNING id, title, completed;

-- name: UpdateTodo :one
UPDATE "todo" SET completed = $1 WHERE id = $2 RETURNING id, title, completed;

-- name: GetTodoById :one
SELECT * FROM "todo"
WHERE id = $1;

-- name: PatchTodo :one
UPDATE todo
SET title = COALESCE($2, title),
    completed = COALESCE($3, completed)
WHERE id = $1 RETURNING id, title, completed;
