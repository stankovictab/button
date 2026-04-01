<script lang="ts">
    import { X, ArrowUpDown, Search, Keyboard, FolderOpen, Monitor } from "lucide-svelte";
    import { BrowserOpenURL } from "../../../wailsjs/runtime/runtime";

    let { onClose }: { onClose: () => void } = $props();

    function handleBackdropClick(e: MouseEvent) {
        if (e.target === e.currentTarget) onClose();
    }

    function handleKeydown(e: KeyboardEvent) {
        if (e.key === "Escape") onClose();
    }
</script>

<svelte:window onkeydown={handleKeydown} />

<!-- svelte-ignore a11y_click_events_have_key_events a11y_no_static_element_interactions a11y_no_noninteractive_element_interactions -->
<div class="backdrop" role="dialog" aria-modal="true" tabindex="-1" onclick={handleBackdropClick}>
    <div class="panel">
        <button class="close-btn" onclick={onClose} aria-label="Close help">
            <X size={14} />
        </button>

        <!-- Identity -->
        <div class="identity">
            <img src="/appicon.png" alt="button logo" class="app-icon" />
            <div class="app-name">button</div>
            <div class="app-version">v0.0.0-alpha</div>
            <p class="app-tagline">A keyboard shortcut reference<br />for your desktop apps.</p>
        </div>

        <div class="divider"></div>

        <!-- Quick reference -->
        <div class="section-label">QUICK REFERENCE</div>
        <ul class="tips">
            <li class="tip">
                <span class="tip-icon"><Search size={13} /></span>
                <span class="tip-text">Type anywhere to search apps and shortcuts</span>
            </li>
            <li class="tip">
                <span class="tip-icon"><ArrowUpDown size={13} /></span>
                <span class="tip-text">Navigate apps with arrow keys</span>
            </li>
            <li class="tip">
                <span class="tip-icon"><Keyboard size={13} /></span>
                <span class="tip-text"><kbd>Esc</kbd> clears the search query</span>
            </li>
            <li class="tip">
                <span class="tip-icon"><Monitor size={13} /></span>
                <span class="tip-text">Toggle Linux / macOS to switch shortcut sets</span>
            </li>
            <li class="tip">
                <span class="tip-icon"><FolderOpen size={13} /></span>
                <span class="tip-text">
                    Add apps by dropping YAML files into<br />
                    <code>~/.config/button/apps/</code><br />
                    <!-- svelte-ignore a11y_invalid_attribute -->
                    <a
                        href="#"
                        class="tip-link"
                        onclick={(e) => { e.preventDefault(); BrowserOpenURL("https://github.com/stankovictab/button/blob/main/examples/template.yaml"); }}
                    >View example config →</a>
                </span>
            </li>
        </ul>
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
        padding: 24px;
        width: 300px;
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
        transition: background 0.1s, color 0.1s;
    }

    .close-btn:hover {
        background: #1c1c1c;
        color: #a1a1a1;
    }

    .identity {
        display: flex;
        flex-direction: column;
        align-items: center;
        gap: 6px;
        padding-bottom: 20px;
    }

    .app-icon {
        width: 64px;
        height: 64px;
        border-radius: 14px;
        margin-bottom: 4px;
    }

    .app-name {
        font-size: 18px;
        font-weight: 700;
        color: #ffffff;
        letter-spacing: 0.02em;
    }

    .app-version {
        font-size: 11px;
        color: #3f3f3f;
    }

    .app-tagline {
        font-size: 12px;
        color: #525252;
        text-align: center;
        margin: 0;
        line-height: 1.5;
    }

    .divider {
        height: 1px;
        background: #1c1c1c;
        margin-bottom: 16px;
    }

    .section-label {
        font-size: 10px;
        font-weight: 600;
        letter-spacing: 0.05em;
        color: #3f3f3f;
        margin-bottom: 10px;
    }

    .tips {
        list-style: none;
        margin: 0;
        padding: 0;
        display: flex;
        flex-direction: column;
        gap: 10px;
    }

    .tip {
        display: flex;
        align-items: flex-start;
        gap: 10px;
    }

    .tip-icon {
        display: flex;
        align-items: center;
        gap: 2px;
        color: #525252;
        flex-shrink: 0;
        width: 22px;
        margin-top: 1px;
    }

    .tip-text {
        font-size: 12px;
        color: #a1a1a1;
        line-height: 1.5;
    }

    .tip-text kbd {
        display: inline-flex;
        align-items: center;
        padding: 1px 5px;
        font-family: ui-monospace, SFMono-Regular, "SF Mono", Menlo, monospace;
        font-size: 10px;
        color: #a1a1a1;
        background: #1c1c1c;
        border: 1px solid #2a2a2a;
        border-radius: 3px;
    }

    .tip-link {
        color: #4597f5;
        text-decoration: none;
        font-size: 11px;
    }

    .tip-link:hover {
        text-decoration: underline;
    }

    .tip-text code {
        font-family: ui-monospace, SFMono-Regular, "SF Mono", Menlo, monospace;
        font-size: 11px;
        color: #525252;
    }
</style>
