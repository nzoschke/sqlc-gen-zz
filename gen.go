package main

import (
	"bytes"
	"context"
	"embed"
	"log/slog"
	"strings"
	"text/template"

	"github.com/gertd/go-pluralize"
	"github.com/iancoleman/strcase"
	"github.com/olekukonko/errors"
	"github.com/sqlc-dev/plugin-sdk-go/plugin"
)

//go:embed *.tmpl
var tmpl embed.FS

var pl = pluralize.NewClient()

func Gen(ctx context.Context, req *plugin.GenerateRequest) (*plugin.GenerateResponse, error) {
	slog.Info("gen", "req", req)

	funcMap := template.FuncMap{
		"camel": strcase.ToCamel,
		"gotype": func(dbtype string) string {
			switch strings.ToLower(dbtype) {
			case "blob":
				return "[]byte"
			case "integer":
				return "int64"
			case "real":
				return "float64"
			case "text":
				return "string"
			}
			return "any"
		},
		"singular": pl.Singular,
	}

	t, err := template.New("catalog.tmpl").Funcs(funcMap).ParseFS(tmpl, "catalog.tmpl")
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
				Name:     "catalog.go",
			},
			{
				Contents: []byte("package zz"),
				Name:     "gen.go",
			},
		},
	}, nil
}
