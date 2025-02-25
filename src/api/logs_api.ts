import { request } from './request'
import type { LogEntry, LogStats, LogFile } from '../types/log'
import axios from 'axios'

export async function getLogFiles(dir: string = ''): Promise<LogFile[]> {
    const { files } = await request.get<{ files: string[] }>('/logs/files', { params: { dir } })
    return files.map(path => ({ path }))
}

export async function getLogEntries(file: string, limit: number = 100): Promise<LogEntry[]> {
    const { entries } = await request.get<{ entries: LogEntry[] }>('/logs/entries', {
        params: { file, limit }
    })
    return entries
}

export async function searchLogs(dir: string, keyword: string): Promise<LogEntry[]> {
    const { data } = await request.get<{ data: LogEntry[] }>('/logs/search', {
        params: { dir, keyword }
    })
    return data
}

export async function getLogStats(dir: string = ''): Promise<LogStats> {
    const { data } = await request.get<{ data: LogStats }>('/logs/stats', { params: { dir } })
    return data
}

export async function cleanOldLogs(days: number = 30): Promise<void> {
    await request.post('/logs/clean', { params: { days } })
}

export async function deleteFile(path: string): Promise<void> {
    await request.post('/logs/delete', { params: { path } })
}

// 清空指定的日志文件
export async function clearLogFile(file: string): Promise<void> {
    await request.post('/logs/clear', { params: { file } })
} 