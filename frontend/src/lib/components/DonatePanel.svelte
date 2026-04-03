<script lang="ts">
    import { X, Heart } from "lucide-svelte";

    let {
        onClose,
    }: {
        onClose: () => void;
    } = $props();

    function handleBackdropClick(e: MouseEvent) {
        if (e.target === e.currentTarget) onClose();
    }

    function handleKeydown(e: KeyboardEvent) {
        if (e.key === "Escape") onClose();
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
        <button class="close-btn" onclick={onClose} aria-label="Close">
            <X size={14} />
        </button>

        <div class="duo">
            <img src="/appicon.png" alt="Button logo" class="duo-logo" />
            <Heart size={42} class="duo-heart" fill="currentColor" />
        </div>

        <h2 class="title">Support Button</h2>

        <p class="message">
            Thank you for considering a donation,<br />it means a lot!<br />
            <br />
            Donations aren't set up yet,<br />but stay tuned.
        </p>
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
        padding: 32px 28px 28px;
        width: 300px;
        box-shadow: 0 24px 48px rgba(0, 0, 0, 0.6);
        display: flex;
        flex-direction: column;
        align-items: center;
        gap: 10px;
        text-align: center;
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

    .duo {
        display: flex;
        align-items: center;
        justify-content: center;
        gap: 10px;
        margin-bottom: 4px;
    }

    .duo-logo {
        width: 52px;
        height: 52px;
        transform: rotate(-12deg);
    }

    :global(.duo-heart) {
        color: #b83919;
        transform: rotate(12deg);
        flex-shrink: 0;
    }

    .title {
        font-size: 17px;
        font-weight: 700;
        color: #ffffff;
        margin: 0;
    }

    .message {
        font-size: 13px;
        color: #6b6b6b;
        line-height: 1.6;
        margin: 0;
    }
</style>
