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

// New makes a new client for nextbus with the specified agency tag.
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
		fmt.Printf("%-8v%v\n", s.Tag, s.Title)
	}

	return nil
}

// ListPredictions lists the predictions for service for the specified route and stop.
func (c *Client) ListPredictions(route string, stop string) error {
	preds, err := c.nb.GetPredictions(c.agency, route, stop)
	if err != nil {
		return fmt.Errorf("error getting predictions for route '%v' at stop '%v': %v",
			route, stop, err.Error())
	}

	if len(preds) == 0 {
		return fmt.Errorf("invalid route '%v' or stop identifier '%v'", route, stop)
	} else if len(preds) > 1 {
		return fmt.Errorf("invalid route '%v' and stop identifier '%v'", route, stop)
	}

	pred := preds[0].PredictionDirectionList
	for _, dir := range pred {
		if len(dir.PredictionList) == 0 {
			continue
		} else if len(dir.PredictionList) == 1 {
			fmt.Printf("%v: %v mins\n", dir.Title, dir.PredictionList[0].Minutes)
		} else {
			fmt.Printf("%v: %v & %v mins\n", dir.Title,
				dir.PredictionList[0].Minutes, dir.PredictionList[1].Minutes)
		}
	}
	return nil
}
