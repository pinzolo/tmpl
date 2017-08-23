package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"text/template"

	"golang.org/x/crypto/ssh/terminal"

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

var jsonOpt bool

func main() {
	w := writer{
		out: os.Stdout,
		err: os.Stderr,
	}
	flag.BoolVar(&jsonOpt, "json", false, "Read value as json")
	flag.Parse()
	code := render(w, flag.Args()...)
	os.Exit(code)
}

func render(w writer, args ...string) int {
	if len(args) == 0 {
		// TODO: print usage
	}
	name := args[0]
	t, err := tmpl(name)
	if err != nil {
		fmt.Fprintln(w.err, err)
		return 2
	}
	t = t.Funcs(funcs())
	p, err := params(args[1:])
	if err != nil {
		fmt.Fprintln(w.err, err)
		return 2
	}
	err = t.Execute(w.out, p)
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

func params(args []string) (Params, error) {
	if jsonOpt {
		b, err := jsonBytes(args)
		if err != nil {
			return Params{}, err
		}
		return jsonParams(b)
	}
	return flatParams(args)
}

func jsonBytes(args []string) ([]byte, error) {
	var err error
	var b []byte
	if terminal.IsTerminal(0) {
		if len(args) == 0 {
			return nil, errors.New("no parameter")
		}
		b = []byte(args[0])
	} else {
		b, err = ioutil.ReadAll(os.Stdin)
	}
	if err != nil {
		return nil, err
	}
	return b, nil
}

func jsonParams(b []byte) (Params, error) {
	var m = make(map[string]interface{})
	err := json.Unmarshal(b, &m)
	if err != nil {
		return Params{}, err
	}
	return Params(m), nil
}

func flatParams(args []string) (Params, error) {
	if len(args) == 0 {
		return Params{}, errors.New("no parameter")
	}
	p := Params(make(map[string]interface{}))
	for _, a := range args {
		sa := strings.Split(a, ":")
		if len(sa) > 1 {
			p[sa[0]] = sa[1]
			continue
		}
		p[sa[0]] = ""
	}
	return p, nil
}
