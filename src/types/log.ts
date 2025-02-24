export interface LogEntry {
    time: string;
    level: string;
    caller: string;
    message: string;
    extra?: Record<string, any>;
}

export interface LogStats {
    error: number;
    warn: number;
    info: number;
    debug: number;
}

export interface LogFile {
    path: string;
    modTime: string;
    size: number;
} 