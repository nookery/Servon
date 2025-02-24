<script setup lang="ts">
import { ref, onMounted, watch, computed } from 'vue'
import { useConfirm } from '../../composables/useConfirm'
import { useToast } from '../../composables/useToast'
import { useError } from '../../composables/useError'
import type { LogEntry, LogStats } from '../../types/log'
import * as logApi from '../../api/logs'
import IconButton from '../IconButton.vue'
import { RiDeleteBinLine } from '@remixicon/vue'

const props = defineProps<{
    currentDir: string
}>()

const confirm = useConfirm()
const toast = useToast()
const { error } = useError()

const logFiles = ref<LogEntry[]>([])
const selectedFile = ref<string>('')
const logEntries = ref<LogEntry[]>([])
const logStats = ref<LogStats | null>(null)
const searchKeyword = ref('')
const loading = ref(false)
const selectedLevels = ref<string[]>(['error', 'warn', 'info', 'debug'])

// 加载日志文件列表
async function loadLogFiles() {
    try {
        loading.value = true
        logFiles.value = await logApi.getLogFiles(props.currentDir)
        toast.success('日志列表已刷新')
    } catch (err: any) {
        error('获取日志文件列表失败: ' + (err.response?.data?.error || err.message || '未知错误'))
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
        toast.success('日志内容已刷新')
    } catch (err: any) {
        error('获取日志内容失败: ' + (err.response?.data?.error || err.message || '未知错误'))
    } finally {
        loading.value = false
    }
}

// 搜索日志
async function handleSearch() {
    if (!searchKeyword.value) return

    try {
        loading.value = true
        logEntries.value = await logApi.searchLogs(props.currentDir, searchKeyword.value)
    } catch (err: any) {
        error('搜索日志失败: ' + (err.response?.data?.error || err.message || '未知错误'))
    } finally {
        loading.value = false
    }
}

// 加载统计信息
async function loadStats() {
    try {
        logStats.value = await logApi.getLogStats(props.currentDir)
    } catch (err: any) {
        error('获取日志统计失败: ' + (err.response?.data?.error || err.message || '未知错误'))
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
            error('清理日志失败: ' + err.message)
        }
    }
}

// 删除当前选中的日志文件
async function handleDeleteCurrentLog() {
    if (!selectedFile.value) {
        toast.warning('请先选择要删除的日志文件')
        return
    }

    if (await confirm.error('删除日志文件', `确定要删除日志文件 "${selectedFile.value}" 吗？此操作不可恢复！`, {
        confirmText: '删除',
        cancelText: '取消'
    })) {
        try {
            await logApi.deleteFile(selectedFile.value)
            toast.success('日志文件删除成功')
            selectedFile.value = ''
            logEntries.value = []
            await loadLogFiles()
        } catch (err: any) {
            error('删除日志文件失败: ' + (err.response?.data?.error || err.message))
        }
    }
}

// 过滤日志条目
const filteredLogEntries = computed(() => {
    return logEntries.value.filter(entry =>
        selectedLevels.value.includes(entry.level.toLowerCase())
    )
})

// 格式化日志级别样式
function getLevelClass(level: string): string {
    const levelMap: Record<string, string> = {
        error: 'badge-error',
        warn: 'badge-warning',
        info: 'badge-info',
        debug: 'badge-neutral'
    }
    return `badge ${levelMap[level.toLowerCase()] || 'badge-ghost'}`
}

// 监听 logFiles 变化，自动选择第一个文件
watch(logFiles, (newFiles) => {
    if (newFiles.length > 0 && !selectedFile.value) {
        selectedFile.value = newFiles[0].path
        loadLogEntries()
    }
})

onMounted(() => {
    loadLogFiles()
    loadStats()
})
</script>

