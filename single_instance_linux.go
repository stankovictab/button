//go:build linux && !bindings

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/godbus/dbus/v5"
	"github.com/wailsapp/wails/v2/pkg/options"
)

func handleExistingLinuxInstance(uniqueID string, action launchAction, args []string) (bool, error) {
	name, path := linuxSingleInstanceAddress(uniqueID)
	conn, err := dbus.ConnectSessionBus()
	if err != nil {
		return false, fmt.Errorf("connect to session bus: %w", err)
	}
	defer conn.Close()

	reply, err := conn.RequestName(name, dbus.NameFlagDoNotQueue)
	if err != nil {
		return false, fmt.Errorf("check running Button instance: %w", err)
	}

	if reply == dbus.RequestNameReplyPrimaryOwner {
		_, _ = conn.ReleaseName(name)
		return action == launchQuit, nil
	}
	if reply != dbus.RequestNameReplyExists {
		return false, fmt.Errorf("unexpected D-Bus name reply: %d", reply)
	}

	workingDir, err := os.Getwd()
	if err != nil {
		return false, fmt.Errorf("get working directory: %w", err)
	}
	data, err := json.Marshal(options.SecondInstanceData{
		Args:             args,
		WorkingDirectory: workingDir,
	})
	if err != nil {
		return false, fmt.Errorf("marshal second-instance data: %w", err)
	}

	call := conn.Object(name, dbus.ObjectPath(path)).Call(name+".SendMessage", 0, string(data))
	if call.Err != nil {
		return false, fmt.Errorf("contact running Button instance: %w", call.Err)
	}
	return true, nil
}

func linuxSingleInstanceAddress(uniqueID string) (string, string) {
	id := "wails_app_" + strings.ReplaceAll(strings.ReplaceAll(uniqueID, "-", "_"), ".", "_")
	return "org." + id + ".SingleInstance", "/org/" + id + "/SingleInstance"
}
