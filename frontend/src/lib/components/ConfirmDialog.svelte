<script lang="ts">
    let {
        title,
        message,
        confirmLabel = "Delete",
        cancelLabel = "Cancel",
        danger = false,
        onConfirm,
        onCancel,
    }: {
        title: string;
        message: string;
        confirmLabel?: string;
        cancelLabel?: string;
        danger?: boolean;
        onConfirm: () => void;
        onCancel: () => void;
    } = $props();

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
        <h3 class="title">{title}</h3>
        <p class="message">{@html message}</p>
        <div class="actions">
            <button class="btn btn--cancel" onclick={onCancel}
                >{cancelLabel}</button
            >
            <button
                class="btn"
                class:btn--danger={danger}
                class:btn--primary={!danger}
                onclick={onConfirm}>{confirmLabel}</button
            >
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
        background: #161616;
        border: 1px solid #2a2a2a;
        border-radius: 12px;
        padding: 20px 24px;
        width: 450px;
        box-shadow: 0 24px 48px rgba(0, 0, 0, 0.6);
    }

    .title {
        font-size: 15px;
        font-weight: 600;
        color: #ffffff;
        margin: 0 0 8px;
    }

    .message {
        font-size: 13px;
        color: #a1a1a1;
        margin: 0 0 20px;
        line-height: 1.5;
    }

    :global(.message code) {
        font-family: ui-monospace, SFMono-Regular, "SF Mono", Menlo, monospace;
        font-size: 12px;
        color: #d4d4d4;
        background: #121212;
        padding: 2px 6px;
        border-radius: 4px;
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

    .btn--cancel {
        background: #1c1c1c;
        color: #a1a1a1;
        border: 1px solid #2a2a2a;
    }

    .btn--cancel:hover {
        background: #252525;
        color: #d4d4d4;
    }

    .btn--danger {
        background: #7f1d1d;
        color: #fca5a5;
    }

    .btn--danger:hover {
        background: #991b1b;
        color: #fecaca;
    }

    .btn--primary {
        background: #1e3a5f;
        color: #93c5fd;
    }

    .btn--primary:hover {
        background: #1e4a7f;
        color: #bfdbfe;
    }
</style>
