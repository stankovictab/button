<script lang="ts">
    import { onMount, tick } from "svelte";
    import { ChevronDown, X } from "lucide-svelte";
    import type { Shortcut } from "../../types";

    let {
        mode,
        initial = null,
        existingCategories = [],
        onSave,
        onCancel,
    }: {
        mode: "create" | "edit";
        initial?: { category: string; shortcut: Shortcut } | null;
        existingCategories?: string[];
        onSave: (payload: { category: string; shortcut: Shortcut }) => void;
        onCancel: () => void;
    } = $props();

    type KeyField = "keys" | "linux" | "macos";
    type DraftField = "keysDraft" | "linuxDraft" | "macosDraft";

    // --- Form state (captured once at mount — initial is intentionally read only at creation) ---
    // svelte-ignore state_referenced_locally
    const initCategory = initial?.category ?? "";
    // svelte-ignore state_referenced_locally
    const initDesc = initial?.shortcut.desc ?? "";
    // svelte-ignore state_referenced_locally
    const initKeys = [...(initial?.shortcut.keys ?? [])];
    // svelte-ignore state_referenced_locally
    const initLinux = [...(initial?.shortcut.linux ?? [])];
    // svelte-ignore state_referenced_locally
    const initMacos = [...(initial?.shortcut.macos ?? [])];

    let category: string = $state(initCategory);
    let desc: string = $state(initDesc);
    let keys: string[] = $state(initKeys);
    let linux: string[] = $state(initLinux);
    let macos: string[] = $state(initMacos);
    let keysDraft: string = $state("");
    let linuxDraft: string = $state("");
    let macosDraft: string = $state("");
    let panelEl: HTMLDivElement | undefined = $state();
    let descInput: HTMLInputElement | undefined = $state();
    let categoryInput: HTMLInputElement | undefined = $state();
    let categoryMenuOpen: boolean = $state(false);

    function draftFieldFor(field: KeyField): DraftField {
        return (field + "Draft") as DraftField;
    }

    function valuesFor(field: KeyField): string[] {
        if (field === "keys") return keys;
        if (field === "linux") return linux;
        return macos;
    }

    function setValues(field: KeyField, value: string[]) {
        if (field === "keys") {
            keys = value;
            return;
        }
        if (field === "linux") {
            linux = value;
            return;
        }
        macos = value;
    }

    function draftValueFor(field: KeyField): string {
        if (field === "keys") return keysDraft;
        if (field === "linux") return linuxDraft;
        return macosDraft;
    }

    function setDraftValue(field: KeyField, value: string) {
        if (field === "keys") {
            keysDraft = value;
            return;
        }
        if (field === "linux") {
            linuxDraft = value;
            return;
        }
        macosDraft = value;
    }

    function commitDraft(field: KeyField) {
        const draft = draftValueFor(field).trim();
        if (!draft) return;

        setValues(field, [...valuesFor(field), draft]);
        setDraftValue(field, "");
    }

    function handleChipKeydown(e: KeyboardEvent, field: KeyField) {
        const draft = draftValueFor(field).trim();

        if (e.key === "Enter" && draft) {
            e.preventDefault();
            commitDraft(field);
            return;
        }

        if (e.key === "Backspace" && !draftValueFor(field)) {
            setValues(field, valuesFor(field).slice(0, -1));
        }
    }

    function removeChip(field: KeyField, chipIndex: number) {
        setValues(
            field,
            valuesFor(field).filter((_, i) => i !== chipIndex),
        );
    }

    function handleSave() {
        commitDraft("keys");
        commitDraft("linux");
        commitDraft("macos");

        const shortcut: Shortcut = {
            desc: desc.trim(),
        };
        if (keys.length) shortcut.keys = keys;
        if (linux.length) shortcut.linux = linux;
        if (macos.length) shortcut.macos = macos;

        onSave({
            category: category.trim(),
            shortcut,
        });
    }

    let canSave = $derived(category.trim().length > 0 && desc.trim().length > 0);
    let filteredCategories = $derived.by(() => {
        const query = category.trim().toLowerCase();
        if (!query) return existingCategories;
        return existingCategories.filter((existingCategory) =>
            existingCategory.toLowerCase().includes(query),
        );
    });

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
        descInput?.focus();
        descInput?.select();
    }

    onMount(() => {
        void focusFirstField();
    });

    function openCategoryMenu() {
        categoryMenuOpen = true;
    }

    function closeCategoryMenu() {
        categoryMenuOpen = false;
    }

    function toggleCategoryMenu() {
        categoryMenuOpen = !categoryMenuOpen;
        if (categoryMenuOpen) {
            categoryInput?.focus();
        }
    }

    function selectCategory(nextCategory: string) {
        category = nextCategory;
        categoryMenuOpen = false;
        categoryInput?.focus();
    }

    function handleCategoryBlur() {
        requestAnimationFrame(() => {
            const active = document.activeElement;
            if (
                active instanceof HTMLElement &&
                panelEl?.contains(active) &&
                active.dataset.categoryOption === "true"
            ) {
                return;
            }
            closeCategoryMenu();
        });
    }

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

    function handleCategoryKeydown(e: KeyboardEvent) {
        if (e.key === "ArrowDown") {
            e.preventDefault();
            openCategoryMenu();
            return;
        }

        if (e.key === "Escape") {
            closeCategoryMenu();
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
                {mode === "create" ? "New Shortcut" : "Edit Shortcut"}
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
            <div class="field-row">
                <label class="field">
                    <span class="field-label">Description</span>
                    <input
                        bind:this={descInput}
                        type="text"
                        class="field-input"
                        bind:value={desc}
                        placeholder="Shortcut description"
                    />
                </label>
                <label class="field field--category">
                    <span class="field-label">Group</span>
                    <div class="combobox">
                        <input
                            bind:this={categoryInput}
                            type="text"
                            class="field-input combobox-input"
                            bind:value={category}
                            placeholder="e.g. Navigation"
                            onfocus={openCategoryMenu}
                            oninput={openCategoryMenu}
                            onblur={handleCategoryBlur}
                            onkeydown={handleCategoryKeydown}
                        />
                        <button
                            type="button"
                            class="combobox-toggle"
                            onclick={toggleCategoryMenu}
                            aria-label="Toggle group suggestions"
                        >
                            <ChevronDown size={15} />
                        </button>
                        {#if categoryMenuOpen && filteredCategories.length > 0}
                            <div class="combobox-menu">
                                {#each filteredCategories as existingCategory}
                                    <button
                                        type="button"
                                        class="combobox-option"
                                        data-category-option="true"
                                        onmousedown={(e) => e.preventDefault()}
                                        onclick={() =>
                                            selectCategory(existingCategory)}
                                    >
                                        {existingCategory}
                                    </button>
                                {/each}
                            </div>
                        {/if}
                    </div>
                </label>
            </div>

            <div class="shortcut-card">
                <div class="shortcut-card-label">Key Bindings</div>
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
                                {#each valuesFor(field) as key, chipIndex}
                                    <span class="chip">
                                        {key}
                                        <button
                                            class="chip-remove"
                                            type="button"
                                            tabindex="-1"
                                            onclick={() =>
                                                removeChip(field, chipIndex)}
                                        >
                                            ×
                                        </button>
                                    </span>
                                {/each}
                                <input
                                    class="chips-input"
                                    type="text"
                                    value={draftValueFor(field)}
                                    oninput={(e) =>
                                        setDraftValue(
                                            field,
                                            (e.currentTarget as HTMLInputElement)
                                                .value,
                                        )}
                                    onkeydown={(e) =>
                                        handleChipKeydown(e, field)}
                                    onblur={() => commitDraft(field)}
                                    placeholder={valuesFor(field).length === 0
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
        width: min(760px, 88vw);
        background: #161616;
        border: 1px solid #2a2a2a;
        border-radius: 12px;
        max-height: 90vh;
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

    .panel-title {
        font-size: 15px;
        font-weight: 600;
        color: #ffffff;
        margin: 0;
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

    .field--category {
        flex: 0 0 220px;
    }

    .combobox {
        position: relative;
    }

    .combobox-input {
        padding-right: 32px;
    }

    .combobox-toggle {
        position: absolute;
        top: 50%;
        right: 6px;
        display: flex;
        align-items: center;
        justify-content: center;
        width: 20px;
        height: 20px;
        padding: 0;
        transform: translateY(-50%);
        border: none;
        border-radius: 4px;
        background: transparent;
        color: #666666;
        cursor: pointer;
        transition:
            background 0.1s,
            color 0.1s;
    }

    .combobox-toggle:hover {
        background: #262626;
        color: #cfcfcf;
    }

    .combobox-menu {
        position: absolute;
        top: calc(100% + 6px);
        left: 0;
        right: 0;
        z-index: 5;
        display: flex;
        flex-direction: column;
        max-height: 180px;
        overflow-y: auto;
        padding: 6px;
        border: 1px solid #2a2a2a;
        border-radius: 8px;
        background: #151515;
        box-shadow: 0 16px 32px rgba(0, 0, 0, 0.45);
    }

    .combobox-option {
        display: flex;
        align-items: center;
        width: 100%;
        padding: 7px 8px;
        border: none;
        border-radius: 6px;
        background: transparent;
        color: #cfcfcf;
        font-size: 13px;
        text-align: left;
        cursor: pointer;
        transition:
            background 0.1s,
            color 0.1s;
    }

    .combobox-option:hover {
        background: #202020;
        color: #ffffff;
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

    .shortcut-card {
        border: 1px solid #222222;
        border-radius: 8px;
        background: #131313;
        overflow: hidden;
    }

    .shortcut-card-label {
        font-size: 13px;
        font-weight: 600;
        color: #555555;
        padding: 10px 12px;
        border-bottom: 1px solid #1c1c1c;
        background: #1a1a1a;
    }

    .shortcut-row-keys {
        display: flex;
        gap: 6px;
        padding: 12px;
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

        .field--category {
            flex-basis: auto;
        }
    }
</style>
