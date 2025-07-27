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
	_ IDer      = WorkspaceItem{}
	_ list.Item = WorkspaceItem{}
)

func (m *RootScreenModel) SwitchToWorkspaces() {
	m.name = "workspaces"

	m.listFn = func(ctx context.Context) ([]list.Item, error) {
		client := m.Client

		savedWks := client.Workspace()

		client.SetWorkspace("") // Can't list workspaces with a workspace set
		defer client.SetWorkspace(savedWks)

		workspaces, err := client.Workspaces.ListAll(ctx)
		if err != nil {
			return nil, err
		}

		res := make([]list.Item, len(workspaces))
		for i := range workspaces {
			res[i] = WorkspaceItem{workspaces[i]}
		}

		return res, nil
	}

	m.getFn = func(ctx context.Context, nameOrID string) (any, error) {
		return m.Client.Workspaces.Get(ctx, &nameOrID)
	}

	m.deleteFn = func(ctx context.Context, nameOrID string) error {
		return m.Client.Workspaces.Delete(ctx, &nameOrID)
	}

	m.updateFn = func(ctx context.Context, content []byte) error {
		workspace := kong.Workspace{}

		err := yaml.Unmarshal(content, &workspace)
		if err != nil {
			return err
		}

		_, err = m.Client.Workspaces.Update(ctx, &workspace)
		if err != nil {
			return err
		}

		return nil
	}

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

type WorkspaceItem struct {
	*kong.Workspace
}

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
