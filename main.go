package main

import (
	"errors"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/mchlumsky/kongvisor/internal/config"
	"github.com/mchlumsky/kongvisor/internal/model"
)

var (
	errUsage              = errors.New("usage: kongvisor <gateway-name>")
	errGatewayCfgNotFound = errors.New("error: gateway config not found")
)

func run() error {
	cfg, err := config.NewConfig()
	if err != nil {
		return err
	}

	args := os.Args
	if len(args) != 2 {
		return errUsage
	}

	if _, found := cfg[args[1]]; !found {
		return fmt.Errorf("%w: %q", errGatewayCfgNotFound, args[1])
	}

	cl, err := cfg[args[1]].GetKongClient()
	if err != nil {
		return err
	}

	initModel, err := model.InitModel(cl)
	if err != nil {
		return err
	}

	p := tea.NewProgram(initModel, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		return err
	}

	return nil
}

func main() {
	err := run()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
