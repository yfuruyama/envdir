package main

import (
	"bytes"
	"testing"
)

func TestEnvdir_success(t *testing.T) {
	t.Run("simple case", func(t *testing.T) {
		outStream, errStream := new(bytes.Buffer), new(bytes.Buffer)
		e := &Envdir{outStream: outStream, errStream: errStream}

		status := e.Run([]string{"envdir", "testenv/simple", "printenv", "FOO"})
		if status != 0 {
			t.Errorf("expected %d, but got %d", 0, status)
		}

		if outStream.String() != "foo\n" {
			t.Errorf("expected %q, but got %q", "foo", outStream.String())
		}
	})

	t.Run("the end of space and tab in env value is removed", func(t *testing.T) {
		outStream, errStream := new(bytes.Buffer), new(bytes.Buffer)
		e := &Envdir{outStream: outStream, errStream: errStream}

		status := e.Run([]string{"envdir", "testenv/space", "printenv", "FOO"})
		if status != 0 {
			t.Errorf("expected %d, but got %d: out=%s, err=%s", 0, status, outStream, errStream)
		}

		if outStream.String() != "hello world\n" {
			t.Errorf("expected %q, but got %q", "hello world\n", outStream.String())
		}
	})

	t.Run("characters after newline are removed", func(t *testing.T) {
		outStream, errStream := new(bytes.Buffer), new(bytes.Buffer)
		e := &Envdir{outStream: outStream, errStream: errStream}

		status := e.Run([]string{"envdir", "testenv/multiline", "printenv", "FOO"})
		if status != 0 {
			t.Errorf("expected %d, but got %d: out=%s, err=%s", 0, status, outStream, errStream)
		}

		if outStream.String() != "multi\n" {
			t.Errorf("expected %q, but got %q", "multi\n", outStream.String())
		}
	})

	t.Run("filename beginning with '.' are ignored", func(t *testing.T) {
		outStream, errStream := new(bytes.Buffer), new(bytes.Buffer)
		e := &Envdir{outStream: outStream, errStream: errStream}

		status := e.Run([]string{"envdir", "testenv/dotfile", "printenv", ".BAR"})
		if status != 1 { // printenv set exit status to 1 if env not found
			t.Errorf("expected %d, but got %d", 1, status)
		}

		if outStream.String() != "" {
			t.Errorf("expected %q, but got %q", "", outStream.String())
		}
	})

	t.Run("embedded nulls are converted to newlines", func(t *testing.T) {
		outStream, errStream := new(bytes.Buffer), new(bytes.Buffer)
		e := &Envdir{outStream: outStream, errStream: errStream}

		status := e.Run([]string{"envdir", "testenv/null", "printenv", "HELLO"})
		if status != 0 {
			t.Errorf("expected %d, but got %d: out=%s, err=%s", 0, status, outStream, errStream)
		}

		if outStream.String() != "hello\nworld\n" {
			t.Errorf("expected %q, but got %q", "hello\nworld\n", outStream.String())
		}
	})

	t.Run("env file size is 0, so remove that env", func(t *testing.T) {
		outStream, errStream := new(bytes.Buffer), new(bytes.Buffer)
		e := &Envdir{outStream: outStream, errStream: errStream, env: []string{"EMPTY=not empty"}}

		status := e.Run([]string{"envdir", "testenv/empty", "printenv", "EMPTY"})
		if status != 1 {
			t.Errorf("expected %d, but got %d", 1, status)
		}

		// if outStream.String() != "" {
		// t.Errorf("expected %q, but got %q", "", outStream.String())
		// }
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

		status := e.Run([]string{"envdir", "testenv/dir", "printenv", "FOO"})
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
