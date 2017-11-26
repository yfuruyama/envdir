package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"
	"syscall"
)

const (
	ExitCodeOk    = 0
	ExitCodeFatal = 111
)

type Envdir struct {
	outStream, errStream io.Writer
	env                  []string
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
		return e.fatal(fmt.Sprintf("%s\n", err.Error()))
	}

	for _, fileInfo := range fileInfos {
		if fileInfo.IsDir() {
			return e.fatal(fmt.Sprintf("%s is not a file, but a directory\n", fileInfo.Name()))
		}

		fileName := fileInfo.Name()
		if strings.HasPrefix(fileName, ".") {
			continue
		}

		filePath := path.Join(dir, fileName)
		file, err := os.Open(filePath)
		if err != nil {
			return e.fatal(fmt.Sprintf("%s\n", err.Error()))
		}

		fsize := fileInfo.Size()
		if fsize == 0 {
			for i, elem := range e.env {
				if strings.HasPrefix(elem, fileName+"=") {
					// remove env
					e.env = append(e.env[:i], e.env[i+1:]...)
					break
				}
			}
			continue
		}

		data := make([]byte, fsize)
		n, err := file.Read(data)
		if err != nil {
			return e.fatal(fmt.Sprintf("%s\n", err.Error()))
		}
		if int64(n) != fsize {
			return e.fatal(fmt.Sprintf("invalid file read size, got: %s, expected: %s, \n", n, fsize))
		}

		v := strings.SplitN(string(data), "\n", 2)[0]
		v = strings.Replace(v, "\x00", "\n", -1) // replace NULL character with newline
		v = strings.TrimRight(v, " \t")          // trim trailing space and tab

		e.env = append(e.env, fileName+"="+v)
	}

	cmd := exec.Command(child, childArgs...)
	cmd.Stdout = e.outStream
	cmd.Stderr = e.errStream
	cmd.Env = e.env
	err = cmd.Run()
	if err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				return status.ExitStatus()
			}
		} else {
			return e.fatal(fmt.Sprintf("%s\n", err.Error()))
		}
	}

	return ExitCodeOk
}
