<script lang="ts">
    import type { AppConfig, SortMode } from "../../types";
    import AppIcon from "./AppIcon.svelte";
    import { ArrowDownAZ, Clock } from "lucide-svelte";

    let {
        apps,
        selectedIndex,
        searchQuery = "",
        matchCounts = {},
        nameMatches = {},
        width = 310,
        sortMode = "alpha",
        onSelect,
        onToggleSort,
    }: {
        apps: AppConfig[];
        selectedIndex: number;
        searchQuery: string;
        matchCounts: Record<number, number>;
        nameMatches: Record<number, boolean>;
        width: number;
        sortMode: SortMode;
        onSelect: (index: number) => void;
        onToggleSort: () => void;
    } = $props();

    function totalShortcuts(app: AppConfig): number {
        return app.groups.reduce((sum, g) => sum + g.shortcuts.length, 0);
    }

    let totalShortcutsAll = $derived(
        apps.reduce((sum, app) => sum + totalShortcuts(app), 0),
    );
</script>

<div class="app-list" style="width: {width}px">
    <div class="app-list-header">
        <span class="app-list-header-label">Apps</span>
        <span class="app-list-header-count">{apps.length}</span>
        <span class="app-list-header-sep">&middot;</span>
        <span class="app-list-header-label">Shortcuts</span>
        <span class="app-list-header-count">{totalShortcutsAll}</span>
        <button
            class="sort-btn"
            class:sort-btn--active={sortMode === "last-updated"}
            onclick={onToggleSort}
            title={sortMode === "alpha"
                ? "Sorted alphabetically"
                : "Sorted by last updated"}
        >
            {#if sortMode === "alpha"}
                <ArrowDownAZ size={13} />
            {:else}
                <Clock size={13} />
            {/if}
        </button>
    </div>

    <div class="app-list-items">
        {#each apps as app, i}
            {@const count = totalShortcuts(app)}
            {@const matches = matchCounts[i] ?? 0}
            {@const isSelected = i === selectedIndex}
            {@const nameMatch = nameMatches[i] ?? false}
            {@const hasNoMatches =
                searchQuery !== "" && matches === 0 && !nameMatch}
            <button
                class="app-row"
                class:app-row--selected={isSelected}
                class:app-row--dimmed={hasNoMatches}
                onclick={() => onSelect(i)}
            >
                <div class="app-row-left">
                    <div class="app-icon">
                        <AppIcon icon={app.icon} name={app.app} size={18} />
                    </div>
                    <span class="app-row-name">{app.app}</span>
                </div>
                <div class="app-row-right">
                    {#if searchQuery && nameMatch}
                        <span class="app-row-meta">
                            {count} shortcuts
                            <span class="app-row-app-match">&middot; app match</span>
                        </span>
                    {:else if searchQuery && matches > 0}
                        <span class="app-row-meta">
                            {count} shortcuts
                            <span class="app-row-match-count"
                                >&middot; {matches}
                                {matches === 1 ? "match" : "matches"}</span
                            >
                        </span>
                    {:else}
                        <span class="app-row-meta">{count} shortcuts</span>
                    {/if}
                </div>
            </button>
        {/each}
    </div>

    <div class="app-list-footer">
        <span class="app-list-version">v0.0.0-alpha</span>
    </div>
</div>

<style>
    .app-list {
        display: flex;
        flex-direction: column;
        min-width: 200px;
        max-width: 500px;
        border-right: 1px solid #222222;
        background: #111111;
        height: 100%;
    }

    .app-list-header {
        display: flex;
        align-items: center;
        gap: 6px;
        padding: 10px 10px 6px 14px;
    }

    .app-list-header-sep {
        font-size: 13px;
        color: #3a3a3a;
    }

    .sort-btn {
        margin-left: auto;
        display: flex;
        align-items: center;
        justify-content: center;
        width: 22px;
        height: 22px;
        border-radius: 5px;
        border: none;
        background: transparent;
        color: #525252;
        cursor: pointer;
        transition:
            background 0.1s,
            color 0.1s;
        flex-shrink: 0;
    }

    .sort-btn:hover {
        background: #1c1c1c;
        color: #a1a1a1;
    }

    .sort-btn--active {
        color: #4597f5;
    }

    .sort-btn--active:hover {
        color: #6aabf7;
    }

    .app-list-header-label {
        font-size: 10px;
        font-weight: 600;
        letter-spacing: 0.05em;
        color: #525252;
    }

    .app-list-header-count {
        font-size: 10px;
        color: #3f3f3f;
    }

    .app-list-items {
        flex: 1;
        overflow-y: auto;
        padding: 2px 6px;
    }

    .app-row {
        display: flex;
        align-items: center;
        justify-content: space-between;
        width: 100%;
        padding: 5px 8px;
        border-radius: 6px;
        border: none;
        background: transparent;
        cursor: pointer;
        transition: background 0.1s;
        text-align: left;
    }

    .app-row:hover {
        background: #1c1c1c;
    }

    .app-row--selected {
        background: linear-gradient(to right, #111111, #1e3a5f) !important;
        position: relative;
        overflow: hidden;
    }

    .app-row--selected::after {
        content: "";
        position: absolute;
        inset: 0;
        background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg'%3E%3Cfilter id='noise'%3E%3CfeTurbulence type='fractalNoise' baseFrequency='0.65' numOctaves='3' stitchTiles='stitch'/%3E%3CfeColorMatrix type='saturate' values='0'/%3E%3C/filter%3E%3Crect width='100%25' height='100%25' filter='url(%23noise)' opacity='0.08'/%3E%3C/svg%3E");
        pointer-events: none;
        border-radius: inherit;
    }

    .app-row--dimmed {
        opacity: 0.35;
    }

    .app-row-left {
        display: flex;
        align-items: center;
        gap: 8px;
        min-width: 0;
    }

    .app-icon {
        width: 30px;
        height: 30px;
        border-radius: 7px;
        background: #1c1c1c;
        border: 1px solid #2a2a2a;
        display: flex;
        align-items: center;
        justify-content: center;
        font-size: 12px;
        font-weight: 600;
        color: #a1a1a1;
        flex-shrink: 0;
    }

    .app-row-name {
        font-size: 13px;
        color: #d4d4d4;
        white-space: nowrap;
        overflow: hidden;
        text-overflow: ellipsis;
    }

    .app-row--selected .app-row-name {
        color: #ffffff;
        font-weight: 600;
    }

    .app-row--selected .app-icon {
        border-color: #2a4a6f;
    }

    .app-row-right {
        display: flex;
        align-items: center;
        gap: 4px;
        flex-shrink: 0;
    }

    .app-row-meta {
        font-size: 11px;
        color: #777777;
        white-space: nowrap;
    }

    .app-row--selected .app-row-meta {
        color: #a8c4e8;
    }

    .app-row-match-count {
        color: #4597f5;
    }

    .app-row-app-match {
        color: #f59e0b;
    }

    .app-list-footer {
        padding: 7px 14px;
        border-top: 1px solid #1c1c1c;
        display: flex;
        align-items: center;
        justify-content: center;
    }

    .app-list-version {
        font-size: 12px;
        color: #3f3f3f;
    }
</style>
