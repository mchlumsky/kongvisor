package client

import (
	"context"

	"github.com/goccy/go-yaml"
	"github.com/kong/go-kong/kong"
)

func (c *Client) GetRoute(ctx context.Context, id string) (interface{}, error) {
	route, err := c.Routes.Get(ctx, &id)
	if err != nil {
		return nil, err
	}

	return route, nil
}

func (c *Client) DeleteRoute(ctx context.Context, id string) error {
	return c.Routes.Delete(ctx, &id)
}

func (c *Client) UpdateRoute(ctx context.Context, content []byte) error {
	route := kong.Route{}

	err := yaml.Unmarshal(content, &route)
	if err != nil {
		return err
	}

	_, err = c.Routes.Update(ctx, &route)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) ListRoutes(ctx context.Context) ([]any, error) {
	var (
		routes []*kong.Route
		err    error
	)

	if *c.FilterService != "" {
		routes, _, err = c.Routes.ListForService(ctx, c.FilterService, nil)
	} else {
		routes, err = c.Routes.ListAll(ctx)
	}

	if err != nil {
		return nil, err
	}

	res := make([]any, len(routes))
	for i := range routes {
		res[i] = routes[i]
	}

	return res, nil
}
