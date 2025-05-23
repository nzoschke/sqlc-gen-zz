package gen

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"text/template"

	"github.com/gertd/go-pluralize"
	"github.com/iancoleman/strcase"
	"github.com/nzoschke/sqlc-gen-zz/pkg/tmpl"
	"github.com/olekukonko/errors"
	"github.com/sqlc-dev/plugin-sdk-go/plugin"
	"google.golang.org/protobuf/encoding/protojson"
)

var (
	overrides = map[string]GoType{}
	pl        = pluralize.NewClient()
)

func init() {
	strcase.ConfigureAcronym("id", "ID")
}

type Options struct {
	Overrides []Override `json:"overrides"`
}

type GoType struct {
	Import  string `json:"import"`
	Package string `json:"package"`
	Type    string `json:"type"`
}

type Override struct {
	Column string `json:"column"`
	GoType GoType `json:"go_type"`
}

type Data struct {
	Catalog   *plugin.Catalog
	Overrides []Override
	Queries   []*plugin.Query
	Settings  *plugin.Settings
}

func Gen(ctx context.Context, req *plugin.GenerateRequest) (*plugin.GenerateResponse, error) {
	opts := Options{}
	if err := json.Unmarshal(req.PluginOptions, &opts); err != nil {
		return nil, errors.WithStack(err)
	}

	for _, o := range opts.Overrides {
		overrides[o.Column] = o.GoType
	}

	data := Data{
		Catalog:   req.Catalog,
		Overrides: opts.Overrides,
		Queries:   req.Queries,
		Settings:  req.Settings,
	}

	res := &plugin.GenerateResponse{
		Files: []*plugin.File{},
	}

	fm := template.FuncMap{
		"camel":    strcase.ToCamel,
		"bindval":  bindval,
		"dbtype":   dbtype,
		"gotype":   gotype,
		"imports":  imports,
		"inarg":    inarg,
		"outarg":   outarg,
		"lower":    strings.ToLower,
		"retempty": retempty,
		"retval":   retval,
		"singular": pl.Singular,
	}

	t, err := template.New("catalog.tmpl").Funcs(fm).ParseFS(tmpl.Tmpl, "*.tmpl")
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return nil, errors.WithStack(err)
	}

	res.Files = append(res.Files, &plugin.File{
		Contents: buf.Bytes(),
		Name:     "catalog.go",
	})

	t, err = template.New("queries.tmpl").Funcs(fm).ParseFS(tmpl.Tmpl, "*.tmpl")
	if err != nil {
		return nil, errors.WithStack(err)
	}

	buf = bytes.Buffer{}
	if err := t.Execute(&buf, data); err != nil {
		return nil, errors.WithStack(err)
	}

	res.Files = append(res.Files, &plugin.File{
		Contents: buf.Bytes(),
		Name:     "queries.go",
	})

	if err := config(false, opts, req, res); err != nil {
		return nil, errors.WithStack(err)
	}

	return res, nil
}

// config generates .json files from the request for debugging purposes
func config(debug bool, opts Options, req *plugin.GenerateRequest, res *plugin.GenerateResponse) error {
	if !debug {
		return nil
	}

	bs, err := json.MarshalIndent(opts, "", "  ")
	if err != nil {
		return errors.WithStack(err)
	}

	res.Files = append(res.Files, &plugin.File{
		Contents: bs,
		Name:     "opts.json",
	})

	m := &protojson.MarshalOptions{
		EmitUnpopulated: true,
		Indent:          "",
		UseProtoNames:   true,
	}
	data, err := m.Marshal(req)
	if err != nil {
		return errors.WithStack(err)
	}

	var rm json.RawMessage = data
	bs, err = json.MarshalIndent(rm, "", "  ")
	if err != nil {
		return errors.WithStack(err)
	}

	res.Files = append(res.Files, &plugin.File{
		Contents: bs,
		Name:     "req.json",
	})

	return nil
}

func bindval(ps []*plugin.Parameter, i int32) string {
	p := ps[i-1]

	cast := func(v string) string {
		if _, t := overridetype(p.Column); t != "" {
			return fmt.Sprintf(`jsonMarshal(%s)`, v)
		}

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
	case "int", "integer":
		return "Int64"
	case "real":
		return "Float"
	case "text":
		return "Text"
	default:
		return "Bytes"
	}
}

func overridetype(c *plugin.Column) (GoType, string) {
	if c.Table == nil {
		return GoType{}, ""
	}

	t, ok := overrides[fmt.Sprintf("%s.%s", c.Table.Name, c.Name)]
	if !ok {
		return GoType{}, ""
	}

	return t, fmt.Sprintf("%s.%s", t.Package, t.Type)
}

func gotype(c *plugin.Column) string {
	if _, t := overridetype(c); t != "" {
		return t
	}

	if strings.HasSuffix(c.Name, "_at") {
		return "time.Time"
	}

	// https://sqlite.org/datatype3.html#affinity_name_examples
	switch strings.ToLower(c.Type.Name) {
	case "int", "integer":
		return "int64"
	case "real":
		return "float64"
	case "text":
		return "string"
	default:
		return "[]byte"
	}
}

func imports(os []Override) string {
	pkgs := map[string]struct{}{}

	for _, o := range os {
		pkgs[fmt.Sprintf("\"%s/%s\"", o.GoType.Import, o.GoType.Package)] = struct{}{}
	}

	ps := []string{}
	for p := range pkgs {
		ps = append(ps, p)
	}

	return strings.Join(ps, "\n")
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
		if g, t := overridetype(c); t != "" {
			return fmt.Sprintf("%s.%s{}", g.Package, g.Type)
		}

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

	if g, t := overridetype(c); t != "" {
		return fmt.Sprintf("jsonUnmarshal%s%s([]byte(stmt.ColumnText(%d)))", strcase.ToCamel(g.Package), strcase.ToCamel(g.Type), i)
	}

	switch gotype(c) {
	case "[]byte":
		return fmt.Sprintf("[]byte(stmt.ColumnText(%d))", i)
	case "json.RawMessage":
		return fmt.Sprintf("[]byte(stmt.ColumnText(%d))", i)
	case "time.Time":
		return fmt.Sprintf("timeParse(stmt.ColumnText(%d))", i)
	default:
		return fmt.Sprintf("stmt.Column%s(%d)", dbtype(c.Type.Name), i)
	}
}
