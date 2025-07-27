package model

import (
	"context"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/goccy/go-yaml"
	"github.com/kong/go-kong/kong"
)

var (
	_ IDer      = PluginItem{}
	_ list.Item = PluginItem{}
)

func (m *RootScreenModel) SwitchToPlugins() {
	m.name = "plugins"

	m.listFn = func(ctx context.Context) ([]list.Item, error) {
		var (
			plugins []*kong.Plugin
			err     error
			client  = m.Client
		)

		switch {
		case *client.FilterRoute != "":
			plugins, err = client.Plugins.ListAllForRoute(ctx, client.FilterRoute)
		case *client.FilterService != "":
			plugins, err = client.Plugins.ListAllForService(ctx, client.FilterService)
		default:
			var ps []*kong.Plugin

			ps, err = client.Plugins.ListAll(ctx)
			for _, p := range ps {
				if p.Service == nil && p.Route == nil {
					plugins = append(plugins, p)
				}
			}
		}

		if err != nil {
			return nil, err
		}

		res := make([]list.Item, len(plugins))
		for i := range plugins {
			res[i] = PluginItem{plugins[i]}
		}

		return res, nil
	}

	m.getFn = func(ctx context.Context, nameOrID string) (any, error) {
		return m.Client.Plugins.Get(ctx, &nameOrID)
	}

	m.deleteFn = func(ctx context.Context, nameOrID string) error {
		return m.Client.Plugins.Delete(ctx, &nameOrID)
	}

	m.updateFn = func(ctx context.Context, content []byte) error {
		plugin := kong.Plugin{}

		err := yaml.Unmarshal(content, &plugin)
		if err != nil {
			return err
		}

		_, err = m.Client.Plugins.Update(ctx, &plugin)

		return err
	}

	m.list = list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	m.list.Title = "Plugins"
	m.list.SetStatusBarItemName("plugin", "plugins")

	keymap := NewPluginKeyMap()

	m.list.AdditionalShortHelpKeys = func() []key.Binding {
		return keymap
	}
	m.list.AdditionalFullHelpKeys = func() []key.Binding {
		return keymap
	}
}

type PluginItem struct {
	*kong.Plugin
}

func (pi PluginItem) FilterValue() string {
	return joinStrPtrs(pi.Name, pi.ID) + joinStrPtrs(pi.Tags...)
}

func (pi PluginItem) Title() string {
	return fmt.Sprintf("Plugin: %s [id=%s]", joinStrPtrs(pi.Name), joinStrPtrs(pi.ID))
}

func (pi PluginItem) Description() string {
	return pluginDesc(pi.Plugin)
}

func (pi PluginItem) GetID() *string {
	return pi.ID
}

func pluginDesc(plugin *kong.Plugin) string {
	tags := make([]string, len(plugin.Tags))
	for i, t := range plugin.Tags {
		tags[i] = *t
	}

	var appliedTo []string

	if plugin.Service != nil {
		appliedTo = []string{"Service"}
	}

	if plugin.Consumer != nil {
		appliedTo = append(appliedTo, "Consumer")
	}

	if plugin.Route != nil {
		appliedTo = append(appliedTo, "Route")
	}

	if len(appliedTo) == 0 {
		appliedTo = []string{"Global"}
	}

	var sb strings.Builder

	sb.WriteString("Enabled: " + ptrToString(plugin.Enabled) + ", ")
	sb.WriteString("Applied to: " + strings.Join(appliedTo, ", ") + ", ")
	sb.WriteString("Tags: [" + strings.Join(tags, ", ") + "]")

	return sb.String()
}
