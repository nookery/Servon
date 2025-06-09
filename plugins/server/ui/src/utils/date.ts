export function formatTimestamp(timestamp: string): string {
    return new Date(timestamp).toLocaleString()
} 