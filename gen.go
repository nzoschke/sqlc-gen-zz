package main

import (
	"bytes"
	"context"
	"embed"
	"log/slog"
	"text/template"

	"github.com/olekukonko/errors"
	"github.com/sqlc-dev/plugin-sdk-go/plugin"
)

//go:embed *.tmpl
var tmpl embed.FS

func Gen(ctx context.Context, req *plugin.GenerateRequest) (*plugin.GenerateResponse, error) {
	slog.Info("gen", "req", req)

	t, err := template.ParseFS(tmpl, "model.tmpl")
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, req); err != nil {
		return nil, errors.WithStack(err)
	}

	return &plugin.GenerateResponse{
		Files: []*plugin.File{
			{
				Contents: buf.Bytes(),
				Name:     "model.go",
			},
			{
				Contents: []byte("package zz"),
				Name:     "gen.go",
			},
		},
	}, nil
}
