<script lang="ts">
    import { X, Download, Loader, TriangleAlert } from "lucide-svelte";
    import {
        GetRegistryApps,
        GetExistingAppFiles,
    } from "../../../wailsjs/go/main/App";
    import type { RegistryEntry } from "../../types";
    import AppIcon from "./AppIcon.svelte";

    let {
        onImport,
        onClose,
    }: {
        onImport: (filenames: string[]) => void;
        onClose: () => void;
    } = $props();

    let entries: RegistryEntry[] = $state([]);
    let existingFiles: Set<string> = $state(new Set());
    let selected: Set<string> = $state(new Set());
    let loading: boolean = $state(true);
    let error: string | null = $state(null);

    // Group entries by their first tag
    let grouped = $derived.by(() => {
        const groups: Record<string, RegistryEntry[]> = {};
        for (const entry of entries) {
            const tag =
                entry.tags && entry.tags.length > 0 ? entry.tags[0] : "Other";
            if (!groups[tag]) groups[tag] = [];
            groups[tag].push(entry);
        }
        const sorted = Object.entries(groups).sort(([a], [b]) => {
            if (a === "Other") return 1;
            if (b === "Other") return -1;
            return a.localeCompare(b);
        });
        return sorted;
    });

    let selectedCount = $derived(selected.size);

    // Count how many selected entries would override existing files
    let overrideCount = $derived.by(() => {
        let count = 0;
        for (const filename of selected) {
            if (existingFiles.has(filename)) count++;
        }
        return count;
    });

    function toggleEntry(filename: string) {
        const next = new Set(selected);
        if (next.has(filename)) {
            next.delete(filename);
        } else {
            next.add(filename);
        }
        selected = next;
    }

    function toggleGroup(groupEntries: RegistryEntry[]) {
        const allSelected = groupEntries.every((e) => selected.has(e.filename));
        const next = new Set(selected);
        for (const entry of groupEntries) {
            if (allSelected) {
                next.delete(entry.filename);
            } else {
                next.add(entry.filename);
            }
        }
        selected = next;
    }

    function selectAll() {
        selected = new Set(entries.map((e) => e.filename));
    }

    function deselectAll() {
        selected = new Set();
    }

    function handleImport() {
        if (selectedCount > 0) {
            onImport(Array.from(selected));
        }
    }

    function handleBackdropClick(e: MouseEvent) {
        if (e.target === e.currentTarget) onClose();
    }

    function handleKeydown(e: KeyboardEvent) {
        if (e.key === "Escape") onClose();
    }

    async function loadRegistry() {
        loading = true;
        error = null;
        try {
            const [registryEntries, existingFilesList] = await Promise.all([
                GetRegistryApps(),
                GetExistingAppFiles(),
            ]);
            entries = registryEntries;
            existingFiles = new Set(existingFilesList ?? []);
            // Pre-select only entries that don't already exist
            selected = new Set(
                entries
                    .filter((e) => !existingFiles.has(e.filename))
                    .map((e) => e.filename),
            );
        } catch (err) {
            error = String(err);
        } finally {
            loading = false;
        }
    }

    $effect(() => {
        loadRegistry();
    });
</script>

<svelte:window onkeydown={handleKeydown} />

<!-- svelte-ignore a11y_click_events_have_key_events a11y_no_static_element_interactions a11y_no_noninteractive_element_interactions -->
<div
    class="backdrop"
    role="dialog"
    aria-modal="true"
    tabindex="-1"
    onclick={handleBackdropClick}
