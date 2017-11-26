package main

import "os"

func main() {
	envdir := &Envdir{outStream: os.Stdout, errStream: os.Stderr}
	os.Exit(envdir.Run(os.Args))
}
