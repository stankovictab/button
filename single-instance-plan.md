# Button Linux Single-Instance, Tray, and Autostart Plan

This plan intentionally ignores `docs/OLD-PLAN.md`. It reflects the agreed Linux-first direction for making Button run in the background, expose a KDE status bar icon, and support a KDE global shortcut through `button --toggle`.

## Goals

- Keep one Button process alive per user session.
- Let the window be hidden/shown without exiting the app.
- Add a small CLI control surface:
  - `button`
  - `button --toggle`
  - `button --quit`
- Let KDE users bind `Meta+B` to `button --toggle` in KDE System Settings.
- Add a Settings UI to Button with a `Run Button on login` checkbox.
    - Create/remove `~/.config/autostart/button.desktop` from that checkbox.
- Add a KDE-compatible status bar/tray icon using `assets/images/button-logo-simple-white.png`.

## Non-Goals For This Phase

- Do not implement native global hotkey registration inside Button.
- Do not use a system-wide `systemctl` service.
- Do not create AUR/pacman packaging yet.
- Do not add `--show`, `--hide`, or `--background` unless implementation proves they are required.
- Do not rely on X11-only global key grabbing, this solution focuses on Wayland.

## Current Repo Context

- Main Wails entry point: `main.go`
- Backend app methods and Linux desktop asset installation: `app.go`
- Existing Linux desktop launcher installation happens in `installLinuxAssets()`.
- Current Linux app icon embed uses `build/appicon.png`.
- Desired tray icon asset is `assets/images/button-logo-simple-white.png`.
- User config lives under the existing `internal/config` package.
- There is an unrelated modified file at `frontend/src/lib/components/ImportPanel.svelte`; do not revert it unless the user explicitly asks.

## Recommended Design

### Process Model

Use Wails `SingleInstanceLock` so only one Button process owns the GUI/session state.

Behavior:

- `button`
  - If no instance exists: start Button normally with the window visible.
  - If an instance already exists: show/focus the existing window.
- `button --toggle`
  - If no instance exists: start Button and show the window.
  - If an instance already exists: toggle the existing window visible/hidden.
- `button --quit`
  - If no instance exists: exit successfully.
  - If an instance already exists: ask the running instance to quit cleanly.

Implementation note:

- Parse CLI arguments before `wails.Run`.
- Pass the requested action to the app instance.
- Use Wails `SingleInstanceLock.OnSecondInstanceLaunch` to receive second-launch arguments and dispatch them to the existing instance.
- Track visible/hidden state inside the Go app because Wails exposes show/hide APIs, but not necessarily a reliable cross-platform visibility query.

### Window Lifecycle

Use Wails window lifecycle options:

- `HideWindowOnClose: true`
- Keep the process alive when the user closes the window.
- Tray menu `Quit` and CLI `--quit` must perform an actual app quit.

Likely runtime APIs:

- `runtime.WindowShow(ctx)`
- `runtime.WindowHide(ctx)`
- `runtime.WindowUnminimise(ctx)`
- `runtime.WindowCenter(ctx)` if needed
- `runtime.Quit(ctx)` for actual quit

Create app methods similar to:

- `showWindow()`
- `hideWindow()`
- `toggleWindow()`
- `quitApp()`

Guard these methods if `a.ctx` is not initialized yet.

### CLI Argument Shape

Keep it intentionally small:

```bash
button
button --toggle
button --quit
```

Invalid arguments should print a concise usage message to stderr and exit non-zero before starting Wails.

Suggested internal representation:

```go
type launchAction string

const (
    launchDefault launchAction = "default"
    launchToggle  launchAction = "toggle"
    launchQuit    launchAction = "quit"
)
```

### KDE Global Shortcut

Do not register a global shortcut in the app yet.

Document that KDE users can configure:

```bash
button --toggle
```

as a custom shortcut, for example `Meta+B`, in KDE System Settings.

Add this to `README.md`, probably under Linux installation or a new `Linux background mode` section.

### Autostart

Add a settings panel in the app, with a checkbox:

```text
Run Button on login
```

When enabled, write:

```text
~/.config/autostart/button.desktop
```

The file should be created only in the current user account.

Suggested contents:

```ini
[Desktop Entry]
Type=Application
Name=Button
Exec=/absolute/path/to/button
Icon=button
Categories=Utility;
StartupNotify=false
X-GNOME-Autostart-enabled=true
```

Notes:

- Use `os.Executable()` for the `Exec=` path.
- Quote or escape the executable path correctly for desktop files if the path contains spaces.
- The initial autostart behavior can start Button normally. If this feels too intrusive during testing, add a persisted `StartHiddenOnAutostart` later, but do not add extra CLI flags in this phase unless necessary.
- On disable, remove only Button's own autostart file at `~/.config/autostart/button.desktop`.
- Before removing, optionally verify the file looks like Button's file, for example `Name=Button` and an `Exec=` path containing `button`.

Backend API suggestions:

- `GetAutostartEnabled() (bool, error)`
- `SetAutostartEnabled(enabled bool) error`

The UI checkbox should reflect the actual autostart file state, not only a stored preference.

