<script setup lang="ts">
import { ref, onMounted, watch, computed } from 'vue'
import { useConfirm } from '../../composables/useConfirm'
import { useToast } from '../../composables/useToast'
import { useError } from '../../composables/useError'
import type { LogEntry, LogFile, LogStats } from '../../types/log'
import * as logApi from '../../api/logs_api'
import IconButton from '../IconButton.vue'
import TableLogView from './TableLogView.vue'
import TerminalLogView from './TerminalLogView.vue'
import ViewModeSelector from './ViewModeSelector.vue'
import FieldSelector from './FieldSelector.vue'
import LevelSelector from './LevelSelector.vue'

const props = defineProps<{
    currentDir: string
}>()

const confirm = useConfirm()
const toast = useToast()
const { error } = useError()

const logFiles = ref<LogFile[]>([])
const selectedFile = ref<string>('')
const logEntries = ref<LogEntry[]>([])
const logStats = ref<LogStats | null>(null)
const searchKeyword = ref('')
const loading = ref(false)
const selectedLevels = ref<string[]>(['error', 'warn', 'info', 'debug'])
const viewMode = ref<'table' | 'terminal'>('table')
const visibleFields = ref<string[]>(['time', 'level', 'caller', 'message'])

// 基础加载函数，不显示提示
async function loadLogFiles(showToast = false) {
    try {
        loading.value = true
        logFiles.value = await logApi.getLogFiles(props.currentDir)
        if (showToast) {
            toast.success('日志列表已刷新')
        }
    } catch (err: any) {
        error('获取日志文件列表失败: ' + (err.response?.data?.error || err.message || '未知错误'))
    } finally {
        loading.value = false
    }
}

// 基础加载函数，不显示提示
async function loadLogEntries(showToast = false) {
    if (!selectedFile.value) return

    try {
        loading.value = true
        logEntries.value = await logApi.getLogEntries(selectedFile.value)
        if (showToast) {
            toast.success('日志内容已刷新')
        }
    } catch (err: any) {
        error('获取日志内容失败: ' + (err.response?.data?.error || err.message || '未知错误'))
    } finally {
        loading.value = false
    }
}

// 手动刷新按钮点击事件，显示提示
async function handleRefresh() {
    await loadLogFiles(true)
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

// 清空当前日志文件
async function handleClearCurrentLog() {
    if (!selectedFile.value) {
        toast.warning('请先选择要清空的日志文件')
        return
    }

    if (await confirm.error('清空日志文件', `确定要清空日志文件 "${selectedFile.value}" 吗？此操作不可恢复！`, {
        confirmText: '清空',
        cancelText: '取消'
    })) {
        try {
            await logApi.clearLogFile(selectedFile.value)
            toast.success('日志文件已清空')
            await loadLogEntries(true)  // 重新加载日志内容
        } catch (err: any) {
            error('清空日志文件失败: ' + (err.response?.data?.error || err.message))
        }
    }
}

// 过滤日志条目
const filteredLogEntries = computed(() => {
    return logEntries.value.filter(entry =>
        selectedLevels.value.includes(entry.level.toLowerCase())
    )
})

// 监听文件列表变化，自动加载不显示提示
watch(logFiles, (newFiles) => {
    if (newFiles.length > 0 && !selectedFile.value) {
        selectedFile.value = newFiles[0].path
        loadLogEntries()  // 自动加载不显示提示
    }
})

onMounted(() => {
    loadLogFiles()  // 初始加载不显示提示
    loadStats()
})

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
</script>

<template>
    <!-- 使用 h-full 和 overflow-hidden 确保组件占满所有可用空间 -->
    <div class="h-full flex flex-col overflow-hidden">
        <!-- 固定的操作栏和统计信息 -->
        <div class="flex-none space-y-4">
            <!-- 操作栏 -->
            <div class="flex justify-between items-center gap-4">
                <div class="flex gap-2 items-center">
                    <input type="text" v-model="props.currentDir" placeholder="日志目录"
                        class="input input-bordered input-sm" />
                    <IconButton icon="ri-refresh-line" size="sm" @click="handleRefresh">刷新</IconButton>
                </div>
                <div class="flex gap-2">
                    <div class="join">
                        <input type="text" v-model="searchKeyword" placeholder="搜索日志"
                            class="input input-bordered input-sm join-item" @keyup.enter="handleSearch" />
                        <IconButton icon="ri-search-line" size="sm" class="join-item" @click="handleSearch">搜索
                        </IconButton>
                    </div>
                    <IconButton icon="ri-delete-bin-line" variant="error" size="sm" @click="handleDeleteCurrentLog"
                        :disabled="!selectedFile" title="删除当前日志文件">
                        删除日志
                    </IconButton>
                    <IconButton icon="ri-delete-bin-line" variant="error" size="sm" @click="handleCleanLogs"
                        title="清理30天前的日志">
                        清理旧日志
                    </IconButton>
                    <IconButton icon="ri-eraser-line" variant="error" size="sm" @click="handleClearCurrentLog"
                        :disabled="!selectedFile" title="清空当前日志内容">
                        清空日志
                    </IconButton>
                </div>
            </div>

            <!-- 日志统计 -->
            <div v-if="logStats" class="stats shadow w-full">
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
        </div>

        <!-- 可滚动的内容区域 - 使用 flex-1 和 min-h-0 确保正确的滚动行为 -->
        <div class="flex-1 min-h-0 grid grid-cols-12 gap-4 mt-4">
            <!-- 左侧区域：视图切换、日志级别筛选和文件列表 -->
            <div class="col-span-3 flex flex-col gap-4">
                <!-- 视图切换卡片 -->
                <ViewModeSelector v-model="viewMode" />

                <!-- 字段显示控制卡片 -->
                <FieldSelector v-model="visibleFields" />

                <!-- 日志级别筛选卡片 -->
                <LevelSelector v-model="selectedLevels" />

                <!-- 日志文件列表卡片 -->
                <div class="card bg-base-200 overflow-hidden flex-1">
                    <div class="h-full flex flex-col">
                        <!-- 文件列表头部 -->
                        <div class="flex-none p-3">
                            <span class="text-sm font-medium">日志文件列表</span>
                        </div>

                        <!-- 可滚动的文件列表 -->
                        <div class="flex-1 min-h-0 overflow-y-auto p-2">
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
            </div>

            <!-- 日志内容 - 使用 overflow-hidden 和 flex 布局 -->
            <div class="col-span-9 card bg-base-200 overflow-hidden">
                <div class="h-full flex flex-col">
                    <div v-if="loading" class="flex-1 flex justify-center items-center">
                        <span class="loading loading-spinner loading-lg"></span>
                    </div>
                    <div v-else class="flex-1 min-h-0 overflow-auto">
                        <!-- 根据视图模式切换组件 -->
                        <TableLogView v-if="viewMode === 'table'" :logEntries="filteredLogEntries"
                            :visibleFields="visibleFields" />
                        <TerminalLogView v-else :logEntries="filteredLogEntries" :visibleFields="visibleFields" />
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>