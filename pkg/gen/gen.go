package gen

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"strings"
	"text/template"

	"github.com/gertd/go-pluralize"
	"github.com/iancoleman/strcase"
	"github.com/nzoschke/sqlc-gen-zz/pkg/tmpl"
	"github.com/olekukonko/errors"
	"github.com/sqlc-dev/plugin-sdk-go/plugin"
)

var pl = pluralize.NewClient()

func Gen(ctx context.Context, req *plugin.GenerateRequest) (*plugin.GenerateResponse, error) {
	slog.Info("gen", "req", req)

	funcMap := template.FuncMap{
		"camel": strcase.ToCamel,
		"dbtype": func(dbtype string) string {
			switch strings.ToLower(dbtype) {
			case "integer":
				return "Int64"
			case "real":
				return "Float"
			case "datetime", "text":
				return "Text"
			default:
				return "Bytes"
			}
		},
		"gotype": gotype,
		"inarg": func(name string, ps []*plugin.Parameter) string {
			switch len(ps) {
			case 0:
				return ""
			case 1:
				// id int
				p := ps[0]
				return fmt.Sprintf(", %s %s", p.Column.Name, gotype(p.Column))
			default:
				// ContactCreateIn
				return fmt.Sprintf(", in %sIn", name)
			}
		},
		"outarg": func(name string, cs []*plugin.Column) string {
			switch len(cs) {
			case 0:
				return ""
			case 1:
				// int
				c := cs[0]
				return fmt.Sprintf("%s, ", gotype(c))
			default:
				// *ContactCreateOut
				return fmt.Sprintf("*%sOut, ", name)
			}
		},
		"outempty": func(name string, cs []*plugin.Column) string {
			switch len(cs) {
			case 0:
				return ""
			case 1:
				c := cs[0]
				switch gotype(c) {
				case "[]byte":
					return "nil"
				case "text":
					return `""`
				default:
					return "0"
				}
			default:
				return "nil"
			}
		},
		"lower":    strings.ToLower,
		"singular": pl.Singular,
		"timeimport": func(c *plugin.Catalog) string {
			for _, s := range c.Schemas {
				for _, t := range s.Tables {
					for _, c := range t.Columns {
						if gotype(c) == "time.Time" {
							return `"time"`
						}
					}
				}
			}
			return ""
		},
	}

	res := &plugin.GenerateResponse{
		Files: []*plugin.File{},
	}

	t, err := template.New("catalog.tmpl").Funcs(funcMap).ParseFS(tmpl.Tmpl, "*.tmpl")
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, req); err != nil {
		return nil, errors.WithStack(err)
	}

	res.Files = append(res.Files, &plugin.File{
		Contents: buf.Bytes(),
		Name:     "catalog.go",
	})

	t, err = template.New("queries.tmpl").Funcs(funcMap).ParseFS(tmpl.Tmpl, "*.tmpl")
	if err != nil {
		return nil, errors.WithStack(err)
	}

	buf = bytes.Buffer{}
	if err := t.Execute(&buf, req); err != nil {
		return nil, errors.WithStack(err)
	}

	res.Files = append(res.Files, &plugin.File{
		Contents: buf.Bytes(),
		Name:     "queries.go",
	})

	return res, nil
}

func gotype(c *plugin.Column) string {
	if strings.HasSuffix(c.Name, "_at") {
		return "time.Time"
	}

	// https://sqlite.org/datatype3.html#affinity_name_examples
	switch strings.ToLower(c.Type.Name) {
	case "integer":
		return "int64"
	case "real":
		return "float64"
	case "text":
		return "string"
	default:
		return "[]byte"
	}
}
