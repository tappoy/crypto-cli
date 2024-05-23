package main

import (
	"github.com/tappoy/env"
)

func parse() *option {
	// args
	// 0       1       2   3
	// program command src dst
	args := env.Args
	ret := &option{}

	// command
	if len(args) < 2 {
		ret.command = ""
	} else {
		ret.command = args[1]
	}

	// src
	if len(args) < 3 {
		ret.src = ""
	} else {
		ret.src = args[2]
	}

	// dst
	if len(args) < 4 {
		ret.dst = ""
	} else {
		ret.dst = args[3]
	}

	return ret
}
