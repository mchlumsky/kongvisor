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

type ServiceTestSuite struct {
	suite.Suite
	model *RootScreenModel
}

func (suite *ServiceTestSuite) SetupTest() {
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

func (suite *ServiceTestSuite) TestServices() {
	require := suite.Require()

	suite.model.SwitchToServices()

	suite.Equal("services", suite.model.name, "Model should be in services mode")

	suite.NotNil(suite.model.listFn, "List function should not be nil")
	suite.NotNil(suite.model.getFn, "Get function should not be nil")
	suite.NotNil(suite.model.deleteFn, "Delete function should not be nil")
	suite.NotNil(suite.model.updateFn, "Update function should not be nil")

	items, err := suite.model.listFn(context.Background())
	require.NoError(err, "List function should not return an error")
	suite.Equal(2, len(items), "List function should return exactly 2 items")
	suite.Equal(map[string]struct{}{"foo": {}, "bar": {}}, extractServiceItemNames(items))

	svc, err := suite.model.getFn(context.Background(), "foo")
	require.NoError(err, "Get function should not return an error")
	ksvc, ok := svc.(*kong.Service)
	suite.True(ok, "Get function should return a *kong.Service type")
	suite.Equal("foo", *ksvc.Name, "Get function should return the correct service name")
	suite.Equal("/foo", *ksvc.Path, "Get function should return the correct service path")
	suite.True(*ksvc.Enabled, "Get function should return the service as enabled")
	suite.Equal("foo-server.dev", *ksvc.Host, "Get function should return the correct service host")
	suite.Equal(80, *ksvc.Port, "Get function should return the correct service port")
	suite.Equal("http", *ksvc.Protocol, "Get function should return the correct service protocol")

	require.Error(suite.model.deleteFn(context.Background(), ""), "Delete function should return an error for empty ID")
}

func TestIntegrationServiceTestSuite(t *testing.T) {
	t.Run("service integration test", func(t *testing.T) {
		suite.Run(t, new(ServiceTestSuite))
	})
}

func extractServiceItemNames(items []list.Item) map[string]struct{} {
	names := make(map[string]struct{}, len(items))
	for _, item := range items {
		if routeItem, ok := item.(ServiceItem); ok && routeItem.Name != nil {
			names[*routeItem.Name] = struct{}{}
		}
	}

	return names
}
