<script lang="ts">
    import { GetApps, GetCurrentOS } from "../wailsjs/go/main/App.js";
    import { EventsOn } from "../wailsjs/runtime/runtime.js";
    import { onMount } from "svelte";
    import type { AppConfig, AppsResponse } from "./types";

    import Toolbar from "./lib/components/Toolbar.svelte";
    import AppList from "./lib/components/AppList.svelte";
    import ShortcutDetail from "./lib/components/ShortcutDetail.svelte";

    // --- State ---
    let apps: AppConfig[] = $state([]);
    let warnings: string[] = $state([]);
    let error: string = $state("");
    let selectedIndex: number = $state(0);
    let searchQuery: string = $state("");
    let currentOS: "linux" | "darwin" = $state("linux");
    let listWidth: number = $state(310);
    let isResizing: boolean = $state(false);

    // --- Derived: search matching ---
    // For each app (by original index), compute which shortcut descriptions match the search.
    let searchResults = $derived.by(() => {
        const q = searchQuery.toLowerCase().trim();
        if (!q) {
            return {
                sortedApps: apps,
                matchCounts: {} as Record<number, number>,
                matchingSets: {} as Record<number, Set<string>>,
                nameMatches: {} as Record<number, boolean>,
            };
        }

        type ScoredApp = {
            app: AppConfig;
            originalIndex: number;
            nameMatch: boolean;
            matchCount: number;
            matchingDescs: Set<string>;
        };
        const scored: ScoredApp[] = apps.map((app, i) => {
            const nameMatch = app.app.toLowerCase().includes(q);
            const matchingDescs = new Set<string>();
            for (const group of app.groups) {
                for (const shortcut of group.shortcuts) {
                    if (shortcut.desc.toLowerCase().includes(q)) {
                        matchingDescs.add(shortcut.desc);
                    }
                }
            }
            return {
                app,
                originalIndex: i,
                nameMatch,
                matchCount: matchingDescs.size,
                matchingDescs,
            };
        });

        // Sort: apps with matches first (by match count desc), then name matches, then the rest
        scored.sort((a, b) => {
            const aHasMatch = a.nameMatch || a.matchCount > 0 ? 1 : 0;
            const bHasMatch = b.nameMatch || b.matchCount > 0 ? 1 : 0;
            if (aHasMatch !== bHasMatch) return bHasMatch - aHasMatch;
            if (a.matchCount !== b.matchCount)
                return b.matchCount - a.matchCount;
            return 0;
        });

        const sortedApps = scored.map((s) => s.app);
        const matchCounts: Record<number, number> = {};
        const matchingSets: Record<number, Set<string>> = {};
        const nameMatches: Record<number, boolean> = {};
        scored.forEach((s, newIdx) => {
            matchCounts[newIdx] = s.matchCount;
            matchingSets[newIdx] = s.matchingDescs;
            nameMatches[newIdx] = s.nameMatch;
        });

        return { sortedApps, matchCounts, matchingSets, nameMatches };
    });

    let displayApps = $derived(searchResults.sortedApps);
    let selectedApp = $derived(displayApps[selectedIndex] ?? null);
    let currentMatchingDescs = $derived(
        searchResults.matchingSets[selectedIndex] ?? new Set<string>(),
    );

    // --- Auto-select first app when search changes ---
    $effect(() => {
        // Access searchQuery to track it
        searchQuery;
        selectedIndex = 0;
    });

    // --- Data loading ---
    function applyResponse(resp: AppsResponse) {
        apps = resp.apps ?? [];
        warnings = resp.warnings ?? [];
        error = "";
        // Keep selection in bounds
        if (selectedIndex >= apps.length) {
            selectedIndex = Math.max(0, apps.length - 1);
        }
    }

    function loadApps() {
        GetApps()
            .then((result: AppsResponse) => applyResponse(result))
            .catch((err: any) => {
                error = String(err);
                apps = [];
                warnings = [];
            });
    }

    // --- Keyboard navigation ---
    function handleKeydown(e: KeyboardEvent) {
        if (e.key === "ArrowDown") {
            e.preventDefault();
            if (selectedIndex < displayApps.length - 1) {
                selectedIndex++;
            }
        } else if (e.key === "ArrowUp") {
            e.preventDefault();
            if (selectedIndex > 0) {
                selectedIndex--;
            }
        } else if (e.key === "Escape") {
            if (searchQuery) {
                searchQuery = "";
            }
        }
    }

    onMount(() => {
        // Detect OS
        GetCurrentOS()
            .then((os: string) => {
                if (os === "darwin" || os === "linux") {
                    currentOS = os;
                }
            })
            .catch(() => {}); // fallback to 'linux'

        loadApps();

        const cleanup = EventsOn("config:changed", (resp: AppsResponse) => {
            applyResponse(resp);
        });

        return cleanup;
    });

    function toggleOS() {
        currentOS = currentOS === "linux" ? "darwin" : "linux";
    }

    // --- Resize handle ---
    function onResizeStart(e: PointerEvent) {
        e.preventDefault();
        isResizing = true;
        const target = e.currentTarget as HTMLElement;
        target.setPointerCapture(e.pointerId);
    }

    function onResizeMove(e: PointerEvent) {
        if (!isResizing) return;
        const newWidth = Math.min(Math.max(e.clientX, 200), 500);
        listWidth = newWidth;
    }

    function onResizeEnd() {
        isResizing = false;
    }
