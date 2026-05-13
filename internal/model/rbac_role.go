package model

import (
	"context"
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/goccy/go-yaml"
	"github.com/kong/go-kong/kong"
)

var (
	_ IDer             = RBACRoleItem{}
	_ list.DefaultItem = RBACRoleItem{}
)

func (m *RootScreenModel) SwitchToRBACRoles() {
	m.name = "roles"

	m.listFn = func(ctx context.Context) ([]list.Item, error) {
		client := m.Client

		savedWks := client.Workspace()

		// TODO: do we really need this for roles?
		client.SetWorkspace("") // Can't list roles with a workspace set
		defer client.SetWorkspace(savedWks)

		roles, err := client.RBACRoles.ListAll(ctx)
		if err != nil {
			return nil, err
		}

		res := make([]list.Item, len(roles))
		for i := range roles {
			res[i] = RBACRoleItem{roles[i]}
		}

		return res, nil
	}

	m.getFn = func(ctx context.Context, nameOrID string) (any, error) {
		return m.Client.RBACRoles.Get(ctx, &nameOrID)
	}

	m.deleteFn = func(ctx context.Context, nameOrID string) error {
		return m.Client.RBACRoles.Delete(ctx, &nameOrID)
	}

	m.updateFn = func(ctx context.Context, content []byte) error {
		role := kong.RBACRole{}

		err := yaml.Unmarshal(content, &role)
		if err != nil {
			return err
		}

		_, err = m.Client.RBACRoles.Update(ctx, &role)
		if err != nil {
			return err
		}

		return nil
	}

	m.list = list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	m.list.Title = "RBAC Roles"
	m.list.SetStatusBarItemName("role", "roles")

	keymap := NewRBACRoleKeyMap()

	m.list.AdditionalShortHelpKeys = func() []key.Binding {
		return keymap
	}
	m.list.AdditionalFullHelpKeys = func() []key.Binding {
		return keymap
	}
}

type RBACRoleItem struct {
	*kong.RBACRole
}

func (ri RBACRoleItem) FilterValue() string {
	return joinStrPtrs(ri.Name, ri.ID, ri.Comment)
}

func (ri RBACRoleItem) Title() string {
	return fmt.Sprintf("Role: %s [id=%s]", joinStrPtrs(ri.Name), joinStrPtrs(ri.ID))
}

func (ri RBACRoleItem) Description() string {
	return "Comment: " + joinStrPtrs(ri.Comment)
}

func (ri RBACRoleItem) GetID() *string {
	return ri.ID
}
