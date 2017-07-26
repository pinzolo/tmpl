package main

import (
	"bytes"
	"os"
	"testing"
)

func TestRender(t *testing.T) {
	reset := setTestEnv()
	defer reset()
	ob := &bytes.Buffer{}
	eb := &bytes.Buffer{}
	w := writer{out: ob, err: eb}
	code := render(w, "doma_domain", "pkg:foo.bar.baz", "name:Code", "type:String")
	if code != 0 {
		t.Error(eb)
		return
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
	reset := setTestEnv()
	defer reset()
	ob := &bytes.Buffer{}
	eb := &bytes.Buffer{}
	w := writer{out: ob, err: eb}
	code := render(w, "unknown", "pkg:foo.bar.baz", "name:Code", "type:String")
	if code == 0 {
		t.Error("render with unknown template name should be error.")
	}
}

func TestBrokenTemplate(t *testing.T) {
	reset := setTestEnv()
	defer reset()
	ob := &bytes.Buffer{}
	eb := &bytes.Buffer{}
	w := writer{out: ob, err: eb}
	code := render(w, "broken", "foo:bar")
	if code == 0 {
		t.Error("render with broken template should be error.")
	}
}

func setTestEnv() func() {
	pre := os.Getenv("XDG_DATA_HOME")
	os.Setenv("XDG_DATA_HOME", "testdata")
	return func() {
		os.Setenv("XDG_DATA_HOME", pre)
	}
}
