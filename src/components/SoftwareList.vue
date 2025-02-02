<!-- SoftwareList.vue -->
<template>
    <div>
        <!-- Tab 切换 -->
        <div class="mb-6 border-b border-gray-200">
            <nav class="-mb-px flex space-x-8" aria-label="Tabs">
                <button
                    v-for="tab in tabs"
                    :key="tab.value"
                    class="whitespace-nowrap py-4 px-1 border-b-2 font-medium text-sm"
                    :class="[
                        currentTab === tab.value
                            ? 'border-blue-500 text-blue-600'
                            : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
                    ]"
                    @click="currentTab = tab.value">
                    {{ tab.label }}
                </button>
            </nav>
        </div>

        <!-- 加载状态 -->
        <div v-if="loading" class="text-gray-600">加载中...</div>
        <div v-if="error" class="text-red-600">{{ error }}</div>

        <!-- 软件列表 -->
        <div v-if="!loading && !error">
            <div class="space-y-4">
                <div v-if="filteredSoftware.length === 0" class="text-center py-8 text-gray-500">
                    {{ currentTab === 'installed' ? '暂无已安装软件' : '暂无未安装软件' }}
                </div>
                <div
                    v-for="software in filteredSoftware"
                    :key="software.name"
                    class="bg-white shadow rounded-lg p-6">
                    <div class="flex justify-between items-start">
                        <div>
                            <h3 class="text-lg font-semibold text-gray-900">
                                {{ software.name }}
                            </h3>
                            <p class="text-sm text-gray-500 mt-1">
                                {{ software.version || '未知版本' }}
                            </p>
                            <p v-if="software.description" class="text-gray-700 mt-2">
                                {{ software.description }}
                            </p>
                        </div>
                        <div class="flex flex-col items-end gap-2">
                            <span
                                class="px-2 py-1 text-sm rounded-full"
                                :class="getStatusClass(software.status)">
                                {{ getStatusText(software.status) }}
                            </span>
                            <button
                                class="px-3 py-1 text-sm rounded-md text-white"
                                :class="software.status === 'not_installed' ? 'bg-blue-500 hover:bg-blue-600' : 'bg-red-500 hover:bg-red-600'"
                                @click="software.status === 'not_installed' ? install(software) : uninstall(software)">
                                {{ software.status === 'not_installed' ? '安装' : '卸载' }}
                            </button>
                        </div>
                    </div>
                    <p v-if="software.path" class="text-sm text-gray-500 mt-2">
                        路径: {{ software.path }}
                    </p>
                </div>
            </div>

            <!-- 原始响应显示 -->
            <div class="mt-4">
                <button
                    class="px-4 py-2 bg-gray-100 hover:bg-gray-200 rounded-md text-sm"
                    @click="showRaw = !showRaw">
                    {{ showRaw ? '隐藏原始响应' : '显示原始响应' }}
                </button>
                <pre
                    v-if="showRaw"
                    class="mt-4 p-4 bg-gray-50 rounded-md overflow-auto">
                    {{ JSON.stringify(softwareList, null, 2) }}
                </pre>
            </div>
        </div>
    </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'

const tabs = [
    { label: '已安装', value: 'installed' },
    { label: '未安装', value: 'uninstalled' }
]

const currentTab = ref('installed')
const loading = ref(true)
const error = ref(null)
const softwareList = ref([])
const showRaw = ref(false)

const filteredSoftware = computed(() => {
    return softwareList.value.filter(software =>
        currentTab.value === 'installed'
            ? software.status === 'running' || software.status === 'stopped'
            : software.status === 'not_installed'
    )
})

function getStatusClass(status) {
    switch (status) {
        case 'running':
            return 'bg-green-100 text-green-800'
        case 'stopped':
            return 'bg-red-100 text-red-800'
        default:
            return 'bg-gray-100 text-gray-800'
    }
}

function getStatusText(status) {
    switch (status) {
        case 'running':
            return '运行中'
        case 'stopped':
            return '已停止'
        default:
            return '未安装'
    }
}

async function fetchSoftwareList() {
    try {
        loading.value = true
        error.value = null
        const response = await fetch('/api/system/software')
        if (!response.ok) {
            throw new Error('获取软件列表失败')
        }
        softwareList.value = await response.json()
    } catch (err) {
        error.value = err.message
    } finally {
        loading.value = false
    }
}

async function install(software) {
    if (!confirm(`确定要安装 ${software.name} 吗？`)) return

    try {
        const response = await fetch(`/api/system/software/${software.name}/install`, {
            method: 'POST'
        })
        if (!response.ok) {
            const data = await response.json()
            throw new Error(data.error || '安装失败')
        }
        alert('安装成功！')
        await fetchSoftwareList()
    } catch (err) {
        alert(err.message)
    }
}

async function uninstall(software) {
    if (!confirm(`确定要卸载 ${software.name} 吗？这可能会删除相关的数据和配置。`)) return

    try {
        const response = await fetch(`/api/system/software/${software.name}/uninstall`, {
            method: 'POST'
        })
        if (!response.ok) {
            const data = await response.json()
            throw new Error(data.error || '卸载失败')
        }
        alert('卸载成功！')
        await fetchSoftwareList()
    } catch (err) {
        alert(err.message)
    }
}

onMounted(fetchSoftwareList)
</script> 