package client

import (
	"context"
	"time"

	"github.com/kong/go-kong/kong"
)

type Client struct {
	*kong.Client
	FilterService *string
	FilterRoute   *string
}

func (c *Client) KongVersion() (kong.Version, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	info, err := c.Info.Get(ctx)
	if err != nil {
		return kong.Version{}, err
	}

	version, err := kong.NewVersion(info.Version)
	if err != nil {
		return kong.Version{}, err
	}

	return version, nil
}
