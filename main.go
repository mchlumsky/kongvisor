package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/mchlumsky/kongvisor/internal/config"
	"github.com/mchlumsky/kongvisor/internal/model"
)

func run() error {
	cfg, err := config.NewConfig()
	if err != nil {
		return err
	}

	args := os.Args

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
