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
	app.Version = "1.1.0"
	app.Description = "An app to get nextbus predictions for the SF Muni."

	var agency string
	var bound int
	var route string
	var stopTag string
	var stopID string

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "agency, a",
			Value:       "sf-muni",
			Usage:       "agency for prediction lookups",
			Destination: &agency,
		},
		cli.IntFlag{
			Name:        "bound, b",
			Value:       30,
			Usage:       "Prediction times greater than this limit will be omitted",
			Destination: &bound,
		},
		cli.StringFlag{
			Name:        "route, r",
			Value:       "",
			Usage:       "specified route for certain commands",
			Destination: &route,
		},
		cli.StringFlag{
			Name:        "stop_tag, t",
			Value:       "",
			Usage:       "specified stop tag for certain commands",
			Destination: &stopTag,
		}, cli.StringFlag{
			Name:        "stop_id, i",
			Value:       "",
			Usage:       "specified stop id for certain commands",
			Destination: &stopID,
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
				return client.New(agency).ListStops(route)
			},
		},
		{
			Name:  "predictions",
			Usage: "list up to two next predictions at the specified stop (and route, if specified)",
			Action: func(c *cli.Context) error {
				if route == "" {
					return client.New(agency).ListStopPredictions(stopID, bound)
				}
				return client.New(agency).ListPredictions(route, stopTag, bound)
			},
		},
	}

	app.Run(os.Args)
}
