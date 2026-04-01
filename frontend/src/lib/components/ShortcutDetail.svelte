<script lang="ts">
    import type { AppConfig } from "../../types";
    import KeyBadge from "./KeyBadge.svelte";
    import AppIcon from "./AppIcon.svelte";

    let {
        app,
        currentOS,
        searchQuery = "",
        matchingDescs = new Set<string>(),
    }: {
        app: AppConfig | null;
        currentOS: "linux" | "darwin";
        searchQuery: string;
        matchingDescs: Set<string>;
    } = $props();

    function resolveKeys(shortcut: {
        keys?: string[];
        linux?: string[];
        macos?: string[];
    }): string[] {
        if (currentOS === "linux" && shortcut.linux?.length)
            return shortcut.linux;
        if (currentOS === "darwin" && shortcut.macos?.length)
            return shortcut.macos;
        return shortcut.keys ?? [];
    }

    function totalShortcuts(a: AppConfig): number {
        return a.groups.reduce((sum, g) => sum + g.shortcuts.length, 0);
    }

    function totalGroups(a: AppConfig): number {
        return a.groups.length;
    }
</script>

{#if app}
    <div class="detail-panel">
        <!-- Header -->
        <div class="detail-header">
            <div class="detail-header-icon">
                <AppIcon icon={app.icon} name={app.app} size={22} />
            </div>
            <div class="detail-header-info">
                <h2 class="detail-header-name">{app.app}</h2>
                <span class="detail-header-meta">
                    {totalShortcuts(app)} shortcuts &middot; {totalGroups(app)}
                    {totalGroups(app) === 1 ? "group" : "groups"}
                </span>
            </div>
            {#if searchQuery && matchingDescs.size > 0}
                <div class="detail-match-badge">
                    {matchingDescs.size}
                    {matchingDescs.size === 1 ? "match" : "matches"} for &ldquo;{searchQuery}&rdquo;
                </div>
            {/if}
        </div>

        <!-- Shortcut groups -->
        <div class="detail-body">
            {#each app.groups as group}
                <div class="shortcut-group">
                    <div class="shortcut-group-label">{group.category}</div>
                    {#each group.shortcuts as shortcut}
                        {@const isMatch = matchingDescs.has(shortcut.desc)}
                        {@const keys = resolveKeys(shortcut)}
                        <div
                            class="shortcut-row"
                            class:shortcut-row--match={searchQuery && isMatch}
                            class:shortcut-row--dimmed={searchQuery !== "" &&
                                !isMatch}
                        >
                            <span class="shortcut-desc">{shortcut.desc}</span>
                            {#if keys.length > 0}
                                <KeyBadge
                                    {keys}
                                    highlight={searchQuery !== "" && isMatch}
                                />
                            {:else}
                                <span class="shortcut-no-keys">--</span>
                            {/if}
                        </div>
                    {/each}
                </div>
            {/each}
        </div>

        <!-- Footer -->
        <div class="detail-footer">
            <span class="detail-footer-hint">
                <kbd class="hint-key">&uarr;&darr;</kbd> navigate
                <kbd class="hint-key">esc</kbd> clear / close
            </span>
        </div>
    </div>
{:else}
    <div class="detail-panel detail-panel--empty">
        <p class="detail-empty-text">No apps loaded</p>
        <p class="detail-empty-hint">
            Add YAML files to <code>~/.config/button/apps/</code>
        </p>
    </div>
{/if}

<style>
    .detail-panel {
        flex: 1;
        display: flex;
        flex-direction: column;
        height: 100%;
        min-width: 0;
        background: #111111;
    }

    .detail-panel--empty {
        align-items: center;
        justify-content: center;
    }

    .detail-empty-text {
        font-size: 14px;
        color: #525252;
    }

    .detail-empty-hint {
        font-size: 12px;
        color: #3f3f3f;
        margin-top: 4px;
    }

    .detail-empty-hint code {
        color: #525252;
        font-family: ui-monospace, SFMono-Regular, "SF Mono", Menlo, monospace;
    }

    .detail-header {
        display: flex;
        align-items: center;
        gap: 10px;
        padding: 12px 16px;
        border-bottom: 1px solid #1c1c1c;
    }

    .detail-header-icon {
        width: 36px;
        height: 36px;
        border-radius: 9px;
        background: #1c1c1c;
        border: 1px solid #2a2a2a;
        display: flex;
        align-items: center;
        justify-content: center;
        font-size: 15px;
        font-weight: 600;
        color: #a1a1a1;
        flex-shrink: 0;
    }

    .detail-header-info {
        flex: 1;
        min-width: 0;
    }

    .detail-header-name {
        font-size: 15px;
        font-weight: 600;
        color: #ffffff;
        margin: 0;
    }

    .detail-header-meta {
        font-size: 11px;
        color: #525252;
    }

    .detail-match-badge {
        font-size: 11px;
        color: #93c5fd;
        background: #172554;
        padding: 3px 8px;
        border-radius: 10px;
        white-space: nowrap;
        flex-shrink: 0;
    }

    .detail-body {
        flex: 1;
        overflow-y: auto;
        padding: 8px 16px;
    }

    .shortcut-group {
        margin-bottom: 16px;
    }

    .shortcut-group-label {
        font-size: 10px;
        font-weight: 600;
        letter-spacing: 0.05em;
        color: #525252;
        text-transform: uppercase;
        padding: 4px 0 6px;
    }

    .shortcut-row {
        display: flex;
        align-items: center;
        justify-content: space-between;
        padding: 6px 8px;
        border-radius: 6px;
        border-left: 2px solid transparent;
        transition: all 0.1s;
    }

    .shortcut-row:hover {
        background: #1c1c1c;
    }

    .shortcut-row--match {
        border-left-color: #3a88ed;
        background: #111827;
    }

    .shortcut-row--match:hover {
        background: #172554;
    }

    .shortcut-row--dimmed {
        opacity: 0.4;
    }

    .shortcut-desc {
        font-size: 13px;
        color: #d4d4d4;
        white-space: nowrap;
        overflow: hidden;
        text-overflow: ellipsis;
        margin-right: 12px;
    }

    .shortcut-no-keys {
        font-size: 11px;
        color: #3f3f3f;
    }

    .detail-footer {
        display: flex;
        align-items: center;
        justify-content: space-between;
        padding: 8px 16px;
        border-top: 1px solid #1c1c1c;
    }

    .detail-footer-hint {
        font-size: 11px;
        color: #3f3f3f;
        display: flex;
        align-items: center;
        gap: 6px;
    }

    .hint-key {
        display: inline-flex;
        align-items: center;
        justify-content: center;
        min-width: 18px;
        height: 16px;
        padding: 0 4px;
        font-family: ui-monospace, SFMono-Regular, "SF Mono", Menlo, monospace;
        font-size: 10px;
        color: #525252;
        background: #1c1c1c;
        border: 1px solid #2a2a2a;
        border-radius: 3px;
    }
</style>
