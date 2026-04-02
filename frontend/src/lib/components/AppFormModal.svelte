<script lang="ts">
    import { X, Plus, Trash2 } from "lucide-svelte";
    import type { AppConfig, Group, Shortcut } from "../../types";

    let {
        mode,
        initial = null,
        onSave,
        onCancel,
    }: {
        mode: "create" | "edit";
        initial?: AppConfig | null;
        onSave: (app: AppConfig) => void;
        onCancel: () => void;
    } = $props();

    type KeyField = "keys" | "linux" | "macos";
    type DraftField = "keysDraft" | "linuxDraft" | "macosDraft";

    type ShortcutForm = {
        desc: string;
        keys: string[];
        linux: string[];
        macos: string[];
        keysDraft: string;
        linuxDraft: string;
        macosDraft: string;
    };

    type GroupForm = {
        category: string;
        shortcuts: ShortcutForm[];
    };

    function emptyShortcut(): ShortcutForm {
        return {
            desc: "",
            keys: [],
            linux: [],
            macos: [],
            keysDraft: "",
            linuxDraft: "",
            macosDraft: "",
        };
    }

    // --- Form state (captured once at mount — initial is intentionally read only at creation) ---
    // svelte-ignore state_referenced_locally
    const initApp = initial?.app ?? "";
    // svelte-ignore state_referenced_locally
    const initIcon = initial?.icon ?? "";
    // svelte-ignore state_referenced_locally
    const initGroups: GroupForm[] = initial
        ? initial.groups.map((g) => ({
              category: g.category,
              shortcuts: g.shortcuts.map((s) => ({
                  desc: s.desc,
                  keys: [...(s.keys ?? [])],
                  linux: [...(s.linux ?? [])],
                  macos: [...(s.macos ?? [])],
                  keysDraft: "",
                  linuxDraft: "",
                  macosDraft: "",
              })),
          }))
        : [{ category: "", shortcuts: [emptyShortcut()] }];

    let appName: string = $state(initApp);
    let iconId: string = $state(initIcon);
    let groups: GroupForm[] = $state(initGroups);

    function addGroup() {
        groups = [...groups, { category: "", shortcuts: [emptyShortcut()] }];
    }

    function removeGroup(gi: number) {
        groups = groups.filter((_, i) => i !== gi);
    }

    function addShortcut(gi: number) {
        groups[gi].shortcuts = [...groups[gi].shortcuts, emptyShortcut()];
    }

    function removeShortcut(gi: number, si: number) {
        groups[gi].shortcuts = groups[gi].shortcuts.filter((_, i) => i !== si);
    }

    function draftFieldFor(field: KeyField): DraftField {
        return (field + "Draft") as DraftField;
    }

    function commitDraft(gi: number, si: number, field: KeyField) {
        const draftField = draftFieldFor(field);
        const draft = groups[gi].shortcuts[si][draftField].trim();
        if (draft) {
            groups[gi].shortcuts[si][field] = [
                ...groups[gi].shortcuts[si][field],
                draft,
            ];
            groups[gi].shortcuts[si][draftField] = "";
        }
    }

    function handleChipKeydown(
        e: KeyboardEvent,
        gi: number,
        si: number,
        field: KeyField,
    ) {
        const draftField = draftFieldFor(field);
        const draft = groups[gi].shortcuts[si][draftField].trim();

        if ((e.key === "Enter" || e.key === "Tab") && draft) {
            e.preventDefault();
            commitDraft(gi, si, field);
        } else if (
            e.key === "Backspace" &&
            !groups[gi].shortcuts[si][draftField]
        ) {
            groups[gi].shortcuts[si][field] = groups[gi].shortcuts[si][
                field
            ].slice(0, -1);
        }
    }

    function removeChip(gi: number, si: number, field: KeyField, ci: number) {
        groups[gi].shortcuts[si][field] = groups[gi].shortcuts[si][
            field
        ].filter((_, i) => i !== ci);
    }

    function handleSave() {
        // Commit any unconfirmed draft text before saving.
        for (let gi = 0; gi < groups.length; gi++) {
            for (let si = 0; si < groups[gi].shortcuts.length; si++) {
                commitDraft(gi, si, "keys");
                commitDraft(gi, si, "linux");
                commitDraft(gi, si, "macos");
            }
        }

        const appConfig: AppConfig = {
            app: appName.trim(),
            icon: iconId.trim(),
            modTime: 0,
            groups: groups
                .filter(
                    (g) =>
                        g.category.trim() ||
                        g.shortcuts.some((s) => s.desc.trim()),
                )
                .map(
                    (g): Group => ({
                        category: g.category.trim(),
                        shortcuts: g.shortcuts
                            .filter((s) => s.desc.trim())
                            .map((s): Shortcut => {
                                const shortcut: Shortcut = {
                                    desc: s.desc.trim(),
                                };
                                if (s.keys.length) shortcut.keys = s.keys;
                                if (s.linux.length) shortcut.linux = s.linux;
                                if (s.macos.length) shortcut.macos = s.macos;
                                return shortcut;
                            }),
                    }),
                ),
        };
        onSave(appConfig);
    }

    let canSave = $derived(appName.trim().length > 0);

    function handleBackdropClick(e: MouseEvent) {
        if (e.target === e.currentTarget) onCancel();
    }

    function handleKeydown(e: KeyboardEvent) {
        if (e.key === "Escape") onCancel();
    }
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
        <button class="close-btn" onclick={onCancel} aria-label="Close">
            <X size={18} />
        </button>

        <h3 class="panel-title">
            {mode === "create" ? "New App" : "Edit App"}
        </h3>

        <!-- App name & icon -->
        <div class="field-row">
            <label class="field">
                <span class="field-label">App Name</span>
                <input
                    type="text"
                    class="field-input"
                    bind:value={appName}
                    placeholder="My App"
                />
            </label>
            <label class="field field--icon">
                <span class="field-label">Icon ID</span>
                <input
                    type="text"
                    class="field-input"
                    bind:value={iconId}
                    placeholder="myapp"
                />
            </label>
        </div>

        <!-- Groups -->
        <div class="groups-section">
            <div class="groups-list">
                {#each groups as group, gi}
                    <div class="group-card">
                        <div class="group-header">
                            <div class="group-name-wrap">
                                <span class="group-name-label">Group</span>
                                <input
                                    type="text"
                                    class="field-input group-name-input"
                                    bind:value={group.category}
                                    placeholder="e.g. Navigation"
                                />
                            </div>
                            <button
                                class="icon-btn icon-btn--danger"
                                onclick={() => removeGroup(gi)}
                                title="Remove group"
                                disabled={groups.length <= 1}
                            >
                                <Trash2 size={17} />
                            </button>
                        </div>
                        <div class="shortcuts-label">Shortcuts</div>

                        <!-- Shortcuts in this group -->
                        {#each group.shortcuts as shortcut, si}
                            <div class="shortcut-entry">
                                <div class="shortcut-row-top">
                                    <input
                                        type="text"
                                        class="field-input shortcut-desc"
                                        bind:value={shortcut.desc}
                                        placeholder="Shortcut description"
                                    />
                                    <button
                                        class="icon-btn icon-btn--danger"
                                        onclick={() => removeShortcut(gi, si)}
                                        title="Remove shortcut"
                                        disabled={group.shortcuts.length <= 1}
                                    >
                                        <Trash2 size={15} />
                                    </button>
                                </div>
                                <div class="shortcut-row-keys">
                                    {#each ["keys", "linux", "macos"] as KeyField[] as field}
                                        <div class="key-field">
                                            <span class="key-field-label">
                                                {field === "keys"
                                                    ? "Keys"
                                                    : field === "linux"
                                                      ? "Linux"
                                                      : "macOS"}
                                            </span>
                                            <div class="chips-wrap">
                                                {#each shortcut[field] as key, ci}
                                                    <span class="chip">
                                                        {key}
                                                        <button
                                                            class="chip-remove"
                                                            type="button"
                                                            tabindex="-1"
                                                            onclick={() =>
                                                                removeChip(
                                                                    gi,
                                                                    si,
                                                                    field,
                                                                    ci,
                                                                )}>×</button
                                                        >
                                                    </span>
                                                {/each}
                                                <input
                                                    class="chips-input"
                                                    type="text"
                                                    bind:value={
                                                        shortcut[
                                                            draftFieldFor(field)
                                                        ]
                                                    }
                                                    onkeydown={(e) =>
                                                        handleChipKeydown(
                                                            e,
                                                            gi,
                                                            si,
                                                            field,
                                                        )}
                                                    onblur={() =>
                                                        commitDraft(
                                                            gi,
                                                            si,
                                                            field,
                                                        )}
                                                    placeholder={shortcut[field]
                                                        .length === 0
                                                        ? field === "keys"
                                                            ? "Ctrl Alt j"
                                                            : "Override"
                                                        : ""}
                                                />
                                            </div>
                                        </div>
                                    {/each}
                                </div>
                            </div>
                        {/each}

                        <button
                            class="add-shortcut-btn"
                            onclick={() => addShortcut(gi)}
                        >
                            <Plus size={12} /> Add Shortcut
                        </button>
                    </div>
                {/each}
            </div>
            <button class="add-group-btn" onclick={addGroup}>
                <Plus size={13} /> Add Group
            </button>
        </div>

        <!-- Actions -->
        <div class="actions">
            <button class="btn btn--cancel" onclick={onCancel}>Cancel</button>
            <button
                class="btn btn--primary"
                onclick={handleSave}
                disabled={!canSave}
            >
                {mode === "create" ? "Create" : "Save"}
            </button>
        </div>
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
        z-index: 200;
    }

    .panel {
        position: relative;
        background: #161616;
        border: 1px solid #2a2a2a;
        border-radius: 12px;
        padding: 24px;
        max-height: 90vh;
        max-width: 88vw;
        overflow-y: auto;
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

    .panel-title {
        font-size: 15px;
        font-weight: 600;
        color: #ffffff;
        margin: 0 0 16px;
    }

    .field-row {
        display: flex;
        gap: 10px;
        margin-bottom: 16px;
    }

    .field {
        display: flex;
        flex-direction: column;
        gap: 4px;
        flex: 1;
    }

    .field--icon {
        flex: 0 0 140px;
    }

    .field-label {
        font-size: 10px;
        font-weight: 600;
        letter-spacing: 0.05em;
        color: #525252;
    }

    .field-input {
        padding: 6px 8px;
        border-radius: 6px;
        border: 1px solid #2a2a2a;
        background: #1c1c1c;
        color: #d4d4d4;
        font-size: 13px;
        font-family: inherit;
        outline: none;
        transition: border-color 0.1s;
    }

    .field-input:focus {
        border-color: #3a88ed;
    }

    .field-input::placeholder {
        color: #3f3f3f;
    }

    .groups-section {
        margin-bottom: 20px;
    }

    .add-group-btn {
        display: flex;
        align-items: center;
        justify-content: center;
        gap: 5px;
        width: 100%;
        margin-top: 8px;
        padding: 7px 0;
        border: 1px dashed #2a2a2a;
        border-radius: 8px;
        background: transparent;
        color: #3a3a3a;
        font-size: 12px;
        cursor: pointer;
        transition:
            border-color 0.1s,
            color 0.1s;
    }

    .add-group-btn:hover {
        border-color: #3f3f3f;
        color: #777777;
    }

    .icon-btn {
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

    .icon-btn:hover {
        background: #1c1c1c;
        color: #a1a1a1;
    }

    .icon-btn:disabled {
        opacity: 0.3;
        cursor: not-allowed;
    }

    .icon-btn--danger:hover:not(:disabled) {
        background: #2a1515;
        color: #f87171;
    }

    .groups-list {
        display: flex;
        flex-direction: column;
        gap: 10px;
    }

    .group-card {
        border: 1px solid #222222;
        border-radius: 8px;
        overflow: hidden;
        background: #131313;
    }

    .group-header {
        display: flex;
        align-items: center;
        gap: 8px;
        padding: 10px 12px;
        background: #1a1a1a;
        border-bottom: 1px solid #222222;
    }

    .group-name-wrap {
        display: flex;
        align-items: center;
        gap: 10px;
        flex: 1;
        min-width: 0;
    }

    .group-name-label {
        font-size: 13px;
        font-weight: 600;
        color: #666666;
        flex-shrink: 0;
    }

    .group-name-input {
        flex: 1;
        background: #141414;
        border-color: #2a2a2a;
        padding: 4px 8px;
        font-size: 13px;
        font-weight: 500;
        color: #d4d4d4;
    }

    .group-name-input:focus {
        border-color: #3a88ed;
    }

    .shortcuts-label {
        font-size: 13px;
        font-weight: 600;
        color: #555555;
        padding: 8px 12px 2px;
    }

    .shortcut-entry {
        display: flex;
        flex-direction: column;
        gap: 6px;
        padding: 8px 12px;
        border-bottom: 1px solid #1c1c1c;
    }

    .shortcut-entry:last-of-type {
        border-bottom: none;
    }

    .shortcut-row-top {
        display: flex;
        align-items: center;
        gap: 6px;
    }

    .shortcut-desc {
        flex: 1;
    }

    .shortcut-row-keys {
        display: flex;
        gap: 6px;
    }

    .key-field {
        display: flex;
        flex-direction: column;
        gap: 3px;
        flex: 1;
    }

    .key-field-label {
        font-size: 12px;
        font-weight: 600;
        color: #555555;
        padding-left: 2px;
    }

    /* Chip input container — styled to match .field-input */
    .chips-wrap {
        display: flex;
        flex-wrap: wrap;
        align-items: center;
        gap: 4px;
        padding: 4px 6px;
        min-height: 30px;
        border-radius: 6px;
        border: 1px solid #2a2a2a;
        background: #1c1c1c;
        cursor: text;
        transition: border-color 0.1s;
    }

    .chips-wrap:focus-within {
        border-color: #3a88ed;
    }

    .chip {
        display: inline-flex;
        align-items: center;
        gap: 4px;
        padding: 1px 5px 1px 6px;
        background: #2a2a2a;
        border: 1px solid #3f3f3f;
        border-bottom: 2px solid #525252;
        border-radius: 4px;
        font-size: 12px;
        font-family: "JetBrains Mono", ui-monospace, SFMono-Regular, "SF Mono",
            Menlo, monospace;
        color: #e0e0e0;
        white-space: nowrap;
        line-height: 1.4;
        overflow: hidden;
    }

    .chip-remove {
        display: flex;
        align-items: center;
        justify-content: center;
        width: 12px;
        height: 12px;
        flex-shrink: 0;
        background: none;
        border: none;
        color: #707070;
        cursor: pointer;
        font-size: 15px;
        line-height: 1;
        padding: 0;
        margin-left: 1px;
        overflow: hidden;
    }

    .chip-remove:hover {
        color: #ff5a5a;
    }

    .chips-input {
        flex: 1;
        min-width: 32px;
        background: none;
        border: none;
        outline: none;
        color: #d4d4d4;
        font-size: 12px;
        font-weight: 500;
        font-family: "JetBrains Mono", ui-monospace, SFMono-Regular, "SF Mono",
            monospace;
        padding: 0;
    }

    .chips-input::placeholder {
        color: #3f3f3f;
    }

    .add-shortcut-btn {
        display: flex;
        align-items: center;
        gap: 4px;
        padding: 6px 12px;
        border: none;
        border-radius: 0;
        background: transparent;
        color: #3a3a3a;
        font-size: 11px;
        cursor: pointer;
        transition:
            background 0.1s,
            color 0.1s;
        width: 100%;
        border-top: 1px solid #1c1c1c;
    }

    .add-shortcut-btn:hover {
        background: #161616;
        color: #777777;
    }

    .actions {
        display: flex;
        justify-content: flex-end;
        gap: 8px;
    }

    .btn {
        padding: 6px 14px;
        border-radius: 6px;
        border: none;
        font-size: 13px;
        font-weight: 500;
        cursor: pointer;
        transition:
            background 0.1s,
            color 0.1s;
    }

    .btn:disabled {
        opacity: 0.4;
        cursor: not-allowed;
    }

    .btn--cancel {
        background: #1c1c1c;
        color: #a1a1a1;
        border: 1px solid #2a2a2a;
    }

    .btn--cancel:hover {
        background: #252525;
        color: #d4d4d4;
    }

    .btn--primary {
        background: #1e3a5f;
        color: #93c5fd;
    }

    .btn--primary:hover:not(:disabled) {
        background: #1e4a7f;
        color: #bfdbfe;
    }
</style>
