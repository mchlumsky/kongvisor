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

func NewWorkspaceKeyMap() []key.Binding {
	return []key.Binding{
		NewServicesBinding(),
		NewRoutesBinding(),
		NewPluginsBinding(),
	}
}

func NewServiceKeyMap() []key.Binding {
	return []key.Binding{
		NewWorkspacesBinding(),
		NewRoutesBinding(),
		NewPluginsBinding(),
	}
}

func NewRouteKeyMap() []key.Binding {
	return []key.Binding{
		NewWorkspacesBinding(),
		NewServicesBinding(),
		NewPluginsBinding(),
	}
}

func NewPluginKeyMap() []key.Binding {
	return []key.Binding{
		NewWorkspacesBinding(),
		NewServicesBinding(),
		NewRoutesBinding(),
	}
}
