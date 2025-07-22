package client

import (
	"context"

	"github.com/goccy/go-yaml"
	"github.com/kong/go-kong/kong"
)

func (c *Client) GetWorkspace(ctx context.Context, nameOrID string) (interface{}, error) {
	workspace, err := c.Workspaces.Get(ctx, &nameOrID)
	if err != nil {
		return nil, err
	}

	return workspace, nil
}

func (c *Client) DeleteWorkspace(ctx context.Context, id string) error {
	return c.Workspaces.Delete(ctx, &id)
}

func (c *Client) UpdateWorkspace(ctx context.Context, content []byte) error {
	workspace := kong.Workspace{}

	err := yaml.Unmarshal(content, &workspace)
	if err != nil {
		return err
	}

	_, err = c.Workspaces.Update(ctx, &workspace)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) ListWorkspaces(ctx context.Context) ([]any, error) {
	savedWks := c.Workspace()

	c.SetWorkspace("") // Can't list workspaces with a workspace set
	defer c.SetWorkspace(savedWks)

	workspaces, err := c.Workspaces.ListAll(ctx)
	if err != nil {
		return nil, err
	}

	res := make([]any, len(workspaces))
	for i := range workspaces {
		res[i] = workspaces[i]
	}

	return res, nil
}
