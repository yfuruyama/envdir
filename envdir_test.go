package main

import (
	"bytes"
	"testing"
)

func TestEnvdir(t *testing.T) {
	outStream, errStream := new(bytes.Buffer), new(bytes.Buffer)
	e := &Envdir{outStream: outStream, errStream: errStream}

	status := e.Run([]string{"envdir", "testenv", "sh", "-c", "printf %s $FOO"})
	if status != ExitCodeOk {
		t.Errorf("expected %d, but got %d", ExitCodeOk, status)
	}

	if outStream.String() != "foo" {
		t.Errorf("expected %q, but got %q", "foo", outStream.String())
	}
}
