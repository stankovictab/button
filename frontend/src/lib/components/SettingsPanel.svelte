<script lang="ts">
    import { onMount } from "svelte";
    import { X, Settings, Info } from "lucide-svelte";
    import {
        GetAutostartEnabled,
        SetAutostartEnabled,
    } from "../../../wailsjs/go/main/App.js";
    import type { NotificationType } from "../../types";

    let {
        onClose,
        onNotify,
    }: {
        onClose: () => void;
        onNotify: (type: NotificationType, message: string) => void;
    } = $props();

    let autostartEnabled: boolean = $state(false);
    let loading: boolean = $state(true);
    let saving: boolean = $state(false);

    function handleBackdropClick(e: MouseEvent) {
        if (e.target === e.currentTarget) onClose();
    }

    function handleKeydown(e: KeyboardEvent) {
        if (e.key === "Escape") onClose();
    }

    function loadAutostart() {
        loading = true;
        GetAutostartEnabled()
            .then((enabled: boolean) => {
                autostartEnabled = enabled;
            })
            .catch((err: any) => {
                onNotify("error", String(err));
            })
            .finally(() => {
                loading = false;
            });
    }

    function toggleAutostart() {
        const next = !autostartEnabled;
        saving = true;
        SetAutostartEnabled(next)
            .then(() => {
                autostartEnabled = next;
            })
            .catch((err: any) => {
                onNotify("error", String(err));
                loadAutostart();
            })
            .finally(() => {
                saving = false;
            });
    }

    onMount(loadAutostart);
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

        <div class="heading">
            <Settings size={17} />
            <h2 class="title">Settings</h2>
        </div>

        <label class="setting-row" class:setting-row--disabled={loading || saving}>
            <input
                type="checkbox"
                checked={autostartEnabled}
                disabled={loading || saving}
                onchange={toggleAutostart}
            />
            <span class="setting-copy">
                <span class="setting-label">
                    Run Button on login
                    <button
                        type="button"
                        class="info-hint"
                        aria-label="Autostart details"
                        onclick={(e) => e.stopPropagation()}
                    >
                        <Info size={12} />
                        <span class="info-tooltip"
                            >Button writes a desktop entry to <code
                                >~/.config/autostart/button.desktop</code
                            > using the stable Button binary path. Turning this
                            off removes that Button autostart file.</span
                        >
                    </button>
                </span>
                <span class="setting-note">
                    Creates a user autostart entry for this Button binary.
                </span>
            </span>
        </label>
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
        padding: 22px;
        width: 360px;
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

    .heading {
        display: flex;
        align-items: center;
        gap: 8px;
        color: #d8d8d8;
        margin-bottom: 18px;
    }

    .title {
        font-size: 16px;
        font-weight: 700;
        color: #ffffff;
        margin: 0;
    }

    .setting-row {
        display: flex;
        align-items: flex-start;
        gap: 10px;
        padding: 12px;
        border: 1px solid #2a2a2a;
        border-radius: 8px;
        background: #1c1c1c;
        cursor: pointer;
    }

    .setting-row--disabled {
        cursor: default;
        opacity: 0.7;
    }

    input {
        margin-top: 2px;
        accent-color: #3a88ed;
    }

    .setting-copy {
        display: flex;
        flex-direction: column;
        gap: 4px;
    }

    .setting-label {
        display: inline-flex;
        align-items: center;
        gap: 4px;
        color: #e5e5e5;
        font-size: 13px;
        font-weight: 600;
    }

    .setting-note {
        color: #777777;
        font-size: 12px;
        line-height: 1.4;
    }

    .info-hint {
        position: relative;
        display: inline-flex;
        align-items: center;
        padding: 0;
        border: none;
        background: transparent;
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
        width: 245px;
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

    .info-tooltip code {
        color: #d4d4d4;
        font-family:
            "SF Mono", Monaco, "Cascadia Code", "Roboto Mono", Consolas,
            "Courier New", monospace;
        font-size: 10px;
    }

    .info-hint:hover .info-tooltip {
        display: block;
    }
</style>
