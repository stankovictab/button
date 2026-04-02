<script lang="ts">
    import type { AppConfig, SortMode } from "../../types";
    import AppIcon from "./AppIcon.svelte";
    import { ArrowDownAZ, Clock, Plus, EllipsisVertical } from "lucide-svelte";

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
        onCreateApp,
        onEditApp,
        onDeleteApp,
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
        onCreateApp: () => void;
        onEditApp: (index: number) => void;
        onDeleteApp: (index: number) => void;
    } = $props();

    let menuOpenIndex: number | null = $state(null);

    function toggleMenu(e: MouseEvent, i: number) {
        e.stopPropagation();
        menuOpenIndex = menuOpenIndex === i ? null : i;
    }

    function handleMenuEdit(e: MouseEvent, i: number) {
        e.stopPropagation();
        menuOpenIndex = null;
        onEditApp(i);
    }

    function handleMenuDelete(e: MouseEvent, i: number) {
        e.stopPropagation();
        menuOpenIndex = null;
        onDeleteApp(i);
    }

    function handleWindowClick() {
        menuOpenIndex = null;
    }

    function totalShortcuts(app: AppConfig): number {
        return app.groups.reduce((sum, g) => sum + g.shortcuts.length, 0);
    }

    let totalShortcutsAll = $derived(
        apps.reduce((sum, app) => sum + totalShortcuts(app), 0),
    );

    let rowEls: HTMLButtonElement[] = $state([]);

    $effect(() => {
        rowEls[selectedIndex]?.scrollIntoView({ block: "nearest" });
    });
</script>

<!-- svelte-ignore a11y_click_events_have_key_events a11y_no_static_element_interactions -->
<svelte:window onclick={handleWindowClick} />

<div class="app-list" style="width: {width}px">
    <div class="app-list-header">
        <span class="app-list-header-count">{apps.length}</span>
        <span class="app-list-header-label">Apps</span>
        <span class="app-list-header-sep">&bull;</span>
        <span class="app-list-header-count">{totalShortcutsAll}</span>
        <span class="app-list-header-label">Shortcuts</span>
        <button class="header-icon-btn" onclick={onCreateApp} title="New app">
            <Plus size={17} />
        </button>
        <button
            class="sort-btn"
            class:sort-btn--active={sortMode === "last-updated"}
            onclick={onToggleSort}
            title={sortMode === "alpha"
                ? "Sorted alphabetically"
                : "Sorted by last updated"}
        >
            {#if sortMode === "alpha"}
                <ArrowDownAZ size={15} />
            {:else}
                <Clock size={15} />
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
            <div class="app-row-wrapper">
                <button
                    bind:this={rowEls[i]}
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
                                <span class="app-row-app-match"
                                    >&middot; app match</span
                                >
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
                <button
                    class="menu-trigger"
                    onclick={(e) => toggleMenu(e, i)}
                    title="App actions"
                >
                    <EllipsisVertical size={15} />
                </button>
                {#if menuOpenIndex === i}
                    <!-- svelte-ignore a11y_click_events_have_key_events a11y_no_static_element_interactions -->
                    <div
                        class="context-menu"
                        onclick={(e) => e.stopPropagation()}
                    >
                        <button
                            class="context-menu-item"
                            onclick={(e) => handleMenuEdit(e, i)}>Edit</button
                        >
                        <button
                            class="context-menu-item context-menu-item--danger"
                            onclick={(e) => handleMenuDelete(e, i)}
                            >Delete</button
                        >
                    </div>
                {/if}
            </div>
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
        font-size: 12px;
        color: #3a3a3a;
    }

    .sort-btn {
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
        color: #cfcfcf;
    }

    .sort-btn--active {
        color: #4597f5;
    }

    .sort-btn--active:hover {
        color: #6aabf7;
    }

    .header-icon-btn {
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

    .header-icon-btn:hover {
        background: #1c1c1c;
        color: #cfcfcf;
    }

    .app-list-header-label {
        font-size: 12px;
        font-weight: 600;
        color: #525252;
    }

    .app-list-header-count {
        font-size: 12px;
        font-weight: 600;
        color: #4a89db;
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
        /* right padding reserves space so shortcuts text clears the trigger */
        padding: 5px 34px 5px 8px;
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

    .app-row-wrapper {
        position: relative;
    }

    .menu-trigger {
        position: absolute;
        right: 6px;
        top: 50%;
        transform: translateY(-50%);
        display: flex;
        align-items: center;
        justify-content: center;
        width: 22px;
        height: 22px;
        border-radius: 4px;
        border: none;
        background: transparent;
        color: #888888;
        cursor: pointer;
        opacity: 0;
        transition:
            opacity 0.1s,
            background 0.1s,
            color 0.1s;
        z-index: 2;
    }

    .app-row-wrapper:hover .menu-trigger {
        opacity: 1;
    }

    .menu-trigger:hover {
        background: rgba(255, 255, 255, 0.07);
        color: #a1a1a1;
    }

    .context-menu {
        position: absolute;
        top: 100%;
        right: 4px;
        background: #1c1c1c;
        border: 1px solid #2a2a2a;
        border-radius: 8px;
        padding: 4px;
        min-width: 100px;
        box-shadow: 0 8px 24px rgba(0, 0, 0, 0.5);
        z-index: 50;
    }

    .context-menu-item {
        display: block;
        width: 100%;
        padding: 6px 10px;
        border: none;
        border-radius: 5px;
        background: transparent;
        color: #d4d4d4;
        font-size: 12px;
        text-align: left;
        cursor: pointer;
        transition: background 0.1s;
    }

    .context-menu-item:hover {
        background: #252525;
    }

    .context-menu-item--danger {
        color: #f87171;
    }

    .context-menu-item--danger:hover {
        background: #2a1515;
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
