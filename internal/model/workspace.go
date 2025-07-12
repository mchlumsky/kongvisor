package model

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/kong/go-kong/kong"
)

var _ IDer = WorkspaceItem{}

func (m *RootScreenModel) SwitchToWorkspaces() { //nolint:dupl
	m.name = "workspaces"

	m.listFn = m.Client.ListWorkspaces
	m.toItemFn = func(workspace any) list.Item {
		return WorkspaceItem(*workspace.(*kong.Workspace)) //nolint:forcetypeassert
	}
	m.getFn = m.Client.GetWorkspace
	m.deleteFn = m.Client.DeleteWorkspace
	m.updateFn = m.Client.UpdateWorkspace

	m.list = list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	m.list.Title = "Workspaces"
	m.list.SetStatusBarItemName("workspace", "workspaces")

	keymap := NewWorkspaceKeyMap()

	m.list.AdditionalShortHelpKeys = func() []key.Binding {
		return keymap
	}
	m.list.AdditionalFullHelpKeys = func() []key.Binding {
		return keymap
	}
}

type WorkspaceItem kong.Workspace

func (wi WorkspaceItem) FilterValue() string {
	return joinStrPtrs(wi.Name, wi.ID, wi.Comment)
}

func (wi WorkspaceItem) Title() string {
	return fmt.Sprintf("Workspace: %s [id=%s]", joinStrPtrs(wi.Name), joinStrPtrs(wi.ID))
}

func (wi WorkspaceItem) Description() string {
	return "Comment: " + joinStrPtrs(wi.Comment)
}

func (wi WorkspaceItem) GetID() *string {
	return wi.ID
}

func joinStrPtrs(ptrs ...*string) string {
	var sb strings.Builder
	for _, ptr := range ptrs {
		sb.WriteString(unNil(ptr))
	}

	return sb.String()
}

func unNil[T any](v *T) T {
	var empty T

	if v != nil {
		return *v
	}

	return empty
}
