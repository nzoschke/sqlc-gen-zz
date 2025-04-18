package sql_test

import (
	"testing"
	"time"

	"github.com/nzoschke/sqlc-gen-zz/pkg/db"
	"github.com/nzoschke/sqlc-gen-zz/pkg/sql/c"
	"github.com/nzoschke/sqlc-gen-zz/pkg/sql/models"
	"github.com/stretchr/testify/assert"
)

func TestCRUD(t *testing.T) {
	ctx := t.Context()
	a := assert.New(t)

	db, err := db.New(ctx, "file::memory:?mode=memory&cache=shared")
	a.NoError(err)

	conn, put, err := db.Take(ctx)
	a.NoError(err)
	defer put()

	c1, err := c.ContactCreate(conn, c.ContactCreateIn{
		Blob: []byte("b"),
		Info: models.Info{
			Age: 21,
		},
		Name: "name",
	})
	a.NoError(err)

	a.Equal(&c.ContactCreateOut{
		Blob:      []byte("b"),
		CreatedAt: c1.CreatedAt,
		Id:        1,
		Info: models.Info{
			Age: 21,
		},
		Name: "name",
	}, c1)

	a.Equal(time.Now().UTC().Format("2006-01-02"), c1.CreatedAt.Format("2006-01-02"))

	n, err := c.ContactCount(conn)
	a.NoError(err)
	a.Equal(int64(1), n)

	at := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	err = c.ContactUpdate(conn, c.ContactUpdateIn{
		CreatedAt: at,
		Id:        1,
		Name:      "new",
	})
	a.NoError(err)

	cr, err := c.ContactRead(conn, 1)
	a.NoError(err)

	a.Equal(&c.ContactReadOut{
		Blob:      []byte("b"),
		CreatedAt: at,
		Id:        1,
		Info: models.Info{
			Age: 21,
		},
		Name: "new",
	}, cr)

	c2, err := c.ContactCreate(conn, c.ContactCreateIn{
		Blob: []byte("b"),
		Name: "name",
	})
	a.NoError(err)

	cs, err := c.ContactList(conn, 10)
	a.NoError(err)

	a.Equal(c.ContactListOut{
		{
			Blob:      []byte("b"),
			CreatedAt: at,
			Id:        1,
			Info: models.Info{
				Age: 21,
			},
			Name: "new",
		},
		{
			Blob:      []byte("b"),
			CreatedAt: c2.CreatedAt,
			Id:        2,
			Info:      models.Info{},
			Name:      "name",
		},
	}, cs)

	ns, err := c.ContactListNames(conn, 10)
	a.NoError(err)

	a.Equal([]string{"new", "name"}, ns)

	err = c.ContactDelete(conn, 1)
	a.NoError(err)

	err = c.ContactDelete(conn, 2)
	a.NoError(err)

	cs, err = c.ContactList(conn, 10)
	a.NoError(err)

	a.Equal(c.ContactListOut{}, cs)
}

func TestJSONB(t *testing.T) {
	ctx := t.Context()
	a := assert.New(t)

	db, err := db.New(ctx, "file::memory:?mode=memory&cache=shared")
	a.NoError(err)

	conn, put, err := db.Take(ctx)
	a.NoError(err)
	defer put()

	c1, err := c.ContactCreateJSONB(conn, c.ContactCreateJSONBIn{
		Blob: []byte("{}"),
		Info: models.Info{},
		Name: "name",
	})
	a.NoError(err)

	a.Equal(&c.ContactCreateJSONBOut{
		Json:      []byte("{}"),
		Blob:      c1.Blob,
		CreatedAt: c1.CreatedAt,
		Id:        1,
		Info:      models.Info{},
		Name:      "name",
	}, c1)

	c2, err := c.ContactReadJSONB(conn, 1)
	a.NoError(err)

	a.Equal([]byte("{}"), c2)

	err = c.ContactDeleteAll(conn)
	a.NoError(err)
}

func TestJSON(t *testing.T) {

}
