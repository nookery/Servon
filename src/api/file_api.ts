import axios from 'axios'
import type { FileInfo, SortBy, SortOrder } from '../models/FileInfo'

export const fileAPI = {
    // 获取文件列表
    getFiles: (path: string, sortBy?: SortBy, order?: SortOrder) =>
        axios.get<FileInfo[]>(`/web_api/files`, {
            params: {
                path,
                sortBy,
                order
            }
        }),

    // 下载文件
    downloadFile: (path: string) =>
        axios.get(`/web_api/files/download`, {
            params: { path },
            responseType: 'blob'
        }),

    // 删除文件
    deleteFile: (path: string) =>
        axios.delete(`/web_api/files/delete`, { params: { path } }),

    // 创建文件或目录
    createFile: (path: string, type: 'file' | 'directory') =>
        axios.post('/web_api/files/create', { path, type }),

    // 搜索文件
    searchFiles: (path: string, query: string) =>
        axios.get<FileInfo[]>(`/web_api/files/search`, {
            params: { path, query }
        }),

    // 获取文件内容
    getFileContent: (path: string) =>
        axios.get<{ content: string }>(`/web_api/files/content`, {
            params: { path }
        }),

    // 保存文件内容
    saveFileContent: (path: string, content: string) =>
        axios.post('/web_api/files/save', { path, content }),

    // 重命名文件
    renameFile: (oldPath: string, newPath: string) =>
        axios.post('/web_api/files/rename', { oldPath, newPath }),
}
