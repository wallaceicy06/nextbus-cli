package main

import (
	"os"

	"github.com/urfave/cli"
	"github.com/wallaceicy06/nextmuni/client"
)

const (
	muniAgencyTag = "sf-muni"
)

func main() {
	app := cli.NewApp()
	app.Name = "nextmuni"
	app.Description = "An app to get nextbus predictions for the SF Muni."

	muniClient := client.New(muniAgencyTag)

	app.Commands = []cli.Command{
		{
			Name:  "routes",
			Usage: "list the routes in the system",
			Action: func(c *cli.Context) error {
				return muniClient.ListRoutes()
			},
		},
		{
			Name:  "stops",
			Usage: "list the stops on the specified route",
			Action: func(c *cli.Context) error {
				route := c.Args().First()
				return muniClient.ListStops(route)
			},
		},
	}

	app.Run(os.Args)
}
