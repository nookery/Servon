<template>
    <div>
        <!-- 错误提示 -->
        <div v-if="error && !refreshing" class="mb-4 text-red-600">{{ error }}</div>

        <!-- 文件列表 -->
        <div>
            <!-- 工具栏 -->
            <div class="flex justify-between items-center mb-4">
                <div class="flex items-center gap-4">
                    <!-- 当前路径 -->
                    <div class="flex items-center gap-1 text-sm text-gray-600">
                        <button @click="navigateToPath('/')" class="hover:text-blue-600"
                            :disabled="currentPath === '/' || refreshing">
                            根目录
                        </button>
                        <template v-for="(segment, index) in pathSegments" :key="index">
                            <span>/</span>
                            <button @click="navigateToPath(getPathUpToIndex(index))" class="hover:text-blue-600"
                                :disabled="refreshing">
                                {{ segment }}
                            </button>
                        </template>
                    </div>
                    <!-- 刷新按钮 -->
                    <button @click="loadFiles(true)"
                        class="p-2 text-gray-600 hover:text-blue-600 rounded-lg hover:bg-gray-100"
                        :disabled="refreshing">
                        <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" :class="{ 'animate-spin': refreshing }"
                            viewBox="0 0 20 20" fill="currentColor">
                            <path fill-rule="evenodd"
                                d="M4 2a1 1 0 011 1v2.101a7.002 7.002 0 0111.601 2.566 1 1 0 11-1.885.666A5.002 5.002 0 005.999 7H9a1 1 0 010 2H4a1 1 0 01-1-1V3a1 1 0 011-1zm.008 9.057a1 1 0 011.276.61A5.002 5.002 0 0014.001 13H11a1 1 0 110-2h5a1 1 0 011 1v5a1 1 0 11-2 0v-2.101a7.002 7.002 0 01-11.601-2.566 1 1 0 01.61-1.276z"
                                clip-rule="evenodd" />
                        </svg>
                    </button>
                </div>
            </div>

            <!-- 文件列表表格 -->
            <div class="bg-white shadow rounded-lg overflow-hidden">
                <table class="min-w-full divide-y divide-gray-200">
                    <thead class="bg-gray-50">
                        <tr>
                            <th scope="col"
                                class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                名称
                            </th>
                            <th scope="col"
                                class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                大小
                            </th>
                            <th scope="col"
                                class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                修改时间
                            </th>
                            <th scope="col"
                                class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                权限
                            </th>
                        </tr>
                    </thead>
                    <tbody class="bg-white divide-y divide-gray-200 relative min-h-[200px]">
                        <!-- 加载遮罩 -->
                        <tr v-if="refreshing">
                            <td colspan="4" class="absolute inset-0 bg-gray-50/50">
                                <div class="flex items-center justify-center h-full">
                                    <div class="text-sm text-gray-500">加载中...</div>
                                </div>
                            </td>
                        </tr>
                        <!-- 返回上级目录 -->
                        <tr v-if="currentPath !== '/'" class="hover:bg-gray-50 cursor-pointer" @click="navigateUp"
                            :class="{ 'pointer-events-none opacity-50': refreshing }">
                            <td class="px-6 py-4 whitespace-nowrap">
                                <div class="flex items-center">
                                    <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-gray-400 mr-2"
                                        viewBox="0 0 20 20" fill="currentColor">
                                        <path fill-rule="evenodd"
                                            d="M7.707 14.707a1 1 0 01-1.414 0l-4-4a1 1 0 010-1.414l4-4a1 1 0 011.414 1.414L5.414 9H17a1 1 0 110 2H5.414l2.293 2.293a1 1 0 010 1.414z"
                                            clip-rule="evenodd" />
                                    </svg>
                                    <span class="text-sm text-gray-900">上级目录</span>
                                </div>
                            </td>
                            <td colspan="3"></td>
                        </tr>
                        <!-- 文件和目录列表 -->
                        <template v-if="!error && files.length > 0">
                            <tr v-for="item in files" :key="item.path" class="hover:bg-gray-50" :class="{
                                'cursor-pointer': item.isDir,
                                'pointer-events-none opacity-50': refreshing
                            }" @click="item.isDir ? navigateToPath(item.path) : null">
                                <td class="px-6 py-4 whitespace-nowrap">
                                    <div class="flex items-center">
                                        <!-- 文件夹图标 -->
                                        <svg v-if="item.isDir" xmlns="http://www.w3.org/2000/svg"
                                            class="h-5 w-5 text-blue-500 mr-2" viewBox="0 0 20 20" fill="currentColor">
                                            <path
                                                d="M2 6a2 2 0 012-2h5l2 2h5a2 2 0 012 2v6a2 2 0 01-2 2H4a2 2 0 01-2-2V6z" />
                                        </svg>
                                        <!-- 文件图标 -->
                                        <svg v-else xmlns="http://www.w3.org/2000/svg"
                                            class="h-5 w-5 text-gray-400 mr-2" viewBox="0 0 20 20" fill="currentColor">
                                            <path fill-rule="evenodd"
                                                d="M4 4a2 2 0 012-2h4.586A2 2 0 0112 2.586L15.414 6A2 2 0 0116 7.414V16a2 2 0 01-2 2H6a2 2 0 01-2-2V4z"
                                                clip-rule="evenodd" />
                                        </svg>
                                        <span class="text-sm text-gray-900">{{ item.name }}</span>
                                    </div>
                                </td>
                                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                                    {{ item.isDir ? '-' : formatSize(item.size) }}
                                </td>
                                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                                    {{ formatDate(item.modTime) }}
                                </td>
                                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                                    {{ item.mode }}
                                </td>
                            </tr>
                        </template>
                        <!-- 空目录提示 -->
                        <tr v-if="!error && files.length === 0 && !refreshing">
                            <td colspan="4" class="px-6 py-4 text-center text-gray-500">
                                当前目录为空
                            </td>
                        </tr>
                    </tbody>
                </table>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed, watch } from 'vue'

