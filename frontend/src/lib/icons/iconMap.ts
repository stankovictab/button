/**
 * Icon map: resolves YAML `icon` field values to SVG strings.
 *
 * - Simple Icons: brand logos imported from the `simple-icons` package
 * - Custom SVGs: raw SVG files from `./custom/` for apps missing from Simple Icons
 * - Lucide: UI chrome icons are handled separately in Toolbar.svelte
 *
 * To add a new icon:
 *   1. If it exists in Simple Icons: import it and add an entry below
 *   2. If not: add an SVG to `./custom/`, import it as `?raw`, and add an entry
 */

import {
    siKofi,
    siLinux,
    siApple,
    si1password,
    siAlacritty,
    siDiscord,
    siFishshell,
    siGhostty,
    siKde,
    siLinear,
    siNeovim,
    siNotion,
    siObsidian,
    siTmux,
    siVivaldi,
    siGnubash,
} from 'simple-icons';

// Custom SVGs (apps not in Simple Icons)
import buttonSvg from './custom/button.svg?raw';
import yaziSvg from './custom/yazi.svg?raw';
import vscodeSvg from './custom/vscode.svg?raw';
import terminalSvg from './custom/terminal.svg?raw';
import slackSvg from './custom/slack.svg?raw';

interface IconEntry {
    svg: string;
    hex: string; // brand color without #
}

/**
 * Map from YAML `icon` field value to icon data.
 * Keys must match the `icon:` string in app YAML files exactly.
 */
export const iconMap: Record<string, IconEntry> = {
    // Simple Icons
    kofi: { svg: siKofi.svg, hex: siKofi.hex },
    linux: { svg: siLinux.svg, hex: 'cccccc' },
    apple: { svg: siApple.svg, hex: 'cccccc' },
    '1password': { svg: si1password.svg, hex: si1password.hex },
    alacritty: { svg: siAlacritty.svg, hex: siAlacritty.hex },
    discord: { svg: siDiscord.svg, hex: siDiscord.hex },
    fish: { svg: siFishshell.svg, hex: siFishshell.hex },
    ghostty: { svg: siGhostty.svg, hex: siGhostty.hex },
    kde: { svg: siKde.svg, hex: siKde.hex },
    linear: { svg: siLinear.svg, hex: siLinear.hex },
    neovim: { svg: siNeovim.svg, hex: siNeovim.hex },
    notion: { svg: siNotion.svg, hex: 'cccccc' },
    obsidian: { svg: siObsidian.svg, hex: siObsidian.hex },
    tmux: { svg: siTmux.svg, hex: siTmux.hex },
    vivaldi: { svg: siVivaldi.svg, hex: 'cccccc' },
    gnubash: { svg: siGnubash.svg, hex: 'ff8400' },

    // Custom SVGs
    button: { svg: buttonSvg, hex: 'cccccc' },
    yazi: { svg: yaziSvg, hex: '6CB644' },
    vscode: { svg: vscodeSvg, hex: '007ACC' },
    terminal: { svg: terminalSvg, hex: 'a1a1a1' },
    slack: { svg: slackSvg, hex: 'ffffff' },
};
