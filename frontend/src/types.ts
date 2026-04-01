export type Shortcut = {
  desc: string
  keys?: string[]
  linux?: string[]
  macos?: string[]
}

export type Group = {
  category: string
  shortcuts: Shortcut[]
}

export type AppConfig = {
  app: string
  icon: string
  groups: Group[]
  modTime: number
}

export type SortMode = "alpha" | "last-updated"

export type AppsResponse = {
  apps: AppConfig[]
  warnings: string[]
}
