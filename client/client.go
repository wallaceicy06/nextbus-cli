package client

import (
	"fmt"
	"os"
	"strconv"
	"text/tabwriter"

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

// ListAgencies lists the all agencies served by nextbus.
func (c *Client) ListAgencies() error {
	agencies, err := c.nb.GetAgencyList()
	if err != nil {
		return fmt.Errorf("error getting agencies: %v", err.Error())
	}

	format := "%v\t%v\n"
	tw := tabwriter.NewWriter(os.Stdout, 0, 8, 2, ' ', 0)
	fmt.Fprintf(tw, format, "Tag", "Title")
	fmt.Fprintf(tw, format, "---", "-----")

	for _, a := range agencies {
		fmt.Fprintf(tw, format, a.Tag, a.Title)
	}

	tw.Flush()

	return nil
}

// ListRoutes lists the all routes in the system.
func (c *Client) ListRoutes() error {
	routes, err := c.nb.GetRouteList(c.agency)
	if err != nil {
		return fmt.Errorf("error getting routes: %v", err.Error())
	}

	format := "%v\t%v\n"
	tw := tabwriter.NewWriter(os.Stdout, 0, 8, 2, ' ', 0)
	fmt.Fprintf(tw, format, "Tag", "Title")
	fmt.Fprintf(tw, format, "---", "-----")

	for _, r := range routes {
		fmt.Fprintf(tw, format, r.Tag, r.Title)
	}

	tw.Flush()

	return nil
}

// ListStops lists the stops on the specified route.
func (c *Client) ListStops(route string) error {
	if len(route) == 0 {
		return fmt.Errorf("route identifier may not be empty")
	}

	rtCfgs, err := c.nb.GetRouteConfig(c.agency, nextbus.RouteConfigTag(route))
	if err != nil {
		return fmt.Errorf("error getting route config for %q: %v", route, err)
	}

	if len(rtCfgs) == 0 {
		return fmt.Errorf("invalid route identifier %q", route)
	} else if len(rtCfgs) > 1 {
		return fmt.Errorf("non-unique route identifier %q", route)
	}

	stops := rtCfgs[0].StopList

	format := "%v\t%v\n"
	tw := tabwriter.NewWriter(os.Stdout, 0, 8, 2, ' ', 0)
	fmt.Fprintf(tw, format, "Tag", "Title")
	fmt.Fprintf(tw, format, "---", "-----")

	for _, s := range stops {
		fmt.Fprintf(tw, format, s.Tag, s.Title)
	}

	tw.Flush()

	return nil
}

// ListPredictions lists the predictions for service for the specified route and stop.
// Predictions with no prediction less than bound will be omitted.
func (c *Client) ListPredictions(route string, stop string, bound int) error {
	preds, err := c.nb.GetPredictions(c.agency, route, stop)
	if err != nil {
		return fmt.Errorf("error getting predictions for route %q at stop %q: %v",
			route, stop, err.Error())
	}

	if len(preds) == 0 {
		return fmt.Errorf("invalid route %q or stop identifier %q", route, stop)
	} else if len(preds) > 1 {
		return fmt.Errorf("invalid route %q and stop identifier %q", route, stop)
	}

	format := "%v\t%v\n"
	tw := tabwriter.NewWriter(os.Stdout, 0, 8, 2, ' ', 0)
	fmt.Fprintf(tw, format, "Route", "Next Arrivals")
	fmt.Fprintf(tw, format, "-----", "-------------")

	pred := preds[0].PredictionDirectionList
	for _, dir := range pred {
		if len(dir.PredictionList) == 0 {
			continue
		}

		// If the first prediction is greater than the bound, then skip this
		// prediction and move on to the next one.
		if firstPred, err := strconv.Atoi(dir.PredictionList[0].Minutes); err != nil {
			return fmt.Errorf("non-numerical minute prediction received: %v", err)
		} else if firstPred > bound {
			continue
		}

		if len(dir.PredictionList) == 1 {
			fmt.Fprintf(tw, format, dir.Title,
				fmt.Sprintf("%s mins", dir.PredictionList[1].Minutes))
		} else {
			fmt.Fprintf(tw, format, dir.Title,
				fmt.Sprintf("%s & %s mins", dir.PredictionList[0].Minutes, dir.PredictionList[1].Minutes))
		}
	}

	tw.Flush()

	return nil
}
