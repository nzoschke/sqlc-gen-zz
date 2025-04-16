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
		"camel":      strcase.ToCamel,
		"bindval":    bindval,
		"dbtype":     dbtype,
		"gotype":     gotype,
		"inarg":      inarg,
		"outarg":     outarg,
		"lower":      strings.ToLower,
		"retempty":   retempty,
		"retval":     retval,
		"singular":   pl.Singular,
		"timeimport": timeimport,
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

func bindval(ps []*plugin.Parameter, i int32) string {
	p := ps[i-1]

	cast := func(v string) string {
		switch gotype(p.Column) {
		case "time.Time":
			return fmt.Sprintf(`%s.Format("2006-01-02 15:04:05")`, v)
		default:
			return v
		}
	}

	switch len(ps) {
	case 1:
		// id
		return cast(p.Column.Name)
	default:
		// in.Id
		return cast(fmt.Sprintf("in.%s", strcase.ToCamel(p.Column.Name)))
	}
}

func dbtype(t string) string {
	switch strings.ToLower(t) {
	case "integer":
		return "Int64"
	case "real":
		return "Float"
	case "text":
		return "Text"
	default:
		return "Bytes"
	}
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

func inarg(name string, ps []*plugin.Parameter) string {
	switch len(ps) {
	case 0:
		return ""
	case 1:
		// id int
		p := ps[0]
		return fmt.Sprintf(", %s %s", p.Column.Name, gotype(p.Column))
	default:
		// in ContactCreateIn
		return fmt.Sprintf(", in %sIn", name)
	}
}

func outarg(name string, cs []*plugin.Column) string {
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
}

func retempty(name string, cs []*plugin.Column) string {
	switch len(cs) {
	case 0:
		return ""
	case 1:
		c := cs[0]
		switch gotype(c) {
		case "[]byte":
			return "nil"
		case "time.Time":
			return "time.Time{}"
		case "text":
			return `""`
		case "int64", "float64":
			return "0"
		default:
			return "nil"
		}
	default:
		return "nil"
	}
}

func retval(cs []*plugin.Column, i int) string {
	c := cs[i]

	switch gotype(c) {
	case "[]byte":
		return fmt.Sprintf("[]byte(stmt.ColumnText(%d))", i)
	case "time.Time":
		return fmt.Sprintf("timeParse(stmt.ColumnText(%d))", i)
	default:
		return fmt.Sprintf("stmt.Column%s(%d)", dbtype(c.Type.Name), i)
	}
}

func timeimport(c *plugin.Catalog) string {
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
}
