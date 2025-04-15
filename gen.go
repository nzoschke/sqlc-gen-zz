package main

import (
	"bytes"
	"context"
	"embed"
	"fmt"
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

func gotype(dbtype string) string {
	switch strings.ToLower(dbtype) {
	case "any", "blob":
		return "[]byte"
	case "integer":
		return "int64"
	case "real":
		return "float64"
	case "text":
		return "string"
	}
	return "[]byte"
}

func Gen(ctx context.Context, req *plugin.GenerateRequest) (*plugin.GenerateResponse, error) {
	slog.Info("gen", "req", req)

	funcMap := template.FuncMap{
		"inarg": func(name string, ps []*plugin.Parameter) string {
			switch len(ps) {
			case 0:
				return ""
			case 1:
				// id int
				p := ps[0]
				return fmt.Sprintf(", %s %s", p.Column.Name, gotype(p.Column.Type.Name))
			default:
				// ContactCreateIn
				return fmt.Sprintf(", in %sIn", name)
			}
		},
		"camel": strcase.ToCamel,
		"dbtype": func(dbtype string) string {
			switch strings.ToLower(dbtype) {
			case "any", "blob":
				return "Bytes"
			case "integer":
				return "Int64"
			case "real":
				return "Float"
			case "text":
				return "Text"
			}
			return "Bytes"
		},
		"gotype":   gotype,
		"lower":    strings.ToLower,
		"singular": pl.Singular,
	}

	c, err := template.New("catalog.tmpl").Funcs(funcMap).ParseFS(tmpl, "*.tmpl")
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var cbuf bytes.Buffer
	if err := c.Execute(&cbuf, req); err != nil {
		return nil, errors.WithStack(err)
	}

	q, err := template.New("queries.tmpl").Funcs(funcMap).ParseFS(tmpl, "*.tmpl")
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var qbuf bytes.Buffer
	if err := q.Execute(&qbuf, req); err != nil {
		return nil, errors.WithStack(err)
	}

	return &plugin.GenerateResponse{
		Files: []*plugin.File{
			{
				Contents: cbuf.Bytes(),
				Name:     "catalog.go",
			},
			{
				Contents: qbuf.Bytes(),
				Name:     "queries.go",
			},
		},
	}, nil
}
