//go:build integration && enterprise

package model

import (
	"context"
	"testing"

	"github.com/goccy/go-yaml"
	"github.com/kong/go-kong/kong"
	"github.com/mchlumsky/kongvisor/internal/config"
	"github.com/stretchr/testify/suite"
)

type WorkspaceTestSuite struct {
	suite.Suite
	model *RootScreenModel
}

func (suite *WorkspaceTestSuite) SetupTest() {
	require := suite.Require()

	gc := config.GatewayConfig{
		URL: "http://localhost:8001",
	}

	kc, err := gc.GetKongClient()
	require.NoError(err, "Failed to create Kong client")
	suite.NotNil(kc.Client, "Kong client should not be nil")

	suite.model, err = InitModel(kc)
	require.NoError(err, "Failed to initialize model")
}

func (suite *WorkspaceTestSuite) TestWorkspaces() {
	require := suite.Require()

	suite.model.SwitchToWorkspaces()

	suite.Equal("workspaces", suite.model.name, "Model should be in workspaces mode")

	suite.NotNil(suite.model.listFn, "List function should not be nil")
	suite.NotNil(suite.model.getFn, "Get function should not be nil")
	suite.NotNil(suite.model.deleteFn, "Delete function should not be nil")
	suite.NotNil(suite.model.updateFn, "Update function should not be nil")

	items, err := suite.model.listFn(context.Background())
	require.NoError(err, "List function should not return an error")
	suite.NotEmpty(items, "List function should return items")

	wks, err := suite.model.getFn(context.Background(), "default")
	require.NoError(err, "Get function should not return an error")

	workspace, ok := wks.(*kong.Workspace)
	require.True(ok, "Get function should return a Workspace type")
	suite.Equal("default", *workspace.Name, "Get function should return the default workspace")

	require.Error(suite.model.deleteFn(context.Background(), "default"), "Delete function should return an error when license missing")

	ws, err := yaml.Marshal(kong.Workspace{ID: kong.String("Default"), Name: kong.String("default"), Comment: kong.String("new comment")})
	require.NoError(err, "Failed to marshal workspace to YAML")
	require.Error(suite.model.updateFn(context.Background(), ws), "Update function should return an error when license missing")
}

func TestIntegrationWorkspaceTestSuite(t *testing.T) {
	t.Run("workspace integration test", func(t *testing.T) {
		suite.Run(t, new(WorkspaceTestSuite))
	})
}
