<script lang="ts">
    import { X } from "lucide-svelte";
    import type { Notification } from "../../types";

    let {
        notifications,
        onDismiss,
    }: {
        notifications: Notification[];
        onDismiss: (id: number) => void;
    } = $props();
</script>

{#if notifications.length > 0}
    <div class="notification-stack">
        {#each notifications as notif (notif.id)}
            <div class="notification notification--{notif.type}">
                <span class="notification-message">{@html notif.message}</span>
                <button
                    class="notification-dismiss"
                    onclick={() => onDismiss(notif.id)}
                    aria-label="Dismiss"
                >
                    <X size={15} />
                </button>
            </div>
        {/each}
    </div>
{/if}

<style>
    .notification-stack {
        display: flex;
        flex-direction: column;
    }

    .notification {
        display: flex;
        align-items: center;
        justify-content: space-between;
        gap: 8px;
        padding: 6px 12px;
        font-size: 12px;
        border-bottom: 2px solid transparent;
    }

    .notification--error {
        background: #450a0a;
        border-bottom-color: #7f1d1d;
        color: #fca5a5;
    }

    .notification--warning {
        background: #451a03;
        border-bottom-color: #78350f;
        color: #fde68a;
    }

    .notification--info {
        background: #0c1a3d;
        border-bottom-color: #1e3a5f;
        color: #93c5fd;
    }

    .notification-message {
        flex: 1;
        line-height: 1.5;
    }

    .notification-message :global(code) {
        font-family: "JetBrains Mono", ui-monospace, SFMono-Regular, "SF Mono",
            Menlo, monospace;
        font-size: 11px;
        background: rgba(255, 255, 255, 0.08);
        padding: 1px 4px;
        border-radius: 3px;
    }

    .notification-dismiss {
        display: flex;
        align-items: center;
        justify-content: center;
        width: 18px;
        height: 18px;
        border: none;
        border-radius: 4px;
        background: transparent;
        color: inherit;
        opacity: 0.5;
        cursor: pointer;
        flex-shrink: 0;
        transition:
            opacity 0.1s,
            background 0.1s;
    }

    .notification-dismiss:hover {
        opacity: 1;
        background: rgba(255, 255, 255, 0.1);
    }
</style>
