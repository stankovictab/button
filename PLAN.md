# PLAN.md: Project "Button"

## 1. Project Overview
**Button** is a cross-platform (Linux and macOS) "Quick-Reference" GUI for personal keyboard shortcuts.
- **The Problem:** DevOps engineers manage dozens of tools (tmux, zellij, NeoVim, Ghostty, WezTerm, KDE) with complex, custom keybinds that are hard to memorize. On top of this, the difference between Linux and macOS shortcuts is not always clear because of the keyboard layout differences.
- **The Solution:** A lightweight, Raycast-inspired GUI app, accessible via a single global hotkey, that functions as a viewer and editor of YAML configuration for custom shortcuts. It provides a searchable grid of application "Cards" that reveal shortcuts upon selection.
It's important to note that this app doesn't change any app's functionality or shortcuts, it's just a "quick reference" for the user.
For the time being it also won't be automatically updated from the app's configuration, instead all changes must be done manually.

## 2. Technical Stack
- **Backend:** Go (Golang) 1.21+
- **Framework:** [Wails v2](https://wails.io/) (Go + Webview Bridge)
- **Frontend:** Svelte + Tailwind CSS
- **Data Format:** YAML (stored in `~/.config/button/`)
- **Target Systems:**
    - **Linux:** CachyOS (Arch-based) + KDE Plasma (Wayland/X11)
    - **macOS:** Latest versions, supporting native vibrancy/blur.

## 3. Core Requirements & Architecture

### A. Data Structure (The "Source of Truth")
The app acts as a GUI layer over a flat-file database.
- **Location:** `$HOME/.config/button/apps/*.yaml`
- **Schema Design:** Must support default/fallback keybinds, and Linux and macOS specific overrides. See example :
```yaml
app: "App Name"
icon: "icon-name"
groups:
  - category: "Navigation"
    shortcuts:
      - desc: "Description"
        keys: ["Alt", "p"]
        linux: ["Ctrl", "Shift", "p"]
        macos: ["Cmd", "p"]
```

### B. The GUI (The "Raycast" Aesthetic)
- **Single Window:** Center-screen, fixed width (~600px-800px), floating.
- **Visuals:**
    - Raycast-style search bar at the top (autofocus on launch).
    - Grid/List of App Cards with icons.
    - Possible blur, transparency and glass-like effects for macOS and KDE down the line.
- **Navigation:**
    - `Fuzzy Find` as you type to filter apps or specific shortcuts.
    - `Enter` to drill into a card.
    - `Esc` to go back or hide the window.

### C. OS Integration
- **Platform Detection:** Go backend detects OS at runtime to serve the correct key-array to the Svelte frontend.
- **KDE Integration:** App should be configured via KWin rules (or Go hints) to stay on top, skip taskbar, and appear on all virtual desktops.
- **macOS Integration:** Support for "Vibrancy" (Native backdrop blur).

## 4. Development Phases

### Phase 1: The Reader (MVP)
- [ ] Initialize Wails project with Go + Svelte.
- [ ] Implement Go logic to read `~/.config/button/apps/` and parse YAML files into a JSON-bridge for Svelte.
- [ ] Build basic Svelte UI: Search Bar + App Cards.
- [ ] Implement "Key Badge" components (rendering `CMD` on Mac and `META` on Linux).

### Phase 2: The Search & Interaction
- [ ] Implement fuzzy search in JS to find apps by name or description.
- [ ] Add keyboard navigation (Arrow keys to select, Esc to clear).
- [ ] Add icons for app cards and UI chrome:
    - **Simple Icons** (`svelte-simple-icons`) for brand/app logos (NeoVim, tmux, KDE, etc.)
    - **Lucide** (`lucide-svelte`) for UI chrome (search, settings, close, back, etc.)
    - Custom fallback SVGs in `frontend/src/lib/icons/custom/` for apps missing from Simple Icons (e.g. Zellij, Yazi). The YAML `icon` field routes unknown names to this folder. Import these SVGs as raw strings (`import YaziIcon from './custom/yazi.svg?raw'`) and inject them inline via Svelte's `{@html}`. This ensures they get bundled into the binary by Vite and no external files are needed at runtime.

### Phase 3: The Editor (Writing)
- [ ] Create a "New Shortcut/App" UI form.
- [ ] Implement Go logic to write/update YAML files safely without destroying comments (if possible) or using a standard clean YAML structure.

### Phase 4: Polish & Deployment
- [ ] Global hotkey implementation.
- [ ] Blur/transparency for macOS and KDE.
- [ ] Build pipeline for `.app` (Mac) and Binary (Linux).

## 5. Agent Instructions (How to help me)
When working on this project:
1. **Keep it "DevOps-y":** Ensure the YAML structure is clean and easy to edit outside the app.
2. **Prioritize Performance:** The window must appear instantly. Avoid heavy JS libraries; keep the Svelte bundle lean.
3. **Cross-Platform Awareness:** Every time a UI element is built, verify how "Meta", "Cmd", "Ctrl", and "Alt" are represented across OSs.
4. **Window Management:** If writing Go code for the window, ensure it handles "Focus Lost" events to auto-hide (or close) the app (optional but preferred).

---

### Suggested First Task for AI:
*"Create the Go Structs and the YAML parser logic that can read a directory of YAML files and return a platform-filtered list of shortcuts to the frontend."*
