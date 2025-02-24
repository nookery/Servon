export interface FileInfo {
    name: string
    path: string
    size: number
    isDir: boolean
    mode: string
    modTime: string
    owner: string
    group: string
}

export type SortBy = 'name' | 'size' | 'modTime'
export type SortOrder = 'asc' | 'desc'