//go:build !linux

package linuxtray

import "context"

type Tray struct{}

func New(icon []byte) *Tray {
	return nil
}

func (t *Tray) Start(ctx context.Context, onToggle func(), onQuit func()) error {
	return nil
}

func (t *Tray) Stop() {}
