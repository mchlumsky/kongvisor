package client

import (
	"context"

	"github.com/goccy/go-yaml"
	"github.com/kong/go-kong/kong"
)

func (c *Client) GetService(ctx context.Context, nameOrID string) (interface{}, error) {
	service, err := c.Services.Get(ctx, &nameOrID)
	if err != nil {
		return nil, err
	}

	return service, nil
}

func (c *Client) DeleteService(ctx context.Context, nameOrID string) error {
	return c.Services.Delete(ctx, &nameOrID)
}

func (c *Client) UpdateService(ctx context.Context, content []byte) error {
	service := kong.Service{}

	err := yaml.Unmarshal(content, &service)
	if err != nil {
		return err
	}

	_, err = c.Services.Update(ctx, &service)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) ListServices(ctx context.Context) ([]any, error) {
	services, err := c.Services.ListAll(ctx)
	if err != nil {
		return nil, err
	}

	res := make([]any, len(services))
	for i := range services {
		res[i] = services[i]
	}

	return res, nil
}
