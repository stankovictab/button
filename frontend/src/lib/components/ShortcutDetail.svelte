<script lang="ts">
    import type { AppConfig, Group } from "../../types";
    import KeyBadge from "./KeyBadge.svelte";
    import AppIcon from "./AppIcon.svelte";
    import {
        SquarePen,
        Trash2,
        Ghost,
        Plus,
        FilePenLine,
        FolderOpen,
    } from "lucide-svelte";
    import { OpenAppFile } from "../../../wailsjs/go/main/App";

    let {
        app,
        currentOS,
        searchQuery = "",
        matchingDescs = new Set<string>(),
        onBodyMount,
        onCreateShortcut,
        onEditShortcut,
        onDeleteShortcut,
        onEdit,
        onDelete,
    }: {
        app: AppConfig | null;
        currentOS: "linux" | "darwin" | "windows";
        searchQuery: string;
        matchingDescs: Set<string>;
        onBodyMount?: (element: HTMLDivElement | undefined) => void;
        onCreateShortcut?: () => void;
        onEditShortcut?: (groupIndex: number, shortcutIndex: number) => void;
        onDeleteShortcut?: (groupIndex: number, shortcutIndex: number) => void;
        onEdit?: () => void;
        onDelete?: () => void;
    } = $props();

    let detailBody: HTMLDivElement | undefined = $state();

    $effect(() => {
        onBodyMount?.(detailBody);
    });

    function resolveKeys(shortcut: {
        keys?: string[][];
        linux?: string[][];
        macos?: string[][];
    }): string[][] {
        if (
            (currentOS === "linux" || currentOS === "windows") &&
            shortcut.linux?.length
        )
            return shortcut.linux;
        if (currentOS === "darwin" && shortcut.macos?.length)
            return shortcut.macos;
        return shortcut.keys ?? [];
    }

    function osDisplayName(os: "linux" | "darwin" | "windows"): string {
        if (os === "darwin") return "macOS";
        if (os === "windows") return "Windows";
        return "Linux";
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

    function visibleShortcuts(group: Group) {
        return group.shortcuts
            .map((shortcut, shortcutIndex) => ({ shortcut, shortcutIndex }))
            .filter(
                ({ shortcut }) =>
                    !searchQuery ||
                    matchingDescs.size === 0 ||
                    matchingDescs.has(shortcut.desc),
            );
    }

    function showFilteredGroups(): boolean {
        return Boolean(searchQuery && matchingDescs.size > 0);
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
            <div class="detail-header-actions">
                {#if onCreateShortcut}
                    <button
                        class="detail-action-btn"
                        onclick={onCreateShortcut}
                        title="New shortcut"
                        aria-label="New shortcut"
                    >
                        <Plus size={14} />
                    </button>
                {/if}
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
                {#if app}
                    <button
                        class="detail-action-btn"
                        onclick={() => OpenAppFile(app.app)}
                        title="Open YAML file"
                        aria-label="Open YAML file"
                    >
                        <FilePenLine size={14} />
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
        <div class="detail-body" bind:this={detailBody}>
            {#if totalShortcuts(app) === 0}
                <div class="detail-body-empty">
                    <Ghost size={50} />
                    <p class="detail-body-empty-text">
                        No shortcuts yet!<br />
                        Add one here or edit the app YAML directly.
                    </p>
                </div>
            {/if}
            {#each app.groups as group, groupIndex}
                {@const shortcuts = visibleShortcuts(group)}
                {#if shortcuts.length > 0}
                    <div
                        class="shortcut-group"
                        class:shortcut-group--match={showFilteredGroups()}
                    >
                        <div class="shortcut-group-label">{group.category}</div>
                        {#each shortcuts as { shortcut, shortcutIndex }}
                            {@const binds = resolveKeys(shortcut)}
                            <div class="shortcut-row">
                                <span class="shortcut-desc"
                                    >{shortcut.desc}</span
                                >
                                <div class="shortcut-row-right">
                                    {#if binds.length > 0}
                                        <div class="shortcut-keys">
                                            {#each binds as bind, i}
                                                {#if i > 0}
                                                    <span
                                                        class="shortcut-bind-or"
                                                        >or</span
                                                    >
                                                {/if}
                                                <KeyBadge
                                                    keys={bind}
                                                    highlight={showFilteredGroups()}
                                                />
                                            {/each}
                                        </div>
                                    {:else}
                                        <span class="shortcut-no-keys"
                                            >Not set for {osDisplayName(
                                                currentOS,
                                            )}</span
                                        >
                                    {/if}
                                    <div class="shortcut-actions">
                                        {#if onEditShortcut}
                                            <button
                                                class="shortcut-action-btn"
                                                onclick={() =>
                                                    onEditShortcut(
                                                        groupIndex,
                                                        shortcutIndex,
                                                    )}
                                                title="Edit shortcut"
                                                aria-label={`Edit ${shortcut.desc}`}
                                            >
                                                <SquarePen size={14} />
                                            </button>
                                        {/if}
                                        {#if onDeleteShortcut}
                                            <button
                                                class="shortcut-action-btn shortcut-action-btn--danger"
                                                onclick={() =>
                                                    onDeleteShortcut(
                                                        groupIndex,
                                                        shortcutIndex,
                                                    )}
                                                title="Delete shortcut"
                                                aria-label={`Delete ${shortcut.desc}`}
                                            >
                                                <Trash2 size={14} />
                                            </button>
                                        {/if}
                                    </div>
                                </div>
                            </div>
                        {/each}
                    </div>
                {/if}
            {/each}
        </div>

        <!-- Footer -->
        <div class="detail-footer">
            <div
                class="detail-footer-shortcuts"
                aria-label="Keyboard shortcuts"
            >
                <span class="detail-footer-shortcut">
                    <span class="detail-footer-keys">
                        <kbd class="hint-key">j/k</kbd>
                    </span>
                    <span class="detail-footer-label">Move</span>
                </span>
                <span class="detail-footer-shortcut">
                    <span class="detail-footer-keys">
                        <kbd class="hint-key">/</kbd>
                    </span>
                    <span class="detail-footer-label">Search</span>
                </span>
                <span class="detail-footer-shortcut">
                    <span class="detail-footer-keys">
                        <kbd class="hint-key">n</kbd>
                    </span>
                    <span class="detail-footer-label">New</span>
                </span>
                <span class="detail-footer-shortcut">
                    <span class="detail-footer-keys">
                        <kbd class="hint-key">e</kbd>
                    </span>
                    <span class="detail-footer-label">Edit</span>
                </span>
                <span class="detail-footer-shortcut">
                    <span class="detail-footer-keys">
                        <kbd class="hint-key">?</kbd>
                    </span>
                    <span class="detail-footer-label">Help</span>
                </span>
            </div>
            <span class="detail-footer-custom">Who still uses the mouse?</span>
        </div>
    </div>
{:else}
    <div class="detail-panel detail-panel--empty">
        <div class="detail-body-empty">
            <FolderOpen size={50} />
            <p class="detail-body-empty-text">
                No apps loaded!<br />
                Press + in the app list, or press <code>n</code> to add an app,<br
                />
                or add YAML files to the config directory.
            </p>
        </div>
    </div>
{/if}

<style>
    .detail-panel {
        position: relative;
        flex: 1;
        display: flex;
        flex-direction: column;
        height: 100%;
        min-width: 0;
        background: #111111;
        overflow: hidden;
        isolation: isolate;
    }

    .detail-panel::before {
        content: "";
        position: absolute;
        inset: 0;
        pointer-events: none;
        background-image: radial-gradient(
            circle,
            rgba(255, 255, 255, 0.06) 1px,
            transparent 1.7px
        );
        background-size: 16px 16px;
        background-position: 0 0;
        opacity: 0.5;
        -webkit-mask-image: linear-gradient(
            to top,
            rgba(0, 0, 0, 1) 0%,
            rgba(0, 0, 0, 0.95) 28%,
            rgba(0, 0, 0, 0.35) 68%,
            rgba(0, 0, 0, 0.05) 100%
        );
        mask-image: linear-gradient(
            to top,
            rgba(0, 0, 0, 1) 0%,
            rgba(0, 0, 0, 0.95) 28%,
            rgba(0, 0, 0, 0.35) 68%,
            rgba(0, 0, 0, 0.05) 100%
        );
    }

    .detail-panel--empty {
        align-items: center;
        justify-content: center;
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
        font-weight: 500;
        color: #525252;
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
        font-size: 14px;
        font-weight: 500;
        color: #383838;
        text-align: center;
    }

    .detail-body-empty-text code {
        color: #4a4a4a;
        font-family: "JetBrains Mono", ui-monospace, SFMono-Regular, "SF Mono",
            Menlo, monospace;
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
        padding: 4px 8px;
        border-radius: 6px;
        border-left: 2px solid transparent;
        transition: all 0.1s;
        min-height: 34px;
    }

    .shortcut-row:hover {
        background: #1c1c1c;
    }

    .shortcut-group--match .shortcut-row:hover {
        background: #112645;
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
        min-width: 0;
    }

    .shortcut-row-right {
        display: flex;
        align-items: center;
        justify-content: flex-end;
        min-width: 0;
        margin-left: auto;
        flex-shrink: 0;
    }

    .shortcut-keys {
        display: flex;
        align-items: center;
        flex-wrap: wrap;
        gap: 4px;
        min-width: 0;
    }

    .shortcut-bind-or {
        font-size: 11px;
        color: #525252;
        padding: 0 2px;
        white-space: nowrap;
    }

    .shortcut-actions {
        display: flex;
        align-items: center;
        gap: 6px;
        width: 0;
        margin-left: 0;
        overflow: hidden;
        opacity: 0;
        pointer-events: none;
        transition:
            width 0.14s ease,
            margin-left 0.14s ease,
            opacity 0.14s ease;
        flex-shrink: 0;
    }

    .shortcut-row:hover .shortcut-actions,
    .shortcut-row:focus-within .shortcut-actions {
        width: 62px;
        margin-left: 10px;
        opacity: 1;
        pointer-events: auto;
    }

    .shortcut-action-btn {
        display: flex;
        align-items: center;
        justify-content: center;
        width: 28px;
        height: 28px;
        padding: 0;
        border: 1px solid #2a2a2a;
        border-radius: 6px;
        background: #171717;
        color: #808080;
        cursor: pointer;
        transition:
            background 0.12s,
            border-color 0.12s,
            color 0.12s;
        flex-shrink: 0;
    }

    .shortcut-action-btn:hover {
        background: #242424;
        border-color: #3a3a3a;
        color: #d4d4d4;
    }

    .shortcut-action-btn--danger:hover {
        background: #3d1f1f;
        border-color: #5a2a2a;
        color: #ff6b6b;
    }

    .shortcut-no-keys {
        font-size: 14px;
        font-weight: 600;
        font-family: "JetBrains Mono", ui-monospace, SFMono-Regular, "SF Mono",
            Menlo, monospace;
        letter-spacing: -0.04em;
        color: #333333;
        white-space: nowrap;
    }

    .detail-footer {
        position: relative;
        z-index: 1;
        display: flex;
        align-items: center;
        justify-content: space-between;
        gap: 12px;
        padding: 8px 8px 8px 6px;
        border-top: 1px solid #1c1c1c;
        height: 33px;
        background: #111111;
    }

    .detail-footer-shortcuts {
        display: flex;
        align-items: center;
        justify-content: flex-start;
        min-width: 0;
    }

    .detail-footer-shortcut {
        display: inline-flex;
        align-items: center;
        justify-content: center;
        gap: 6px;
        padding: 2px 4px;
        border-radius: 5px;
        /* background: #151515; */
        /* border: 1px solid #1c1c1c; */
    }

    .detail-footer-keys {
        display: inline-flex;
        align-items: center;
        gap: 3px;
    }

    .detail-footer-label {
        font-size: 11px;
        font-weight: 400;
        color: #525252;
    }

    .detail-footer-custom {
        margin-left: auto;
        font-size: 12px;
        font-weight: 500;
        color: #3f3f3f;
        white-space: nowrap;
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