</script>

<svelte:window onkeydown={handleKeydown} />

<main class="app-shell">
    <Toolbar bind:searchQuery {currentOS} onToggleOS={toggleOS} />

    {#if error}
        <div class="error-bar">
            <span>{error}</span>
        </div>
    {/if}

    {#if warnings.length > 0}
        <div class="warning-bar">
            {#each warnings as warning}
                <span>{warning}</span>
            {/each}
        </div>
    {/if}

    <div class="panels" class:panels--resizing={isResizing}>
        <AppList
            apps={displayApps}
            {selectedIndex}
            {searchQuery}
            matchCounts={searchResults.matchCounts}
            nameMatches={searchResults.nameMatches}
            width={listWidth}
            onSelect={(i) => {
                selectedIndex = i;
            }}
        />
        <!-- svelte-ignore a11y_no_static_element_interactions -->
        <div
            class="resize-handle"
            onpointerdown={onResizeStart}
            onpointermove={onResizeMove}
            onpointerup={onResizeEnd}
            onpointercancel={onResizeEnd}
        ></div>
        <ShortcutDetail
            app={selectedApp}
            {currentOS}
            {searchQuery}
            matchingDescs={currentMatchingDescs}
        />
    </div>
</main>

<style>
    .app-shell {
        display: flex;
        flex-direction: column;
        height: 100vh;
        background: #121212;
        color: #e5e5e5;
        overflow: hidden;
        font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", system-ui,
            sans-serif;
    }

    .error-bar {
        padding: 6px 12px;
        background: #450a0a;
        border-bottom: 1px solid #7f1d1d;
        font-size: 12px;
        color: #fca5a5;
    }

    .warning-bar {
        padding: 6px 12px;
        background: #451a03;
        border-bottom: 1px solid #78350f;
        font-size: 12px;
        color: #fde68a;
        display: flex;
        flex-direction: column;
        gap: 2px;
    }

    .panels {
        flex: 1;
        display: flex;
        min-height: 0;
    }

    .panels--resizing {
        user-select: none;
        cursor: col-resize;
    }

    .resize-handle {
        width: 4px;
        cursor: col-resize;
        background: transparent;
        transition: background 0.15s;
        flex-shrink: 0;
    }

    .resize-handle:hover,
    .panels--resizing .resize-handle {
        background: #3a88ed;
    }
</style>
