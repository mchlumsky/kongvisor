//go:build integration

package client_test

import (
	"os"
	"strings"
	"testing"

	"github.com/mchlumsky/kongvisor/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestIntegrationKongVersion(t *testing.T) {
	assert := assert.New(t)

	gc := config.GatewayConfig{
		URL: "http://localhost:8001",
	}

	kc, err := gc.GetKongClient()
	assert.NoError(err, "Failed to create Kong client")
	assert.NotNil(kc.Client, "Kong client should not be nil")

	actualVersion, err := kc.KongVersion()
	assert.NoError(err, "Failed to get Kong version")

	expectedVersion, ok := os.LookupEnv("TEST_KONG_VERSION")
	assert.True(ok, "TEST_KONG_VERSION environment variable must be set for this test to pass")

	assert.Truef(strings.HasPrefix(actualVersion.String(), expectedVersion+"."),
		"Returned Kong version `%s` does not match expected prefix `%s.`", actualVersion, expectedVersion)
}