<template>
    <div>
        <!-- 操作栏 -->
        <div class="flex justify-between items-center mb-4">
            <div class="flex gap-2 items-center">
                <input type="text" v-model="props.currentDir" placeholder="日志目录"
                    class="input input-bordered input-sm" />
                <IconButton icon="ri-refresh-line" size="sm" @click="loadLogFiles">刷新</IconButton>
            </div>
            <div class="flex gap-2">
                <div class="join">
                    <input type="text" v-model="searchKeyword" placeholder="搜索日志"
                        class="input input-bordered input-sm join-item" @keyup.enter="handleSearch" />
                    <IconButton icon="ri-search-line" size="sm" class="join-item" @click="handleSearch">搜索</IconButton>
                </div>
                <div class="join">
                    <label class="join-item btn btn-sm" :class="{ 'btn-error': selectedLevels.includes('error') }">
                        <input type="checkbox" class="hidden" v-model="selectedLevels" value="error" />
                        <i class="ri-error-warning-fill mr-1"></i>
                        错误
                    </label>
                    <label class="join-item btn btn-sm" :class="{ 'btn-warning': selectedLevels.includes('warn') }">
                        <input type="checkbox" class="hidden" v-model="selectedLevels" value="warn" />
                        <i class="ri-alert-fill mr-1"></i>
                        警告
                    </label>
                    <label class="join-item btn btn-sm" :class="{ 'btn-info': selectedLevels.includes('info') }">
                        <input type="checkbox" class="hidden" v-model="selectedLevels" value="info" />
                        <i class="ri-information-fill mr-1"></i>
                        信息
                    </label>
                    <label class="join-item btn btn-sm" :class="{ 'btn-neutral': selectedLevels.includes('debug') }">
                        <input type="checkbox" class="hidden" v-model="selectedLevels" value="debug" />
                        <i class="ri-bug-fill mr-1"></i>
                        调试
                    </label>
                </div>
                <IconButton icon="ri-delete-bin-line" variant="error" size="sm" @click="handleCleanLogs">
                    清理旧日志
                </IconButton>
            </div>
        </div>

        <!-- 日志统计 -->
        <div v-if="logStats" class="stats shadow mb-4">
            <div class="stat">
                <div class="stat-title">错误</div>
                <div class="stat-value text-error">{{ logStats.error }}</div>
                <div class="stat-desc"><i class="ri-error-warning-fill"></i></div>
            </div>
            <div class="stat">
                <div class="stat-title">警告</div>
                <div class="stat-value text-warning">{{ logStats.warn }}</div>
                <div class="stat-desc"><i class="ri-alert-fill"></i></div>
            </div>
            <div class="stat">
                <div class="stat-title">信息</div>
                <div class="stat-value text-info">{{ logStats.info }}</div>
                <div class="stat-desc"><i class="ri-information-fill"></i></div>
            </div>
            <div class="stat">
                <div class="stat-title">调试</div>
                <div class="stat-value text-neutral">{{ logStats.debug }}</div>
                <div class="stat-desc"><i class="ri-bug-fill"></i></div>
            </div>
        </div>

        <!-- 日志内容 -->
        <div class="grid grid-cols-12 gap-4">
            <!-- 日志文件列表 -->
            <div class="col-span-3">
                <div class="card bg-base-200">
                    <div class="card-body p-2">
                        <div class="flex items-center justify-between mb-2 px-2">
                            <span class="text-sm font-medium">日志文件列表</span>
                            <IconButton v-if="selectedFile" icon="ri-delete-bin-line" variant="error" size="xs"
                                @click="handleDeleteCurrentLog" title="删除当前日志文件" />
                        </div>
                        <ul class="menu bg-base-200 w-full">
                            <li v-for="file in logFiles" :key="file.path">
                                <a class="flex items-center gap-2 transition-colors duration-200 hover:bg-base-300"
                                    :class="{
                                        'bg-primary/10 text-primary border-l-4 border-primary': selectedFile === file.path,
                                        'border-l-4 border-transparent': selectedFile !== file.path
                                    }" @click="selectedFile = file.path; loadLogEntries()">
                                    <i class="ri-file-text-line" />
                                    <span class="truncate">{{ file.path }}</span>
                                </a>
                            </li>
                            <li v-if="logFiles.length === 0" class="p-4 text-center text-base-content/50">
                                暂无日志文件
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
                                    <tr v-for="entry in filteredLogEntries" :key="entry.time">
                                        <td class="whitespace-nowrap">{{ entry.time }}</td>
                                        <td>
                                            <span :class="getLevelClass(entry.level)">
                                                {{ entry.level }}
                                            </span>
                                        </td>
                                        <td class="text-xs">{{ entry.caller }}</td>
                                        <td class="whitespace-pre-wrap">{{ entry.message }}</td>
                                    </tr>
                                    <tr v-if="filteredLogEntries.length === 0">
                                        <td colspan="4" class="text-center text-base-content/50 py-4">
                                            暂无符合条件的日志
                                        </td>
                                    </tr>
                                </tbody>
                            </table>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>