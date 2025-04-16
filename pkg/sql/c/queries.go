// Code generated by "sqlc-gen-zz". DO NOT EDIT.

package c

import (
	"database/sql"
	"time"

	"zombiezen.com/go/sqlite"
)

type ContactCreateIn struct {
	Blob []byte `json:"blob"`
	Name string `json:"name"`
}

type ContactCreateOut struct {
	Blob      []byte    `json:"blob"`
	CreatedAt time.Time `json:"created_at"`
	Id        int64     `json:"id"`
	Name      string    `json:"name"`
}

func ContactCreate(tx *sqlite.Conn, in ContactCreateIn) (*ContactCreateOut, error) {
	stmt := tx.Prep(`INSERT INTO
  contacts (blob, name)
VALUES
  (?, ?)
RETURNING
  blob, created_at, id, name`)
	defer stmt.Reset()

	stmt.BindBytes(1, in.Blob)
	stmt.BindText(2, in.Name)

	ok, err := stmt.Step()
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, sql.ErrNoRows
	}

	out := ContactCreateOut{}
	out.Blob = []byte(stmt.ColumnText(0))
	out.CreatedAt = timeParse(stmt.ColumnText(1))
	out.Id = stmt.ColumnInt64(2)
	out.Name = stmt.ColumnText(3)

	return &out, nil

}

type ContactReadOut struct {
	Blob      []byte    `json:"blob"`
	CreatedAt time.Time `json:"created_at"`
	Id        int64     `json:"id"`
	Name      string    `json:"name"`
}

func ContactRead(tx *sqlite.Conn, id int64) (*ContactReadOut, error) {
	stmt := tx.Prep(`SELECT blob, created_at, id, name FROM contacts WHERE id = ? LIMIT 1`)
	defer stmt.Reset()

	stmt.BindInt64(1, id)

	ok, err := stmt.Step()
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, sql.ErrNoRows
	}

	out := ContactReadOut{}
	out.Blob = []byte(stmt.ColumnText(0))
	out.CreatedAt = timeParse(stmt.ColumnText(1))
	out.Id = stmt.ColumnInt64(2)
	out.Name = stmt.ColumnText(3)

	return &out, nil

}

type ContactCountOut struct {
	Count int64 `json:"count"`
}

func ContactCount(tx *sqlite.Conn) (int64, error) {
	stmt := tx.Prep(`SELECT COUNT(*) FROM contacts`)
	defer stmt.Reset()

	ok, err := stmt.Step()
	if err != nil {
		return 0, err
	}
	if !ok {
		return 0, sql.ErrNoRows
	}

	return stmt.ColumnInt64(0), nil

}

type ContactListOut []ContactListRow

type ContactListRow struct {
	Blob      []byte    `json:"blob"`
	CreatedAt time.Time `json:"created_at"`
	Id        int64     `json:"id"`
	Name      string    `json:"name"`
}

func ContactList(tx *sqlite.Conn, limit int64) (ContactListOut, error) {
	stmt := tx.Prep(`SELECT
  blob, created_at, id, name
FROM
  contacts
LIMIT
  ?`)
	defer stmt.Reset()

	stmt.BindInt64(1, limit)

	out := ContactListOut{}
	for {
		ok, err := stmt.Step()
		if err != nil {
			return nil, err
		}
		if !ok {
			break
		}

		row := ContactListRow{}
		row.Blob = []byte(stmt.ColumnText(0))
		row.CreatedAt = timeParse(stmt.ColumnText(1))
		row.Id = stmt.ColumnInt64(2)
		row.Name = stmt.ColumnText(3)

		out = append(out, row)
	}

	return out, nil
}

func ContactListNames(tx *sqlite.Conn, limit int64) ([]string, error) {
	stmt := tx.Prep(`SELECT
  name
FROM
  contacts
LIMIT
  ?`)
	defer stmt.Reset()

	stmt.BindInt64(1, limit)

	out := []string{}
	for {
		ok, err := stmt.Step()
		if err != nil {
			return nil, err
		}
		if !ok {
			break
		}

		c := stmt.ColumnText(0)

		out = append(out, c)
	}

	return out, nil
}

type ContactUpdateIn struct {
	CreatedAt time.Time `json:"created_at"`
	Name      string    `json:"name"`
	Id        int64     `json:"id"`
}

func ContactUpdate(tx *sqlite.Conn, in ContactUpdateIn) error {
	stmt := tx.Prep(`UPDATE
  contacts
SET
  created_at = ?,
  name = ?
WHERE
  id = ?`)
	defer stmt.Reset()

	stmt.BindText(1, in.CreatedAt.Format("2006-01-02 15:04:05"))
	stmt.BindText(2, in.Name)
	stmt.BindInt64(3, in.Id)

	_, err := stmt.Step()
	if err != nil {
		return err
	}

	return nil
}

func ContactDelete(tx *sqlite.Conn, id int64) error {
	stmt := tx.Prep(`DELETE FROM
  contacts
WHERE
  id = ?`)
	defer stmt.Reset()

	stmt.BindInt64(1, id)

	_, err := stmt.Step()
	if err != nil {
		return err
	}

	return nil
}

func ContactDeleteAll(tx *sqlite.Conn) error {
	stmt := tx.Prep(`DELETE FROM
  contacts`)
	defer stmt.Reset()

	_, err := stmt.Step()
	if err != nil {
		return err
	}

	return nil
}

type ContactCreateJSONBIn struct {
	Blob []byte `json:"blob"`
	Name string `json:"name"`
}

type ContactCreateJSONBOut struct {
	Json      []byte    `json:"json"`
	Blob      []byte    `json:"blob"`
	CreatedAt time.Time `json:"created_at"`
	Id        int64     `json:"id"`
	Name      string    `json:"name"`
}

func ContactCreateJSONB(tx *sqlite.Conn, in ContactCreateJSONBIn) (*ContactCreateJSONBOut, error) {
	stmt := tx.Prep(`INSERT INTO
  contacts (blob, name)
VALUES
  (JSONB(?1), ?2) -- JSONB requires functional named param
RETURNING
  JSON(blob),  -- and requires functional return param in position 1
  blob, created_at, id, name`)
	defer stmt.Reset()

	stmt.BindBytes(1, in.Blob)
	stmt.BindText(2, in.Name)

	ok, err := stmt.Step()
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, sql.ErrNoRows
	}

	out := ContactCreateJSONBOut{}
	out.Json = []byte(stmt.ColumnText(0))
	out.Blob = []byte(stmt.ColumnText(1))
	out.CreatedAt = timeParse(stmt.ColumnText(2))
	out.Id = stmt.ColumnInt64(3)
	out.Name = stmt.ColumnText(4)

	return &out, nil

}

type ContactReadJSONBOut struct {
	Blob []byte `json:"blob"`
}

func ContactReadJSONB(tx *sqlite.Conn, id int64) ([]byte, error) {
	stmt := tx.Prep(`SELECT
  JSON(blob) AS blob
FROM
  contacts
WHERE
  id = ?
LIMIT
  1`)
	defer stmt.Reset()

	stmt.BindInt64(1, id)

	ok, err := stmt.Step()
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, sql.ErrNoRows
	}

	return []byte(stmt.ColumnText(0)), nil

}

func timeParse(s string) time.Time {
	t, _ := time.Parse("2006-01-02 15:04:05", s)
	return t
}
