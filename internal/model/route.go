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
	_ IDer             = RouteItem{}
	_ list.DefaultItem = RouteItem{}
)

func (m *RootScreenModel) SwitchToRoutes() {
	m.name = "routes"

	m.listFn = func(ctx context.Context) ([]list.Item, error) {
		var (
			routes []*kong.Route
			err    error
			client = m.Client
		)

		if *client.FilterService != "" {
			routes, _, err = client.Routes.ListForService(ctx, client.FilterService, nil)
		} else {
			routes, err = client.Routes.ListAll(ctx)
		}

		if err != nil {
			return nil, err
		}

		res := make([]list.Item, len(routes))
		for i := range routes {
			res[i] = RouteItem{routes[i]}
		}

		return res, nil
	}

	m.getFn = func(ctx context.Context, nameOrID string) (any, error) {
		return m.Client.Routes.Get(ctx, &nameOrID)
	}

	m.deleteFn = func(ctx context.Context, nameOrID string) error {
		return m.Client.Routes.Delete(ctx, &nameOrID)
	}

	m.updateFn = func(ctx context.Context, content []byte) error {
		route := kong.Route{}

		err := yaml.Unmarshal(content, &route)
		if err != nil {
			return err
		}

		_, err = m.Client.Routes.Update(ctx, &route)
		if err != nil {
			return err
		}

		return nil
	}

	m.list = list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	m.list.SetFilteringEnabled(false)
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

type RouteItem struct {
	*kong.Route
}

func (ri RouteItem) FilterValue() string {
	return joinStrPtrs(ri.Name, ri.ID) + joinStrPtrs(ri.Paths...)
}

func (ri RouteItem) Title() string {
	return fmt.Sprintf("Route: %s [id=%s]", joinStrPtrs(ri.Name), joinStrPtrs(ri.ID))
}

func (ri RouteItem) Description() string {
	return routeDesc(ri.Route)
}

func (ri RouteItem) GetID() *string {
	return ri.ID
}

func routeDesc(route *kong.Route) string {
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
