package sql_test

import (
	"testing"

	"github.com/nzoschke/sqlc-gen-zz/db"
	"github.com/nzoschke/sqlc-gen-zz/sql"
	"github.com/nzoschke/sqlc-gen-zz/zz"
	"github.com/stretchr/testify/assert"
)

func TestCRUD(t *testing.T) {
	ctx := t.Context()
	a := assert.New(t)

	db, err := db.New(ctx, sql.SQL, "file::memory:?mode=memory&cache=shared")
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

	n, err := zz.ContactCount(conn)
	a.NoError(err)
	a.Equal(int64(1), n)

	err = zz.ContactUpdate(conn, zz.ContactUpdateIn{
		Id:   1,
		Name: "new",
	})
	a.NoError(err)

	c, err := zz.ContactRead(conn, 1)
	a.NoError(err)

	a.Equal(&zz.ContactReadOut{
		Blob:      []byte("b"),
		CreatedAt: c.CreatedAt,
		Id:        1,
		Name:      "new",
	}, c)

	c2, err := zz.ContactCreate(conn, zz.ContactCreateIn{
		Blob: []byte("b"),
		Name: "name",
	})
	a.NoError(err)

	cs, err := zz.ContactList(conn, 10)
	a.NoError(err)

	a.Equal(zz.ContactListOut{
		{
			Blob:      []byte("b"),
			CreatedAt: c1.CreatedAt,
			Id:        1,
			Name:      "new",
		},
		{
			Blob:      []byte("b"),
			CreatedAt: c2.CreatedAt,
			Id:        2,
			Name:      "name",
		},
	}, cs)

	ns, err := zz.ContactListNames(conn, 10)
	a.NoError(err)

	a.Equal([]string{"new", "name"}, ns)

	err = zz.ContactDelete(conn, 1)
	a.NoError(err)

	err = zz.ContactDelete(conn, 2)
	a.NoError(err)

	cs, err = zz.ContactList(conn, 10)
	a.NoError(err)

	a.Equal(zz.ContactListOut{}, cs)
}

func TestJSONB(t *testing.T) {
	ctx := t.Context()
	a := assert.New(t)

	db, err := db.New(ctx, sql.SQL, "file::memory:?mode=memory&cache=shared")
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
		Id:        1,
		Name:      "name",
	}, c1)

	c2, err := zz.ContactReadJSONB(conn, 1)
	a.NoError(err)

	a.Equal([]byte("{}"), c2)

	err = zz.ContactDeleteAll(conn)
	a.NoError(err)
}
