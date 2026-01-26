# PLAN.md: Project "Button"

## 1. Project Overview
**Button** is a cross-platform (Linux/macOS) "Quick-Reference" GUI for personal keyboard shortcuts.
- **The Problem:** DevOps engineers manage dozens of tools (KDE, Ghostty, Zellij, Neovim, Tmux) with complex, custom keybinds that are hard to memorize.
- **The Solution:** A lightweight, Raycast-inspired "Prettier YAML Viewer and Editor" accessible via a single global hotkey. It provides a searchable grid of application "Cards" that reveal shortcuts upon selection.

## 2. Technical Stack
- **Backend:** Go (Golang) 1.21+
- **Framework:** [Wails v2/v3](https://wails.io/) (Go + Webview Bridge)
- **Frontend:** Svelte + Tailwind CSS
- **Data Format:** YAML (stored in `~/.config/button/`)
- **Target Systems:**
    - **Linux:** CachyOS (Arch-based) + KDE Plasma (Wayland/X11)
    - **macOS:** Latest versions, supporting native vibrancy/blur.

## 3. Core Requirements & Architecture

### A. Data Structure (The "Source of Truth")
The app acts as a GUI layer over a flat-file database.
- **Location:** `$HOME/.config/button/apps/*.yaml`
- **Schema Design:** Must support platform-specific overrides.
```yaml
app: "App Name"
icon: "icon-name"
groups:
  - category: "Navigation"
    shortcuts:
      - desc: "Description"
        keys: ["Alt", "p"]           # Default/Fallback
        linux: ["Ctrl", "Shift", "p"] # Linux Override
        macos: ["Cmd", "p"]          # macOS Override
```

### B. The GUI (The "Raycast" Aesthetic)
- **Single Window:** Center-screen, fixed width (~600px-800px), floating.
- **Visuals:**
    - "Glassmorphism" (Blurred transparency).
    - Raycast-style search bar at the top (autofocus on launch).
    - Grid/List of App Cards with icons.
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
- [ ] Implement fuzzy search (logic in Go or JS) to find apps by name or description.
- [ ] Add keyboard navigation (Arrow keys to select, Esc to clear).

### Phase 3: The Editor (Writing)
- [ ] Create a "New Shortcut/App" UI form.
- [ ] Implement Go logic to write/update YAML files safely without destroying comments (if possible) or using a standard clean YAML structure.

### Phase 4: Polish & Deployment
- [ ] Global hotkey implementation.
- [ ] Native blur/transparency for macOS and KDE.
- [ ] Build pipeline for `.app` (Mac) and Binary (Linux).

## 5. Agent Instructions (How to help me)
When working on this project:
1. **Keep it "DevOps-y":** Ensure the YAML structure is clean and easy to edit via CLI (Neovim) outside the app.
2. **Prioritize Performance:** The window must appear instantly. Avoid heavy JS libraries; keep the Svelte bundle lean.
3. **Cross-Platform Awareness:** Every time a UI element is built, verify how "Meta", "Cmd", "Ctrl", and "Alt" are represented across OSs.
4. **Window Management:** If writing Go code for the window, ensure it handles "Focus Lost" events to auto-hide the app (optional but preferred).

---

### Suggested First Task for AI:
*“Create the Go Structs and the YAML parser logic that can read a directory of YAML files and return a platform-filtered list of shortcuts to the frontend.”*
