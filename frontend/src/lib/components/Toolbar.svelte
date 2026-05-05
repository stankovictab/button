<script lang="ts">
    import {
        Search,
        X,
        CircleHelp,
        Heart,
        ChevronDown,
        Download,
        Settings,
    } from "lucide-svelte";
    import { siGithub } from "simple-icons";
    import { BrowserOpenURL } from "../../../wailsjs/runtime/runtime";

    let {
        searchQuery = $bindable(""),
        showHelp = $bindable(false),
        showDonate = $bindable(false),
        showImport = $bindable(false),
        showSettings = $bindable(false),
        currentOS,
        matchingDescs = new Set<string>(),
        onSetOS,
        onSearchInput,
    }: {
        searchQuery: string;
        showHelp: boolean;
        showDonate: boolean;
        showImport: boolean;
        showSettings: boolean;
        currentOS: "linux" | "darwin" | "windows";
        matchingDescs: Set<string>;
        onSetOS: (os: "linux" | "darwin" | "windows") => void;
        onSearchInput?: (element: HTMLInputElement | undefined) => void;
    } = $props();

    let searchInput: HTMLInputElement | undefined = $state();
    let osDropdownOpen: boolean = $state(false);
    let dropdownRef: HTMLDivElement | undefined = $state();

    const osOptions: {
        value: "linux" | "darwin" | "windows";
        label: string;
    }[] = [
        { value: "linux", label: "Linux" },
        { value: "windows", label: "Windows" },
        { value: "darwin", label: "macOS" },
    ];

    let currentLabel = $derived(
        osOptions.find((o) => o.value === currentOS)?.label ?? "Linux",
    );

    function handleClickOutside(e: MouseEvent) {
        if (
            dropdownRef &&
            !dropdownRef.contains(e.target as Node) &&
            osDropdownOpen
        ) {
            osDropdownOpen = false;
        }
    }

    $effect(() => {
        // Autofocus the search bar on mount
        searchInput?.focus();
    });

    $effect(() => {
        onSearchInput?.(searchInput);
    });
</script>

<svelte:window onclick={handleClickOutside} />

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

    <!-- OS selector -->
    <div class="os-dropdown" bind:this={dropdownRef}>
        <button
            class="os-dropdown-btn"
            class:os-dropdown-btn--linux={currentOS === "linux"}
            class:os-dropdown-btn--windows={currentOS === "windows"}
            class:os-dropdown-btn--macos={currentOS === "darwin"}
            onclick={(e) => {
                e.stopPropagation();
                osDropdownOpen = !osDropdownOpen;
            }}
            title="Switch platform"
        >
            <span class="os-dropdown-label">{currentLabel}</span>
            <ChevronDown size={13} />
        </button>
        {#if osDropdownOpen}
            <div class="os-dropdown-menu">
                {#each osOptions as opt}
                    <button
                        class="os-dropdown-option"
                        class:os-dropdown-option--active={currentOS ===
                            opt.value}
                        onclick={(e) => {
                            e.stopPropagation();
                            onSetOS(opt.value);
                            osDropdownOpen = false;
                        }}
                    >
                        {opt.label}
                    </button>
                {/each}
            </div>
        {/if}
    </div>

    <!-- Import button -->
    <button
        class="icon-btn"
        class:icon-btn--active={showImport}
        aria-label="Import apps from registry"
        title="Import apps"
        onclick={() => (showImport = !showImport)}
    >
        <Download size={16} />
    </button>

    <!-- Settings button -->
    <button
        class="icon-btn"
        class:icon-btn--active={showSettings}
        aria-label="Settings"
        title="Settings"
        onclick={() => (showSettings = !showSettings)}
    >
        <Settings size={16} />
    </button>

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

    .os-dropdown {
        position: relative;
        flex-shrink: 0;
    }

    .os-dropdown-btn {
        display: flex;
        align-items: center;
        gap: 4px;
        padding: 6px 10px;
        font-size: 13px;
        font-weight: 500;
        color: #525252;
        background: #1c1c1c;
        border: 1px solid #2a2a2a;
        border-radius: 8px;
        cursor: pointer;
        transition: all 0.15s;
        line-height: 1;
        height: 34px;
        white-space: nowrap;
    }

    .os-dropdown-btn:hover {
        /* border-color: #3f3f3f; */
        color: #a1a1a1;
    }

    .os-dropdown-btn--linux {
        background: #172c47;
        color: #59a4ff;
        border-color: #1e3a5f;
    }

    .os-dropdown-btn--linux:hover {
        background: #1e3a5f;
        color: #9bcaff;
    }

    .os-dropdown-btn--windows {
        background: #172c47;
        color: #59a4ff;
        border-color: #1e3a5f;
    }

    .os-dropdown-btn--windows:hover {
        background: #1e3a5f;
        color: #9bcaff;
    }

    .os-dropdown-btn--macos {
        background: #352716;
        color: #ff8928;
        border-color: #3d2e1a;
    }

    .os-dropdown-btn--macos:hover {
        background: #3d2e1a;
        color: #f3a620;
    }

    .os-dropdown-label {
        display: inline-block;
        width: 60px;
        text-align: left;
    }

    .os-dropdown-menu {
        position: absolute;
        top: calc(100% + 6px);
        right: 0;
        z-index: 50;
        display: flex;
        flex-direction: column;
        padding: 4px;
        border: 1px solid #2a2a2a;
        border-radius: 8px;
        background: #161616;
        box-shadow: 0 12px 32px rgba(0, 0, 0, 0.5);
        min-width: 120px;
    }

    .os-dropdown-option {
        display: flex;
        align-items: center;
        width: 100%;
        padding: 7px 10px;
        border: none;
        border-radius: 6px;
        background: transparent;
        color: #a1a1a1;
        font-size: 13px;
        font-weight: 500;
        text-align: left;
        cursor: pointer;
        transition:
            background 0.1s,
            color 0.1s;
    }

    .os-dropdown-option:hover {
        background: #1c1c1c;
        color: #ffffff;
    }

    .os-dropdown-option--active {
        color: #ffffff;
        background: #1c1c1c;
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
        color: #777777;
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
