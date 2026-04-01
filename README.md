<p align="center">
  <img src="assets/images/button-logo-a-long.png" alt="Button" width="400" />
</p>

<p align="center">
A cross-platform (Linux and macOS) quick-reference GUI for personal keyboard shortcuts.
</p>

<p align="center">
  <img src="assets/images/button-ui.png" alt="Button UI" width="900" />
</p>

See `PLAN.md` for full project details and roadmap.

---

## Prerequisites

- [Go](https://go.dev/) 1.21+
- [Wails v2](https://wails.io/docs/gettingstarted/installation) (`go install github.com/wailsapp/wails/v2/cmd/wails@latest`)
- [Node.js](https://nodejs.org/) (for the frontend)

---

## Build Types

There are three distinct build targets in this project, each serving a different purpose:

### 1. Go only — `go build ./...`

Compiles the Go source code only. Does **not** produce a runnable application — the frontend is not included. Useful for quickly checking that Go code compiles without errors (e.g. after editing backend logic).

```bash
go build ./...
```

### 2. Frontend only — `npm run build`

Compiles the Svelte/TypeScript frontend via Vite into `frontend/dist/`. Does **not** produce a runnable application — the Go binary is not included. Useful for checking frontend compilation errors in isolation.

```bash
cd frontend
npm install
npm run build
```

### 3. Full app — `wails build`

Compiles everything into a single self-contained binary: builds the frontend, embeds it into the Go binary, and outputs the result to `build/bin/`. This is the distributable artifact.

```bash
wails build
# Output: build/bin/button (Linux) or build/bin/button.app (macOS)
```

---

## Development

Run the app in dev mode with hot-reload:

```bash
wails dev
```

This starts:
- A Vite watcher for the frontend (changes to `.svelte`/`.ts`/`.css` files reflect immediately via HMR)
- A Go recompiler (changes to `.go` files trigger an automatic rebuild and restart)

> **Note:** The config directory `~/.config/button/apps/` is created automatically on first launch. Drop `.yaml` or `.yml` files there to populate the app. Changes are picked up live without needing a restart.

---

## Config

App shortcuts are defined as YAML files in `~/.config/button/apps/`. See `PLAN.md` Section 3A for the schema.
