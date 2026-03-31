<script lang="ts">
  let { keys, highlight = false }: { keys: string[]; highlight?: boolean } = $props()
</script>

<div class="flex items-center gap-1">
  {#each keys as key, i}
    <kbd
      class="key-badge"
      class:key-badge--highlight={highlight}
    >
      {key}
    </kbd>
    {#if i < keys.length - 1}
      <span class="text-[11px] text-[#525252]">+</span>
    {/if}
  {/each}
</div>

<style>
  .key-badge {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    min-width: 24px;
    height: 22px;
    padding: 0 6px;
    font-family: ui-monospace, SFMono-Regular, 'SF Mono', Menlo, monospace;
    font-size: 11px;
    line-height: 1;
    color: #a1a1a1;
    background: #1c1c1c;
    border: 1px solid #3f3f3f;
    border-bottom-width: 2px;
    border-bottom-color: #525252;
    border-radius: 5px;
    white-space: nowrap;
  }

  .key-badge--highlight {
    color: #c4b5fd;
    background: #1e1533;
    border-color: #5b21b6;
    border-bottom-color: #7c3aed;
  }
</style>
