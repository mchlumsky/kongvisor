package config

import (
	"fmt"
	"net/http"

	"github.com/kong/go-kong/kong"
	kclient "github.com/mchlumsky/kongvisor/internal/client"
	"github.com/spf13/viper"
)

type Config map[string]GatewayConfig

type GatewayConfig struct {
	URL            string
	KongAdminToken string
}

func (gc GatewayConfig) GetKongClient() (*kclient.Client, error) {
	headers := http.Header{}
	headers.Set("Kong-Admin-Token", gc.KongAdminToken)

	client, err := kong.NewClient(kong.String(gc.URL), kong.HTTPClientWithHeaders(nil, headers))
	if err != nil {
		return nil, fmt.Errorf("error creating kong client: %w", err)
	}

	return &kclient.Client{Client: client}, nil
}

func NewConfig() (Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/.config/kongvisor")
	viper.AddConfigPath("/etc/kongvisor/")

	err := viper.ReadInConfig()
	if err != nil {
		return Config{}, fmt.Errorf("error reading config: %w", err)
	}

	var config Config

	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling config: %w", err)
	}

	return config, nil
}
