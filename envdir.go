package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"syscall"
)

const (
	ExitCodeOk    = 0
	ExitCodeFatal = 111
)

type Envdir struct {
	outStream, errStream io.Writer
}

func (e *Envdir) fatal(msg string) int {
	fmt.Fprint(e.errStream, msg)
	return ExitCodeFatal
}

func (e *Envdir) Run(args []string) int {
	if len(args) < 3 {
		return e.fatal("usage: envdir dir command\n")
	}

	dir := args[1]
	child := args[2]
	childArgs := args[3:]

	fileInfos, err := ioutil.ReadDir(dir)
	if err != nil {
		e.fatal(fmt.Sprintf("%s\n", err.Error()))
	}

	env := os.Environ()
	for _, fileInfo := range fileInfos {
		if fileInfo.IsDir() {
			continue // TODO
		}
		fileName := fileInfo.Name()
		filePath := path.Join(dir, fileName)

		file, err := os.Open(filePath)
		if err != nil {
			e.fatal(fmt.Sprintf("%s\n", err.Error()))
		}

		fsize := fileInfo.Size()
		data := make([]byte, fsize)
		n, err := file.Read(data)
		if err != nil {
			e.fatal(fmt.Sprintf("%s\n", err.Error()))
		}
		if int64(n) != fsize {
			e.fatal(fmt.Sprintf("invalid file read size, got: %s, expected: %s, \n", n, fsize))
		}
		// TODO: remove until newline

		env = append(env, fileName+"="+string(data))
	}

	cmd := exec.Command(child, childArgs...)
	// TODO: how about Stdin?
	cmd.Stdout = e.outStream
	cmd.Stderr = e.errStream
	cmd.Env = env
	err = cmd.Run()
	if err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				return status.ExitStatus()
			}
		} else {
			e.fatal(fmt.Sprintf("%s\n", err.Error()))
		}
	}

	return ExitCodeOk
}
