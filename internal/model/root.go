package model

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/davecgh/go-spew/spew"
	"github.com/goccy/go-yaml"
	"github.com/mchlumsky/kongvisor/internal/client"
)

const (
	period            = 2 * time.Second
	timeout           = 15 * time.Second
	workspaceName     = "workspaces"
	serviceName       = "services"
	routeName         = "routes"
	pluginName        = "plugins"
	EditionEnterprise = "Enterprise"
	EditionOSS        = "OSS"
)

var (
	ErrNilID      = errors.New("error: id for selected item is nil")
	loadingStatus = lipgloss.NewStyle().Blink(true).Faint(true).Render("loading...") //nolint:gochecknoglobals
)

type (
	ListFn     func(context.Context) ([]any, error)
	ToItemFn   func(any) list.Item
	ListItemFn func(context.Context) ([]list.Item, error)
	GetFn      func(context.Context, string) (interface{}, error)
	DeleteFn   func(context.Context, string) error
	UpdateFn   func(context.Context, []byte) error
)

type ItemsMsg struct {
	Items []list.Item
	Name  string
}

type editorFinishedMsg struct {
	filename string
	err      error
}

func openEditor(filePath string) tea.Cmd {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vim"
	}

	c := exec.Command(editor, filePath)

	return tea.ExecProcess(c, func(err error) tea.Msg {
		return editorFinishedMsg{filePath, err}
	})
}

type UnexpectedTypeError struct {
	value any
}

func (e UnexpectedTypeError) Error() string {
	return fmt.Sprintf("unexpected type %T", e.value)
}

func ErrCmd(e error) tea.Cmd {
	return func() tea.Msg {
		return e
	}
}

type RootScreenModel struct {
	Client        *client.Client
	dump          io.Writer
	list          list.Model
	help          help.Model
	viewport      viewport.Model
	viewportHelp  help.Model
	yaml          bool
	baseURL       string
	kongVersion   string
	edition       string
	listItemFn    ListItemFn
	listFn        ListFn
	toItemFn      ToItemFn
	getFn         GetFn
	deleteFn      DeleteFn
	updateFn      UpdateFn
	name          string
	status        string
	FilterService string
	FilterRoute   string
	width         int
	height        int
	err           error
}

func (m *RootScreenModel) Init() tea.Cmd {
	return tickCmd()
}

func (m *RootScreenModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) { //nolint: ireturn
	if m.dump != nil {
		spew.Fdump(m.dump, msg)
	}

	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "W":
			m.Client.SetWorkspace("")
		case "S":
			m.FilterService = ""
		case "R":
			m.FilterRoute = ""
		case "p", "s", "w", "r":
			if msg.String() == "w" {
				if m.edition != EditionEnterprise {
					return m, nil
				}
			}

			m.status = loadingStatus

			switch msg.String() {
			case "p":
				m.SwitchToPlugins()
			case "s":
				m.SwitchToServices()
			case "w":
				m.SwitchToWorkspaces()
			case "r":
				m.SwitchToRoutes()
			}

			return m, nil
		case "y":
			m.yaml = true

			id, err := m.SelectedID()
			if err != nil {
				return m, ErrCmd(err)
			}

			if id == nil {
				return m, ErrCmd(ErrNilID)
			}

			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()

			resource, err := m.getFn(ctx, *id)
			if err != nil {
				return m, ErrCmd(err)
			}

			var buf bytes.Buffer
			enc := yaml.NewEncoder(&buf, yaml.IndentSequence(true))

			err = enc.Encode(resource)
			if err != nil {
				return m, ErrCmd(err)
			}

			m.viewport.SetContent(buf.String())
		case "e":
			temp, err := os.CreateTemp("", "kongvisor")
			if err != nil {
				return m, ErrCmd(err)
			}

			id, err := m.SelectedID()
			if err != nil {
				return m, ErrCmd(err)
			}

			if id == nil {
				return m, ErrCmd(ErrNilID)
			}

			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()

			resource, err := m.getFn(ctx, *id)
			if err != nil {
				return m, ErrCmd(err)
			}

			err = yaml.NewEncoder(temp).Encode(resource)
			if err != nil {
				return m, ErrCmd(err)
			}

			defer closeIgnoringErr(temp)

			return m, openEditor(temp.Name())
		case tea.KeyCtrlD.String():
			id, err := m.SelectedID()
			if err != nil {
				return m, ErrCmd(err)
			}

			if id == nil {
				return m, ErrCmd(ErrNilID)
			}

			cmd = func() tea.Msg {
				ctx, cancel := context.WithTimeout(context.Background(), timeout)
				defer cancel()

				err := m.deleteFn(ctx, *id)
				if err != nil {
					return err
				}

				return nil
			}

			return m, cmd
		case tea.KeyEsc.String():
			if m.yaml {
				m.yaml = false
			}

			return m, cmd
		case tea.KeyDown.String(), "j", tea.KeyUp.String(), "k":
			m.list, cmd = m.list.Update(msg)

			m.updateFilters()
		default:
			if !m.yaml {
				m.list, cmd = m.list.Update(msg)
			} else {
				m.viewport, cmd = m.viewport.Update(msg)
			}
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case editorFinishedMsg:
		if msg.err != nil {
			return m, ErrCmd(msg.err)
		}

		content, err := os.ReadFile(msg.filename)
		if err != nil {
			return m, ErrCmd(err)
		}
		defer func(name string) {
			_ = os.Remove(name)
		}(msg.filename)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		err = m.updateFn(ctx, content)
		if err != nil {
			return m, ErrCmd(err)
		}
	case TickMsg:
		return m, func() tea.Msg {
			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()

			if m.listFn == nil {
				items, err := m.listItemFn(ctx)
				if err != nil {
					return err
				}

				return ItemsMsg{items, m.name}
			}

			kongItems, err := m.listFn(ctx)
			if err != nil {
				return err
			}

			items := make([]list.Item, len(kongItems))
			for i := range kongItems {
				items[i] = m.toItemFn(kongItems[i])
			}

			return ItemsMsg{items, m.name}
		}
	case ItemsMsg:
		if m.name == msg.Name {
			m.list.SetItems(msg.Items)
			m.updateFilters()
			m.status = ""
		}

		return m, tickCmd()
	case error:
		m.err = msg

		return m, nil
	}

	return m, cmd
}

