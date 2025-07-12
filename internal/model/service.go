package model

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/kong/go-kong/kong"
)

var _ IDer = ServiceItem{}

func (m *RootScreenModel) SwitchToServices() { //nolint:dupl
	m.name = "services"

	m.listFn = m.Client.ListServices
	m.toItemFn = func(service any) list.Item {
		return ServiceItem(*service.(*kong.Service)) //nolint:forcetypeassert
	}
	m.getFn = m.Client.GetService
	m.deleteFn = m.Client.DeleteService
	m.updateFn = m.Client.UpdateService

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

type ServiceItem kong.Service

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
