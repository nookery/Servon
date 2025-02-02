<!-- ProcessList.vue -->
<template>
    <div>
        <!-- 初始加载状态 -->
        <div v-if="loading" class="text-gray-600">加载中...</div>
        <div v-if="error" class="text-red-600">{{ error }}</div>

        <!-- 进程列表 -->
        <div v-if="!loading">
            <!-- 搜索和刷新 -->
            <div class="flex justify-between items-center mb-4">
                <div class="relative">
                    <input type="text" v-model="searchQuery" placeholder="搜索进程..."
                        class="pl-10 pr-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500" />
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-gray-400 absolute left-3 top-2.5"
                        viewBox="0 0 20 20" fill="currentColor">
                        <path fill-rule="evenodd"
                            d="M8 4a4 4 0 100 8 4 4 0 000-8zM2 8a6 6 0 1110.89 3.476l4.817 4.817a1 1 0 01-1.414 1.414l-4.816-4.816A6 6 0 012 8z"
                            clip-rule="evenodd" />
                    </svg>
                </div>
                <button @click="fetchProcesses(true)"
                    class="px-4 py-2 bg-blue-500 text-white rounded-lg hover:bg-blue-600 focus:outline-none focus:ring-2 focus:ring-blue-500 flex items-center gap-2"
                    :disabled="refreshing">
                    <svg v-if="refreshing" class="animate-spin h-4 w-4" xmlns="http://www.w3.org/2000/svg" fill="none"
                        viewBox="0 0 24 24">
                        <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4">
                        </circle>
                        <path class="opacity-75" fill="currentColor"
                            d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z">
                        </path>
                    </svg>
                    {{ refreshing ? '刷新中...' : '刷新' }}
                </button>
            </div>

            <!-- 进程表格 -->
            <div class="bg-white shadow rounded-lg overflow-hidden relative">
                <table class="min-w-full divide-y divide-gray-200">
                    <thead class="bg-gray-50">
                        <tr>
                            <th scope="col"
                                class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                PID
                            </th>
                            <th scope="col"
                                class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                用户
                            </th>
                            <th scope="col"
                                class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider cursor-pointer select-none"
                                @click="toggleSort('cpu')">
                                <div class="flex items-center gap-1">
                                    CPU (%)
                                    <span v-if="sortConfig.field === 'cpu'" class="text-blue-500">
                                        {{ sortConfig.direction === 'desc' ? '↓' : '↑' }}
                                    </span>
                                </div>
                            </th>
                            <th scope="col"
                                class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider cursor-pointer select-none"
                                @click="toggleSort('memory')">
                                <div class="flex items-center gap-1">
                                    内存 (%)
                                    <span v-if="sortConfig.field === 'memory'" class="text-blue-500">
                                        {{ sortConfig.direction === 'desc' ? '↓' : '↑' }}
                                    </span>
                                </div>
                            </th>
                            <th scope="col"
                                class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                命令
                            </th>
                        </tr>
                    </thead>
                    <tbody class="bg-white divide-y divide-gray-200 relative">
                        <!-- 刷新遮罩 -->
                        <tr v-if="refreshing">
                            <td colspan="5" class="absolute inset-0 bg-gray-50/50">
                                <div class="flex items-center justify-center h-full">
                                    <div class="text-sm text-gray-500">正在更新数据...</div>
                                </div>
                            </td>
                        </tr>
                        <tr v-for="process in paginatedProcesses" :key="process.pid">
                            <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                                {{ process.pid }}
                            </td>
                            <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                                {{ process.user }}
                            </td>
                            <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                                {{ process.cpu.toFixed(1) }}
                            </td>
                            <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                                {{ process.memory.toFixed(1) }}
                            </td>
                            <td class="px-6 py-4 text-sm text-gray-900 max-w-md truncate">
                                {{ process.command }}
                            </td>
                        </tr>
                    </tbody>
                </table>
            </div>

            <!-- 分页控件 -->
            <div v-if="totalPages > 1" class="mt-4 flex items-center justify-between">
                <div class="text-sm text-gray-700">
                    共 {{ filteredAndSortedProcesses.length }} 个进程，
                    当前显示第 {{ (currentPage - 1) * pageSize + 1 }} - {{ Math.min(currentPage * pageSize,
                        filteredAndSortedProcesses.length) }} 个
                </div>
                <div class="flex items-center space-x-2">
                    <button @click="changePage(currentPage - 1)" :disabled="currentPage === 1"
                        class="px-3 py-1 rounded border hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed">
                        上一页
                    </button>
                    <div class="flex items-center space-x-1">
                        <template v-for="page in totalPages" :key="page">
                            <!-- 显示首页、末页、当前页及其前后页 -->
                            <button
                                v-if="page === 1 || page === totalPages || (page >= currentPage - 1 && page <= currentPage + 1)"
                                @click="changePage(page)" :class="[
                                    'px-3 py-1 rounded',
                                    currentPage === page
                                        ? 'bg-blue-500 text-white'
                                        : 'hover:bg-gray-50'
                                ]">
                                {{ page }}
                            </button>
                            <!-- 显示省略号 -->
                            <span v-else-if="page === currentPage - 2 || page === currentPage + 2" class="px-2">
                                ...
                            </span>
                        </template>
                    </div>
                    <button @click="changePage(currentPage + 1)" :disabled="currentPage === totalPages"
                        class="px-3 py-1 rounded border hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed">
                        下一页
                    </button>
                </div>
            </div>

            <!-- 原始响应显示 -->
            <div class="mt-4">
                <button class="px-4 py-2 bg-gray-100 hover:bg-gray-200 rounded-md text-sm" @click="showRaw = !showRaw">
                    {{ showRaw ? '隐藏原始响应' : '显示原始响应' }}
                </button>
                <pre v-if="showRaw" class="mt-4 p-4 bg-gray-50 rounded-md overflow-auto">
        {{ JSON.stringify(processes, null, 2) }}
    </pre>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, watch } from 'vue'

