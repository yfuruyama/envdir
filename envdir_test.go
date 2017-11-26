package main

import (
	"bytes"
	"testing"
)

func TestEnvdir_success(t *testing.T) {
	t.Run("simple case", func(t *testing.T) {
		outStream, errStream := new(bytes.Buffer), new(bytes.Buffer)
		e := &Envdir{outStream: outStream, errStream: errStream}

		status := e.Run([]string{"envdir", "testenv/simple", "sh", "-c", "printf %s $FOO"})
		if status != 0 {
			t.Errorf("expected %d, but got %d", 0, status)
		}

		if outStream.String() != "foo" {
			t.Errorf("expected %q, but got %q", "foo", outStream.String())
		}
	})

	t.Run("env file size is 0, so remove that env", func(t *testing.T) {
		// outStream, errStream := new(bytes.Buffer), new(bytes.Buffer)
		// e := &Envdir{outStream: outStream, errStream: errStream}

		// status := e.Run([]string{"envdir", "testenv", "sh", "-c", "printf %s $NO_VALUE"})
		// if status != 0 {
		// t.Errorf("expected %d, but got %d", 0, status)
		// }

		// if outStream.String() != "" {
		// t.Errorf("expected %q, but got %q", "", outStream.String())
		// }
	})

	t.Run("the end of space in env value is removed", func(t *testing.T) {
		// TODO
	})

	t.Run("the end of tab in env value is removed", func(t *testing.T) {
		// TODO
	})
}

func TestEnvdir_error(t *testing.T) {
	t.Run("less arguments provided", func(t *testing.T) {
		outStream, errStream := new(bytes.Buffer), new(bytes.Buffer)
		e := &Envdir{outStream: outStream, errStream: errStream}

		status := e.Run([]string{"envdir"})
		if status != 111 {
			t.Errorf("expected %d, but got %d", 111, status)
		}

		status = e.Run([]string{"envdir", "testenv/simple"})
		if status != 111 {
			t.Errorf("expected %d, but got %d", 111, status)
		}
	})

	t.Run("dir must not contain inner directory", func(t *testing.T) {
		outStream, errStream := new(bytes.Buffer), new(bytes.Buffer)
		e := &Envdir{outStream: outStream, errStream: errStream}

		status := e.Run([]string{"envdir", "testenv/include_dir", "sh", "-c", "printf %s $FOO"})
		if status != 111 {
			t.Errorf("expected %d, but got %d", 111, status)
		}
	})

	t.Run("command exit with error", func(t *testing.T) {
		outStream, errStream := new(bytes.Buffer), new(bytes.Buffer)
		e := &Envdir{outStream: outStream, errStream: errStream}

		status := e.Run([]string{"envdir", "testenv/simple", "sh", "-c", "exit 3"})
		if status != 3 {
			t.Errorf("expected %d, but got %d", 3, status)
		}
	})
}
