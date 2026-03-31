<script lang="ts">
  import { GetApps } from '../wailsjs/go/main/App.js'
  import { EventsOn } from '../wailsjs/runtime/runtime.js'
  import { onMount } from 'svelte'

  type Shortcut = {
    desc: string
    keys: string[]
  }

  type Group = {
    category: string
    shortcuts: Shortcut[]
  }

  type AppConfig = {
    app: string
    icon: string
    groups: Group[]
  }

  type AppsResponse = {
    apps: AppConfig[]
    warnings: string[]
  }

  let apps: AppConfig[] = $state([])
  let warnings: string[] = $state([])
  let error: string = $state("")

  function applyResponse(resp: AppsResponse) {
    apps = resp.apps ?? []
    warnings = resp.warnings ?? []
    error = ""
  }

  function loadApps() {
    GetApps()
      .then((result: AppsResponse) => applyResponse(result))
      .catch((err: any) => {
        error = String(err)
        apps = []
        warnings = []
      })
  }

  onMount(() => {
    loadApps()

    // Listen for live config changes from the file watcher
    const cleanup = EventsOn("config:changed", (resp: AppsResponse) => {
      applyResponse(resp)
    })

    return cleanup
  })
</script>

<main class="flex flex-col items-center h-screen bg-gray-900 text-white p-6 overflow-y-auto">
  <h1 class="text-3xl font-bold mb-2">Button</h1>
  <p class="text-sm text-gray-400 mb-6">Quick-Reference for Keyboard Shortcuts</p>

  {#if error}
    <div class="bg-red-900/50 border border-red-700 rounded-lg p-4 max-w-2xl w-full mb-4">
      <p class="text-red-300 text-sm font-mono">{error}</p>
    </div>
  {/if}

  {#if warnings.length > 0}
    <div class="bg-yellow-900/30 border border-yellow-700/50 rounded-lg p-4 max-w-2xl w-full mb-4">
      {#each warnings as warning}
        <p class="text-yellow-300 text-sm font-mono">{warning}</p>
      {/each}
    </div>
  {/if}

  {#if apps.length === 0 && !error}
    <div class="bg-gray-800 rounded-lg p-6 max-w-2xl w-full text-center">
      <p class="text-gray-400">No YAML files found in <code class="text-gray-300">~/.config/button/apps/</code></p>
      <p class="text-gray-500 text-sm mt-2">Drop <code class="text-gray-400">.yaml</code> files there — they'll appear here automatically.</p>
    </div>
  {/if}

  <div class="flex flex-col gap-4 max-w-2xl w-full">
    {#each apps as app}
      <div class="bg-gray-800 rounded-lg p-4">
        <h2 class="text-xl font-semibold mb-3">{app.app}</h2>

        {#each app.groups as group}
          <div class="mb-3">
            <h3 class="text-xs font-medium text-gray-400 uppercase tracking-wider mb-2">{group.category}</h3>
            <div class="flex flex-col gap-1">
              {#each group.shortcuts as shortcut}
                <div class="flex items-center justify-between py-1.5 px-2 rounded hover:bg-gray-700/50">
                  <span class="text-sm text-gray-200">{shortcut.desc}</span>
                  <div class="flex gap-1">
                    {#each shortcut.keys as key}
                      <kbd class="px-2 py-0.5 text-xs font-mono bg-gray-700 border border-gray-600 rounded text-gray-300">{key}</kbd>
                    {/each}
                  </div>
                </div>
              {/each}
            </div>
          </div>
        {/each}
      </div>
    {/each}
  </div>
</main>
