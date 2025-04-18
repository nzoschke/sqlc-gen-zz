-- name: ContactCreate :one
INSERT INTO
  contacts (blob, info, name)
VALUES
  (?, ?, ?)
RETURNING
  *;

-- name: ContactRead :one
SELECT * FROM contacts WHERE id = ? LIMIT 1;

-- name: ContactCount :one
SELECT COUNT(*) FROM contacts;

-- name: ContactList :many
SELECT
  *
FROM
  contacts
LIMIT
  ?;

-- name: ContactListNames :many
SELECT
  name
FROM
  contacts
LIMIT
  ?;

-- name: ContactUpdate :exec
UPDATE
  contacts
SET
  created_at = ?,
  name = ?
WHERE
  id = ?;

-- name: ContactDelete :exec
DELETE FROM
  contacts
WHERE
  id = ?;

-- name: ContactDeleteAll :exec
DELETE FROM
  contacts;

-- name: ContactCreateJSONB :one
INSERT INTO
  contacts (blob, info, name)
VALUES
  (JSONB(:blob), :info, :name) -- JSONB requires functional named param
RETURNING
  JSON(blob),  -- and requires functional return param in position 1
  *;

-- name: ContactReadJSONB :one
SELECT
  JSON(blob) AS blob
FROM
  contacts
WHERE
  id = ?
LIMIT
  1;
