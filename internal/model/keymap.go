package model

import "github.com/charmbracelet/bubbles/key"

func NewWorkspacesBinding() key.Binding {
	return key.NewBinding(
		key.WithKeys("w"),
		key.WithHelp("w", "workspaces"),
	)
}

func NewServicesBinding() key.Binding {
	return key.NewBinding(
		key.WithKeys("s"),
		key.WithHelp("s", "services"),
	)
}

func NewRoutesBinding() key.Binding {
	return key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "routes"),
	)
}

func NewPluginsBinding() key.Binding {
	return key.NewBinding(
		key.WithKeys("p"),
		key.WithHelp("p", "plugins"),
	)
}

func NewYAMLViewBinding() key.Binding {
	return key.NewBinding(
		key.WithKeys("y"),
		key.WithHelp("y", "yaml view"),
	)
}

func NewWorkspaceKeyMap() []key.Binding {
	return []key.Binding{
		NewServicesBinding(),
		NewRoutesBinding(),
		NewPluginsBinding(),
		NewYAMLViewBinding(),
	}
}

func NewServiceKeyMap() []key.Binding {
	return []key.Binding{
		NewWorkspacesBinding(),
		NewRoutesBinding(),
		NewPluginsBinding(),
		NewYAMLViewBinding(),
	}
}

func NewRouteKeyMap() []key.Binding {
	return []key.Binding{
		NewWorkspacesBinding(),
		NewServicesBinding(),
		NewPluginsBinding(),
		NewYAMLViewBinding(),
	}
}

func NewPluginKeyMap() []key.Binding {
	return []key.Binding{
		NewWorkspacesBinding(),
		NewServicesBinding(),
		NewRoutesBinding(),
		NewYAMLViewBinding(),
	}
}
