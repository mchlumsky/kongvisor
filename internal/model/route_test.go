//go:build integration

package model

import (
	"context"
	"testing"

	"github.com/charmbracelet/bubbles/list"
	"github.com/kong/go-kong/kong"
	"github.com/mchlumsky/kongvisor/internal/config"
	"github.com/stretchr/testify/suite"
)

type RouteTestSuite struct {
	suite.Suite
	model *RootScreenModel
}

func (suite *RouteTestSuite) SetupTest() {
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

func (suite *RouteTestSuite) TestRoutes() {
	require := suite.Require()

	suite.model.SwitchToRoutes()
	suite.Equal("routes", suite.model.name, "Model should be in routes mode")

	suite.NotNil(suite.model.listFn, "List function should not be nil")
	suite.NotNil(suite.model.getFn, "Get function should not be nil")
	suite.NotNil(suite.model.deleteFn, "Delete function should not be nil")
	suite.NotNil(suite.model.updateFn, "Update function should not be nil")

	// list all routes in default workspace
	items, err := suite.model.listFn(context.Background())
	require.NoError(err, "List function should not return an error")
	suite.Equal(3, len(items), "List function should return exactly 3 items")
	suite.Equal(
		map[string]struct{}{"firstFoo": {}, "secondFoo": {}, "bar": {}},
		extractRouteItemNames(items),
		"List function should return expected route names",
	)

	// list all routes in foo service
	suite.model.FilterService = "foo"
	items, err = suite.model.listFn(context.Background())
	require.NoError(err, "List function should not return an error")
	suite.Equal(2, len(items), "List function should return exactly 2 items")
	suite.Equal(
		map[string]struct{}{"firstFoo": {}, "secondFoo": {}},
		extractRouteItemNames(items),
		"List function should return expected route names",
	)

	// get route by name
	route, err := suite.model.getFn(context.Background(), "firstFoo")
	require.NoError(err, "Get function should not return an error")

	kroute := route.(*kong.Route)
	suite.Equal("firstFoo", *kroute.Name, "Get function should return the correct route name")

	// delete a route
	require.Error(suite.model.deleteFn(context.Background(), ""), "Delete function should return an error for empty name")
}

func TestIntegrationRouteTestSuite(t *testing.T) {
	t.Run("route integration test", func(t *testing.T) {
		suite.Run(t, new(RouteTestSuite))
	})
}

func extractRouteItemNames(items []list.Item) map[string]struct{} {
	names := make(map[string]struct{}, len(items))
	for _, item := range items {
		if routeItem, ok := item.(RouteItem); ok && routeItem.Name != nil {
			names[*routeItem.Name] = struct{}{}
		}
	}

	return names
}
