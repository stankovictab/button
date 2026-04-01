<script lang="ts">
    import type { AppConfig } from "../../types";
    import AppIcon from "./AppIcon.svelte";

    let {
        apps,
        selectedIndex,
        searchQuery = "",
        matchCounts = {},
        nameMatches = {},
        width = 310,
        onSelect,
    }: {
        apps: AppConfig[];
        selectedIndex: number;
        searchQuery: string;
        matchCounts: Record<number, number>;
        nameMatches: Record<number, boolean>;
        width: number;
        onSelect: (index: number) => void;
    } = $props();

    function totalShortcuts(app: AppConfig): number {
        return app.groups.reduce((sum, g) => sum + g.shortcuts.length, 0);
    }
</script>

<div class="app-list" style="width: {width}px">
    <div class="app-list-header">
        <span class="app-list-header-label">APPS</span>
        <span class="app-list-header-count">{apps.length}</span>
    </div>

    <div class="app-list-items">
        {#each apps as app, i}
            {@const count = totalShortcuts(app)}
            {@const matches = matchCounts[i] ?? 0}
            {@const isSelected = i === selectedIndex}
            {@const nameMatch = nameMatches[i] ?? false}
            {@const hasNoMatches = searchQuery !== "" && matches === 0 && !nameMatch}
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
                    {#if searchQuery && matches > 0}
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
        padding: 10px 14px 6px;
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
        background: #172554 !important;
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
        color: #4597f5;
        font-weight: 600;
    }

    .app-row--selected .app-icon {
        background: #1e3a5f;
        border-color: #2a4a6f;
    }

    .app-row-right {
        display: flex;
        align-items: center;
        gap: 4px;
        flex-shrink: 0;
    }

    .app-row-meta {
        font-size: 12px;
        color: #777777;
        white-space: nowrap;
    }

    .app-row-match-count {
        color: #4597f5;
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
