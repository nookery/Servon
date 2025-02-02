<template>
    <div>
        <!-- 错误提示 -->
        <div v-if="error && !refreshing" class="mb-4 text-red-600">{{ error }}</div>

        <!-- 端口列表 -->
        <div>
            <!-- 工具栏 -->
            <div class="flex justify-between items-center mb-4">
                <div class="relative">
                    <input type="text" v-model="searchQuery" placeholder="搜索端口、进程..."
                        class="pl-10 pr-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500" />
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-gray-400 absolute left-3 top-2.5"
                        viewBox="0 0 20 20" fill="currentColor">
                        <path fill-rule="evenodd"
                            d="M8 4a4 4 0 100 8 4 4 0 000-8zM2 8a6 6 0 1110.89 3.476l4.817 4.817a1 1 0 01-1.414 1.414l-4.816-4.816A6 6 0 012 8z"
                            clip-rule="evenodd" />
                    </svg>
                </div>
                <button @click="loadPorts(true)"
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

            <!-- 端口列表表格 -->
            <div class="bg-white shadow rounded-lg overflow-hidden">
                <table class="min-w-full divide-y divide-gray-200">
                    <thead class="bg-gray-50">
                        <tr>
                            <th scope="col"
                                class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                端口
                            </th>
                            <th scope="col"
                                class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                协议
                            </th>
                            <th scope="col"
                                class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                状态
                            </th>
                            <th scope="col"
                                class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                进程
                            </th>
                            <th scope="col"
                                class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                用户
                            </th>
                            <th scope="col"
                                class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                监听地址
                            </th>
                        </tr>
                    </thead>
                    <tbody class="bg-white divide-y divide-gray-200 relative min-h-[200px]">
                        <!-- 加载遮罩 -->
                        <tr v-if="refreshing">
                            <td colspan="6" class="absolute inset-0 bg-gray-50/50">
                                <div class="flex items-center justify-center h-full">
                                    <div class="text-sm text-gray-500">加载中...</div>
                                </div>
                            </td>
                        </tr>
                        <!-- 端口列表 -->
                        <template v-if="!error && filteredPorts.length > 0">
                            <tr v-for="port in filteredPorts" :key="port.port + port.protocol" class="hover:bg-gray-50">
                                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                                    {{ port.port }}
                                </td>
                                <td class="px-6 py-4 whitespace-nowrap text-sm">
                                    <span :class="{
                                        'px-2 py-1 rounded text-xs font-medium': true,
                                        'bg-green-100 text-green-800': port.protocol === 'TCP',
                                        'bg-blue-100 text-blue-800': port.protocol === 'UDP'
                                    }">
                                        {{ port.protocol }}
                                    </span>
                                </td>
                                <td class="px-6 py-4 whitespace-nowrap text-sm">
                                    <span :class="{
                                        'px-2 py-1 rounded text-xs font-medium': true,
                                        'bg-green-100 text-green-800': port.state === 'LISTEN',
                                        'bg-yellow-100 text-yellow-800': port.state === 'ESTABLISHED',
                                        'bg-gray-100 text-gray-800': !['LISTEN', 'ESTABLISHED'].includes(port.state)
                                    }">
                                        {{ port.state }}
                                    </span>
                                </td>
                                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                                    <div class="flex flex-col">
                                        <span>{{ port.process }}</span>
                                        <span class="text-xs text-gray-500">PID: {{ port.pid }}</span>
                                        <span class="text-xs text-gray-500 truncate max-w-xs" :title="port.command">
                                            {{ port.command }}
                                        </span>
                                    </div>
                                </td>
                                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                                    {{ port.user }}
                                </td>
                                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                                    {{ port.ipAddress }}
                                </td>
                            </tr>
                        </template>
                        <!-- 空列表提示 -->
                        <tr v-if="!error && filteredPorts.length === 0 && !refreshing">
                            <td colspan="6" class="px-6 py-4 text-center text-gray-500">
                                {{ searchQuery ? '没有找到匹配的端口' : '没有正在使用的端口' }}
                            </td>
                        </tr>
                    </tbody>
                </table>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'

interface PortInfo {
    port: number
    protocol: string
    state: string
    pid: number
    process: string
    command: string
    user: string
    ipAddress: string
}

const refreshing = ref(false)
const error = ref<string | null>(null)
const ports = ref<PortInfo[]>([])
const searchQuery = ref('')
let refreshInterval: number | null = null

// 过滤端口列表
const filteredPorts = computed(() => {
    if (!searchQuery.value) {
        return ports.value
    }
    const query = searchQuery.value.toLowerCase()
    return ports.value.filter(port =>
        port.port.toString().includes(query) ||
        port.process.toLowerCase().includes(query) ||
        port.command.toLowerCase().includes(query) ||
        port.user.toLowerCase().includes(query) ||
        port.ipAddress.toLowerCase().includes(query)
    )
})

// 加载端口列表
async function loadPorts(isRefresh = false) {
    try {
        refreshing.value = true
        error.value = null
        const response = await fetch('/api/system/ports')
        if (!response.ok) {
            throw new Error('获取端口列表失败')
        }
        ports.value = await response.json()
    } catch (err) {
        error.value = err instanceof Error ? err.message : '未知错误'
    } finally {
        refreshing.value = false
    }
}

// 自动刷新（每10秒）
function startAutoRefresh() {
    refreshInterval = window.setInterval(() => {
        loadPorts(true)
    }, 10000)
}

onMounted(() => {
    loadPorts()
    startAutoRefresh()
})

onUnmounted(() => {
    if (refreshInterval !== null) {
        clearInterval(refreshInterval)
    }
})
</script>