//go:build linux

package linuxtray

import (
	"bytes"
	"context"
	"fmt"
	"image/png"
	"os"
	"path/filepath"
	"sync"

	"github.com/godbus/dbus/v5"
)

const (
	itemInterface = "org.kde.StatusNotifierItem"
	menuInterface = "com.canonical.dbusmenu"
	itemPath      = dbus.ObjectPath("/StatusNotifierItem")
	menuPath      = dbus.ObjectPath("/StatusNotifierItem/Menu")
)

type Tray struct {
	iconPNG  []byte
	conn     *dbus.Conn
	item     *statusNotifierItem
	menu     *dbusMenu
	cancel   context.CancelFunc
	stopOnce sync.Once
}

func New(icon []byte) *Tray {
	return &Tray{iconPNG: icon}
}

func (t *Tray) Start(ctx context.Context, onToggle func(), onQuit func()) error {
	conn, err := dbus.SessionBus()
	if err != nil {
		return err
	}

	serviceName := fmt.Sprintf("org.kde.StatusNotifierItem-%d-1", os.Getpid())
	reply, err := conn.RequestName(serviceName, dbus.NameFlagDoNotQueue)
	if err != nil {
		conn.Close()
		return err
	}
	if reply != dbus.RequestNameReplyPrimaryOwner {
		conn.Close()
		return fmt.Errorf("status notifier service name is already owned")
	}

	if err := installTrayIcon(t.iconPNG); err != nil {
		conn.Close()
		return err
	}

	trayCtx, cancel := context.WithCancel(ctx)
	t.conn = conn
	t.cancel = cancel
	t.item = &statusNotifierItem{onActivate: onToggle}
	t.menu = &dbusMenu{onToggle: onToggle, onQuit: onQuit}

	if err := conn.Export(t.item, itemPath, itemInterface); err != nil {
		t.Stop()
		return err
	}
	if err := conn.Export(t.item, itemPath, "org.freedesktop.DBus.Properties"); err != nil {
		t.Stop()
		return err
	}
	if err := conn.Export(t.menu, menuPath, menuInterface); err != nil {
		t.Stop()
		return err
	}
	if err := conn.Export(t.menu, menuPath, "org.freedesktop.DBus.Properties"); err != nil {
		t.Stop()
		return err
	}

	watcher := conn.Object("org.kde.StatusNotifierWatcher", "/StatusNotifierWatcher")
	if call := watcher.Call("org.kde.StatusNotifierWatcher.RegisterStatusNotifierItem", 0, serviceName); call.Err != nil {
		t.Stop()
		return call.Err
	}

	go func() {
		<-trayCtx.Done()
		t.Stop()
	}()

	return nil
}

func (t *Tray) Stop() {
	t.stopOnce.Do(func() {
		if t.cancel != nil {
			t.cancel()
		}
		if t.conn != nil {
			t.conn.Close()
		}
	})
}

type statusNotifierItem struct {
	onActivate func()
}

func (i *statusNotifierItem) Activate(x int32, y int32) *dbus.Error {
	if i.onActivate != nil {
		go i.onActivate()
	}
	return nil
}

func (i *statusNotifierItem) SecondaryActivate(x int32, y int32) *dbus.Error {
	return i.Activate(x, y)
}

func (i *statusNotifierItem) ContextMenu(x int32, y int32) *dbus.Error {
	return nil
}

func (i *statusNotifierItem) Scroll(delta int32, orientation string) *dbus.Error {
	return nil
}

func (i *statusNotifierItem) Get(iface string, property string) (dbus.Variant, *dbus.Error) {
	props := i.properties()
	value, ok := props[property]
	if !ok {
		return dbus.MakeVariant(nil), dbus.MakeFailedError(fmt.Errorf("unknown property %s", property))
	}
	return value, nil
}

func (i *statusNotifierItem) GetAll(iface string) (map[string]dbus.Variant, *dbus.Error) {
	return i.properties(), nil
}

func (i *statusNotifierItem) Set(iface string, property string, value dbus.Variant) *dbus.Error {
	return dbus.MakeFailedError(fmt.Errorf("property %s is read-only", property))
}

func (i *statusNotifierItem) properties() map[string]dbus.Variant {
	return map[string]dbus.Variant{
		"Category":              dbus.MakeVariant("ApplicationStatus"),
		"Id":                    dbus.MakeVariant("button"),
		"Title":                 dbus.MakeVariant("Button"),
		"Status":                dbus.MakeVariant("Active"),
		"WindowId":              dbus.MakeVariant(uint32(0)),
		"IconName":              dbus.MakeVariant("button-tray"),
		"IconThemePath":         dbus.MakeVariant(""),
		"AttentionIconName":     dbus.MakeVariant(""),
		"AttentionMovieName":    dbus.MakeVariant(""),
		"OverlayIconName":       dbus.MakeVariant(""),
		"ToolTip":               dbus.MakeVariant(statusNotifierToolTip{}),
		"ItemIsMenu":            dbus.MakeVariant(false),
		"Menu":                  dbus.MakeVariant(menuPath),
		"XAyatanaOrderingIndex": dbus.MakeVariant(uint32(0)),
	}
}

type statusNotifierToolTip struct {
	IconName string
	Pixmaps  []iconPixmap
	Title    string
	Text     string
}

type iconPixmap struct {
	Width  int32
	Height int32
	Bytes  []byte
}

