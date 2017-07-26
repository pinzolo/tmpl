package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"text/template"

	"github.com/pinzolo/casee"
	"github.com/pinzolo/xdgdir"
)

// Params is paramters for template.
// Values are parsed from commandline arguments.
type Params map[string]interface{}

type writer struct {
	out io.Writer
	err io.Writer
}

func main() {
	w := writer{
		out: os.Stdout,
		err: os.Stderr,
	}
	flag.Parse()
	code := render(w, flag.Args()...)
	os.Exit(code)
}

func render(w writer, args ...string) int {
	if len(args) == 0 {
		// TODO: print usage
	}
	name := args[0]
	if len(args) == 1 {
		// TODO: json from pipe
	}
	t, err := tmpl(name)
	if err != nil {
		fmt.Fprintln(w.err, err)
		return 2
	}
	t = t.Funcs(funcs())
	err = t.Execute(w.out, params(args[1:]))
	if err != nil {
		fmt.Fprintln(w.err, err)
		return 2
	}
	return 0
}

func tmpl(name string) (*template.Template, error) {
	app := xdgdir.NewApp("tmpl")
	tp, err := app.DataFile(name + ".tmpl")
	if err != nil {
		return nil, err
	}
	return template.ParseFiles(tp)
}

func funcs() template.FuncMap {
	return template.FuncMap{
		"snakecase":  casee.ToSnakeCase,
		"chaincase":  casee.ToChainCase,
		"camelcase":  casee.ToCamelCase,
		"pascalcase": casee.ToPascalCase,
		"upper":      strings.ToUpper,
		"lower":      strings.ToLower,
		"contains":   strings.Contains,
		"hasprefix":  strings.HasPrefix,
		"hassuffix":  strings.HasSuffix,
	}
}

func params(args []string) Params {
	p := Params(make(map[string]interface{}))
	for _, a := range args {
		sa := strings.Split(a, ":")
		if len(sa) > 1 {
			p[sa[0]] = sa[1]
			continue
		}
		p[sa[0]] = ""
	}
	return p
}
