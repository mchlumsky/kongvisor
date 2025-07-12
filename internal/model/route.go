package model

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/kong/go-kong/kong"
)

var _ IDer = RouteItem{}

func (m *RootScreenModel) SwitchToRoutes() { //nolint:dupl
	m.name = "routes"

	m.listFn = m.Client.ListRoutes
	m.toItemFn = func(route any) list.Item {
		return RouteItem(*route.(*kong.Route)) //nolint:forcetypeassert
	}
	m.getFn = m.Client.GetRoute
	m.deleteFn = m.Client.DeleteRoute
	m.updateFn = m.Client.UpdateRoute

	m.list = list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	m.list.Title = "Routes"
	m.list.SetStatusBarItemName("route", "routes")

	keymap := NewRouteKeyMap()

	m.list.AdditionalShortHelpKeys = func() []key.Binding {
		return keymap
	}
	m.list.AdditionalFullHelpKeys = func() []key.Binding {
		return keymap
	}
}

type RouteItem kong.Route

func (ri RouteItem) FilterValue() string {
	return joinStrPtrs(ri.Name, ri.ID) + joinStrPtrs(ri.Paths...)
}

func (ri RouteItem) Title() string {
	return fmt.Sprintf("Route: %s [id=%s]", joinStrPtrs(ri.Name), joinStrPtrs(ri.ID))
}

func (ri RouteItem) Description() string {
	return routeDesc(kong.Route(ri))
}

func (ri RouteItem) GetID() *string {
	return ri.ID
}

func routeDesc(route kong.Route) string {
	protocols := make([]string, len(route.Protocols))
	for i, p := range route.Protocols {
		protocols[i] = *p
	}

	methods := make([]string, len(route.Methods))
	for i, m := range route.Methods {
		methods[i] = *m
	}

	hosts := make([]string, len(route.Hosts))
	for i, h := range route.Hosts {
		hosts[i] = *h
	}

	paths := make([]string, len(route.Paths))
	for i, p := range route.Paths {
		paths[i] = *p
	}

	tags := make([]string, len(route.Tags))
	for i, t := range route.Tags {
		tags[i] = *t
	}

	var sb strings.Builder

	sb.WriteString("Proto: " + strings.Join(protocols, ", ") + ", ")
	sb.WriteString("Methods: " + strings.Join(methods, ",") + ", ")
	sb.WriteString("Hosts: " + strings.Join(hosts, ", ") + ", ")
	sb.WriteString("Paths: " + strings.Join(paths, ", ") + ", ")
	sb.WriteString("Tags: [" + strings.Join(tags, ", ") + "]")

	return sb.String()
}
