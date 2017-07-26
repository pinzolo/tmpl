package main

import (
	"bytes"
	"os"
	"testing"
)

func TestRender(t *testing.T) {
	os.Setenv("XDG_DATA_HOME", "testdata")
	ob := &bytes.Buffer{}
	eb := &bytes.Buffer{}
	w := writer{out: ob, err: eb}
	code := render(w, "doma_domain", "pkg:foo.bar.baz", "name:Code", "type:String")
	if code != 0 {
		t.Fatal(eb)
	}
	expected := `package foo.bar.baz;

import org.seasar.doma.Domain;

@Domain(valueType = String.class, factoryMethod = "of")
public class Code {
  private Code(final String value) {
    this.value = value;
  }

  private String value;
  public String getValue() {
    return value;
  }

  public static Code of(final String value) {
    return new Code(value);
  }
}
`
	if actual := ob.String(); actual != expected {
		t.Errorf("expected:\n%s\n\nactual:\n%s", expected, actual)
	}
}

func TestTemplateNotFound(t *testing.T) {
	os.Setenv("XDG_DATA_HOME", "testdata")
	ob := &bytes.Buffer{}
	eb := &bytes.Buffer{}
	w := writer{out: ob, err: eb}
	code := render(w, "unknown", "pkg:foo.bar.baz", "name:Code", "type:String")
	if code == 0 {
		t.Fatal("render with unknown template name should be error.")
	}
}

func TestBrokenTemplate(t *testing.T) {
	os.Setenv("XDG_DATA_HOME", "testdata")
	ob := &bytes.Buffer{}
	eb := &bytes.Buffer{}
	w := writer{out: ob, err: eb}
	code := render(w, "broken", "foo:bar")
	if code == 0 {
		t.Fatal("render with broken template should be error.")
	}
}