>
    <div class="panel">
        <button class="close-btn" onclick={onClose} aria-label="Close">
            <X size={14} />
        </button>

        <div class="header">
            <h2 class="title">Import Apps</h2>
            <p class="subtitle">
                Choose apps from the registry to add to your collection. You can
                edit them anytime after importing.
            </p>
        </div>

        <div class="divider"></div>

        <div class="body">
            {#if loading}
                <div class="state-message">
                    <Loader size={18} class="spinner" />
                    <span>Loading registry...</span>
                </div>
            {:else if error}
                <div class="state-message state-error">
                    <span>Failed to load registry: {error}</span>
                    <button class="btn-retry" onclick={() => loadRegistry()}>
                        Retry
                    </button>
                </div>
            {:else if entries.length === 0}
                <div class="state-message">
                    <span>No apps available in the registry yet.</span>
                </div>
            {:else}
                <div class="select-actions">
                    <button class="btn-text" onclick={selectAll}
                        >Select all</button
                    >
                    <span class="select-sep">/</span>
                    <button class="btn-text" onclick={deselectAll}
                        >Deselect all</button
                    >
                </div>

                {#each grouped as [tag, groupEntries]}
                    <div class="group">
                        <button
                            class="group-header"
                            onclick={() => toggleGroup(groupEntries)}
                        >
                            <span class="group-tag">{tag}</span>
                            <span class="group-count"
                                >{groupEntries.filter((e) =>
                                    selected.has(e.filename),
                                ).length}/{groupEntries.length}</span
                            >
                        </button>
                        <div class="group-items">
                            {#each groupEntries as entry}
                                <label
                                    class="entry-row"
                                    class:entry-row--exists={existingFiles.has(
                                        entry.filename,
                                    )}
                                >
                                    <input
                                        type="checkbox"
                                        checked={selected.has(entry.filename)}
                                        onchange={() =>
                                            toggleEntry(entry.filename)}
                                    />
                                    <AppIcon
                                        icon={entry.icon}
                                        name={entry.app}
                                        size={20}
                                    />
                                    <span class="entry-name">{entry.app}</span>
                                    {#if existingFiles.has(entry.filename)}
                                        <span class="exists-badge"
                                            >Installed</span
                                        >
                                    {/if}
                                </label>
                            {/each}
                        </div>
                    </div>
                {/each}
            {/if}
        </div>

        {#if !loading && !error && entries.length > 0}
            <div class="divider"></div>

            {#if overrideCount > 0}
                <div class="override-warning">
                    <TriangleAlert size={13} />
                    <span
                        >{overrideCount} selected {overrideCount === 1
                            ? "app already exists"
                            : "apps already exist"} and will be overwritten.</span
                    >
                </div>
            {/if}

            <div class="footer">
                <button
                    class="btn-primary"
                    disabled={selectedCount === 0}
                    onclick={handleImport}
                >
                    <Download size={14} />
                    Import{selectedCount > 0 ? ` (${selectedCount})` : ""}
                </button>
                <button class="btn-cancel" onclick={onClose}>Cancel</button>
            </div>
        {/if}
    </div>
</div>

<style>
    .backdrop {
        position: fixed;
        inset: 0;
        background: rgba(0, 0, 0, 0.6);
        display: flex;
        align-items: center;
        justify-content: center;
        z-index: 100;
    }

    .panel {
        position: relative;
        background: #161616;
        border: 1px solid #2a2a2a;
        border-radius: 12px;
        padding: 24px 24px 20px;
        width: 400px;
        max-height: 80vh;
        display: flex;
        flex-direction: column;
        box-shadow: 0 24px 48px rgba(0, 0, 0, 0.6);
    }

    .close-btn {
        position: absolute;
        top: 12px;
        right: 12px;
        display: flex;
        align-items: center;
        justify-content: center;
        width: 22px;
        height: 22px;
        background: transparent;
        border: none;
        border-radius: 5px;
        color: #525252;
        cursor: pointer;
        transition:
            background 0.1s,
            color 0.1s;
    }

    .close-btn:hover {
        background: #1c1c1c;
        color: #a1a1a1;
    }

    .header {
        padding-bottom: 12px;
    }

    .title {
        font-size: 16px;
        font-weight: 700;
        color: #ffffff;
        margin: 0 0 6px;
    }

    .subtitle {
        font-size: 12px;
        color: #777777;
        margin: 0;
        line-height: 1.4;
    }

    .divider {
        height: 1px;
        background: #1c1c1c;
        margin: 0 0 12px;
        flex-shrink: 0;
    }

    .body {
        flex: 1;
        overflow-y: auto;
        min-height: 0;
        padding-right: 4px;
    }

    .body::-webkit-scrollbar {
        width: 6px;
    }

    .body::-webkit-scrollbar-track {
        background: transparent;
    }

    .body::-webkit-scrollbar-thumb {
        background: #2a2a2a;
        border-radius: 3px;
    }

    .state-message {
        display: flex;
        flex-direction: column;
        align-items: center;
        gap: 10px;
        padding: 24px 0;
        font-size: 13px;
        color: #777777;
    }

    .state-error {
        color: #ef4444;
    }

    .btn-retry {
        padding: 6px 14px;
        background: #1c1c1c;
        border: 1px solid #2a2a2a;
        border-radius: 6px;
        color: #a1a1a1;
        font-size: 12px;
        cursor: pointer;
        transition: all 0.15s;
    }

    .btn-retry:hover {
        background: #222222;
        color: #ffffff;
    }

    .select-actions {
        display: flex;
        align-items: center;
        gap: 4px;
        margin-bottom: 12px;
    }

    .btn-text {
        background: none;
        border: none;
        color: #4597f5;
        font-size: 11px;
        cursor: pointer;
        padding: 0;
        transition: color 0.15s;
    }

    .btn-text:hover {
        color: #93c5fd;
    }

    .select-sep {
        color: #3f3f3f;
        font-size: 11px;
    }

    .group {
        margin-bottom: 12px;
    }

    .group-header {
        display: flex;
        align-items: center;
        justify-content: space-between;
        width: 100%;
        padding: 6px 8px;
        background: #1a1a1a;
        border: 1px solid #222222;
        border-radius: 6px;
        cursor: pointer;
        transition: background 0.1s;
        margin-bottom: 4px;
    }

    .group-header:hover {
        background: #1e1e1e;
    }

    .group-tag {
        font-size: 11px;
        font-weight: 600;
        color: #a1a1a1;
        text-transform: uppercase;
        letter-spacing: 0.04em;
    }

    .group-count {
        font-size: 10px;
        color: #525252;
    }

    .group-items {
        display: flex;
        flex-direction: column;
        gap: 1px;
    }

    .entry-row {
        display: flex;
        align-items: center;
        gap: 8px;
        padding: 7px 8px;
        border-radius: 6px;
        cursor: pointer;
        transition: background 0.1s;
    }

    .entry-row:hover {
        background: #1c1c1c;
    }

    .entry-row input[type="checkbox"] {
        width: 14px;
        height: 14px;
        accent-color: #3a88ed;
        cursor: pointer;
        flex-shrink: 0;
    }

    .entry-name {
        font-size: 13px;
        color: #e5e5e5;
        font-weight: 500;
        flex: 1;
    }

    .entry-row--exists .entry-name {
        color: #777777;
    }

    .exists-badge {
        font-size: 10px;
        font-weight: 500;
        color: #525252;
        background: #1c1c1c;
        border: 1px solid #2a2a2a;
        padding: 1px 6px;
        border-radius: 8px;
        flex-shrink: 0;
    }

    .override-warning {
        display: flex;
        align-items: center;
        gap: 6px;
        padding: 8px 10px;
        margin-bottom: 10px;
        background: #2a1a0a;
        border: 1px solid #3d2e1a;
        border-radius: 6px;
        font-size: 11px;
        color: #f59e0b;
    }

    .footer {
        display: flex;
        align-items: center;
        gap: 10px;
        padding-top: 4px;
    }

    .btn-primary {
        display: flex;
        align-items: center;
        justify-content: center;
        gap: 6px;
        flex: 1;
        padding: 10px 16px;
        background: #1e3a5f;
        border: 1px solid #2a5a8f;
        border-radius: 8px;
        color: #93c5fd;
        font-size: 13px;
        font-weight: 600;
        cursor: pointer;
        transition: all 0.15s;
    }

    .btn-primary:hover:not(:disabled) {
        background: #264b77;
        color: #bfdbfe;
    }

    .btn-primary:disabled {
        opacity: 0.4;
        cursor: not-allowed;
    }

    .btn-cancel {
        background: none;
        border: none;
        color: #525252;
        font-size: 12px;
        cursor: pointer;
        padding: 4px 8px;
        transition: color 0.15s;
    }

    .btn-cancel:hover {
        color: #a1a1a1;
    }

    @keyframes spin {
        from {
            transform: rotate(0deg);
        }
        to {
            transform: rotate(360deg);
        }
    }

    .state-message :global(.spinner) {
        animation: spin 1s linear infinite;
    }
</style>
