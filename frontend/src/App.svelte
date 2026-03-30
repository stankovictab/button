<script lang="ts">
  import {Greet} from '../wailsjs/go/main/App.js'

  let resultText: string = $state("Click the button to test the Go bridge")
  let name: string = $state("")

  function greet(): void {
    Greet(name).then((result: string) => resultText = result)
  }
</script>

<main class="flex flex-col items-center justify-center h-screen bg-gray-900 text-white">
  <h1 class="text-4xl font-bold mb-8">Button</h1>
  <p class="text-lg text-gray-400 mb-6">Quick-Reference for Keyboard Shortcuts</p>

  <div class="bg-gray-800 rounded-lg p-6 max-w-md w-full">
    <p class="text-sm text-gray-300 mb-4">{resultText}</p>
    <div class="flex gap-2">
      <input
        bind:value={name}
        class="flex-1 px-3 py-2 bg-gray-700 rounded text-white placeholder-gray-500 outline-none focus:ring-2 focus:ring-blue-500"
        placeholder="Enter your name"
        type="text"
      />
      <button
        onclick={greet}
        class="px-4 py-2 bg-blue-600 hover:bg-blue-500 rounded font-medium transition-colors"
      >
        Greet
      </button>
    </div>
  </div>
</main>
