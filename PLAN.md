# PLAN.md: Project "Button"

## 1. Project Overview
**Button** is a cross-platform (Linux and macOS) "Quick-Reference" GUI for personal keyboard shortcuts.
- **The Problem:** DevOps engineers manage dozens of tools (tmux, zellij, NeoVim, Ghostty, WezTerm, KDE) with complex, custom keybinds that are hard to memorize. On top of this, the difference between Linux and macOS shortcuts is not always clear because of the keyboard layout differences.
- **The Solution:** A lightweight, Raycast-inspired GUI app, accessible via a single global hotkey, that functions as a viewer and editor of YAML configuration for custom shortcuts. It provides a searchable list of applications that reveal shortcuts in a persistent side panel upon selection.
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
- **Schema Design:** Must support default/fallback keybinds, and Linux and macOS specific overrides. See examples:
```yaml
# Generic app example
app: "App Name"
icon: "icon-name"
groups:
  - category: "Navigation"
    shortcuts:
      - desc: "Description"
        keys: [Alt, p]
        linux: [Ctrl, Shift, p]
        macos: [Cmd, p]
```
```yaml
# Linear example — shows mixed usage of keys, linux, and macos
app: "Linear"
icon: "linear"
groups:
  - category: "Navigation"
    shortcuts:
      - desc: "Open Settings"
        keys: [G, S]
        linux: [G, S]
        macos: [G, S]
      - desc: "Go to My Issues"
        keys: [G, I]
      - desc: "Command Palette"
        linux: [Ctrl, K]
        macos: [Cmd, K]
  - category: "Issues"
    shortcuts:
      - desc: "New Issue"
        keys: [C]
      - desc: "Assign to Me"
        linux: [Ctrl, Shift, M]
        macos: [Cmd, Shift, M]
```

### B. The GUI (The "Raycast" Aesthetic)
- **Layout:** Two-panel split layout — fixed-width app list on the left (~240px), persistent shortcut detail panel on the right. The detail panel is always visible and updates instantly as the user navigates the app list. There is no "drill-in" navigation; both panels are always shown.
- **Single Window:** Center-screen, ~900px wide, floating, dark theme (`#111111` base).
- **Toolbar (top bar):** Contains, left to right:
    - **Search bar** — autofocused on launch, full-width, searches across both app names and shortcut descriptions simultaneously.
    - **OS toggle** — segmented control (Linux | macOS) to the right of the search bar. One-click switch; no dropdown. Visually uses an active/inactive state (blue accent for Linux, orange accent for macOS). The backend detects OS at launch and sets the default, but the user can override it in the UI.
    - **Help button (`?`)** — icon-only button, opens a mini panel with app version info, GitHub repo link, and keyboard navigation hints.
    - **Donation button** — subtle amber/orange tint, Ko-fi icon with no label.
- **Visuals:**
    - Dark base: `#111111` window, `#1c1c1c` inputs and secondary surfaces, `#272727` hover states.
    - Blue accent color for selection state and match highlights.
    - Key badges rendered in a "raised" style: dark background, `1px solid #3f3f3f` border, `2px solid #525252` bottom border for depth. Highlighted key badges (on match) use blue tints.
    - Possible blur, transparency, and glass-like effects for macOS and KDE down the line.

### C. App List Panel (Left)
- **Style:** Vertical list, not a card grid.
- **Row contents:** App icon (28×28px, rounded 7px), app name, shortcut count. When a search is active, rows that have matching shortcuts show an annotated count (e.g. "18 shortcuts · 2 matches" with the match count).
- **Selection state:** Selected row has a dark blue background. Clicking any row (or selecting it with arrow keys) updates the detail panel immediately.
- **Section header:** Small all-caps label showing total app count (e.g. "Apps · 8").
- **Version:** App version (e.g. "v1.0.0") in the footer, small gray text.