type dbusMenu struct {
	onToggle func()
	onQuit   func()
}

type layoutItem struct {
	ID         int32
	Props      map[string]dbus.Variant
	ChildNodes []dbus.Variant
}

func (m *dbusMenu) GetLayout(parentID int32, recursionDepth int32, propertyNames []string) (uint32, layoutItem, *dbus.Error) {
	children := m.menuItems()
	childVariants := make([]dbus.Variant, 0, len(children))
	for _, child := range children {
		childVariants = append(childVariants, dbus.MakeVariant(child))
	}

	return 1, layoutItem{
		ID: -1,
		Props: map[string]dbus.Variant{
			"children-display": dbus.MakeVariant("submenu"),
		},
		ChildNodes: childVariants,
	}, nil
}

func (m *dbusMenu) menuItems() []layoutItem {
	return []layoutItem{
		{
			ID: 1,
			Props: map[string]dbus.Variant{
				"label":   dbus.MakeVariant("Show/Hide Button"),
				"enabled": dbus.MakeVariant(true),
				"visible": dbus.MakeVariant(true),
			},
			ChildNodes: []dbus.Variant{},
		},
		{
			ID: 2,
			Props: map[string]dbus.Variant{
				"type":    dbus.MakeVariant("separator"),
				"enabled": dbus.MakeVariant(true),
				"visible": dbus.MakeVariant(true),
			},
			ChildNodes: []dbus.Variant{},
		},
		{
			ID: 3,
			Props: map[string]dbus.Variant{
				"label":   dbus.MakeVariant("Quit"),
				"enabled": dbus.MakeVariant(true),
				"visible": dbus.MakeVariant(true),
			},
			ChildNodes: []dbus.Variant{},
		},
	}
}

func (m *dbusMenu) GetGroupProperties(ids []int32, propertyNames []string) ([]struct {
	ID    int32
	Props map[string]dbus.Variant
}, *dbus.Error) {
	result := make([]struct {
		ID    int32
		Props map[string]dbus.Variant
	}, 0, len(ids))

	for _, id := range ids {
		for _, child := range m.menuItems() {
			if child.ID == id {
				result = append(result, struct {
					ID    int32
					Props map[string]dbus.Variant
				}{ID: id, Props: child.Props})
				break
			}
		}
	}
	return result, nil
}

func (m *dbusMenu) GetProperty(id int32, name string) (dbus.Variant, *dbus.Error) {
	items, _ := m.GetGroupProperties([]int32{id}, nil)
	if len(items) == 0 {
		return dbus.MakeVariant(nil), dbus.MakeFailedError(fmt.Errorf("unknown menu id %d", id))
	}
	value, ok := items[0].Props[name]
	if !ok {
		return dbus.MakeVariant(nil), dbus.MakeFailedError(fmt.Errorf("unknown property %s", name))
	}
	return value, nil
}

func (m *dbusMenu) Event(id int32, eventID string, data dbus.Variant, timestamp uint32) *dbus.Error {
	if eventID != "clicked" {
		return nil
	}
	switch id {
	case 1:
		if m.onToggle != nil {
			go m.onToggle()
		}
	case 3:
		if m.onQuit != nil {
			go m.onQuit()
		}
	}
	return nil
}

func (m *dbusMenu) EventGroup(events []struct {
	ID        int32
	EventID   string
	Data      dbus.Variant
	Timestamp uint32
}) ([]int32, *dbus.Error) {
	for _, event := range events {
		_ = m.Event(event.ID, event.EventID, event.Data, event.Timestamp)
	}
	return nil, nil
}

func (m *dbusMenu) AboutToShow(id int32) (bool, *dbus.Error) {
	return false, nil
}

func (m *dbusMenu) AboutToShowGroup(ids []int32) ([]int32, []int32, *dbus.Error) {
	return nil, nil, nil
}

func (m *dbusMenu) Get(iface string, property string) (dbus.Variant, *dbus.Error) {
	props := map[string]dbus.Variant{
		"Version":       dbus.MakeVariant(uint32(3)),
		"TextDirection": dbus.MakeVariant("ltr"),
		"Status":        dbus.MakeVariant("normal"),
		"IconThemePath": dbus.MakeVariant([]string{}),
	}
	value, ok := props[property]
	if !ok {
		return dbus.MakeVariant(nil), dbus.MakeFailedError(fmt.Errorf("unknown property %s", property))
	}
	return value, nil
}

func (m *dbusMenu) GetAll(iface string) (map[string]dbus.Variant, *dbus.Error) {
	return map[string]dbus.Variant{
		"Version":       dbus.MakeVariant(uint32(3)),
		"TextDirection": dbus.MakeVariant("ltr"),
		"Status":        dbus.MakeVariant("normal"),
		"IconThemePath": dbus.MakeVariant([]string{}),
	}, nil
}

func (m *dbusMenu) Set(iface string, property string, value dbus.Variant) *dbus.Error {
	return dbus.MakeFailedError(fmt.Errorf("property %s is read-only", property))
}

func installTrayIcon(iconPNG []byte) error {
	if len(iconPNG) == 0 {
		return nil
	}
	if _, err := png.Decode(bytes.NewReader(iconPNG)); err != nil {
		return err
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	iconDir := filepath.Join(home, ".local", "share", "icons", "hicolor", "256x256", "status")
	if err := os.MkdirAll(iconDir, 0755); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(iconDir, "button-tray.png"), iconPNG, 0644)
}
