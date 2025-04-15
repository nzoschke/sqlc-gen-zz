package db

import (
	"context"
	"io/fs"

	"github.com/nzoschke/sqlc-gen-zz/pkg/sql"
	"github.com/olekukonko/errors"
	"zombiezen.com/go/sqlite/sqlitemigration"
)

func (d *DB) migrate(ctx context.Context) error {
	bs, err := fs.ReadFile(sql.SQL, "schema.sql")
	if err != nil {
		return errors.WithStack(err)
	}

	conn, err := d.pool.Take(ctx)
	if err != nil {
		return errors.WithStack(err)
	}
	defer d.pool.Put(conn)

	if err := sqlitemigration.Migrate(ctx, conn,
		sqlitemigration.Schema{
			Migrations: []string{string(bs)},
		},
	); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
