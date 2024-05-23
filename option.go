package main

import (
	"io"
)

type option struct {
	// parse
	command string
	src     string
	dst     string

	// setSrcAndDst
	in       io.Reader
	out      io.Writer
	isStdout bool
}
