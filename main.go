package main

import (
	"fmt"
	"os"

	"github.com/dinedal/nextbus"
	"github.com/urfave/cli"
)

const (
	muniAgencyTag = "sf-muni"
)

func main() {
	app := cli.NewApp()
	app.Name = "nextmuni"
	app.Description = "An app to get nextbus predictions for the SF Muni."

	app.Commands = []cli.Command{
		{
			Name:  "routes",
			Usage: "list the routes in the system",
			Action: func(c *cli.Context) error {
				return listRoutes()
			},
		},
		{
			Name:  "stops",
			Usage: "list the stops on the specified route",
			Action: func(c *cli.Context) error {
				route := c.Args().First()
				return listStops(route)
			},
		},
	}

	app.Run(os.Args)

	// c := nextbus.DefaultClient
	// fmt.Println(c)

}

func listRoutes() error {
	nb := nextbus.DefaultClient

	routes, err := nb.GetRouteList(muniAgencyTag)
	if err != nil {
		return fmt.Errorf("error getting routes: %v", err.Error())
	}

	for _, r := range routes {
		fmt.Printf("%-8v%v\n", r.Tag, r.Title)
	}

	return nil
}

func listStops(route string) error {
	if len(route) == 0 {
		return fmt.Errorf("route identifier may not be empty")
	}

	nb := nextbus.DefaultClient

	rtCfgs, err := nb.GetRouteConfig(muniAgencyTag, nextbus.RouteConfigTag(route))
	if err != nil {
		return fmt.Errorf("error getting route config for '%v': %v", route, err)
	}

	if len(rtCfgs) == 0 {
		return fmt.Errorf("invalid route identifier '%v'", route)
	} else if len(rtCfgs) > 1 {
		return fmt.Errorf("non-unique route identifier '%v'", route)
	}

	stops := rtCfgs[0].StopList
	for _, s := range stops {
		fmt.Printf("%-8v%v\n", s.StopId, s.Title)
	}

	return nil
}
