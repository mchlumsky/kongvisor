//go:build (integration && enterprise)

package model

import (
	"context"
	"testing"

	"github.com/kong/go-kong/kong"
	"github.com/mchlumsky/kongvisor/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntegrationRootModel(t *testing.T) {
	t.Run("", func(t *testing.T) {
		assert := assert.New(t)
		require := require.New(t)

		gc := config.GatewayConfig{
			URL: "http://localhost:8001",
		}

		kc, err := gc.GetKongClient()
		require.NoError(err, "Failed to create Kong client")
		assert.NotNil(kc.Client, "Kong client should not be nil")

		model, err := InitModel(kc)
		require.NoError(err, "Failed to initialize model")

		model.SwitchToWorkspaces()

		assert.Equal("workspaces", model.name, "Model should be in workspaces mode")

		assert.NotNil(model.listFn, "List function should not be nil")
		items, err := model.listFn(context.Background())
		require.NoError(err, "List function should not return an error")
		assert.NotEmpty(items, "List function should return items")

		assert.NotNil(model.getFn, "Get function should not be nil")
		wks, err := model.getFn(context.Background(), "default")
		workspace, ok := wks.(*kong.Workspace)
		require.True(ok, "Get function should return a Workspace type")

		require.NoError(err, "Get function should not return an error")
		assert.Equal("default", *workspace.Name, "Get function should return the default workspace")

		assert.NotNil(model.deleteFn, "Delete function should not be nil")
		require.Error(model.deleteFn(context.Background(), "default"), "Delete function should return an error when license missing")


		assert.NotNil(model.updateFn, "Update function should not be nil")
	})
}
