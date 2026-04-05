export type Shortcut = {
    desc: string
    keys?: string[][]
    linux?: string[][]
    macos?: string[][]
}

export type Group = {
    category: string
    shortcuts: Shortcut[]
}

export type AppConfig = {
    app: string
    icon: string
    tags?: string[]
    default?: boolean
    groups: Group[]
    modTime: number
}

export type RegistryEntry = {
    filename: string
    app: string
    icon: string
    tags: string[]
}

export type UserConfig = {
    hasSeenWelcome: boolean
}

export type SortMode = "alpha" | "last-updated"

export type NotificationType = "error" | "warning" | "info"

export type Notification = {
    id: number
    type: NotificationType
    message: string
}

export type AppsResponse = {
    apps: AppConfig[]
    warnings: string[]
}

export type AppInfo = {
    name: string
    version: string
}