### Settings UI

Create a Button settings surface in the existing Svelte app.

Implementation should follow existing UI conventions. Inspect these files before editing:

- `frontend/src/App.svelte`
- `frontend/src/lib/components/Toolbar.svelte`
- `frontend/src/lib/components/WelcomePanel.svelte`
- `frontend/src/lib/components/DonatePanel.svelte`
- `frontend/src/lib/components/ConfirmDialog.svelte`

Recommended UX:

- Add a settings button to the toolbar, ideally with a lucide settings icon if the project already has lucide available; otherwise follow the existing icon approach.
- Open a settings modal/panel.
- Include the `Run Button on login` checkbox.
- Show errors through the existing notification mechanism if available.
- Do not overbuild settings beyond the autostart checkbox in this phase.

After adding backend methods, regenerate Wails bindings if needed:

```bash
wails generate module
```

### Tray / KDE Status Bar Icon

Use `assets/images/button-logo-simple-white.png` for the tray/status icon.

Embed it separately from `build/appicon.png`:

```go
//go:embed assets/images/button-logo-simple-white.png
var trayIcon []byte
```

Preferred Linux implementation:

- Implement KDE-compatible StatusNotifierItem over D-Bus.
- Use the session bus.
- Register a status notifier item with a stable service/path.
- Provide at minimum:
  - app/status title: `Button`
  - icon
  - left-click: toggle Button window
  - menu:
    - `Show/Hide Button`
    - `Quit`

Look for a maintained Go StatusNotifierItem library first. If none fits cleanly, implement a minimal internal package using `github.com/godbus/dbus/v5`, which is already present indirectly through Wails.

Possible package layout:

```text
internal/linuxtray/
  tray.go
  tray_linux.go
  tray_unsupported.go
```

Make non-Linux builds no-op cleanly.

Important KDE note:

- KDE Plasma uses the StatusNotifierItem protocol for modern tray/status icons.
- Avoid old XEmbed tray behavior unless a library provides it as harmless fallback.

If StatusNotifierItem takes longer than expected, it is acceptable to land the single-instance and autostart work first, then do tray in a follow-up. Do not block the entire feature on tray complexity.

### Existing Linux Desktop Asset Installation

`installLinuxAssets()` currently writes:

- `~/.local/share/icons/hicolor/256x256/apps/button.png`
- `~/.local/share/applications/button.desktop`

Keep this behavior, but consider updating the launcher `.desktop` file to include:

```ini
StartupNotify=false
```

Do not use the white tray icon as the desktop launcher icon; keep using `build/appicon.png` for launcher/taskbar identity.

### Documentation

Update `README.md` with a Linux section covering:

- Button can stay running in the background.
- Closing the window hides it.
- The tray icon can show/hide Button and quit it.
- How to enable `Run Button on login` from Button settings.
- How to configure KDE global shortcut:
  1. Open KDE System Settings.
  2. Go to Keyboard > Shortcuts > Custom Shortcuts or equivalent Plasma version path.
  3. Add a command shortcut.
  4. Set command to `/absolute/path/to/button --toggle` or `button --toggle` if installed on `PATH`.
  5. Bind it to `Meta+B`.

Mention that native global hotkey registration is intentionally not implemented yet because Wayland requires compositor/portal-mediated shortcuts.

## Implementation Checklist

1. Add launch action parsing for `button`, `--toggle`, and `--quit`.
2. Add `SingleInstanceLock` to Wails options.
3. Add app-level show/hide/toggle/quit behavior.
4. Set close-to-hide behavior with `HideWindowOnClose`.
5. Add backend autostart methods.
6. Add Settings UI with `Run Button on login`.
7. Generate Wails frontend bindings.
8. Add Linux tray/status notifier integration using the white logo asset.
9. Update README with KDE shortcut and background behavior docs.
10. Verify builds and key behavior.

## Verification

Run:

```bash
go build ./...
cd frontend && npm run build
wails build
```

Manual Linux/KDE checks:

- Start `build/bin/button`.
- Close the window; process should remain alive.
- Run `build/bin/button --toggle`; window should show/hide.
- Run `build/bin/button --quit`; process should exit.
- Launch a second `build/bin/button`; it should not create a second long-lived GUI process.
- Enable `Run Button on login`; verify `~/.config/autostart/button.desktop` exists.
- Disable it; verify the autostart file is removed.
- Verify tray icon appears in KDE panel using `assets/images/button-logo-simple-white.png`.
- Tray left-click toggles window.
- Tray menu quit exits the process.
- Configure KDE `Meta+B` to run `build/bin/button --toggle`; verify it toggles Button.

## Risk Notes

- Wayland global shortcuts are intentionally avoided in this phase. KDE user-configured shortcuts are more reliable and transparent.
- Tray support may be the trickiest part because Wails v2 does not provide a mature built-in Linux tray API.
- A white PNG tray icon may be less visible on light KDE panels. If that becomes a problem, consider a symbolic/themed icon follow-up.
- If the binary path changes after enabling autostart, the autostart file may point at the old path. Packaging can solve this later; for tarball users, document that toggling the setting off/on refreshes the path.
