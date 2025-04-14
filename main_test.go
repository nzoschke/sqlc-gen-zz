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

func TestOneMany(t *testing.T) {
	ctx := t.Context()
	a := assert.New(t)

	db, err := db.New(ctx, SQL, "file::memory:?mode=memory&cache=shared")
	a.NoError(err)

	conn, put, err := db.Take(ctx)
	a.NoError(err)
	defer put()

	c1, err := zz.ContactCreate(conn, zz.ContactCreateIn{
		Blob: []byte("b"),
		Name: "name",
	})
	a.NoError(err)

	a.Equal(&zz.ContactCreateOut{
		Blob:      []byte("b"),
		CreatedAt: c1.CreatedAt,
		Id:        1,
		Name:      "name",
	}, c1)

	c2, err := zz.ContactCreate(conn, zz.ContactCreateIn{
		Blob: []byte("b"),
		Name: "name",
	})
	a.NoError(err)

	cs, err := zz.ContactList(conn, zz.ContactListIn{
		Limit: 10,
	})
	a.NoError(err)

	a.Equal(zz.ContactListOut{
		{
			Blob:      []byte("b"),
			CreatedAt: c1.CreatedAt,
			Id:        1,
			Name:      "name",
		},
		{
			Blob:      []byte("b"),
			CreatedAt: c2.CreatedAt,
			Id:        2,
			Name:      "name",
		},
	}, cs)
}

func TestJSONB(t *testing.T) {
	ctx := t.Context()
	a := assert.New(t)

	db, err := db.New(ctx, SQL, "file::memory:?mode=memory&cache=shared")
	a.NoError(err)

	conn, put, err := db.Take(ctx)
	a.NoError(err)
	defer put()

	c1, err := zz.ContactCreateJSONB(conn, zz.ContactCreateJSONBIn{
		Blob: []byte("{}"),
		Name: "name",
	})
	a.NoError(err)

	a.Equal(&zz.ContactCreateJSONBOut{
		Json:      []byte("{}"),
		Blob:      c1.Blob,
		CreatedAt: c1.CreatedAt,
		Id:        3,
		Name:      "name",
	}, c1)

	c2, err := zz.ContactReadJSONB(conn, zz.ContactReadJSONBIn{
		Id: 3,
	})
	a.NoError(err)

	a.Equal(&zz.ContactReadJSONBOut{
		Blob: []byte("{}"),
	}, c2)

}
