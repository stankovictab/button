<script lang="ts">
    import type { AppConfig } from "../../types";
    import KeyBadge from "./KeyBadge.svelte";
    import AppIcon from "./AppIcon.svelte";
    import { SquarePen, Trash2, Ghost } from "lucide-svelte";

    let {
        app,
        currentOS,
        searchQuery = "",
        matchingDescs = new Set<string>(),
        onEdit,
        onDelete,
    }: {
        app: AppConfig | null;
        currentOS: "linux" | "darwin";
        searchQuery: string;
        matchingDescs: Set<string>;
        onEdit?: () => void;
        onDelete?: () => void;
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

    function formatModTime(unixSecs: number): string {
        const now = Date.now();
        const diffMs = now - unixSecs * 1000;
        const diffMins = Math.floor(diffMs / 60_000);
        const diffHours = Math.floor(diffMs / 3_600_000);

        if (diffMins < 1) return "just now";
        if (diffMins < 60) return `${diffMins}m ago`;
        if (diffHours < 24) return `${diffHours}h ago`;

        const d = new Date(unixSecs * 1000);
        const day = d.getDate();
        const month = d.toLocaleString("en", { month: "short" });
        const hh = String(d.getHours()).padStart(2, "0");
        const mm = String(d.getMinutes()).padStart(2, "0");
        return `${day} ${month} ${hh}:${mm}`;
    }
</script>

{#if app}
    <div class="detail-panel">
        <!-- Header -->
        <div class="detail-header">
            <div class="detail-header-icon">
                <AppIcon icon={app.icon} name={app.app} size={26} />
            </div>
            <div class="detail-header-info">
                <h2 class="detail-header-name">{app.app}</h2>
                <span class="detail-header-meta">
                    {totalShortcuts(app)} shortcuts &middot; {totalGroups(app)}
                    {totalGroups(app) === 1 ? "group" : "groups"}
                    {#if app.modTime}
                        &middot; updated {formatModTime(app.modTime)}
                    {/if}
                </span>
            </div>
            {#if searchQuery && matchingDescs.size > 0}
                <div class="detail-match-badge">
                    {matchingDescs.size}
                    {matchingDescs.size === 1 ? "match" : "matches"} for &ldquo;{searchQuery}&rdquo;
                </div>
            {/if}
            <div class="detail-header-actions">
                {#if onEdit}
                    <button
                        class="detail-action-btn"
                        onclick={onEdit}
                        title="Edit app"
                        aria-label="Edit app"
                    >
                        <SquarePen size={14} />
                    </button>
                {/if}
                {#if onDelete}
                    <button
                        class="detail-action-btn detail-action-btn--danger"
                        onclick={onDelete}
                        title="Delete app"
                        aria-label="Delete app"
                    >
                        <Trash2 size={14} />
                    </button>
                {/if}
            </div>
        </div>

        <!-- Shortcut groups -->
        <div class="detail-body">
            {#if totalShortcuts(app) === 0}
                <div class="detail-body-empty">
                    <Ghost size={50} />
                    <p class="detail-body-empty-text">
                        No shortcuts yet!<br />
                        Click on Edit, or add some in the app YAML config file.
                    </p>
                </div>
            {/if}
            {#each app.groups as group}
                {#if !searchQuery || matchingDescs.size === 0}
                    <!-- No active shortcut filter: show all normally -->
                    <div class="shortcut-group">
                        <div class="shortcut-group-label">{group.category}</div>
                        {#each group.shortcuts as shortcut}
                            {@const keys = resolveKeys(shortcut)}
                            <div class="shortcut-row">
                                <span class="shortcut-desc"
                                    >{shortcut.desc}</span
                                >
                                {#if keys.length > 0}
                                    <KeyBadge {keys} highlight={false} />
                                {:else}
                                    <span class="shortcut-no-keys"
                                        >Not set for {currentOS === "darwin"
                                            ? "macOS"
                                            : "Linux"}</span
                                    >
                                {/if}
                            </div>
                        {/each}
                    </div>
                {:else}
                    <!-- Shortcut filter active: show only matching shortcuts, highlight group -->
                    {@const matchingShortcuts = group.shortcuts.filter((s) =>
                        matchingDescs.has(s.desc),
                    )}
                    {#if matchingShortcuts.length > 0}
                        <div class="shortcut-group shortcut-group--match">
                            <div class="shortcut-group-label">
                                {group.category}
                            </div>
                            {#each matchingShortcuts as shortcut}
                                {@const keys = resolveKeys(shortcut)}
                                <div class="shortcut-row">
                                    <span class="shortcut-desc"
                                        >{shortcut.desc}</span
                                    >
                                    {#if keys.length > 0}
                                        <KeyBadge {keys} highlight={true} />
                                    {:else}
                                        <span class="shortcut-no-keys"
                                            >Not set for {currentOS === "darwin"
                                                ? "macOS"
                                                : "Linux"}</span
                                        >
                                    {/if}
                                </div>
                            {/each}
                        </div>
                    {/if}
                {/if}
            {/each}
        </div>

        <!-- Footer -->
        <div class="detail-footer">
            <span class="detail-footer-hint">
                <kbd class="hint-key">&uarr;&darr;</kbd> navigate
                <kbd class="hint-key">Esc</kbd> clear
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
        font-family: "JetBrains Mono", ui-monospace, SFMono-Regular, "SF Mono",
            Menlo, monospace;
    }

    .detail-header {
        display: flex;
        align-items: center;
        gap: 10px;
        padding: 12px 16px;
        border-bottom: 1px solid #1c1c1c;
        background: #0d0d0d;
    }

    .detail-header-icon {
        width: 42px;
        height: 42px;
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

    .detail-header-actions {
        display: flex;
        gap: 6px;
        flex-shrink: 0;
        margin-left: auto;
    }

    .detail-action-btn {
        display: flex;
        align-items: center;
        justify-content: center;
        width: 28px;
        height: 28px;
        padding: 0;
        border: 1px solid #2a2a2a;
        border-radius: 6px;
        background: #1c1c1c;
        color: #a1a1a1;
        cursor: pointer;
        transition: all 0.15s;
    }

    .detail-action-btn:hover {
        background: #262626;
        border-color: #3a3a3a;
        color: #d4d4d4;
    }

    .detail-action-btn:active {
        background: #1c1c1c;
        border-color: #2a2a2a;
    }

    .detail-action-btn--danger:hover {
        background: #3d1f1f;
        border-color: #5a2a2a;
        color: #ff6b6b;
    }

    .detail-body {
        flex: 1;
        overflow-y: auto;
        padding: 8px 16px;
    }

    .detail-body-empty {
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        gap: 8px;
        height: 100%;
        color: #383838;
    }

    .detail-body-empty-text {
        font-size: 13px;
        font-weight: 500;
        color: #383838;
        text-align: center;
    }

    .shortcut-group {
        margin-bottom: 16px;
    }

    .shortcut-group-label {
        font-size: 12px;
        font-weight: 500;
        color: #696969;
        padding: 0px 0 6px;
    }

    .shortcut-row {
        display: flex;
        align-items: center;
        justify-content: space-between;
        padding: 6px 8px;
        border-radius: 6px;
        border-left: 2px solid transparent;
        transition: all 0.1s;
        height: 34px;
    }

    .shortcut-row:hover {
        background: #1c1c1c;
    }

    .shortcut-group--match {
        border-left: 3px solid #3a88ed;
        border-right: 0px;
        border-top: 0px;
        background: #0b1424;
        border-radius: 9px 4px 4px 9px;
        padding: 4px 8px 4px 12px;
        margin-left: -10px;
    }

    .shortcut-group--match .shortcut-group-label {
        color: #4a82c0;
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
        font-size: 14px;
        font-weight: 600;
        font-family: "JetBrains Mono", ui-monospace, SFMono-Regular, "SF Mono",
            Menlo, monospace;
        letter-spacing: -0.04em;
        color: #333333;
    }

    .detail-footer {
        display: flex;
        align-items: center;
        justify-content: space-between;
        padding: 8px 16px;
        border-top: 1px solid #1c1c1c;
        height: 33px;
    }

    .detail-footer-hint {
        font-size: 12px;
        font-weight: 500;
        /* font-family: "JetBrains Mono", ui-monospace, SFMono-Regular, "SF Mono", */
        /* Menlo, monospace; */
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
        font-family: "JetBrains Mono", ui-monospace, SFMono-Regular, "SF Mono",
            Menlo, monospace;
        font-size: 10px;
        color: #525252;
        background: #1c1c1c;
        border: 1px solid #2a2a2a;
        border-radius: 3px;
    }
</style>
