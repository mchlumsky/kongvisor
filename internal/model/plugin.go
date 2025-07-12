package model

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/kong/go-kong/kong"
)

var _ IDer = PluginItem{}

func (m *RootScreenModel) SwitchToPlugins() { //nolint:dupl
	m.name = "plugins"

	m.listFn = m.Client.ListPlugins
	m.toItemFn = func(plugin any) list.Item {
		return PluginItem(*plugin.(*kong.Plugin)) //nolint:forcetypeassert
	}
	m.getFn = m.Client.GetPlugin
	m.deleteFn = m.Client.DeletePlugin
	m.updateFn = m.Client.UpdatePlugin

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

type PluginItem kong.Plugin

func (pi PluginItem) FilterValue() string {
	return joinStrPtrs(pi.Name, pi.ID) + joinStrPtrs(pi.Tags...)
}

func (pi PluginItem) Title() string {
	return fmt.Sprintf("Plugin: %s [id=%s]", joinStrPtrs(pi.Name), joinStrPtrs(pi.ID))
}

func (pi PluginItem) Description() string {
	return pluginDesc(kong.Plugin(pi))
}

func (pi PluginItem) GetID() *string {
	return pi.ID
}

func pluginDesc(plugin kong.Plugin) string {
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
