package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestEnvdir(t *testing.T) {
	outStream, errStream := new(bytes.Buffer), new(bytes.Buffer)
	e := &Envdir{outStream: outStream, errStream: errStream}
	args := strings.Split("envdir testenv perl -e 'print($ENV{FOO})'", " ")

	status := e.Run(args)
	if status != ExitCodeOk {
		t.Errorf("expected %d to eq %d", status, ExitCodeOk)
	}

	// if !strings.Contains(outStream.String(), "foo") {
	// t.Errorf("expected %q to contain %q", errStream.String(), "foo")
	// }
}
