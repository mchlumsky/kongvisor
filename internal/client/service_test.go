//go:build integration

package client_test

import (
	"context"
	"testing"

	"github.com/kong/go-kong/kong"
	"github.com/mchlumsky/kongvisor/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestIntegrationListServices(t *testing.T) {
	assert := assert.New(t)
	gc := config.GatewayConfig{
		URL: "http://localhost:8001",
	}

	kc, err := gc.GetKongClient()
	assert.NoError(err, "Failed to create Kong client")
	assert.NotNil(kc.Client, "Kong client should not be nil")

	t.Run("List Services", func(t *testing.T) {
		actual, err := kc.ListServices(context.Background())
		assert.NoError(err, "Failed to list services")

		assert.NotNil(actual, "Returned services should not be nil")
		assert.Len(actual, 2, "Expected one service to be returned in the test environment")

		expectedServices := []string{"foo", "bar"}
		for _, svc := range actual {
			service := svc.(*kong.Service)
			assert.NotNil(service.ID, "Service ID should not be nil")
			assert.Contains(expectedServices, *service.Name, "Service name should be one of the expected names")
		}
	})

	t.Run("Get Service by name", func(t *testing.T) {
		serviceName := "foo"
		service, err := kc.GetService(context.Background(), serviceName)
		assert.NoError(err, "Failed to get service by name")
		assert.NotNil(service, "Returned service should not be nil")

		svc := service.(*kong.Service)

		assert.Equal(serviceName, *svc.Name, "Expected to get the service with the specified name")
		assert.Equal("/foo", *svc.Path, "Expected service path to be '/foo'")
		assert.True(*svc.Enabled, "Expected service to be enabled")
		assert.Equal("foo-server.dev", *svc.Host, "Expected service host to be 'foo-server.dev'")
		assert.Equal(int(80), *svc.Port, "Expected service port to be 80")
		assert.Equal("http", *svc.Protocol, "Expected service protocol to be 'http'")
	})

	t.Run("Get Service by ID", func(t *testing.T) {
		svcByName, err := kc.GetService(context.Background(), "foo")
		assert.NoError(err, "Failed to get service by name")
		assert.NotNil(svcByName, "Returned service should not be nil")

		svcID := *svcByName.(*kong.Service).ID
		svcByID, err := kc.GetService(context.Background(), svcID)
		assert.NoError(err, "Failed to get service by ID")
		assert.NotNil(svcByID, "Returned service should not be nil")

		svc := svcByID.(*kong.Service)

		assert.Equal(svcID, *svc.ID, "Expected to get the service with the specified ID")
		assert.Equal("/foo", *svc.Path, "Expected service path to be '/foo'")
		assert.True(*svc.Enabled, "Expected service to be enabled")
		assert.Equal("foo-server.dev", *svc.Host, "Expected service host to be 'foo-server.dev'")
		assert.Equal(int(80), *svc.Port, "Expected service port to be 80")
		assert.Equal("http", *svc.Protocol, "Expected service protocol to be 'http'")
	})
}
