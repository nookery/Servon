<script setup lang="ts">
import { ref, onMounted } from 'vue'
import PageContainer from '../layouts/PageContainer.vue'
import { useConfirm } from '../composables/useConfirm'
import { useToast } from '../composables/useToast'
import type { LogEntry, LogStats, LogFile } from '../types/log'
import * as logApi from '../api/logs'
import { RiFileListLine } from '@remixicon/vue'

const confirm = useConfirm()
const toast = useToast()

const currentDir = ref('')
const logFiles = ref<LogFile[]>([])
const selectedFile = ref<string>('')
const logEntries = ref<LogEntry[]>([])
const logStats = ref<LogStats | null>(null)
const searchKeyword = ref('')
const loading = ref(false)
const error = ref<string | null>(null)

// 加载日志文件列表
async function loadLogFiles() {
    try {
        loading.value = true
        logFiles.value = await logApi.getLogFiles(currentDir.value)
        error.value = null
    } catch (err: any) {
        error.value = '获取日志文件列表失败: ' + (err.response?.data?.error || err.message || '未知错误')
    } finally {
        loading.value = false
    }
}

// 加载日志内容
async function loadLogEntries() {
    if (!selectedFile.value) return

    try {
        loading.value = true
        logEntries.value = await logApi.getLogEntries(selectedFile.value)
        error.value = null
    } catch (err: any) {
        error.value = '获取日志内容失败: ' + (err.response?.data?.error || err.message || '未知错误')
    } finally {
        loading.value = false
    }
}

// 搜索日志
async function handleSearch() {
    if (!searchKeyword.value) return

    try {
        loading.value = true
        logEntries.value = await logApi.searchLogs(currentDir.value, searchKeyword.value)
        error.value = null
    } catch (err: any) {
        error.value = '搜索日志失败: ' + (err.response?.data?.error || err.message || '未知错误')
    } finally {
        loading.value = false
    }
}

// 加载统计信息
async function loadStats() {
    try {
        logStats.value = await logApi.getLogStats(currentDir.value)
        error.value = null
    } catch (err: any) {
        error.value = '获取日志统计失败: ' + (err.response?.data?.error || err.message || '未知错误')
    }
}

// 清理旧日志
async function handleCleanLogs() {
    if (await confirm.warning('清理日志', '确定要清理30天前的日志吗？此操作不可撤销。', {
        confirmText: '清理'
    })) {
        try {
            await logApi.cleanOldLogs(30)
            toast.success('清理日志成功')
            await loadLogFiles()
        } catch (err: any) {
            toast.error('清理日志失败: ' + err.message)
        }
    }
}

// 格式化日志级别样式
function getLevelClass(level: string): string {
    const levelMap: Record<string, string> = {
        error: 'badge-error',
        warn: 'badge-warning',
        info: 'badge-info',
        debug: 'badge-ghost'
    }
    return `badge ${levelMap[level.toLowerCase()] || 'badge-ghost'}`
}

onMounted(() => {
    loadLogFiles()
    loadStats()
})
</script>

<template>
    <PageContainer title="日志管理" :error="error" :empty="!logFiles.length" empty-text="暂无日志文件"
        empty-description="系统运行后将自动生成日志文件" :empty-icon="RiFileListLine">
        <template #header>
            <div class="flex justify-between items-center mb-4">
                <div class="flex gap-2 items-center">
                    <input type="text" v-model="currentDir" placeholder="日志目录" class="input input-bordered input-sm" />
                    <button class="btn btn-sm" @click="loadLogFiles">刷新</button>
                </div>
                <div class="flex gap-2">
                    <div class="join">
                        <input type="text" v-model="searchKeyword" placeholder="搜索日志"
                            class="input input-bordered input-sm join-item" @keyup.enter="handleSearch" />
                        <button class="btn btn-sm join-item" @click="handleSearch">搜索</button>
                    </div>
                    <button class="btn btn-error btn-sm" @click="handleCleanLogs">
                        清理旧日志
                    </button>
                </div>
            </div>

            <!-- 日志统计 -->
            <div v-if="logStats" class="stats shadow mb-4">
                <div class="stat">
                    <div class="stat-title">错误</div>
                    <div class="stat-value text-error">{{ logStats.error }}</div>
                </div>
                <div class="stat">
                    <div class="stat-title">警告</div>
                    <div class="stat-value text-warning">{{ logStats.warn }}</div>
                </div>
                <div class="stat">
                    <div class="stat-title">信息</div>
                    <div class="stat-value text-info">{{ logStats.info }}</div>
                </div>
                <div class="stat">
                    <div class="stat-title">调试</div>
                    <div class="stat-value text-base-content">{{ logStats.debug }}</div>
                </div>
            </div>
        </template>

        <!-- 日志内容 -->
        <div class="grid grid-cols-12 gap-4">
            <!-- 日志文件列表 -->
            <div class="col-span-3">
                <div class="card bg-base-200">
                    <div class="card-body p-2">
                        <ul class="menu bg-base-200 w-full">
                            <li v-for="file in logFiles" :key="file.path">
                                <a :class="{ 'active': selectedFile === file.path }"
                                    @click="selectedFile = file.path; loadLogEntries()">
                                    {{ file.path }}
                                </a>
                            </li>
                        </ul>
                    </div>
                </div>
            </div>

            <!-- 日志内容 -->
            <div class="col-span-9">
                <div class="card bg-base-200">
                    <div class="card-body p-4">
                        <div v-if="loading" class="flex justify-center">
                            <span class="loading loading-spinner loading-lg"></span>
                        </div>
                        <div v-else class="overflow-x-auto">
                            <table class="table table-xs">
                                <thead>
                                    <tr>
                                        <th>时间</th>
                                        <th>级别</th>
                                        <th>调用位置</th>
                                        <th>消息</th>
                                    </tr>
                                </thead>
                                <tbody>
                                    <tr v-for="entry in logEntries" :key="entry.time">
                                        <td class="whitespace-nowrap">{{ entry.time }}</td>
                                        <td>
                                            <span :class="getLevelClass(entry.level)">
                                                {{ entry.level }}
                                            </span>
                                        </td>
                                        <td class="text-xs">{{ entry.caller }}</td>
                                        <td class="whitespace-pre-wrap">{{ entry.message }}</td>
                                    </tr>
                                </tbody>
                            </table>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </PageContainer>
</template>

<style scoped>
.stats {
    @apply grid grid-cols-4 w-full;
}
</style>