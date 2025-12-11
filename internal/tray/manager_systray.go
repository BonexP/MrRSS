//go:build windows

package tray

import (
	"context"
	"sync"
	"sync/atomic"

	"github.com/getlantern/systray"

	"MrRSS/internal/handlers/core"
)

// Manager provides a thin wrapper around the system tray menu.
type Manager struct {
	handler *core.Handler
	icon    []byte
	running atomic.Bool
	stopCh  chan struct{}
	mu      sync.Mutex
}

// NewManager creates a new tray manager for supported platforms.
func NewManager(handler *core.Handler, icon []byte) *Manager {
	return &Manager{
		handler: handler,
		icon:    icon,
	}
}

// Start initialises the system tray if it isn't already running.
// onQuit should trigger an application shutdown, and onShow should restore the main window.
func (m *Manager) Start(ctx context.Context, onQuit func(), onShow func()) {
	m.mu.Lock()
	if m.running.Load() {
		m.mu.Unlock()
		return
	}

	m.stopCh = make(chan struct{})
	m.running.Store(true)
	m.mu.Unlock()

	go systray.Run(func() {
		m.run(ctx, onQuit, onShow)
	}, func() {
		m.running.Store(false)
	})
}

func (m *Manager) run(ctx context.Context, onQuit func(), onShow func()) {
	if len(m.icon) > 0 {
		systray.SetIcon(m.icon)
	}
	systray.SetTitle("MrRSS")
	systray.SetTooltip("MrRSS")

	showItem := systray.AddMenuItem("Show MrRSS", "Restore window")
	refreshItem := systray.AddMenuItem("Refresh now", "Refresh all feeds")
	systray.AddSeparator()
	quitItem := systray.AddMenuItem("Quit", "Quit MrRSS")

	go func() {
		for {
			select {
			case <-showItem.ClickedCh:
				if onShow != nil {
					onShow()
				}
			case <-refreshItem.ClickedCh:
				if m.handler != nil && m.handler.Fetcher != nil {
					go m.handler.Fetcher.FetchAll(context.Background())
				}
			case <-quitItem.ClickedCh:
				if onQuit != nil {
					onQuit()
				}
			case <-ctx.Done():
				systray.Quit()
				return
			case <-m.stopCh:
				systray.Quit()
				return
			}
		}
	}()
}

// Stop tears down the system tray if it is running.
func (m *Manager) Stop() {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.running.Load() {
		return
	}
	if m.stopCh != nil {
		select {
		case <-m.stopCh:
		default:
			close(m.stopCh)
		}
		m.stopCh = nil
	}
	m.running.Store(false)
}

// IsRunning returns true if the tray has been started.
func (m *Manager) IsRunning() bool {
	return m.running.Load()
}
