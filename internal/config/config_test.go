package config_test

import (
	"os"
	"reflect"
	"testing"

	"github.com/mchlumsky/kongvisor/internal/config"
)

func TestNewConfig(t *testing.T) { //nolint:paralleltest
	tests := []struct {
		name          string
		configContent []byte
		want          config.Config
		wantErr       bool
	}{
		{
			"happy-path",
			[]byte(`foobar:
  url: http://foobar
  kongAdminToken: foobar`),
			config.Config{
				"foobar": config.GatewayConfig{
					URL:            "http://foobar",
					KongAdminToken: "foobar",
				},
			},
			false,
		},
		{
			"empty-config",
			[]byte{},
			config.Config{},
			false,
		},
	}

	for _, tt := range tests { //nolint:paralleltest
		t.Run(tt.name, func(t *testing.T) {
			_ = os.WriteFile("config.yaml", tt.configContent, 0o600)

			defer os.Remove("config.yaml")

			got, err := config.NewConfig()
			if (err != nil) != tt.wantErr {
				t.Errorf("NewConfig() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewConfig() got = %v, want %v", got, tt.want)
			}
		})
	}
}
