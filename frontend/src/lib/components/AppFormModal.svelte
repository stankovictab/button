<script lang="ts">
    import { onMount, tick } from "svelte";
    import {
        X,
        Plus,
        Trash2,
        ChevronUp,
        ChevronDown,
        Info,
    } from "lucide-svelte";
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
    type ShortcutForm = {
        desc: string;
        keys: string[][];
        linux: string[][];
        macos: string[][];
        keysDrafts: string[];
        linuxDrafts: string[];
        macosDrafts: string[];
    };
    type DraftField = "keysDrafts" | "linuxDrafts" | "macosDrafts";

    type GroupForm = {
        category: string;
        shortcuts: ShortcutForm[];
    };

    function emptyShortcut(): ShortcutForm {
        return {
            desc: "",
            keys: [[]],
            linux: [[]],
            macos: [[]],
            keysDrafts: [""],
            linuxDrafts: [""],
            macosDrafts: [""],
        };
    }

    // --- Form state (captured once at mount — initial is intentionally read only at creation) ---
    // svelte-ignore state_referenced_locally
    const initApp = initial?.app ?? "";
    // svelte-ignore state_referenced_locally
    const initIcon = initial?.icon ?? "";
    // svelte-ignore state_referenced_locally
    const initTags = initial?.tags ? [...initial.tags] : [];
    // svelte-ignore state_referenced_locally
    const initGroups: GroupForm[] = initial
        ? initial.groups.map((g) => ({
              category: g.category,
              shortcuts: g.shortcuts.map((s) => {
                  const k = s.keys?.length ? s.keys.map((b) => [...b]) : [[]];
                  const l = s.linux?.length ? s.linux.map((b) => [...b]) : [[]];
                  const m = s.macos?.length ? s.macos.map((b) => [...b]) : [[]];
                  return {
                      desc: s.desc,
                      keys: k,
                      linux: l,
                      macos: m,
                      keysDrafts: k.map(() => ""),
                      linuxDrafts: l.map(() => ""),
                      macosDrafts: m.map(() => ""),
                  };
              }),
          }))
        : [{ category: "", shortcuts: [emptyShortcut()] }];

    let appName: string = $state(initApp);
    let iconId: string = $state(initIcon);
    let tags: string[] = $state(initTags);
    let tagDraft: string = $state("");
    let groups: GroupForm[] = $state(initGroups);
    let recordingSlot: string | null = $state(null);
    let panelEl: HTMLDivElement | undefined = $state();
    let appNameInput: HTMLInputElement | undefined = $state();
    let tagsInput: HTMLInputElement | undefined = $state();

    function addGroup() {
        groups = [...groups, { category: "", shortcuts: [emptyShortcut()] }];
    }

    function removeGroup(gi: number) {
        groups = groups.filter((_, i) => i !== gi);
    }

    function appendTags(rawValue: string) {
        const nextTags = rawValue
            .split(",")
            .map((tag) => tag.trim())
            .filter(Boolean);

        if (nextTags.length === 0) return;

        const seen = new Set(tags.map((tag) => tag.toLowerCase()));
        const uniqueNewTags = nextTags.filter((tag) => {
            const normalized = tag.toLowerCase();
            if (seen.has(normalized)) return false;
            seen.add(normalized);
            return true;
        });

        if (uniqueNewTags.length > 0) {
            tags = [...tags, ...uniqueNewTags];
        }
    }

    function commitTagDraft() {
        if (!tagDraft.trim()) return;
        appendTags(tagDraft);
        tagDraft = "";
    }

    function handleTagKeydown(e: KeyboardEvent) {
        if ((e.key === "Enter" || e.key === ",") && tagDraft.trim()) {
            e.preventDefault();
            commitTagDraft();
            return;
        }

        if (e.key === "Backspace" && !tagDraft && tags.length > 0) {
            tags = tags.slice(0, -1);
        }
    }

    function handleTagBlur() {
        commitTagDraft();
    }

    function removeTag(tagIndex: number) {
        tags = tags.filter((_, index) => index !== tagIndex);
    }

    async function addShortcut(gi: number) {
        groups[gi].shortcuts = [...groups[gi].shortcuts, emptyShortcut()];
        await tick();
        const newSi = groups[gi].shortcuts.length - 1;
        panelEl
            ?.querySelector<HTMLInputElement>(
                `[data-gi="${gi}"][data-si="${newSi}"].shortcut-desc`,
            )
            ?.focus();
    }

    function removeShortcut(gi: number, si: number) {
        groups[gi].shortcuts = groups[gi].shortcuts.filter((_, i) => i !== si);
    }

    function moveShortcutUp(gi: number, si: number) {
        if (si > 0) {
            const s = groups[gi].shortcuts;
            const updated = [...s];
            [updated[si - 1], updated[si]] = [updated[si], updated[si - 1]];
            groups[gi].shortcuts = updated;
        } else if (gi > 0) {
            const shortcut = groups[gi].shortcuts[0];
            groups[gi].shortcuts = groups[gi].shortcuts.slice(1);
            groups[gi - 1].shortcuts = [...groups[gi - 1].shortcuts, shortcut];
        }
    }

    function moveShortcutDown(gi: number, si: number) {
        const len = groups[gi].shortcuts.length;
        if (si < len - 1) {
            const s = groups[gi].shortcuts;
            const updated = [...s];
            [updated[si], updated[si + 1]] = [updated[si + 1], updated[si]];
            groups[gi].shortcuts = updated;
        } else if (gi < groups.length - 1) {
            const shortcut = groups[gi].shortcuts[len - 1];
            groups[gi].shortcuts = groups[gi].shortcuts.slice(0, -1);
            groups[gi + 1].shortcuts = [shortcut, ...groups[gi + 1].shortcuts];
        }
    }

    function draftFieldFor(field: KeyField): DraftField {
        return (field + "Drafts") as DraftField;
    }

    // --- Recording mode helpers ---

    function slotKey(
        gi: number,
        si: number,
        field: KeyField,
        bindIndex: number,
    ): string {
        return `${gi}-${si}-${field}-${bindIndex}`;
    }

    function isRecording(
        gi: number,
        si: number,
        field: KeyField,
        bindIndex: number,
    ): boolean {
        return recordingSlot === slotKey(gi, si, field, bindIndex);
    }

    async function toggleRecording(
        gi: number,
        si: number,
        field: KeyField,
        bindIndex: number,
    ) {
        const key = slotKey(gi, si, field, bindIndex);
        if (recordingSlot === key) {
            recordingSlot = null;
        } else {
            recordingSlot = key;
            await tick();
            panelEl
                ?.querySelector<HTMLInputElement>(
                    `[data-gi="${gi}"][data-si="${si}"][data-field="${field}"][data-bind-index="${bindIndex}"]`,
                )
                ?.focus();
        }
    }

    // --- Key handling (dual mode) ---

    const MODIFIER_KEYS = new Set(["Control", "Shift", "Alt", "Meta"]);

    function normalizeKeyName(e: KeyboardEvent): string {
        const key = e.key;
        if (key === " ") return "Space";
        if (key.length === 1) return key.toLowerCase();
        return key;
    }

    function commitDraft(
        gi: number,
        si: number,
        field: KeyField,
        bindIndex: number,
    ) {
        const draftField = draftFieldFor(field);
        const drafts = groups[gi].shortcuts[si][draftField] as string[];
        const draft = drafts[bindIndex].trim();
        if (!draft) return;
        groups[gi].shortcuts[si][field] = groups[gi].shortcuts[si][field].map(
            (b, i) => (i === bindIndex ? [...b, draft] : b),
        );
        const newDrafts = [...drafts];
        newDrafts[bindIndex] = "";
        groups[gi].shortcuts[si][draftField] = newDrafts as never;
    }

    function handleChipKeydown(
        e: KeyboardEvent,
        gi: number,
        si: number,
        field: KeyField,
        bindIndex: number,
    ) {
        if (isRecording(gi, si, field, bindIndex)) {
            // --- Record mode ---
            if (e.key === "Tab") return;

            if (
                e.key === "Escape" &&
                !e.ctrlKey &&
                !e.altKey &&
                !e.shiftKey &&
                !e.metaKey
            ) {
                recordingSlot = null;
                e.preventDefault();
                e.stopPropagation();
                return;
            }

            e.preventDefault();
            if (MODIFIER_KEYS.has(e.key)) return;

            if (
                e.key === "Backspace" &&
                !e.ctrlKey &&
                !e.altKey &&
                !e.shiftKey &&
                !e.metaKey
            ) {
                groups[gi].shortcuts[si][field] = groups[gi].shortcuts[si][
                    field
                ].map((b, i) => (i === bindIndex ? [] : b));
                return;
            }

            const combo: string[] = [];
            if (e.ctrlKey) combo.push("Ctrl");
            if (e.altKey) combo.push("Alt");
            if (e.shiftKey) combo.push("Shift");
            if (e.metaKey) combo.push("Cmd");
            combo.push(normalizeKeyName(e));

            groups[gi].shortcuts[si][field] = groups[gi].shortcuts[si][
                field
            ].map((b, i) => (i === bindIndex ? combo : b));
            return;
        }

        // --- Text mode ---
        const draftField = draftFieldFor(field);
        const drafts = groups[gi].shortcuts[si][draftField] as string[];
        const draft = drafts[bindIndex].trim();

        if (e.key === "Enter" && draft) {
            e.preventDefault();
            commitDraft(gi, si, field, bindIndex);
        } else if (e.key === "Backspace" && !drafts[bindIndex]) {
            groups[gi].shortcuts[si][field] = groups[gi].shortcuts[si][
                field
            ].map((b, i) => (i === bindIndex ? b.slice(0, -1) : b));
        }
    }

    function handleChipBlur(
        gi: number,
        si: number,
        field: KeyField,
        bindIndex: number,
    ) {
        if (isRecording(gi, si, field, bindIndex)) {
            recordingSlot = null;
        } else {
            commitDraft(gi, si, field, bindIndex);
        }
    }

    function removeChip(
        gi: number,
        si: number,
        field: KeyField,
        bindIndex: number,
        ci: number,
    ) {
        groups[gi].shortcuts[si][field] = groups[gi].shortcuts[si][field].map(
            (b, i) => (i === bindIndex ? b.filter((_, j) => j !== ci) : b),
        );
    }

    async function addBind(gi: number, si: number, field: KeyField) {
        const newBindIndex = groups[gi].shortcuts[si][field].length;
        groups[gi].shortcuts[si][field] = [
            ...groups[gi].shortcuts[si][field],
            [],
        ];
        const draftField = draftFieldFor(field);
        groups[gi].shortcuts[si][draftField] = [
            ...(groups[gi].shortcuts[si][draftField] as string[]),
            "",
        ] as never;
        await tick();
        panelEl
            ?.querySelector<HTMLInputElement>(
                `[data-gi="${gi}"][data-si="${si}"][data-field="${field}"][data-bind-index="${newBindIndex}"]`,
            )
            ?.focus();
    }

    function removeBind(
        gi: number,
        si: number,
        field: KeyField,
        bindIndex: number,
    ) {
        recordingSlot = null;
        const binds = groups[gi].shortcuts[si][field];
        const draftField = draftFieldFor(field);
        const drafts = groups[gi].shortcuts[si][draftField] as string[];
        if (binds.length <= 1) {
            groups[gi].shortcuts[si][field] = [[]];
            groups[gi].shortcuts[si][draftField] = [""] as never;
            return;
        }
        groups[gi].shortcuts[si][field] = binds.filter(
            (_, i) => i !== bindIndex,
        );
        groups[gi].shortcuts[si][draftField] = drafts.filter(
            (_, i) => i !== bindIndex,
        ) as never;
    }

    function handleSave() {
        commitTagDraft();

        for (let gi = 0; gi < groups.length; gi++) {
            for (let si = 0; si < groups[gi].shortcuts.length; si++) {
                for (const field of ["keys", "linux", "macos"] as KeyField[]) {
                    groups[gi].shortcuts[si][field].forEach((_, bindIndex) => {
                        commitDraft(gi, si, field, bindIndex);
                    });
                }
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
                                const k = s.keys.filter((b) => b.length > 0);
                                const l = s.linux.filter((b) => b.length > 0);
                                const m = s.macos.filter((b) => b.length > 0);
                                if (k.length) shortcut.keys = k;
                                if (l.length) shortcut.linux = l;
                                if (m.length) shortcut.macos = m;
                                return shortcut;
                            }),
                    }),
                ),
        };
        if (tags.length > 0) appConfig.tags = tags;
        onSave(appConfig);
    }

    let canSave = $derived(appName.trim().length > 0);

    function handleBackdropClick(e: MouseEvent) {
        if (e.target === e.currentTarget) onCancel();
    }

    function isEditableTarget(target: EventTarget | null): boolean {
        if (!(target instanceof HTMLElement)) return false;
        return (
            target instanceof HTMLInputElement ||
            target instanceof HTMLTextAreaElement ||
            target instanceof HTMLSelectElement ||
            target.isContentEditable
        );
    }

    function blurActiveEditable() {
        const active = document.activeElement;
        if (!isEditableTarget(active)) return false;
        (active as HTMLElement).blur();
        return true;
    }

    function getFocusableElements(): HTMLElement[] {
        if (!panelEl) return [];

        return Array.from(
            panelEl.querySelectorAll<HTMLElement>(
                'button:not([disabled]):not([tabindex="-1"]), input:not([disabled]), textarea:not([disabled]), select:not([disabled]), a[href], [tabindex]:not([tabindex="-1"])',
            ),
        ).filter(
            (el) =>
                !el.hasAttribute("disabled") &&
                el.tabIndex !== -1 &&
                el.offsetParent !== null,
        );
    }

    async function focusFirstField() {
        await tick();
        appNameInput?.focus();
        appNameInput?.select();
    }

    onMount(() => {
        void focusFirstField();
    });

    function handleKeydown(e: KeyboardEvent) {
        const key = e.key.toLowerCase();
        const editableTarget = isEditableTarget(e.target);

        if (e.key === "Escape") {
            if (blurActiveEditable()) {
                e.preventDefault();
                return;
            }
            onCancel();
            return;
        }

        if (e.key === "Tab") {
            const focusable = getFocusableElements();
            if (focusable.length === 0) return;

            const active = document.activeElement as HTMLElement | null;
            const currentIndex = active ? focusable.indexOf(active) : -1;

            if (currentIndex === -1) {
                e.preventDefault();
                focusable[0]?.focus();
                return;
            }

            if (!e.shiftKey && currentIndex === focusable.length - 1) {
                e.preventDefault();
                focusable[0]?.focus();
                return;
            }

            if (e.shiftKey && currentIndex === 0) {
                e.preventDefault();
                focusable[focusable.length - 1]?.focus();
            }
            return;
        }

        if (
            !editableTarget &&
            !e.ctrlKey &&
            !e.metaKey &&
            !e.altKey &&
            key === "j"
        ) {
            e.preventDefault();
            panelEl?.scrollBy({ top: 180, behavior: "smooth" });
            return;
        }

        if (
            !editableTarget &&
            !e.ctrlKey &&
            !e.metaKey &&
            !e.altKey &&
            key === "k"
        ) {
            e.preventDefault();
            panelEl?.scrollBy({ top: -180, behavior: "smooth" });
        }
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
    <div class="panel" bind:this={panelEl}>
        <div class="panel-topbar">
            <h3 class="panel-title">
                {mode === "create" ? "New App" : "Edit App"}
            </h3>
            <div class="panel-topbar-actions">
                <button class="btn btn--cancel" onclick={onCancel}>
                    Cancel
                </button>
                <button
                    class="btn btn--primary"
                    onclick={handleSave}
                    disabled={!canSave}
                >
                    {mode === "create" ? "Create" : "Save"}
                </button>
                <button class="close-btn" onclick={onCancel} aria-label="Close">
                    <X size={18} />
                </button>
            </div>
        </div>

        <div class="panel-content">
            <!-- App name & icon -->
            <div class="field-row">
                <label class="field">
                    <span class="field-label">App Name</span>
                    <input
                        bind:this={appNameInput}
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

            <div class="field field--full">
                <label class="field-label" for="app-tags-input">Tags</label>
                <div class="tag-editor">
                    {#each tags as tag, index}
                        <span class="tag-chip">
                            {tag}
                            <button
                                class="tag-chip-remove"
                                type="button"
                                tabindex="-1"
                                onclick={() => removeTag(index)}
                                aria-label={`Remove ${tag} tag`}
                            >
                                ×
                            </button>
                        </span>
                    {/each}
                    <input
                        bind:this={tagsInput}
                        id="app-tags-input"
                        type="text"
                        class="tag-input"
                        bind:value={tagDraft}
                        placeholder={tags.length === 0
                            ? "Type a tag and press Enter"
                            : "Add tag"}
                        onkeydown={handleTagKeydown}
                        onblur={handleTagBlur}
                    />
                </div>
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
                            <div class="shortcuts-label">
                                Shortcuts
                                <span class="info-hint">
                                    <Info size={12} />
                                    <span class="info-tooltip"
                                        >Type each key name and press Enter to
                                        add it, or click the record button to
                                        capture a shortcut from your keyboard.</span
                                    >
                                </span>
                            </div>

                            <!-- Shortcuts in this group -->
                            {#each group.shortcuts as shortcut, si}
                                <div class="shortcut-entry">
                                    <div class="shortcut-row-top">
                                        <input
                                            type="text"
                                            class="field-input shortcut-desc"
                                            data-gi={gi}
                                            data-si={si}
                                            bind:value={shortcut.desc}
                                            placeholder="Shortcut description"
                                        />
                                        <button
                                            class="icon-btn"
                                            onclick={() =>
                                                moveShortcutUp(gi, si)}
                                            title="Move up"
                                            disabled={gi === 0 && si === 0}
                                            tabindex="-1"
                                        >
                                            <ChevronUp size={15} />
                                        </button>
                                        <button
                                            class="icon-btn"
                                            onclick={() =>
                                                moveShortcutDown(gi, si)}
                                            title="Move down"
                                            disabled={gi ===
                                                groups.length - 1 &&
                                                si ===
                                                    group.shortcuts.length - 1}
                                            tabindex="-1"
                                        >
                                            <ChevronDown size={15} />
                                        </button>
                                        <button
                                            class="icon-btn icon-btn--danger"
                                            onclick={() =>
                                                removeShortcut(gi, si)}
                                            title="Remove shortcut"
                                            disabled={group.shortcuts.length <=
                                                1}
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
                                                          ? "Linux / Windows"
                                                          : "macOS"}
                                                </span>
                                                <div class="key-binds">
                                                    {#each shortcut[field] as bind, bindIndex}
                                                        {#if bindIndex > 0}
                                                            <div
                                                                class="bind-or-label"
                                                            >
                                                                or
                                                            </div>
                                                        {/if}
                                                        <div class="bind-row">
                                                            <div
                                                                class="chips-wrap"
                                                                class:recording={isRecording(
                                                                    gi,
                                                                    si,
                                                                    field,
                                                                    bindIndex,
                                                                )}
                                                            >
                                                                {#each bind as key, ci}
                                                                    <span
                                                                        class="chip"
                                                                    >
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
                                                                                    bindIndex,
                                                                                    ci,
                                                                                )}
                                                                        >
                                                                            ×
                                                                        </button>
                                                                    </span>
                                                                {/each}
                                                                <input
                                                                    class="chips-input"
                                                                    type="text"
                                                                    readonly={isRecording(
                                                                        gi,
                                                                        si,
                                                                        field,
                                                                        bindIndex,
                                                                    )}
                                                                    data-gi={gi}
                                                                    data-si={si}
                                                                    data-field={field}
                                                                    data-bind-index={bindIndex}
                                                                    value={isRecording(
                                                                        gi,
                                                                        si,
                                                                        field,
                                                                        bindIndex,
                                                                    )
                                                                        ? ""
                                                                        : (
                                                                              shortcut[
                                                                                  draftFieldFor(
                                                                                      field,
                                                                                  )
                                                                              ] as string[]
                                                                          )[
                                                                              bindIndex
                                                                          ]}
                                                                    oninput={(
                                                                        e,
                                                                    ) => {
                                                                        const newDrafts =
                                                                            [
                                                                                ...(groups[
                                                                                    gi
                                                                                ]
                                                                                    .shortcuts[
                                                                                    si
                                                                                ][
                                                                                    draftFieldFor(
                                                                                        field,
                                                                                    )
                                                                                ] as string[]),
                                                                            ];
                                                                        newDrafts[
                                                                            bindIndex
                                                                        ] = (
                                                                            e.currentTarget as HTMLInputElement
                                                                        ).value;
                                                                        groups[
                                                                            gi
                                                                        ].shortcuts[
                                                                            si
                                                                        ][
                                                                            draftFieldFor(
                                                                                field,
                                                                            )
                                                                        ] =
                                                                            newDrafts as never;
                                                                    }}
                                                                    onkeydown={(
                                                                        e,
                                                                    ) =>
                                                                        handleChipKeydown(
                                                                            e,
                                                                            gi,
                                                                            si,
                                                                            field,
                                                                            bindIndex,
                                                                        )}
                                                                    onblur={() =>
                                                                        handleChipBlur(
                                                                            gi,
                                                                            si,
                                                                            field,
                                                                            bindIndex,
                                                                        )}
                                                                    placeholder={isRecording(
                                                                        gi,
                                                                        si,
                                                                        field,
                                                                        bindIndex,
                                                                    )
                                                                        ? "Press shortcut..."
                                                                        : bind.length ===
                                                                            0
                                                                          ? "Type key and Enter"
                                                                          : ""}
                                                                />
                                                            </div>
                                                            <button
                                                                class="bind-record"
                                                                class:active={isRecording(
                                                                    gi,
                                                                    si,
                                                                    field,
                                                                    bindIndex,
                                                                )}
                                                                type="button"
                                                                tabindex="-1"
                                                                title={isRecording(
                                                                    gi,
                                                                    si,
                                                                    field,
                                                                    bindIndex,
                                                                )
                                                                    ? "Stop recording"
                                                                    : "Record shortcut"}
                                                                onmousedown={(
                                                                    e,
                                                                ) =>
                                                                    e.preventDefault()}
                                                                onclick={() =>
                                                                    toggleRecording(
                                                                        gi,
                                                                        si,
                                                                        field,
                                                                        bindIndex,
                                                                    )}
                                                            >
                                                                <span
                                                                    class="record-dot"
                                                                ></span>
                                                            </button>
                                                            {#if shortcut[field].length > 1 || bind.length > 0}
                                                                <button
                                                                    class="bind-remove"
                                                                    type="button"
                                                                    tabindex="-1"
                                                                    onclick={() =>
                                                                        removeBind(
                                                                            gi,
                                                                            si,
                                                                            field,
                                                                            bindIndex,
                                                                        )}
                                                                    title="Remove this bind"
                                                                >
                                                                    ×
                                                                </button>
                                                            {/if}
                                                        </div>
                                                    {/each}
                                                    <button
                                                        class="bind-add"
                                                        type="button"
                                                        tabindex="-1"
                                                        onclick={() =>
                                                            addBind(
                                                                gi,
                                                                si,
                                                                field,
                                                            )}
                                                    >
                                                        + Add Bind
                                                    </button>
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
        max-height: 90vh;
        max-width: 88vw;
        overflow-y: auto;
        box-shadow: 0 24px 48px rgba(0, 0, 0, 0.6);
    }

    .panel-topbar {
        position: sticky;
        top: 0;
        z-index: 2;
        display: flex;
        align-items: center;
        justify-content: space-between;
        gap: 16px;
        padding: 14px 20px;
        background: rgba(22, 22, 22, 0.96);
        border-bottom: 1px solid #222222;
        backdrop-filter: blur(10px);
    }

    .panel-topbar-actions {
        display: flex;
        align-items: center;
        gap: 8px;
        flex-shrink: 0;
    }

    .close-btn {
        display: flex;
        align-items: center;
        justify-content: center;
        width: 28px;
        height: 28px;
        background: transparent;
        border: 1px solid #2a2a2a;
        border-radius: 6px;
        color: #525252;
        cursor: pointer;
        transition:
            background 0.1s,
            color 0.1s,
            border-color 0.1s;
    }

    .close-btn:hover {
        background: #1c1c1c;
        color: #a1a1a1;
        border-color: #3a3a3a;
    }

    .panel-title {
        font-size: 15px;
        font-weight: 600;
        color: #ffffff;
        margin: 0;
    }

    .panel-content {
        padding: 20px;
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

    .field--full {
        margin-bottom: 16px;
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

    .tag-editor {
        display: flex;
        flex-wrap: wrap;
        align-items: center;
        gap: 6px;
        min-height: 38px;
        padding: 6px 8px;
        border-radius: 6px;
        border: 1px solid #2a2a2a;
        background: #1c1c1c;
        cursor: text;
        transition: border-color 0.1s;
    }

    .tag-editor:focus-within {
        border-color: #3a88ed;
    }

    .tag-chip {
        display: inline-flex;
        align-items: center;
        gap: 5px;
        padding: 3px 8px;
        border-radius: 999px;
        border: 1px solid #2e3340;
        background: #171b24;
        color: #b8c2d9;
        font-size: 12px;
        font-weight: 500;
        line-height: 1;
        white-space: nowrap;
    }

    .tag-chip-remove {
        display: flex;
        align-items: center;
        justify-content: center;
        width: 14px;
        height: 14px;
        padding: 0;
        border: none;
        background: none;
        color: #6b7280;
        cursor: pointer;
        font-size: 14px;
        line-height: 1;
    }

    .tag-chip-remove:hover {
        color: #f87171;
    }

    .tag-input {
        flex: 1;
        min-width: 140px;
        padding: 0;
        border: none;
        background: none;
        color: #d4d4d4;
        font-size: 13px;
        font-family: inherit;
        outline: none;
    }

    .tag-input::placeholder {
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
        display: flex;
        align-items: center;
        gap: 6px;
        font-size: 13px;
        font-weight: 600;
        color: #555555;
        padding: 8px 12px 2px;
    }

    .info-hint {
        position: relative;
        display: inline-flex;
        align-items: center;
        color: #3f3f3f;
        cursor: help;
    }

    .info-hint:hover {
        color: #666666;
    }

    .info-tooltip {
        display: none;
        position: absolute;
        top: calc(100% + 6px);
        left: -4px;
        z-index: 10;
        width: 220px;
        padding: 8px 10px;
        border-radius: 6px;
        background: #1a1a1a;
        border: 1px solid #2a2a2a;
        box-shadow: 0 8px 24px rgba(0, 0, 0, 0.4);
        color: #a1a1a1;
        font-size: 11px;
        font-weight: 400;
        line-height: 1.5;
        white-space: normal;
    }

    .info-hint:hover .info-tooltip {
        display: block;
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

    .key-binds {
        display: flex;
        flex-direction: column;
        gap: 4px;
    }

    .bind-row {
        display: flex;
        align-items: center;
        gap: 4px;
    }

    .bind-or-label {
        font-size: 10px;
        font-weight: 600;
        color: #525252;
        padding-left: 2px;
    }

    .bind-remove {
        display: flex;
        align-items: center;
        justify-content: center;
        width: 18px;
        height: 18px;
        flex-shrink: 0;
        background: none;
        border: none;
        color: #525252;
        cursor: pointer;
        font-size: 15px;
        line-height: 1;
        padding: 0;
        border-radius: 3px;
        transition: color 0.1s;
    }

    .bind-remove:hover {
        color: #ff5a5a;
    }

    .bind-record {
        display: flex;
        align-items: center;
        justify-content: center;
        width: 18px;
        height: 18px;
        flex-shrink: 0;
        background: none;
        border: 1px solid #2a2a2a;
        border-radius: 50%;
        cursor: pointer;
        padding: 0;
        transition:
            border-color 0.15s,
            background 0.15s;
    }

    .record-dot {
        display: block;
        width: 6px;
        height: 6px;
        border-radius: 50%;
        background: #525252;
        transition: background 0.15s;
    }

    .bind-record:hover {
        border-color: #b33b3b;
    }

    .bind-record:hover .record-dot {
        background: #d94444;
    }

    .bind-record.active {
        border-color: #d94444;
        background: #2a1515;
    }

    .bind-record.active .record-dot {
        background: #ef4444;
        animation: pulse-record 1.2s ease-in-out infinite;
    }

    @keyframes pulse-record {
        0%,
        100% {
            opacity: 1;
        }
        50% {
            opacity: 0.4;
        }
    }

    .bind-add {
        align-self: flex-start;
        background: none;
        border: 1px dashed #2a2a2a;
        border-radius: 5px;
        color: #525252;
        font-size: 11px;
        font-weight: 500;
        cursor: pointer;
        padding: 2px 8px;
        transition:
            color 0.1s,
            border-color 0.1s;
    }

    .bind-add:hover {
        color: #a1a1a1;
        border-color: #3a3a3a;
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

    .chips-wrap.recording {
        border-color: #b33b3b;
        background: #251717;
    }

    .chips-wrap.recording:focus-within {
        border-color: #d94444;
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

    .chips-input[readonly] {
        cursor: default;
        caret-color: transparent;
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

    @media (max-width: 720px) {
        .panel {
            width: min(88vw, 920px);
        }

        .panel-topbar {
            flex-wrap: wrap;
            align-items: flex-start;
        }

        .panel-topbar-actions {
            width: 100%;
            justify-content: flex-end;
        }

        .field-row,
        .shortcut-row-keys {
            flex-direction: column;
        }

        .field--icon {
            flex-basis: auto;
        }
    }
</style>
