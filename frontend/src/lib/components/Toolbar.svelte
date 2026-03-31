<script lang="ts">
  let {
    searchQuery = $bindable(''),
    currentOS,
    onToggleOS,
  }: {
    searchQuery: string
    currentOS: 'linux' | 'darwin'
    onToggleOS: () => void
  } = $props()

  let searchInput: HTMLInputElement | undefined = $state()

  $effect(() => {
    // Autofocus the search bar on mount
    searchInput?.focus()
  })
</script>

<div class="toolbar">
  <!-- Search bar -->
  <div class="search-wrapper">
    <svg class="search-icon" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
      <circle cx="11" cy="11" r="8"/>
      <path d="m21 21-4.35-4.35"/>
    </svg>
    <input
      bind:this={searchInput}
      bind:value={searchQuery}
      type="text"
      placeholder="Search apps and shortcuts..."
      class="search-input"
    />
    {#if searchQuery}
      <button
        class="search-clear"
        onclick={() => { searchQuery = ''; searchInput?.focus() }}
        aria-label="Clear search"
      >
        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <path d="M18 6 6 18"/><path d="m6 6 12 12"/>
        </svg>
      </button>
    {/if}
  </div>

  <!-- OS toggle -->
  <div class="os-toggle">
    <button
      class="os-toggle-btn"
      class:os-toggle-btn--active-linux={currentOS === 'linux'}
      onclick={() => { if (currentOS !== 'linux') onToggleOS() }}
    >
      Linux
    </button>
    <button
      class="os-toggle-btn"
      class:os-toggle-btn--active-macos={currentOS === 'darwin'}
      onclick={() => { if (currentOS !== 'darwin') onToggleOS() }}
    >
      macOS
    </button>
  </div>

  <!-- Help button -->
  <button class="icon-btn" aria-label="Help" title="Help">
    <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
      <circle cx="12" cy="12" r="10"/>
      <path d="M9.09 9a3 3 0 0 1 5.83 1c0 2-3 3-3 3"/>
      <path d="M12 17h.01"/>
    </svg>
  </button>

  <!-- Donate button -->
  <button class="icon-btn icon-btn--donate" aria-label="Donate" title="Support on Ko-fi">
    <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
      <path d="M19 14c1.49-1.46 3-3.21 3-5.5A5.5 5.5 0 0 0 16.5 3c-1.76 0-3 .5-4.5 2-1.5-1.5-2.74-2-4.5-2A5.5 5.5 0 0 0 2 8.5c0 2.3 1.5 4.05 3 5.5l7 7Z"/>
    </svg>
  </button>
</div>

<style>
  .toolbar {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 10px 12px;
    border-bottom: 1px solid #222222;
  }

  .search-wrapper {
    flex: 1;
    display: flex;
    align-items: center;
    gap: 8px;
    background: #1c1c1c;
    border: 1px solid #2a2a2a;
    border-radius: 8px;
    padding: 0 10px;
    height: 34px;
    transition: border-color 0.15s;
  }

  .search-wrapper:focus-within {
    border-color: #3a88ed;
  }

  .search-icon {
    color: #525252;
    flex-shrink: 0;
  }

  .search-input {
    flex: 1;
    background: transparent;
    border: none;
    outline: none;
    color: #e5e5e5;
    font-size: 13px;
    line-height: 1;
  }

  .search-input::placeholder {
    color: #525252;
  }

  .search-clear {
    display: flex;
    align-items: center;
    justify-content: center;
    color: #525252;
    background: none;
    border: none;
    cursor: pointer;
    padding: 2px;
    border-radius: 4px;
    transition: color 0.15s;
  }

  .search-clear:hover {
    color: #a1a1a1;
  }

  .os-toggle {
    display: flex;
    background: #1c1c1c;
    border: 1px solid #2a2a2a;
    border-radius: 8px;
    overflow: hidden;
    flex-shrink: 0;
  }

  .os-toggle-btn {
    padding: 6px 12px;
    font-size: 12px;
    font-weight: 500;
    color: #525252;
    background: transparent;
    border: none;
    cursor: pointer;
    transition: all 0.15s;
    line-height: 1;
  }

  .os-toggle-btn:hover {
    color: #a1a1a1;
  }

  .os-toggle-btn--active-linux {
    background: #1e3a5f;
    color: #60a5fa;
  }

  .os-toggle-btn--active-macos {
    background: #3d2e1a;
    color: #f59e0b;
  }

  .icon-btn {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 34px;
    height: 34px;
    background: #1c1c1c;
    border: 1px solid #2a2a2a;
    border-radius: 8px;
    color: #525252;
    cursor: pointer;
    flex-shrink: 0;
    transition: all 0.15s;
  }

  .icon-btn:hover {
    color: #a1a1a1;
    border-color: #3f3f3f;
  }

  .icon-btn--donate {
    color: #92400e;
  }

  .icon-btn--donate:hover {
    color: #f59e0b;
    border-color: #92400e;
  }
</style>