### D. Shortcut Detail Panel (Right)
- **Always visible** — no empty state needed once an app is selected; defaults to the first app in the list on launch.
- **Header:** App icon, app name, shortcut count + group count. When a search is active, a badge shows the match count (e.g. "2 matches for 'code actions'").
- **Shortcut rows:** Grouped by category (small all-caps group label). Each row shows the shortcut description on the left, key badges on the right.
- **Search match highlighting:** Matching shortcut rows are highlighted with a blue left-border accent (`2px solid #3a88ed`) and a blue-tinted background. Non-matching rows remain visible but dimmed — context is preserved, not stripped. The panel does not filter out non-matches.
- **Footer:** Keyboard hint bar — `↑↓ navigate`, `esc clear / close`. Right-aligned match summary when search is active (e.g. "2 shortcut matches in NeoVim").

### E. Search Behaviour
- **Unified fuzzy search:** A single search input queries both app names and shortcut descriptions at the same time.
- **App list response:** Apps with no matches are de-ranked to the bottom of the list and dimmed. 
- **Detail panel response:** The currently selected app's shortcut panel updates live, highlighting rows that match. Non-matching rows stay visible for context.
- **Auto-selection:** As the user types, if the top-ranked app changes, the detail panel switches to show that app automatically.

### F. OS Integration
- **Platform Detection:** Go backend detects OS at runtime (`runtime.GOOS`) to serve the correct key-array to the Svelte frontend. The UI OS toggle overrides this for the current session (useful for cross-referencing Linux vs. macOS binds without restarting).
- **KDE Integration:** App should be configured via KWin rules (or Go hints) to stay on top, skip taskbar, and appear on all virtual desktops.
- **macOS Integration:** Support for "Vibrancy" (Native backdrop blur).

## 4. Development Phases

### Phase 1: The Reader (MVP)
- [x] Initialize Wails project with Go + Svelte.
- [x] Implement Go logic to read `~/.config/button/apps/` and parse YAML files into a JSON-bridge for Svelte.
    - **Path note:** Never use a raw `~` string — Go does not expand it. Always resolve the config dir via `os.UserHomeDir()` and `filepath.Join`. Both Linux and macOS use `~/.config/button/apps/`.
    - **Code lives in:** `internal/config/` — `models.go` (structs), `reader.go` (parsing + platform detection), `watcher.go` (live file watching).
    - [x] In `startup()`, ensure the config directory exists before any reads: call `os.MkdirAll(filepath.Join(home, ".config", "button", "apps"), 0755)`. This is safe to call even if the directory already exists.
    - [x] Define Go structs (`AppConfig`, `Group`, `Shortcut`) matching the YAML schema, with `yaml` and `json` struct tags.
    - [x] Add `gopkg.in/yaml.v3` dependency for YAML parsing.
    - [x] Implement directory reader — scan `~/.config/button/apps/*.yaml` and `*.yml`, parse each file into an `AppConfig` struct. Empty or malformed files are reported as warnings instead of hard errors.
    - [x] Implement platform detection (`runtime.GOOS`) — resolve each shortcut's key array: use `linux`/`macos` override if present, otherwise fall back to `keys`.
    - [x] Expose a Wails-bound method `GetApps() (AppsResponse, error)` that returns the platform-filtered list (and any warnings) to the frontend. Always return both values — Wails surfaces the error as a rejected promise in the frontend.
    - [x] Create sample YAML files in `~/.config/button/apps/` for testing (e.g. Linear, NeoVim, Ghostty).
    - [x] **Bonus:** Implemented live file watching via `fsnotify` — the frontend updates automatically when YAML files are added, modified, or removed (no app restart needed).
- [x] Build Svelte UI — two-panel layout.
    - [x] Call `GetApps()` on mount and store the result in Svelte state.
    - [x] Build the top toolbar: search bar (autofocused on launch), OS segmented toggle (set on launch), `?` help button, donation button.
    - [x] Build the App List panel (left, ~240px): vertical list rows with icon, name, shortcut count, and match annotation. Wire selection state.
    - [x] Build the Shortcut Detail panel (right): always-visible, updates on app selection. Header with icon + match badge. Grouped shortcut rows with key badges.
    - [x] Wire unified search: filter/rank app list and highlight matching shortcut rows in the detail panel simultaneously.
    - [x] Default selection to the first app in the list on launch.
    - [x] Wire keyboard navigation: `↑↓` moves selection in the app list, `esc` clears search or closes the window.
