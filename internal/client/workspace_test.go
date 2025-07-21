//go:build (integration && enterprise)

package client_test

import (
	"context"
	"testing"

	"github.com/kong/go-kong/kong"
	"github.com/mchlumsky/kongvisor/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/goccy/go-yaml"
)

func TestIntegrationListWorkspaces(t *testing.T) {
	assert := assert.New(t)
	gc := config.GatewayConfig{
		URL: "http://localhost:8001",
	}

	kc, err := gc.GetKongClient()
	assert.NoError(err, "Failed to create Kong client")
	assert.NotNil(kc.Client, "Kong client should not be nil")

	var workspace *kong.Workspace

	t.Run("List Workspaces", func(t *testing.T) {
		actual, err := kc.ListWorkspaces(context.Background())
		assert.NoError(err, "Failed to list workspaces")

		assert.NotNil(actual, "Returned workspaces should not be nil")
		assert.Equal(1, len(actual), "Expected one workspace to be returned in the test environment")

		workspace = actual[0].(*kong.Workspace)
		assert.Equal("default", *workspace.Name, "Expected workspace name to be 'default'")
	})

	t.Run("Get Workspace", func(t *testing.T) {
		got, err := kc.GetWorkspace(context.Background(), *workspace.ID)
		assert.NoError(err, "Failed to get workspace by ID")
		assert.Equal(workspace, got.(*kong.Workspace), workspace, "Expected to get the same workspace as listed")
	})

	t.Run("Delete default Workspace", func(t *testing.T) {
		// TODO: Need license to delete the default workspace, so we skip this test
		t.Skip("Requires Kong Enterprise license to delete the default workspace")

		err := kc.DeleteWorkspace(context.Background(), *workspace.ID)
		assert.NoError(err)
	})

	t.Run("Update default Workspace", func(t *testing.T) {
		// TODO: Need license to update the default workspace, so we skip this test
		t.Skip("Requires Kong Enterprise license to update the default workspace")

		workspace.Comment = kong.String("test123")

		content, err := yaml.Marshal(workspace)
		assert.NoError(err, "Failed to marshal workspace to YAML")

		err = kc.UpdateWorkspace(context.Background(), content)
		assert.NoError(err)

		// TODO: check workspace was updated
	})
}
