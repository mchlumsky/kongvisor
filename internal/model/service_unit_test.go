package model

// import (
// 	"testing"
//
// 	"github.com/charmbracelet/bubbles/list"
// 	tea "github.com/charmbracelet/bubbletea"
// 	"github.com/mchlumsky/kongvisor/internal/config"
// 	"github.com/stretchr/testify/assert"
// )
//
// func TestServiceFiltering(t *testing.T) {
// 	t.Run("Enable filtering", func(t *testing.T) {
// 		gc := config.GatewayConfig{
// 			URL: "http://localhost:8001",
// 		}
//
// 		kc, err := gc.GetKongClient()
// 		assert.NoError(t, err, "Failed to get client")
// 		assert.NotNil(t, kc, "Kong client should not be nil")
//
// 		model, err := InitModel(kc)
// 		assert.NoError(t, err, "Failed to get model")
//
// 		model.SwitchToServices()
//
// 		// switch to routes by pressing 'r'
// 		k := tea.Key{Type: tea.KeyRunes, Runes: []rune{'r'}}
// 		model.Update(tea.KeyMsg(k))
//
// 		assert.Equal(t, "routes", model.name, "Model should be in routes mode")
//
// 		// reset to services
// 		model.SwitchToServices()
// 		assert.Equal(t, "services", model.name, "Model should be in services mode")
//
// 		// start filtering
// 		k = tea.Key{Type: tea.KeyRunes, Runes: []rune{'/'}}
// 		m, _ := model.Update(tea.KeyMsg(k))
//
// 		rm, _ := m.(*RootScreenModel)
//
// 		assert.True(t, rm.list.FilterState() == list.Filtering)
// 	})
// }
