package main

import (
	"bufio"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"os"
	"strings"
)

const CRLF = "\r\n"
const EOF = ""

type Sent struct {
	RequestLine bool
	Headers     bool
	Verbose     bool
}

func (s *Sent) log(line string) {
	if s.Verbose {
		fmt.Fprintln(os.Stderr, ">", line)
	}
}

func (s *Sent) writeLine(line string) bool {
	if !s.RequestLine {
		s.RequestLine = true
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "/") {
			line = "GET " + line
		}
		if !strings.Contains(line, "HTTP/") {
			line += " HTTP/1.1"
		}

		return s.writeLine(line)
	} else if !s.Headers {
		if !strings.HasSuffix(line, CRLF) {
			line += CRLF
		}
	} else if line == EOF {
		return false
	}

	fmt.Print(line)

	if s.Headers {
		s.log(line)
	} else {
		s.log(fmt.Sprintf("%q", line))
		if line == CRLF {
			s.Headers = true
		}
	}

	return true
}

func before(c *cli.Context) error {
	if c.GlobalBool("debug") {
		log.SetLevel(log.DebugLevel)
	}
	return nil
}

func action(c *cli.Context) error {
	verbose := c.GlobalBool("verbose")
	sent := Sent{
		Verbose: verbose,
	}

	lines := c.Args()
	for _, line := range lines {
		sent.writeLine(line)
	}

	noStdin := c.GlobalBool("no-stdin")
	if !noStdin {
		s := bufio.NewScanner(os.Stdin)

		for s.Scan() {
			sent.writeLine(s.Text())
		}
	}

	return nil
}

func main() {
	app := cli.NewApp()
	app.Name = "httpcat"
	app.Usage = "create raw HTTP requests on the command line"
	app.Version = "0.0.1"
	app.Before = before
	app.Action = action

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "debug, d",
			Usage: "enable debug mode",
		},
		cli.BoolFlag{
			Name:  "verbose, V",
			Usage: "print info about output lines to stderr",
		},
		cli.BoolFlag{
			Name:  "no-stdin, n",
			Usage: "disable reading of lines from stdin",
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