func (m *RootScreenModel) headerView() string {
	var builder strings.Builder

	builder.WriteString("Kong version: ")
	builder.WriteString(lipgloss.NewStyle().Bold(true).Render(m.kongVersion))
	builder.WriteString("\nKong Admin URL: ")
	builder.WriteString(lipgloss.NewStyle().Bold(true).Render(m.baseURL))
	builder.WriteString("\nKong Edition: ")
	builder.WriteString(lipgloss.NewStyle().Bold(true).Render(m.edition))

	return builder.String()
}

func (m *RootScreenModel) title() string {
	switch m.name {
	case workspaceName:
		return "Workspaces"
	case serviceName:
		return "Workspace[" + m.Client.Workspace() + "] > Services"
	case routeName:
		var sb strings.Builder

		sb.WriteString("Workspace[" + m.Client.Workspace() + "]")

		if m.FilterService == "" {
			sb.WriteString(" > Routes")
		} else {
			sb.WriteString(" > Service[" + m.FilterService + "] > Routes")
		}

		return sb.String()
	case pluginName:
		var sb strings.Builder

		sb.WriteString("Workspace[" + m.Client.Workspace() + "]")

		if m.FilterService != "" {
			sb.WriteString(" > Service[" + m.FilterService + "]")
		}

		if m.FilterRoute != "" {
			sb.WriteString(" > Route[" + m.FilterRoute + "]")
		}

		sb.WriteString(" > Plugins")

		return sb.String()
	}

	return ""
}

func (m *RootScreenModel) View() string {
	header := lipgloss.NewStyle().Align(lipgloss.Left).Render(m.headerView())

	if m.yaml {
		m.viewport.Height = m.height - lipgloss.Height(header) - 1
		m.viewport.Width = m.width

		return lipgloss.JoinVertical(
			lipgloss.Left,
			header,
			"",
			m.viewport.View(),
		)
	}

	m.list.SetHeight(m.height - lipgloss.Height(header) - 2)
	m.list.SetWidth(m.width)

	var lastError string

	if m.err != nil {
		lastError = m.err.Error()
	}

	m.list.Title = m.title()

	return lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		fmt.Sprint("Last error: ", lipgloss.NewStyle().Bold(true).Render(lastError), "\n"),
		m.list.View(),
	)
}

type TickMsg time.Time

func tickCmd() tea.Cmd {
	return tea.Tick(period, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

func InitModel(client *client.Client) (*RootScreenModel, error) {
	model := new(RootScreenModel)

	model.Client = client
	model.Client.FilterService = &model.FilterService
	model.Client.FilterRoute = &model.FilterRoute

	version, err := client.KongVersion()
	if err != nil {
		return nil, err
	}

	model.kongVersion = version.String()

	if version.IsKongGatewayEnterprise() {
		model.SwitchToWorkspaces()

		model.edition = EditionEnterprise
	} else {
		model.SwitchToServices()

		model.edition = EditionOSS
	}

	var dump *os.File

	if _, ok := os.LookupEnv("DEBUG"); ok {
		var err error

		dump, err = os.OpenFile("messages.log", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o644)
		if err != nil {
			os.Exit(1)
		}
	}

	model.dump = dump

	model.help = help.New()

	vport := viewport.New(0, 0)
	vport.KeyMap = viewport.DefaultKeyMap()
	model.viewport = vport

	model.viewportHelp = help.New()

	model.baseURL = client.BaseRootURL()

	model.status = loadingStatus

	return model, nil
}

type IDer interface {
	GetID() *string
}

func (m *RootScreenModel) SelectedID() (*string, error) {
	item, ok := m.list.SelectedItem().(IDer)
	if !ok {
		return nil, UnexpectedTypeError{m.list.SelectedItem()}
	}

	return item.GetID(), nil
}

func (m *RootScreenModel) updateFilters() {
	switch m.name {
	case workspaceName:
		workspaceItem, ok := m.list.SelectedItem().(WorkspaceItem)
		if !ok {
			return
		}

		m.Client.SetWorkspace(*workspaceItem.Name)
	case serviceName:
		serviceItem, ok := m.list.SelectedItem().(ServiceItem)
		if !ok {
			return
		}

		m.FilterService = *serviceItem.Name
	case routeName:
		routeItem, ok := m.list.SelectedItem().(RouteItem)
		if !ok {
			return
		}

		m.FilterRoute = *routeItem.Name
	}
}

func closeIgnoringErr(closer io.Closer) {
	_ = closer.Close()
}
