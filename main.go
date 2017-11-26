package main

import "os"

func main() {
	envdir := &Envdir{outStream: os.Stdout, errStream: os.Stderr, env: os.Environ()}
	os.Exit(envdir.Run(os.Args))
}
