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
	_ IDer      = ServiceItem{}
	_ list.Item = ServiceItem{}
)

func (m *RootScreenModel) SwitchToServices() { //nolint:dupl
	m.name = "services"

	m.listFn = func(ctx context.Context) ([]list.Item, error) {
		services, err := m.Client.Services.ListAll(ctx)
		if err != nil {
			return nil, err
		}

		items := make([]list.Item, len(services))
		for i, service := range services {
			si := ServiceItem{service}
			items[i] = list.Item(si)
		}
		return items, nil
	}

	m.getFn = func(ctx context.Context, nameOrID string) (any, error) {
		return m.Client.Services.Get(ctx, &nameOrID)
	}

	m.deleteFn = func(ctx context.Context, nameOrID string) error {
		return m.Client.Services.Delete(ctx, &nameOrID)
	}

	m.updateFn = func(ctx context.Context, content []byte) error {
		service := kong.Service{}

		err := yaml.Unmarshal(content, &service)
		if err != nil {
			return err
		}

		_, err = m.Client.Services.Update(ctx, &service)

		return err
	}

	m.list = list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	m.list.Title = "Services"
	m.list.SetStatusBarItemName("service", "services")

	keymap := NewServiceKeyMap()

	m.list.AdditionalShortHelpKeys = func() []key.Binding {
		return keymap
	}
	m.list.AdditionalFullHelpKeys = func() []key.Binding {
		return keymap
	}
}

type ServiceItem struct {
	*kong.Service
}

func (si ServiceItem) FilterValue() string {
	return joinStrPtrs(si.Name, si.ID, si.Path)
}

func (si ServiceItem) Title() string {
	return fmt.Sprintf("Service: %s [id=%s]", joinStrPtrs(si.Name), joinStrPtrs(si.ID))
}

func (si ServiceItem) Description() string {
	tags := make([]string, len(si.Tags))
	for i, t := range si.Tags {
		tags[i] = *t
	}

	return fmt.Sprintf(
		"Enabled: %s, Protocol: %s, Host: %s, Port: %s, Path: %s, Tags: %s",
		ptrToString(si.Enabled),
		ptrToString(si.Protocol),
		ptrToString(si.Host),
		ptrToString(si.Port),
		ptrToString(si.Path),
		"["+strings.Join(tags, ", ")+"]",
	)
}

func (si ServiceItem) GetID() *string {
	return si.ID
}

func ptrToString[T any](v *T) string {
	if v == nil {
		return ""
	}

	return fmt.Sprintf("%v", *v)
}
