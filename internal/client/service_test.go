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

		assert.Equal("bar", *actual[0].(*kong.Service).Name, "Expected first service name to be 'bar'")
		assert.Equal("foo", *actual[1].(*kong.Service).Name, "Expected first service name to be 'foo'")

		for _, svc := range actual {
			service := svc.(*kong.Service)
			assert.NotNil(service.ID, "Service ID should not be nil")
			assert.NotEmpty(*service.Name, "Service name should not be empty")
		}
	})
}