interface Process {
    pid: number
    user: string
    cpu: number
    memory: number
    command: string
}

interface SortConfig {
    field: 'cpu' | 'memory' | null
    direction: 'asc' | 'desc'
}

const loading = ref(true)
const refreshing = ref(false)
const error = ref<string | null>(null)
const processes = ref<Process[]>([])
const showRaw = ref(false)
const searchQuery = ref('')
const sortConfig = ref<SortConfig>({
    field: null,
    direction: 'desc'
})
const currentPage = ref(1)
const pageSize = ref(20)
let refreshInterval: number | null = null

const filteredAndSortedProcesses = computed(() => {
    let result = [...processes.value]

    if (searchQuery.value) {
        const query = searchQuery.value.toLowerCase()
        result = result.filter(process =>
            process.command.toLowerCase().includes(query) ||
            process.user.toLowerCase().includes(query) ||
            process.pid.toString().includes(query)
        )
    }

    if (sortConfig.value.field) {
        result.sort((a, b) => {
            const multiplier = sortConfig.value.direction === 'desc' ? -1 : 1
            return (a[sortConfig.value.field!] - b[sortConfig.value.field!]) * multiplier
        })
    }

    return result
})

const totalPages = computed(() => Math.ceil(filteredAndSortedProcesses.value.length / pageSize.value))

const paginatedProcesses = computed(() => {
    const start = (currentPage.value - 1) * pageSize.value
    const end = start + pageSize.value
    return filteredAndSortedProcesses.value.slice(start, end)
})

function changePage(page: number) {
    if (page >= 1 && page <= totalPages.value) {
        currentPage.value = page
    }
}

// 当过滤或排序改变时，重置到第一页
watch([searchQuery, sortConfig], () => {
    currentPage.value = 1
})

function toggleSort(field: 'cpu' | 'memory') {
    if (sortConfig.value.field === field) {
        // 如果已经在按这个字段排序，切换排序方向
        sortConfig.value.direction = sortConfig.value.direction === 'desc' ? 'asc' : 'desc'
    } else {
        // 如果是新的排序字段，默认降序（显示最高的在前面）
        sortConfig.value.field = field
        sortConfig.value.direction = 'desc'
    }
}

async function fetchProcesses(isRefresh = false) {
    try {
        if (isRefresh) {
            refreshing.value = true
        } else {
            loading.value = true
        }
        error.value = null
        const response = await fetch('/api/system/processes')
        if (!response.ok) {
            throw new Error('获取进程列表失败')
        }
        processes.value = await response.json()
    } catch (err) {
        error.value = err instanceof Error ? err.message : '未知错误'
    } finally {
        loading.value = false
        refreshing.value = false
    }
}

// 自动刷新进程列表（每5秒）
function startAutoRefresh() {
    refreshInterval = window.setInterval(() => {
        fetchProcesses(true)
    }, 5000)
}

onMounted(() => {
    fetchProcesses(false)
    startAutoRefresh()
})

onUnmounted(() => {
    if (refreshInterval !== null) {
        clearInterval(refreshInterval)
    }
})
</script>