package main

import (
	"PhoenixOracle/gophoenix/core/logger"
	"PhoenixOracle/gophoenix/core/services"
	"PhoenixOracle/gophoenix/core/store"
	"PhoenixOracle/gophoenix/core/web"
	"gopkg.in/urfave/cli.v1"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Usage = "CLI for Chainlink"
	app.Commands = []cli.Command{
		{
			Name:    "node",
			Aliases: []string{"n"},
			Usage:   "Run the chainlink node",
			Action:  runNode,
		},
	}
	app.Run(os.Args)
}


func runNode(c *cli.Context) error {
	cl := services.NewApplication(store.NewConfig())
	services.Authenticate(cl.Store)
	r := web.Router(cl)

	if err := cl.Start(); err != nil {
		logger.Fatal(err)
	}
	defer cl.Stop()
	logger.Fatal(r.Run())
	return nil
}