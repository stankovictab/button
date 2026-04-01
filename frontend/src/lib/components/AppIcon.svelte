<script lang="ts">
    import { iconMap } from '../icons/iconMap';

    let {
        icon,
        name,
        size = 18,
    }: {
        icon: string;
        name: string;
        size?: number;
    } = $props();

    const entry = $derived(iconMap[icon]);
</script>

{#if entry}
    <span
        class="app-icon-svg"
        style="width: {size}px; height: {size}px; color: #{entry.hex};"
    >
        {@html entry.svg}
    </span>
{:else}
    <span class="app-icon-letter" style="width: {size}px; height: {size}px; font-size: {size * 0.55}px;">
        {name.charAt(0).toUpperCase()}
    </span>
{/if}

<style>
    .app-icon-svg {
        display: flex;
        align-items: center;
        justify-content: center;
        flex-shrink: 0;
    }

    .app-icon-svg :global(svg) {
        width: 100%;
        height: 100%;
        fill: currentColor;
    }

    .app-icon-letter {
        display: flex;
        align-items: center;
        justify-content: center;
        font-weight: 600;
        color: #a1a1a1;
        flex-shrink: 0;
    }
</style>