interface FileInfo {
    name: string
    path: string
    size: number
    isDir: boolean
    mode: string
    modTime: string
}

const refreshing = ref(false)
const error = ref<string | null>(null)
const files = ref<FileInfo[]>([])
const currentPath = ref(getInitialPath())

// 从 URL 获取初始路径
function getInitialPath(): string {
    if (typeof window === 'undefined') return '/'
    const params = new URLSearchParams(window.location.search)
    return params.get('path') || '/'
}

// 更新 URL 但不刷新页面
function updateUrlPath(path: string) {
    if (typeof window === 'undefined') return
    const url = new URL(window.location.href)
    url.searchParams.set('path', path)
    window.history.pushState({}, '', url.toString())
}

// 计算当前路径的分段
const pathSegments = computed(() => {
    return currentPath.value.split('/').filter(Boolean)
})

// 获取到指定索引的路径
function getPathUpToIndex(index: number): string {
    return '/' + pathSegments.value.slice(0, index + 1).join('/')
}

// 导航到上级目录
function navigateUp() {
    if (refreshing.value) return
    const parentPath = currentPath.value.substring(0, currentPath.value.lastIndexOf('/'))
    navigateToPath(parentPath || '/')
}

// 导航到指定路径
function navigateToPath(path: string) {
    if (refreshing.value) return
    currentPath.value = path
    updateUrlPath(path)
    loadFiles()
}

// 格式化文件大小
function formatSize(bytes: number): string {
    if (bytes === 0) return '0 B'
    const k = 1024
    const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
    const i = Math.floor(Math.log(bytes) / Math.log(k))
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

// 格式化日期
function formatDate(dateStr: string): string {
    const date = new Date(dateStr)
    return date.toLocaleString()
}

// 加载文件列表
async function loadFiles(isRefresh = false) {
    try {
        refreshing.value = true
        error.value = null
        const response = await fetch(`/api/system/files?path=${encodeURIComponent(currentPath.value)}`)
        if (!response.ok) {
            throw new Error('获取文件列表失败')
        }
        files.value = await response.json()
    } catch (err) {
        error.value = err instanceof Error ? err.message : '未知错误'
    } finally {
        refreshing.value = false
    }
}

// 监听浏览器的前进/后退
onMounted(() => {
    window.addEventListener('popstate', () => {
        currentPath.value = getInitialPath()
        loadFiles()
    })
    loadFiles()
})
</script>