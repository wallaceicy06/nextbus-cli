package main

import (
	"os"

	"github.com/urfave/cli"
	"github.com/wallaceicy06/nextbus-cli/client"
)

func main() {
	app := cli.NewApp()
	app.Name = "nextbus-cli"
	app.Usage = "Retrieve muni arrival time predictions."
	app.Version = "1.0.1"
	app.Description = "An app to get nextbus predictions for the SF Muni."

	var agency string

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "agency, a",
			Value:       "sf-muni",
			Usage:       "agency for prediction lookups",
			Destination: &agency,
		},
	}
	app.Commands = []cli.Command{
		{
			Name:  "agencies",
			Usage: "list all agency tags and names available from nextbus",
			Action: func(c *cli.Context) error {
				return client.New(agency).ListAgencies()
			},
		},
		{
			Name:  "routes",
			Usage: "list all route tags and route names in the specified agency",
			Action: func(c *cli.Context) error {
				return client.New(agency).ListRoutes()
			},
		},
		{
			Name:  "stops",
			Usage: "list all stop tags and names on the specified route",
			Action: func(c *cli.Context) error {
				route := c.Args().First()
				return client.New(agency).ListStops(route)
			},
		},
		{
			Name:  "predictions",
			Usage: "list up to two next predictions at the specified route and stop",
			Action: func(c *cli.Context) error {
				route := c.Args().Get(0)
				stop := c.Args().Get(1)
				return client.New(agency).ListPredictions(route, stop)
			},
		},
	}

	app.Run(os.Args)
}
