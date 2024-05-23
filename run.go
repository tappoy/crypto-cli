package main

import (
	"github.com/tappoy/env"
	ver "github.com/tappoy/version"

	"github.com/tappoy/crypto"

	"io"
	"os"
)

func (o *option) run() int {
	switch o.command {
	case "help":
		return o.usage()
	case "version":
		return o.version()
	case "xc":
		if !o.setStdio() {
			return 1
		}
		return o.withInOut(runXc)
	case "c":
		if !o.setStdio() {
			return 1
		}
		return o.withInOut(runC)
	default:
		env.Errf("Unknown command: %s\n", o.command)
		runHelpMessage()
		return 1
	}
}

// print usage
func (o *option) usage() int {
	env.Outf(usage)
	return 0
}

// print version
func (o *option) version() int {
	env.Outf("crypto-cli version %s\n", ver.Version())
	return 0
}

func (o *option) setStdio() bool {
	if o.src == "" || o.dst == "" {
		env.Errf("src and dst are required.\n")
		return false
	}

	if o.src == "-" {
		env.Errf("src cannot be stdin.\n")
		return false
	}

	if o.dst == "-" {
		o.out = env.Out
		o.isStdout = true
	}

	return true
}

func (o *option) withInOut(f func(r io.Reader, w io.Writer, c *crypto.Crypto) int) int {
	var out io.Writer
	var in io.Reader

	// in
	stat, err := os.Stat(o.src)
	if err != nil {
		env.Errf("Failed to get stat: %s\n", err)
		return 1
	} else {
		if stat.IsDir() {
			env.Errf("src is directory.\n")
			return 1
		}

		f, err := os.Open(o.src)
		if err != nil {
			env.Errf("Failed to open file: %s\n", err)
			return 1
		}
		defer f.Close()
		in = f
	}

	// password
	env.Errf("Password: ")
	password, err := env.InputPassword()
	if err != nil {
		env.Errf("Failed to read password: %s\n", err)
		return 1
	}

	// crypto
	c, err := crypto.NewCrypto(password)
	if err != nil {
		env.Errf("Failed to create crypto: %s\n", err)
		return 1
	}

	// out
	isCreated := false
	if o.isStdout {
		out = env.Out
	} else {
		_, err := os.Stat(o.dst)
		if err != nil {
			f, err := os.Create(o.dst)
			if err != nil {
				env.Errf("Failed to create file: %s\n", err)
				return 1
			}
			defer f.Close()
			isCreated = true
			out = f
		} else {
			env.Errf("File already exists.\n")
			return 1
		}
	}

	ret := f(in, out, c)
	if ret != 0 {
		if isCreated {
			os.Remove(o.dst)
		}
	}

	return ret
}

func runXc(r io.Reader, w io.Writer, c *crypto.Crypto) int {
	cr, err := c.Reader(r)
	if err != nil {
		env.Errf("Failed to create reader: %s\n", err)
		return 1
	}

	_, err = io.Copy(w, cr)
	if err != nil {
		env.Errf("Failed to decrypt: %s\n", err)
		return 1
	}

	return 0
}

func runC(r io.Reader, w io.Writer, c *crypto.Crypto) int {
	cw, err := c.Writer(w)
	if err != nil {
		env.Errf("Failed to create writer: %s\n", err)
		return 1
	}

	_, err = io.Copy(cw, r)
	if err != nil {
		env.Errf("Failed to encrypt: %s\n", err)
		return 1
	}

	return 0
}
