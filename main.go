package main

import (
	"os"
	"os/signal"

	"github.com/urfave/cli/v2"
	"reichard.io/libgen-opds/server"
)

func main() {
	app := &cli.App{
		Name:  "LibGen OPDS Bridge",
		Usage: "A Library Genesis OPDS Bridge",
		Commands: []*cli.Command{
			{
				Name:    "serve",
				Aliases: []string{"s"},
				Usage:   "Start Server",
				Action:  cmdServer,
			},
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		// log.Fatal(err)
	}
}

func cmdServer(ctx *cli.Context) error {
	server := server.NewServer()
	server.StartServer()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	server.StopServer()
	os.Exit(0)

	return nil
}
