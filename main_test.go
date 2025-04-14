package main_test

import (
	"embed"
	"testing"

	"github.com/nzoschke/sqlc-gen-zz/db"
	"github.com/nzoschke/sqlc-gen-zz/zz"
	"github.com/stretchr/testify/assert"
)

//go:embed *.sql
var SQL embed.FS

func TestFoo(t *testing.T) {
	ctx := t.Context()
	a := assert.New(t)

	db, err := db.New(ctx, SQL, "file::memory:?mode=memory&cache=shared")
	a.NoError(err)

	conn, put, err := db.Take(ctx)
	a.NoError(err)
	defer put()

	out, err := zz.ContactCreate(conn, zz.ContactCreateIn{
		Blob: []byte("b"),
		Name: "name",
	})
	a.NoError(err)

	a.Equal(&zz.ContactCreateOut{
		Blob:      []byte("b"),
		CreatedAt: out.CreatedAt,
		Id:        1,
		Name:      "name",
	}, out)
}
