-- name: ContactCreate :one
INSERT INTO
  contacts (blob, name)
VALUES
  (?, ?)
RETURNING
  *;

-- name: ContactList :many
SELECT * FROM contacts LIMIT ?;