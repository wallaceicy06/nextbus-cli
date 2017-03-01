package client

import (
	"fmt"

	"github.com/dinedal/nextbus"
)

// Client is state for the nextbus client.
type Client struct {
	nb     *nextbus.Client
	agency string
}

func New(agency string) *Client {
	return &Client{
		nb:     nextbus.DefaultClient,
		agency: agency,
	}
}

// ListRoutes lists the all routes in the system.
func (c *Client) ListRoutes() error {
	routes, err := c.nb.GetRouteList(c.agency)
	if err != nil {
		return fmt.Errorf("error getting routes: %v", err.Error())
	}

	for _, r := range routes {
		fmt.Printf("%-8v%v\n", r.Tag, r.Title)
	}

	return nil
}

// ListStops lists the stops on the specified route.
func (c *Client) ListStops(route string) error {
	if len(route) == 0 {
		return fmt.Errorf("route identifier may not be empty")
	}

	rtCfgs, err := c.nb.GetRouteConfig(c.agency, nextbus.RouteConfigTag(route))
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
