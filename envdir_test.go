package main

import (
	"bytes"
	"testing"
)

func TestEnvdir(t *testing.T) {
	t.Run("simple case", func(t *testing.T) {
		outStream, errStream := new(bytes.Buffer), new(bytes.Buffer)
		e := &Envdir{outStream: outStream, errStream: errStream}

		status := e.Run([]string{"envdir", "testenv", "sh", "-c", "printf %s $FOO"})
		if status != 0 {
			t.Errorf("expected %d, but got %d", 0, status)
		}

		if outStream.String() != "foo" {
			t.Errorf("expected %q, but got %q", "foo", outStream.String())
		}
	})

	t.Run("command exit with error", func(t *testing.T) {
		outStream, errStream := new(bytes.Buffer), new(bytes.Buffer)
		e := &Envdir{outStream: outStream, errStream: errStream}

		status := e.Run([]string{"envdir", "testenv", "sh", "-c", "exit 3"})
		if status != 3 {
			t.Errorf("expected %d, but got %d", 3, status)
		}
	})
}
