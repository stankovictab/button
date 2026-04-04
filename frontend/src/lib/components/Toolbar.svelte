<script lang="ts">
    import { Search, X, CircleHelp, Heart } from "lucide-svelte";
    import { siGithub } from "simple-icons";
    import { BrowserOpenURL } from "../../../wailsjs/runtime/runtime";

    let {
        searchQuery = $bindable(""),
        showHelp = $bindable(false),
        showDonate = $bindable(false),
        currentOS,
        matchingDescs = new Set<string>(),
        onToggleOS,
        onSearchInput,
    }: {
        searchQuery: string;
        showHelp: boolean;
        showDonate: boolean;
        currentOS: "linux" | "darwin";
        matchingDescs: Set<string>;
        onToggleOS: () => void;
        onSearchInput?: (element: HTMLInputElement | undefined) => void;
    } = $props();

    let searchInput: HTMLInputElement | undefined = $state();

    $effect(() => {
        // Autofocus the search bar on mount
        searchInput?.focus();
    });

    $effect(() => {
        onSearchInput?.(searchInput);
    });
</script>

<div class="toolbar">
    <!-- Search bar -->
    <div class="search-wrapper">
        <span class="search-icon">
            <Search size={16} />
        </span>
        <input
            bind:this={searchInput}
            bind:value={searchQuery}
            type="text"
            placeholder="Search apps and shortcuts..."
            class="search-input"
        />
        {#if searchQuery && matchingDescs.size > 0}
            <div class="search-match-badge">
                {matchingDescs.size}
                {matchingDescs.size === 1 ? "match" : "matches"}
            </div>
        {/if}
        {#if searchQuery}
            <button
                class="search-clear"
                onclick={() => {
                    searchQuery = "";
                    searchInput?.focus();
                }}
                aria-label="Clear search"
                title="Clear search"
            >
                <X size={14} />
            </button>
        {/if}
    </div>

    <!-- OS toggle -->
    <div class="os-toggle">
        <button
            class="os-toggle-btn"
            class:os-toggle-btn--active-linux={currentOS === "linux"}
            onclick={() => {
                if (currentOS !== "linux") onToggleOS();
            }}
            title="Switch to Linux shortcuts"
        >
            Linux
        </button>
        <button
            class="os-toggle-btn"
            class:os-toggle-btn--active-macos={currentOS === "darwin"}
            onclick={() => {
                if (currentOS !== "darwin") onToggleOS();
            }}
            title="Switch to macOS shortcuts"
        >
            macOS
        </button>
    </div>

    <!-- Help button -->
    <button
        class="icon-btn"
        class:icon-btn--active={showHelp}
        aria-label="Help"
        title="Help"
        onclick={() => (showHelp = !showHelp)}
    >
        <CircleHelp size={16} />
    </button>

    <!-- GitHub button -->
    <button
        class="icon-btn icon-btn--github"
        onclick={() => BrowserOpenURL("https://github.com/stankovictab/button")}
        aria-label="GitHub repository"
        title="GitHub repository"
    >
        {@html siGithub.svg}
    </button>

    <!-- Donate button -->
    <button
        class="icon-btn icon-btn--donate"
        class:icon-btn--active={showDonate}
        onclick={() => (showDonate = !showDonate)}
        aria-label="Donate"
        title="Support Button"
    >
        <Heart size={16} />
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
        align-self: stretch;
        background: transparent;
        border: none;
        outline: none;
        color: #e5e5e5;
        font-size: 13px;
    }

    .search-input::placeholder {
        color: #525252;
    }

    .search-match-badge {
        font-size: 11px;
        color: #93c5fd;
        background: #172554;
        padding: 3px 8px;
        border-radius: 10px;
        white-space: nowrap;
        flex-shrink: 0;
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
        /* height: 34px; */ /* Use this to fix height. */
        background: #1c1c1c;
        border: 1px solid #2a2a2a;
        border-radius: 6px;
        overflow: hidden;
        flex-shrink: 0;
    }

    .os-toggle-btn {
        padding: 6px 12px;
        font-size: 13px;
        font-weight: 600;
        color: #525252;
        background: transparent;
        border: none;
        cursor: pointer;
        transition: all 0.15s;
        line-height: 1;
    }

    .os-toggle-btn:hover {
        background: #252525;
        color: #a1a1a1;
    }

    .os-toggle-btn--active-linux {
        background: #172c47;
        color: #59a4ff;
    }

    .os-toggle-btn--active-linux:hover {
        background: #1e3a5f;
        color: #9bcaff;
    }

    .os-toggle-btn--active-macos {
        background: #352716;
        color: #ff8928;
    }

    .os-toggle-btn--active-macos:hover {
        background: #3d2e1a;
        color: #f3a620;
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
        color: #d8d8d8;
        border-color: #3f3f3f;
    }

    .icon-btn--active {
        color: #d8d8d8;
        border-color: #3f3f3f;
    }

    .icon-btn--donate {
        color: #92400e;
    }

    .icon-btn--donate:hover {
        color: #ffa200;
        border-color: #92400e;
    }

    .icon-btn--github :global(svg) {
        width: 15px;
        height: 15px;
        fill: currentColor;
    }
</style>
