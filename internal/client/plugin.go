package client

import (
	"context"

	"github.com/goccy/go-yaml"
	"github.com/kong/go-kong/kong"
)

func (c *Client) ListPlugins(ctx context.Context) ([]any, error) {
	var (
		plugins []*kong.Plugin
		err     error
	)

	switch {
	case *c.FilterRoute != "":
		plugins, err = c.Plugins.ListAllForRoute(ctx, c.FilterRoute)
	case *c.FilterService != "":
		plugins, err = c.Plugins.ListAllForService(ctx, c.FilterService)
	default:
		var ps []*kong.Plugin

		ps, err = c.Plugins.ListAll(ctx)
		for _, p := range ps {
			if p.Service == nil && p.Route == nil {
				plugins = append(plugins, p)
			}
		}
	}

	if err != nil {
		return nil, err
	}

	res := make([]any, len(plugins))
	for i := range plugins {
		res[i] = plugins[i]
	}

	return res, nil
}

func (c *Client) GetPlugin(ctx context.Context, usernameOrID string) (interface{}, error) {
	plugin, err := c.Plugins.Get(ctx, &usernameOrID)
	if err != nil {
		return nil, err
	}

	return plugin, err
}

func (c *Client) DeletePlugin(ctx context.Context, id string) error {
	return c.Plugins.Delete(ctx, &id)
}

func (c *Client) UpdatePlugin(ctx context.Context, content []byte) error {
	plugin := kong.Plugin{}

	err := yaml.Unmarshal(content, &plugin)
	if err != nil {
		return err
	}

	_, err = c.Plugins.Update(ctx, &plugin)
	if err != nil {
		return err
	}

	return nil
}
