-- name: ContactCreate :one
INSERT INTO
  contacts (blob, name)
VALUES
  (?, ?)
RETURNING
  *;

-- name: ContactList :many
SELECT
  *
FROM
  contacts
LIMIT
  ?;

-- name: ContactUpdate :exec
UPDATE
  contacts
SET
  name = ?
WHERE
  id = ?;

-- name: ContactDelete :exec
DELETE FROM
  contacts
WHERE
  id = ?;

-- name: ContactCreateJSONB :one
INSERT INTO
  contacts (blob, name)
VALUES
  (JSONB(:blob), :name) -- JSONB requires functional named param
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
