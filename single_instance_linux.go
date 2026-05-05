//go:build linux

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/godbus/dbus/v5"
	"github.com/wailsapp/wails/v2/pkg/options"
)

func handleExistingLinuxInstance(uniqueID string, action launchAction, args []string) bool {
	name, path := linuxSingleInstanceAddress(uniqueID)
	conn, err := dbus.ConnectSessionBus()
	if err != nil {
		return action == launchQuit
	}
	defer conn.Close()

	reply, err := conn.RequestName(name, dbus.NameFlagDoNotQueue)
	if err != nil {
		return action == launchQuit
	}

	if reply == dbus.RequestNameReplyPrimaryOwner {
		_, _ = conn.ReleaseName(name)
		return action == launchQuit
	}
	if reply != dbus.RequestNameReplyExists {
		return action == launchQuit
	}

	workingDir, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to get working directory:", err)
		return action == launchQuit
	}
	data, err := json.Marshal(options.SecondInstanceData{
		Args:             args,
		WorkingDirectory: workingDir,
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to marshal second-instance data:", err)
		return action == launchQuit
	}

	call := conn.Object(name, dbus.ObjectPath(path)).Call(name+".SendMessage", 0, string(data))
	if call.Err != nil {
		fmt.Fprintln(os.Stderr, "failed to contact running Button instance:", call.Err)
		return action == launchQuit
	}
	return true
}

func linuxSingleInstanceAddress(uniqueID string) (string, string) {
	id := "wails_app_" + strings.ReplaceAll(strings.ReplaceAll(uniqueID, "-", "_"), ".", "_")
	return "org." + id + ".SingleInstance", "/org/" + id + "/SingleInstance"
}
