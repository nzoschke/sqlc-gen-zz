package db

import (
	"context"
	"io/fs"

	"github.com/olekukonko/errors"
	"zombiezen.com/go/sqlite/sqlitemigration"
)

func (d *DB) migrate(ctx context.Context, fsys fs.FS) error {
	bs, err := fs.ReadFile(fsys, "schema.sql")
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
