package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"syscall"
)

func fatal(msg string) {
	fmt.Fprint(os.Stderr, msg)
	os.Exit(111)
}

func main() {
	if len(os.Args) < 3 {
		fatal("usage: envdir dir command\n")
	}

	dir := os.Args[1]
	child := os.Args[2]
	childArgs := os.Args[3:]

	fileInfos, err := ioutil.ReadDir(dir)
	if err != nil {
		fatal(fmt.Sprintf("%s\n", err.Error()))
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
			fatal(fmt.Sprintf("%s\n", err.Error()))
		}

		fsize := fileInfo.Size()
		data := make([]byte, fsize)
		n, err := file.Read(data)
		if err != nil {
			fatal(fmt.Sprintf("%s\n", err.Error()))
		}
		if int64(n) != fsize {
			fatal(fmt.Sprintf("invalid file read size, got: %s, expected: %s, \n", n, fsize))
		}
		// TODO: remove until newline

		env = append(env, fileName+"="+string(data))
	}

	cmd := exec.Command(child, childArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = env
	err = cmd.Run()
	if err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				os.Exit(status.ExitStatus())
			}
		} else {
			fatal(fmt.Sprintf("%s\n", err.Error()))
		}
	}
}