- [x] Implement "Key Badge" component — renders key names with raised styling. Handle special cases:
    - `Cmd` / `⌘` on macOS, `Ctrl` on Linux for the primary modifier.
    - `Meta` / `Super` on Linux, no direct equivalent on macOS.
    - Symbols: `⌃` for Ctrl, `⌥` for Alt/Option, `⇧` for Shift, `⌘` for Cmd.
    - Plain-text fallbacks for multi-char keys: `leader`, `space`, `:vs`, etc.

### Phase 2: The Search & Interaction
- [x] Implement fuzzy search in JS (no heavy libraries — a lightweight scorer or manual implementation) to rank apps by name and match shortcut descriptions.
- [x] Finalise search result behaviour in app list - non-matching apps should be dimmed. 
- [x] Auto-switch detail panel to top-ranked app as search query changes.
- [x] Remove the "x matches" annotation from the detail panel footer, as it's redundant info.
- [x] If searching for "NeoVim", all apps in the list are dimmed, when only neovim should be non-dimmed.
- [x] Non-search shortcut count — In the list, rows without matches show just a bare number (35, 26). Search state shows 35 shortcuts · 13 matches. The default state would feel more consistent as 35 shortcuts instead of a lone number.
- [x] Selected app icon tint — The icon placeholder background stays #1c1c1c even when the row is selected (blue #172554). It could shift to a subtle blue tint to feel more cohesive.
- [x] Add icons for app cards and UI chrome:
    - **Simple Icons** (`simple-icons`) for brand/app logos — SVG strings extracted at build time via a static icon map (`frontend/src/lib/icons/iconMap.ts`). Currently covers: 1Password, Alacritty, Discord, Fish, Ghostty, KDE, Linear, Neovim, Notion, Obsidian, tmux, Vivaldi, Ko-fi.
    - **Lucide** (`lucide-svelte`) for UI chrome (search icon, clear X, `?` help icon, heart donate icon) — replaced all inline SVGs in `Toolbar.svelte`.
    - Custom fallback SVGs in `frontend/src/lib/icons/custom/` for apps missing from Simple Icons (e.g. Zellij, Yazi). The YAML `icon` field is resolved via the icon map; unknown names fall back to a first-letter placeholder. Import custom SVGs as raw strings (`import YaziIcon from './custom/yazi.svg?raw'`) and inject them inline via Svelte's `{@html}`. This ensures they get bundled into the binary by Vite and no external files are needed at runtime.
    - **`AppIcon.svelte`** component handles resolution: checks the icon map for a match (renders brand-colored SVG), otherwise renders a letter fallback. Used in both `AppList.svelte` and `ShortcutDetail.svelte`.
    - **Still need custom SVGs for:** LazyGit, MangoHUD, musikcube, qimgv, Qtile, VS Code, Zellij (these currently show letter fallbacks).
- [x] Implement the `?` help panel: long logo, app version. 

### Phase 3: The Editor (Writing)
- [x] Create a "New Shortcut/App" UI form.
- [x] Implement Go logic to write/update YAML files safely without destroying comments (if possible) or using a standard clean YAML structure.

### Phase 4: Polish & Deployment
- [x] Global hotkey implementation.
- [ ] Build pipeline for `.app` (Mac) and Binary (Linux).

## 5. Agent Instructions (How to help me)
When working on this project:
1. **Keep it "DevOps-y":** Ensure the YAML structure is clean and easy to edit outside the app.
2. **Prioritize Performance:** The window must appear instantly. Avoid heavy JS libraries; keep the Svelte bundle lean.
3. **Cross-Platform Awareness:** Every time a UI element is built, verify how "Meta", "Cmd", "Ctrl", and "Alt" are represented across OSs. The OS toggle in the UI must correctly switch the key arrays served to the frontend without a backend round-trip.
4. **Window Management:** If writing Go code for the window, ensure it handles "Focus Lost" events to auto-hide (or close) the app (optional but preferred).
5. **UI Reference:** The agreed UI design is a two-panel dark-theme layout. Do not revert to a card grid layout — the app list is always a vertical list.