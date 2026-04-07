<script lang="ts">
    import {
        GetApps,
        GetAppInfo,
        GetCurrentOS,
        CreateApp,
        UpdateApp,
        DeleteApp,
        GetMigrationResult,
        GetUserConfig,
        SetHasSeenWelcome,
        SetListPreferences,
        ImportRegistryApps,
    } from "../wailsjs/go/main/App.js";
    import { EventsOn } from "../wailsjs/runtime/runtime.js";
    import { onMount } from "svelte";
    import type {
        AppConfig,
        AppsResponse,
        AppInfo,
        SortMode,
        Notification,
        NotificationType,
        Group,
        Shortcut,
    } from "./types";

    import Toolbar from "./lib/components/Toolbar.svelte";
    import AppList from "./lib/components/AppList.svelte";
    import ShortcutDetail from "./lib/components/ShortcutDetail.svelte";
    import ConfirmDialog from "./lib/components/ConfirmDialog.svelte";
    import AppFormModal from "./lib/components/AppFormModal.svelte";
    import ShortcutFormModal from "./lib/components/ShortcutFormModal.svelte";
    import NotificationBar from "./lib/components/NotificationBar.svelte";
    import HelpPanel from "./lib/components/HelpPanel.svelte";
    import DonatePanel from "./lib/components/DonatePanel.svelte";
    import WelcomePanel from "./lib/components/WelcomePanel.svelte";
    import ImportPanel from "./lib/components/ImportPanel.svelte";
    import FlamingoEasterEgg from "./lib/components/FlamingoEasterEgg.svelte";

    // --- State ---
    let apps: AppConfig[] = $state([]);
    let notifications: Notification[] = $state([]);
    let nextNotificationId = 0;
    let selectedIndex: number = $state(0);
    let searchQuery: string = $state("");
    let currentOS: "linux" | "darwin" | "windows" = $state("linux");
    let listWidth: number = $state(310);
    let isResizing: boolean = $state(false);
    let sortMode: SortMode = $state("alpha");
    let groupByTag: boolean = $state(false);
    let appInfo: AppInfo = $state({
        name: "Button",
        version: "",
    });

    // --- CRUD modal state ---
    let showDeleteConfirm: boolean = $state(false);
    let deleteTargetIndex: number = $state(-1);
    let showAppForm: boolean = $state(false);
    let appFormMode: "create" | "edit" = $state("create");
    let editTargetIndex: number = $state(-1);
    let showShortcutForm: boolean = $state(false);
    let shortcutFormMode: "create" | "edit" = $state("create");
    let showShortcutDeleteConfirm: boolean = $state(false);
    let shortcutTarget: {
        appIndex: number;
        groupIndex: number;
        shortcutIndex: number;
    } | null = $state(null);
    let showOverwriteConfirm: boolean = $state(false);
    let pendingOverwriteApp: AppConfig | null = $state(null);
    let pendingOverwriteOldName: string = $state("");
    let showHelp: boolean = $state(false);
    let showDonate: boolean = $state(false);
    let showWelcome: boolean = $state(false);
    let showImport: boolean = $state(false);
    let cameFromWelcome: boolean = $state(false);
    let searchInput: HTMLInputElement | undefined = $state();
    let detailBody: HTMLDivElement | undefined = $state();
    let flamingoTrigger = $state(0);

    // --- Derived: sorted base apps ---
    function primaryTag(app: AppConfig): string | null {
        const tag = app.tags?.find((value) => value.trim().length > 0)?.trim();
        return tag || null;
    }

    function compareApps(a: AppConfig, b: AppConfig): number {
        if (sortMode === "alpha") {
            return a.app.localeCompare(b.app);
        }

        return b.modTime - a.modTime || a.app.localeCompare(b.app);
    }

    function comparePrimaryTags(a: AppConfig, b: AppConfig): number {
        const aTag = primaryTag(a);
        const bTag = primaryTag(b);

        if (!aTag && bTag) return 1;
        if (aTag && !bTag) return -1;
        if (aTag && bTag) {
            const tagCompare = aTag.localeCompare(bTag);
            if (tagCompare !== 0) return tagCompare;
        }

        return 0;
    }

    function compareTagBuckets(a: AppConfig, b: AppConfig): number {
        const tagCompare = comparePrimaryTags(a, b);
        if (tagCompare !== 0) return tagCompare;

        return compareApps(a, b);
    }

    let sortedBaseApps = $derived.by(() => {
        const copy = apps.slice();
        copy.sort(groupByTag ? compareTagBuckets : compareApps);
        return copy;
    });

    // --- Derived: search matching ---
    // For each app (by original index), compute which shortcut descriptions match the search.
    let searchResults = $derived.by(() => {
        const q = searchQuery.toLowerCase().trim();
        if (!q) {
            return {
                sortedApps: sortedBaseApps,
                matchCounts: {} as Record<number, number>,
                matchingSets: {} as Record<number, Set<string>>,
                nameMatches: {} as Record<number, boolean>,
            };
        }

        type ScoredApp = {
            app: AppConfig;
            originalIndex: number;
            nameMatch: boolean;
            matchCount: number;
            matchingDescs: Set<string>;
        };
        const scored: ScoredApp[] = sortedBaseApps.map((app, i) => {
            const nameMatch = app.app.toLowerCase().includes(q);
            const matchingDescs = new Set<string>();
            for (const group of app.groups) {
                for (const shortcut of group.shortcuts) {
                    if (shortcut.desc.toLowerCase().includes(q)) {
                        matchingDescs.add(shortcut.desc);
                    }
                }
            }
            return {
                app,
                originalIndex: i,
                nameMatch,
                matchCount: matchingDescs.size,
                matchingDescs,
            };
        });

        // Sort search results by relevance, while preserving tag buckets when the
        // category grouping view is enabled.
        scored.sort((a, b) => {
            if (groupByTag) {
                const tagCompare = comparePrimaryTags(a.app, b.app);
                if (tagCompare !== 0) return tagCompare;
            }

            const aHasMatch = a.nameMatch || a.matchCount > 0 ? 1 : 0;
            const bHasMatch = b.nameMatch || b.matchCount > 0 ? 1 : 0;
            if (aHasMatch !== bHasMatch) return bHasMatch - aHasMatch;
            if (a.matchCount !== b.matchCount)
                return b.matchCount - a.matchCount;
            return compareApps(a.app, b.app);
        });

        const sortedApps = scored.map((s) => s.app);
        const matchCounts: Record<number, number> = {};
        const matchingSets: Record<number, Set<string>> = {};
        const nameMatches: Record<number, boolean> = {};
        scored.forEach((s, newIdx) => {
            matchCounts[newIdx] = s.matchCount;
            matchingSets[newIdx] = s.matchingDescs;
            nameMatches[newIdx] = s.nameMatch;
        });

        return { sortedApps, matchCounts, matchingSets, nameMatches };
    });

    let displayApps = $derived(searchResults.sortedApps);
    let selectedApp = $derived(displayApps[selectedIndex] ?? null);
    let currentMatchingDescs = $derived(
        searchResults.matchingSets[selectedIndex] ?? new Set<string>(),
    );
    let configPathDisplay = $derived.by(() => {
        const os: string = currentOS;
        return os === "windows"
            ? "%LOCALAPPDATA%\\button\\apps\\"
            : "~/.config/button/apps/";
    });
    let appVersionLabel = $derived.by(() => {
        const version = appInfo.version.trim();
        if (!version) {
            return "dev";
        }
        return version.startsWith("v") ? version : `v${version}`;
    });

    // --- Auto-select first app when search changes ---
    $effect(() => {
        // Access searchQuery to track it
        searchQuery;
        selectedIndex = 0;
    });

    // --- Notifications ---
    function addNotification(type: NotificationType, message: string) {
        notifications = [
            ...notifications,
            { id: nextNotificationId++, type, message },
        ];
    }

    // --- Utilities ---
    function escapeHtml(text: string): string {
        const map: Record<string, string> = {
            "&": "&amp;",
            "<": "&lt;",
            ">": "&gt;",
            '"': "&quot;",
            "'": "&#039;",
        };
        return text.replace(/[&<>"']/g, (char) => map[char]);
    }

    function dismissNotification(id: number) {
        notifications = notifications.filter((n) => n.id !== id);
    }

    function cloneShortcut(shortcut: Shortcut): Shortcut {
        const clone: Shortcut = {
            desc: shortcut.desc,
        };
        if (shortcut.keys) clone.keys = shortcut.keys.map((b) => [...b]);
        if (shortcut.linux) clone.linux = shortcut.linux.map((b) => [...b]);
        if (shortcut.macos) clone.macos = shortcut.macos.map((b) => [...b]);
        return clone;
    }

    function cloneAppConfig(app: AppConfig): AppConfig {
        return {
            ...app,
            tags: app.tags ? [...app.tags] : undefined,
            groups: app.groups.map(
                (group): Group => ({
                    category: group.category,
                    shortcuts: group.shortcuts.map(cloneShortcut),
                }),
            ),
        };
    }

    function closeShortcutForm() {
        showShortcutForm = false;
        shortcutTarget = null;
    }

    function closeShortcutDeleteConfirm() {
        showShortcutDeleteConfirm = false;
        shortcutTarget = null;
    }

    function persistShortcutApp(appConfig: AppConfig, onSuccess: () => void) {
        UpdateApp(appConfig.app, appConfig as any, false)
            .then((warning: string) => {
                if (warning) {
                    addNotification("warning", warning);
                    return;
                }

                onSuccess();
                loadApps();
            })
            .catch((err: any) => {
                addNotification("error", String(err));
            });
    }

    function uniqueCategories(app: AppConfig): string[] {
        const seen = new Set<string>();
        const categories: string[] = [];

        for (const group of app.groups) {
            const category = group.category.trim();
            if (!category || seen.has(category)) continue;
            seen.add(category);
            categories.push(category);
        }

        return categories;
    }

    // --- Data loading ---
    function applyResponse(resp: AppsResponse) {
        apps = resp.apps ?? [];
        // Clear stale warning notifications from the previous load, then add fresh ones
        notifications = notifications.filter((n) => n.type !== "warning");
        for (const w of resp.warnings ?? []) {
            addNotification("warning", w);
        }
        // Keep selection in bounds
        if (selectedIndex >= apps.length) {
            selectedIndex = Math.max(0, apps.length - 1);
        }
    }

    function loadApps() {
        GetApps()
            .then((result: AppsResponse) => applyResponse(result))
            .catch((err: any) => {
                addNotification("error", String(err));
                apps = [];
            });
    }

    // --- CRUD handlers ---
    function handleCreateApp() {
        appFormMode = "create";
        editTargetIndex = -1;
        showAppForm = true;
    }

    function handleEditApp(displayIndex: number) {
        appFormMode = "edit";
        editTargetIndex = displayIndex;
        showAppForm = true;
    }

    function handleDeleteApp(displayIndex: number) {
        deleteTargetIndex = displayIndex;
        showDeleteConfirm = true;
    }

    function handleCreateShortcut() {
        if (!selectedApp) return;
        shortcutFormMode = "create";
        shortcutTarget = {
            appIndex: selectedIndex,
            groupIndex: -1,
            shortcutIndex: -1,
        };
        showShortcutForm = true;
    }

    function handleEditShortcut(groupIndex: number, shortcutIndex: number) {
        if (!selectedApp) return;
        shortcutFormMode = "edit";
        shortcutTarget = {
            appIndex: selectedIndex,
            groupIndex,
            shortcutIndex,
        };
        showShortcutForm = true;
    }

    function handleDeleteShortcut(groupIndex: number, shortcutIndex: number) {
        if (!selectedApp) return;
        shortcutTarget = {
            appIndex: selectedIndex,
            groupIndex,
            shortcutIndex,
        };
        showShortcutDeleteConfirm = true;
    }

    function confirmDelete() {
        const app = displayApps[deleteTargetIndex];
        if (!app) return;
        DeleteApp(app.app)
            .then(() => {
                showDeleteConfirm = false;
                deleteTargetIndex = -1;
                loadApps();
            })
            .catch((err: any) => {
                addNotification("error", String(err));
                showDeleteConfirm = false;
            });
    }

    function handleShortcutSave(payload: {
        category: string;
        shortcut: Shortcut;
    }) {
        const target = shortcutTarget;
        if (!target) return;

        const sourceApp = displayApps[target.appIndex];
        if (!sourceApp) return;

        const updatedApp = cloneAppConfig(sourceApp);
        const targetCategory = payload.category.trim();

        if (shortcutFormMode === "create") {
            const existingGroup = updatedApp.groups.find(
                (group) => group.category.trim() === targetCategory,
            );

            if (existingGroup) {
                existingGroup.shortcuts = [
                    ...existingGroup.shortcuts,
                    cloneShortcut(payload.shortcut),
                ];
            } else {
                updatedApp.groups = [
                    ...updatedApp.groups,
                    {
                        category: targetCategory,
                        shortcuts: [cloneShortcut(payload.shortcut)],
                    },
                ];
            }

            persistShortcutApp(updatedApp, closeShortcutForm);
            return;
        }

        const sourceGroup = updatedApp.groups[target.groupIndex];
        const sourceShortcut = sourceGroup?.shortcuts[target.shortcutIndex];
        if (!sourceGroup || !sourceShortcut) return;

        const sourceCategory = sourceGroup.category.trim();
        const nextShortcut = cloneShortcut(payload.shortcut);

        if (sourceCategory === targetCategory) {
            sourceGroup.shortcuts[target.shortcutIndex] = nextShortcut;
            persistShortcutApp(updatedApp, closeShortcutForm);
            return;
        }

        sourceGroup.shortcuts = sourceGroup.shortcuts.filter(
            (_, index) => index !== target.shortcutIndex,
        );
        if (sourceGroup.shortcuts.length === 0) {
            updatedApp.groups = updatedApp.groups.filter(
                (_, index) => index !== target.groupIndex,
            );
        }

        const destinationGroup = updatedApp.groups.find(
            (group) => group.category.trim() === targetCategory,
        );
        if (destinationGroup) {
            destinationGroup.shortcuts = [
                ...destinationGroup.shortcuts,
                nextShortcut,
            ];
        } else {
            updatedApp.groups = [
                ...updatedApp.groups,
                {
                    category: targetCategory,
                    shortcuts: [nextShortcut],
                },
            ];
        }

        persistShortcutApp(updatedApp, closeShortcutForm);
    }

    function confirmShortcutDelete() {
        const target = shortcutTarget;
        if (!target) return;

        const sourceApp = displayApps[target.appIndex];
        if (!sourceApp) return;

        const updatedApp = cloneAppConfig(sourceApp);
        const sourceGroup = updatedApp.groups[target.groupIndex];
        if (!sourceGroup) return;

        sourceGroup.shortcuts = sourceGroup.shortcuts.filter(
            (_, index) => index !== target.shortcutIndex,
        );

        if (sourceGroup.shortcuts.length === 0) {
            updatedApp.groups = updatedApp.groups.filter(
                (_, index) => index !== target.groupIndex,
            );
        }

        persistShortcutApp(updatedApp, closeShortcutDeleteConfirm);
    }

    function handleFormSave(appConfig: AppConfig) {
        if (appFormMode === "create") {
            CreateApp(appConfig as any)
                .then(() => {
                    showAppForm = false;
                    loadApps();
                })
                .catch((err: any) => {
                    addNotification("error", String(err));
                });
        } else {
            const oldName = displayApps[editTargetIndex]?.app ?? "";
            UpdateApp(oldName, appConfig as any, false)
                .then((warning: string) => {
                    if (warning) {
                        // There's a file collision — ask user to confirm overwrite
                        pendingOverwriteApp = appConfig;
                        pendingOverwriteOldName = oldName;
                        showOverwriteConfirm = true;
                        return;
                    }
                    showAppForm = false;
                    loadApps();
                })
                .catch((err: any) => {
                    addNotification("error", String(err));
                });
        }
    }

    function confirmOverwrite() {
        if (!pendingOverwriteApp) return;
        UpdateApp(pendingOverwriteOldName, pendingOverwriteApp as any, true)
            .then(() => {
                showOverwriteConfirm = false;
                showAppForm = false;
                pendingOverwriteApp = null;
                pendingOverwriteOldName = "";
                loadApps();
            })
            .catch((err: any) => {
                addNotification("error", String(err));
                showOverwriteConfirm = false;
            });
    }

    // --- Keyboard navigation ---
    function hasBlockingOverlay(): boolean {
        return (
            showDeleteConfirm ||
            showShortcutDeleteConfirm ||
            showAppForm ||
            showShortcutForm ||
            showOverwriteConfirm ||
            showHelp ||
            showDonate ||
            showWelcome ||
            showImport
        );
    }

    function isEditableTarget(target: EventTarget | null): boolean {
        if (!(target instanceof HTMLElement)) return false;
        return (
            target instanceof HTMLInputElement ||
            target instanceof HTMLTextAreaElement ||
            target instanceof HTMLSelectElement ||
            target.isContentEditable
        );
    }

    function setSearchInput(element: HTMLInputElement | undefined) {
        searchInput = element;
    }

    function setDetailBody(element: HTMLDivElement | undefined) {
        detailBody = element;
    }

    function focusSearch(selectText = true) {
        searchInput?.focus();
        if (selectText) {
            searchInput?.select();
        }
    }

    function blurSearch() {
        searchInput?.blur();
    }

    function selectNextApp() {
        if (selectedIndex < displayApps.length - 1) {
            selectedIndex++;
        }
    }

    function selectPreviousApp() {
        if (selectedIndex > 0) {
            selectedIndex--;
        }
    }

    function cycleSortMode() {
        sortMode = sortMode === "alpha" ? "last-updated" : "alpha";
        selectedIndex = 0;
        persistListPreferences();
    }

    function toggleGroupByTag() {
        groupByTag = !groupByTag;
        selectedIndex = 0;
        persistListPreferences();
    }

    function persistListPreferences() {
        SetListPreferences(sortMode, groupByTag).catch((err: any) => {
            addNotification("error", String(err));
        });
    }

    function scrollDetailPane(direction: "up" | "down") {
        const amount = Math.max((detailBody?.clientHeight ?? 320) * 0.75, 160);
        detailBody?.scrollBy({
            top: direction === "down" ? amount : -amount,
            behavior: "smooth",
        });
    }

    function launchFlamingo() {
        flamingoTrigger += 1;
    }

    function isSearchInputTarget(target: EventTarget | null): boolean {
        return target === searchInput;
    }

    function handleKeydown(e: KeyboardEvent) {
        const key = e.key.toLowerCase();
        const editableTarget = isEditableTarget(e.target);
        const searchTarget = isSearchInputTarget(e.target);

        if (hasBlockingOverlay()) {
            return;
        }

        if (
            e.ctrlKey &&
            !e.metaKey &&
            !e.altKey &&
            !e.shiftKey &&
            key === "f"
        ) {
            e.preventDefault();
            focusSearch();
            return;
        }

        if (
            !editableTarget &&
            !e.ctrlKey &&
            !e.metaKey &&
            !e.altKey &&
            e.key === "/"
        ) {
            e.preventDefault();
            focusSearch(false);
            return;
        }

        if (
            !editableTarget &&
            !e.ctrlKey &&
            !e.metaKey &&
            !e.altKey &&
            e.key === "?"
        ) {
            e.preventDefault();
            showHelp = true;
            return;
        }

        if (
            !editableTarget &&
            e.ctrlKey &&
            !e.metaKey &&
            !e.altKey &&
            !e.shiftKey &&
            key === "j"
        ) {
            e.preventDefault();
            scrollDetailPane("down");
            return;
        }

        if (
            !editableTarget &&
            e.ctrlKey &&
            !e.metaKey &&
            !e.altKey &&
            !e.shiftKey &&
            key === "k"
        ) {
            e.preventDefault();
            scrollDetailPane("up");
            return;
        }

        if (
            !editableTarget &&
            !e.ctrlKey &&
            !e.metaKey &&
            !e.altKey &&
            key === "f"
        ) {
            e.preventDefault();
            launchFlamingo();
            return;
        }

        if (
            !editableTarget &&
            !e.ctrlKey &&
            !e.metaKey &&
            !e.altKey &&
            key === "e" &&
            selectedApp
        ) {
            e.preventDefault();
            handleEditApp(selectedIndex);
            return;
        }

        if (
            !editableTarget &&
            !e.ctrlKey &&
            !e.metaKey &&
            !e.altKey &&
            key === "n"
        ) {
            e.preventDefault();
            handleCreateApp();
            return;
        }

        if (
            !editableTarget &&
            !e.ctrlKey &&
            !e.metaKey &&
            !e.altKey &&
            key === "d" &&
            selectedApp
        ) {
            e.preventDefault();
            handleDeleteApp(selectedIndex);
            return;
        }

        if (
            !editableTarget &&
            !e.ctrlKey &&
            !e.metaKey &&
            !e.altKey &&
            (key === "h" || key === "l")
        ) {
            e.preventDefault();
            cycleOS();
            return;
        }

        if (
            !editableTarget &&
            !e.ctrlKey &&
            !e.metaKey &&
            !e.altKey &&
            e.key === "Backspace" &&
            searchQuery
        ) {
            e.preventDefault();
            searchQuery = "";
            return;
        }

        if (
            !editableTarget &&
            !e.ctrlKey &&
            !e.metaKey &&
            !e.altKey &&
            key === "s"
        ) {
            e.preventDefault();
            cycleSortMode();
            return;
        }

        if (
            !editableTarget &&
            !e.ctrlKey &&
            !e.metaKey &&
            !e.altKey &&
            key === "j"
        ) {
            e.preventDefault();
            selectNextApp();
            return;
        }

        if (
            !editableTarget &&
            !e.ctrlKey &&
            !e.metaKey &&
            !e.altKey &&
            key === "k"
        ) {
            e.preventDefault();
            selectPreviousApp();
            return;
        }

        if (!editableTarget && e.key === "ArrowDown") {
            e.preventDefault();
            selectNextApp();
        } else if (!editableTarget && e.key === "ArrowUp") {
            e.preventDefault();
            selectPreviousApp();
        } else if (searchTarget && (e.key === "Escape" || e.key === "Enter")) {
            e.preventDefault();
            blurSearch();
        }
    }

    // --- Welcome / Import handlers ---
    function handleWelcomeContinue() {
        showWelcome = false;
        cameFromWelcome = true;
        showImport = true;
    }

    function handleWelcomeSkip() {
        showWelcome = false;
        SetHasSeenWelcome().catch(() => {});
    }

    function handleImport(filenames: string[]) {
        ImportRegistryApps(filenames)
            .then((count: number) => {
                showImport = false;
                if (cameFromWelcome) {
                    SetHasSeenWelcome().catch(() => {});
                    cameFromWelcome = false;
                }
                if (count > 0) {
                    addNotification(
                        "info",
                        `Imported ${count} ${count === 1 ? "app" : "apps"} from the registry.`,
                    );
                }
                loadApps();
            })
            .catch((err: any) => {
                addNotification("error", String(err));
            });
    }

    function handleImportClose() {
        showImport = false;
        if (cameFromWelcome) {
            SetHasSeenWelcome().catch(() => {});
            cameFromWelcome = false;
        }
    }

    onMount(() => {
        GetAppInfo()
            .then((info: AppInfo) => {
                appInfo = info;
            })
            .catch(() => {});

        // Detect OS
        GetCurrentOS()
            .then((os: string) => {
                if (os === "darwin" || os === "linux" || os === "windows") {
                    currentOS = os;
                }
            })
            .catch(() => {}); // fallback to 'linux'

        // Check for .yml → .yaml migration results
        GetMigrationResult()
            .then((result: any) => {
                if (result.warnings) {
                    for (const w of result.warnings) {
                        addNotification("warning", w);
                    }
                }
                if (result.migrated > 0) {
                    const n = result.migrated;
                    addNotification(
                        "info",
                        `Migrated ${n} ${n === 1 ? "file" : "files"} from <code>.yml</code> to <code>.yaml</code>.`,
                    );
                }
            })
            .catch(() => {}); // non-critical

        // Check if welcome panel should be shown
        GetUserConfig()
            .then((cfg) => {
                if (
                    cfg.lastSortMode === "alpha" ||
                    cfg.lastSortMode === "last-updated"
                ) {
                    sortMode = cfg.lastSortMode;
                }
                groupByTag = cfg.groupByTag ?? false;
                if (!cfg.hasSeenWelcome) {
                    showWelcome = true;
                }
            })
            .catch(() => {});

        loadApps();

        const cleanup = EventsOn("config:changed", (resp: AppsResponse) => {
            applyResponse(resp);
        });

        return cleanup;
    });

    function setOS(os: "linux" | "darwin" | "windows") {
        currentOS = os;
    }

    function cycleOS() {
        const order: ("linux" | "darwin" | "windows")[] = [
            "linux",
            "windows",
            "darwin",
        ];
        const idx = order.indexOf(currentOS);
        currentOS = order[(idx + 1) % order.length];
    }

    // --- Resize handle ---
    function onResizeStart(e: PointerEvent) {
        e.preventDefault();
        isResizing = true;
        const target = e.currentTarget as HTMLElement;
        target.setPointerCapture(e.pointerId);
    }

    function onResizeMove(e: PointerEvent) {
        if (!isResizing) return;
        const newWidth = Math.min(Math.max(e.clientX, 200), 500);
        listWidth = newWidth;
    }

    function onResizeEnd() {
        isResizing = false;
    }
</script>

<svelte:window onkeydown={handleKeydown} />

<main class="app-shell">
    <FlamingoEasterEgg trigger={flamingoTrigger} />

    <Toolbar
        bind:searchQuery
        bind:showHelp
        bind:showDonate
        bind:showImport
        {currentOS}
        matchingDescs={currentMatchingDescs}
        onSetOS={setOS}
        onSearchInput={setSearchInput}
    />

    <NotificationBar {notifications} onDismiss={dismissNotification} />

    <div class="panels" class:panels--resizing={isResizing}>
        <AppList
            apps={displayApps}
            {selectedIndex}
            {searchQuery}
            matchCounts={searchResults.matchCounts}
            nameMatches={searchResults.nameMatches}
            width={listWidth}
            {sortMode}
            {groupByTag}
            versionLabel={appVersionLabel}
            onSelect={(i) => {
                selectedIndex = i;
            }}
            onToggleSort={() => {
                cycleSortMode();
            }}
            onToggleGroupByTag={toggleGroupByTag}
            onCreateApp={handleCreateApp}
            onEditApp={handleEditApp}
            onDeleteApp={handleDeleteApp}
        />
        <!-- svelte-ignore a11y_no_static_element_interactions -->
        <div
            class="resize-handle"
            onpointerdown={onResizeStart}
            onpointermove={onResizeMove}
            onpointerup={onResizeEnd}
            onpointercancel={onResizeEnd}
        ></div>
        <ShortcutDetail
            app={selectedApp}
            {currentOS}
            {searchQuery}
            matchingDescs={currentMatchingDescs}
            onBodyMount={setDetailBody}
            onCreateShortcut={handleCreateShortcut}
            onEditShortcut={handleEditShortcut}
            onDeleteShortcut={handleDeleteShortcut}
            onEdit={() => handleEditApp(selectedIndex)}
            onDelete={() => handleDeleteApp(selectedIndex)}
        />
    </div>

    {#if showDeleteConfirm && displayApps[deleteTargetIndex]}
        <ConfirmDialog
            title="Delete App"
            message={`Are you sure you want to delete app "${escapeHtml(displayApps[deleteTargetIndex].app)}"?<br>This will remove its config file from <code>${configPathDisplay}</code>.`}
            confirmLabel="Delete"
            danger={true}
            onConfirm={confirmDelete}
            onCancel={() => {
                showDeleteConfirm = false;
                deleteTargetIndex = -1;
            }}
        />
    {/if}

    {#if showAppForm}
        <AppFormModal
            mode={appFormMode}
            initial={appFormMode === "edit"
                ? (displayApps[editTargetIndex] ?? null)
                : null}
            onSave={handleFormSave}
            onCancel={() => {
                showAppForm = false;
            }}
        />
    {/if}

    {#if showShortcutForm && shortcutTarget && displayApps[shortcutTarget.appIndex]}
        <ShortcutFormModal
            mode={shortcutFormMode}
            initial={shortcutFormMode === "edit"
                ? {
                      category:
                          displayApps[shortcutTarget.appIndex].groups[
                              shortcutTarget.groupIndex
                          ]?.category ?? "",
                      shortcut: displayApps[shortcutTarget.appIndex].groups[
                          shortcutTarget.groupIndex
                      ]?.shortcuts[shortcutTarget.shortcutIndex] ?? {
                          desc: "",
                      },
                  }
                : null}
            existingCategories={uniqueCategories(
                displayApps[shortcutTarget.appIndex],
            )}
            onSave={handleShortcutSave}
            onCancel={closeShortcutForm}
        />
    {/if}

    {#if showHelp}
        <HelpPanel
            appName={appInfo.name}
            versionLabel={appVersionLabel}
            {currentOS}
            onClose={() => (showHelp = false)}
        />
    {/if}

    {#if showDonate}
        <DonatePanel onClose={() => (showDonate = false)} />
    {/if}

    {#if showWelcome}
        <WelcomePanel
            onContinue={handleWelcomeContinue}
            onSkip={handleWelcomeSkip}
        />
    {/if}

    {#if showImport}
        <ImportPanel onImport={handleImport} onClose={handleImportClose} />
    {/if}

    {#if showShortcutDeleteConfirm && shortcutTarget && displayApps[shortcutTarget.appIndex]}
        <ConfirmDialog
            title="Delete Shortcut"
            message={`Are you sure you want to delete "${escapeHtml(displayApps[shortcutTarget.appIndex].groups[shortcutTarget.groupIndex]?.shortcuts[shortcutTarget.shortcutIndex]?.desc ?? "")}" from app "${escapeHtml(displayApps[shortcutTarget.appIndex].app)}"?`}
            confirmLabel="Delete"
            danger={true}
            onConfirm={confirmShortcutDelete}
            onCancel={closeShortcutDeleteConfirm}
        />
    {/if}

    {#if showOverwriteConfirm}
        <ConfirmDialog
            title="Overwrite Existing File"
            message="A config file for this app name already exists. Saving will overwrite it."
            confirmLabel="Overwrite"
            danger={true}
            onConfirm={confirmOverwrite}
            onCancel={() => {
                showOverwriteConfirm = false;
                pendingOverwriteApp = null;
            }}
        />
    {/if}
</main>

<style>
    .app-shell {
        position: relative;
        display: flex;
        flex-direction: column;
        height: 100vh;
        background: #121212;
        color: #e5e5e5;
        overflow: hidden;
        font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", system-ui,
            sans-serif;
    }

    .panels {
        flex: 1;
        display: flex;
        min-height: 0;
    }

    .panels--resizing {
        user-select: none;
        cursor: col-resize;
    }

    .resize-handle {
        width: 2px;
        cursor: col-resize;
        background: transparent;
        transition: background 0.15s;
        flex-shrink: 0;
        margin: 0 -1px;
    }

    .resize-handle:hover,
    .panels--resizing .resize-handle {
        background: #3a88ed;
    }
</style>
