package main

import (
	"os"

	"github.com/urfave/cli"
	"github.com/wallaceicy06/nextbus-cli/client"
)

func main() {
	app := cli.NewApp()
	app.Name = "nextmuni"
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
			Usage: "list the agencies available from nextbus",
			Action: func(c *cli.Context) error {
				return client.New(agency).ListAgencies()
			},
		},
		{
			Name:  "routes",
			Usage: "list the routes in the system",
			Action: func(c *cli.Context) error {
				return client.New(agency).ListRoutes()
			},
		},
		{
			Name:  "stops",
			Usage: "list the stops on the specified route",
			Action: func(c *cli.Context) error {
				route := c.Args().First()
				return client.New(agency).ListStops(route)
			},
		},
		{
			Name:  "predictions",
			Usage: "list the predicitons at the specified route and stop",
			Action: func(c *cli.Context) error {
				route := c.Args().Get(0)
				stop := c.Args().Get(1)
				return client.New(agency).ListPredictions(route, stop)
			},
		},
	}

	app.Run(os.Args)
}
