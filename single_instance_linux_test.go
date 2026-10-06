//go:build linux && !bindings

package main

import (
	"bufio"
	"encoding/json"
	"os/exec"
	"testing"
	"time"

	"github.com/godbus/dbus/v5"
	"github.com/wailsapp/wails/v2/pkg/options"
)

func TestQuitReportsSessionBusFailure(t *testing.T) {
	t.Setenv("DBUS_SESSION_BUS_ADDRESS", "unix:path="+t.TempDir()+"/missing-bus")
	handled, err := handleExistingLinuxInstance("test.button", launchQuit, []string{"--quit"})
	if err == nil || handled {
		t.Fatalf("handled=%v, err=%v; expected communication failure", handled, err)
	}
}

type testInstance struct{ messages chan string }

func (s *testInstance) SendMessage(message string) *dbus.Error {
	s.messages <- message
	return nil
}

func TestQuitDelivery(t *testing.T) {
	bin, err := exec.LookPath("dbus-daemon")
	if err != nil {
		t.Skip("dbus-daemon unavailable")
	}
	cmd := exec.Command(bin, "--session", "--nofork", "--print-address=1")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		t.Fatal(err)
	}
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = cmd.Process.Kill(); _ = cmd.Wait() })
	scanner := bufio.NewScanner(stdout)
	if !scanner.Scan() {
		t.Fatal("missing private bus address")
	}
	t.Setenv("DBUS_SESSION_BUS_ADDRESS", scanner.Text())

	const id = "test.button.quit"
	if handled, err := handleExistingLinuxInstance(id, launchQuit, []string{"--quit"}); err != nil || !handled {
		t.Fatalf("no running instance should succeed: handled=%v, err=%v", handled, err)
	}
	if handled, err := handleExistingLinuxInstance(id, launchDefault, nil); err != nil || handled {
		t.Fatalf("normal launch should continue: handled=%v, err=%v", handled, err)
	}
	owner, err := dbus.ConnectSessionBus()
	if err != nil {
		t.Fatal(err)
	}
	defer owner.Close()
	name, path := linuxSingleInstanceAddress(id)
	if _, err := owner.RequestName(name, dbus.NameFlagDoNotQueue); err != nil {
		t.Fatal(err)
	}
	// The name exists but there is no message receiver: delivery must fail.
	if handled, err := handleExistingLinuxInstance(id, launchQuit, []string{"--quit"}); err == nil || handled {
		t.Fatalf("delivery failure should be reported: handled=%v, err=%v", handled, err)
	}
	receiver := &testInstance{messages: make(chan string, 1)}
	if err := owner.Export(receiver, dbus.ObjectPath(path), name); err != nil {
		t.Fatal(err)
	}
	if handled, err := handleExistingLinuxInstance(id, launchQuit, []string{"--quit"}); err != nil || !handled {
		t.Fatalf("delivery should succeed: handled=%v, err=%v", handled, err)
	}
	select {
	case message := <-receiver.messages:
		var data options.SecondInstanceData
		if err := json.Unmarshal([]byte(message), &data); err != nil {
			t.Fatal(err)
		}
		if len(data.Args) != 1 || data.Args[0] != "--quit" {
			t.Fatalf("unexpected args: %v", data.Args)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("quit message not received")
	}
}
