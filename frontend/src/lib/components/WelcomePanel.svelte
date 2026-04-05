<script lang="ts">
    import { X, ArrowRight } from "lucide-svelte";

    let {
        onContinue,
        onSkip,
    }: {
        onContinue: () => void;
        onSkip: () => void;
    } = $props();

    function handleBackdropClick(e: MouseEvent) {
        if (e.target === e.currentTarget) onSkip();
    }

    function handleKeydown(e: KeyboardEvent) {
        if (e.key === "Escape") onSkip();
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
        <button class="close-btn" onclick={onSkip} aria-label="Close">
            <X size={14} />
        </button>

        <div class="identity">
            <img src="/appicon.png" alt="Button logo" class="app-logo" />
            <div class="app-name">Welcome to Button</div>
            <p class="app-tagline">
                A cross-platform keyboard shortcut<br />reference for your
                desktop apps.
            </p>
        </div>

        <div class="divider"></div>

        <div class="description">
            <p>
                Button lets you keep track of keyboard shortcuts for all your
                apps in one place. You can create your own app configs or import
                pre-made ones to get started quickly.
            </p>
            <p>
                Would you like to browse the app registry and import some
                shortcuts?
            </p>
        </div>

        <div class="actions">
            <button class="btn-primary" onclick={onContinue}>
                Browse Apps
                <ArrowRight size={14} />
            </button>
            <button class="btn-skip" onclick={onSkip}> Skip for now </button>
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
        z-index: 100;
    }

    .panel {
        position: relative;
        background: #161616;
        border: 1px solid #2a2a2a;
        border-radius: 12px;
        padding: 24px 24px 20px;
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

    .identity {
        display: flex;
        flex-direction: column;
        align-items: center;
        gap: 6px;
        padding-bottom: 16px;
    }

    .app-logo {
        width: 68px;
        height: 68px;
        margin-bottom: 4px;
    }

    .app-name {
        font-size: 18px;
        font-weight: 700;
        color: #ffffff;
        letter-spacing: 0.02em;
    }

    .app-tagline {
        font-size: 13px;
        color: #555555;
        text-align: center;
        margin: 0;
        line-height: 1.3;
        font-weight: 400;
    }

    .divider {
        height: 1px;
        background: #1c1c1c;
        margin-bottom: 16px;
    }

    .description {
        margin-bottom: 20px;
    }

    .description p {
        font-size: 13px;
        color: #a1a1a1;
        line-height: 1.5;
        margin: 0 0 10px;
    }

    .description p:last-child {
        margin-bottom: 0;
    }

    .actions {
        display: flex;
        flex-direction: column;
        align-items: center;
        gap: 10px;
    }

    .btn-primary {
        display: flex;
        align-items: center;
        justify-content: center;
        gap: 6px;
        width: 100%;
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

    .btn-primary:hover {
        background: #264b77;
        color: #bfdbfe;
    }

    .btn-skip {
        background: none;
        border: none;
        color: #525252;
        font-size: 12px;
        cursor: pointer;
        padding: 4px 8px;
        transition: color 0.15s;
    }

    .btn-skip:hover {
        color: #a1a1a1;
    }
</style>
