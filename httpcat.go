package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"os"
)

func before(c *cli.Context) error {
	if c.GlobalBool("debug") {
		log.SetLevel(log.DebugLevel)
	}
	return nil
}

func action(c *cli.Context) error {
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
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
