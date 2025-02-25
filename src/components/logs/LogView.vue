<script setup lang="ts">
import { ref, onMounted, watch, computed, onUnmounted } from 'vue'
import { useConfirm } from '../../composables/useConfirm'
import { useToast } from '../../composables/useToast'
import { useError } from '../../composables/useError'
import type { LogEntry, LogFile, LogStats } from '../../types/log'
import * as logApi from '../../api/logs_api'
import TableLogView from './TableLogView.vue'
import TerminalLogView from './TerminalLogView.vue'
import ViewModeSelector from './ViewModeSelector.vue'
import FieldSelector from './FieldSelector.vue'
import LevelSelector from './LevelSelector.vue'
import LogFileList from './LogFileList.vue'
import LogToolbar from './LogToolbar.vue'

const props = defineProps<{
    currentDir: string
}>()

const emit = defineEmits<{
    'update:currentDir': [value: string]
}>()

const confirm = useConfirm()
const toast = useToast()
const { error } = useError()

const logFiles = ref<LogFile[]>([])
const selectedFile = ref<string>('')
const logEntries = ref<LogEntry[]>([])
const logStats = ref<LogStats | null>(null)
const loading = ref(false)
const selectedLevels = ref<string[]>(['error', 'warn', 'info', 'debug'])
const viewMode = ref<'table' | 'terminal'>('table')
const visibleFields = ref<string[]>(['time', 'level', 'caller', 'message'])

// 添加一个状态来控制左侧面板的显示
const showSidebar = ref(true)

// 切换左侧面板显示
function toggleSidebar() {
    showSidebar.value = !showSidebar.value
}

// 添加一个计算属性来检测容器宽度
const isNarrow = ref(false)

// 监听窗口大小变化
onMounted(() => {
    // 初始检查
    checkContainerWidth()

    // 添加resize事件监听
    window.addEventListener('resize', checkContainerWidth)
})

// 组件卸载时移除事件监听
onUnmounted(() => {
    window.removeEventListener('resize', checkContainerWidth)
})

// 检查容器宽度
function checkContainerWidth() {
    // 获取组件容器宽度
    const container = document.querySelector('.log-view-container')
    if (container) {
        isNarrow.value = container.clientWidth < 768
    }
}

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

// 搜索日志
async function handleSearch(keyword: string) {
    if (!keyword) return

    try {
        loading.value = true
        logEntries.value = await logApi.searchLogs(props.currentDir, keyword)
    } catch (err: any) {
        error('搜索日志失败: ' + (err.response?.data?.error || err.message || '未知错误'))
    } finally {
        loading.value = false
    }
}

// 处理文件选择
function handleFileSelect(filePath: string) {
    selectedFile.value = filePath
    loadLogEntries()
}

// 更新目录
function updateCurrentDir(dir: string) {
    emit('update:currentDir', dir)
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
</script>

<template>
    <!-- 使用 h-full 和 overflow-hidden 确保组件占满所有可用空间 -->
    <div class="h-full flex flex-col overflow-hidden log-view-container">
        <!-- 固定的操作栏和统计信息 -->
        <div class="flex-none space-y-4">
            <!-- 操作栏 -->
            <LogToolbar :currentDir="props.currentDir" :selectedFile="selectedFile"
                @update:currentDir="updateCurrentDir" @refresh="handleRefresh" @search="handleSearch"
                @delete-log="handleDeleteCurrentLog" @clean-logs="handleCleanLogs" @clear-log="handleClearCurrentLog" />

            <!-- 日志统计 -->
            <div v-if="logStats" class="stats shadow w-full overflow-x-auto">
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

        <!-- 侧边栏切换按钮 -->
        <div class="flex justify-end mb-2">
            <button @click="toggleSidebar" class="btn btn-sm btn-ghost">
                <i :class="showSidebar ? (isNarrow ? 'ri-arrow-up-s-line' : 'ri-arrow-left-s-line') : (isNarrow ? 'ri-arrow-down-s-line' : 'ri-arrow-right-s-line')"
                    class="mr-1"></i>
                {{ showSidebar ? '收起控制面板' : '展开控制面板' }}
            </button>
        </div>

        <!-- 可滚动的内容区域 - 根据宽度切换布局 -->
        <div class="flex-1 min-h-0" :class="isNarrow ? 'flex flex-col gap-4' : 'flex gap-4'">
            <!-- 左侧区域：视图切换、日志级别筛选和文件列表 -->
            <div v-show="showSidebar" class="transition-all duration-300 overflow-hidden" :class="{
                'flex-none': true,
                'h-0': !showSidebar,
                'w-0': !showSidebar && !isNarrow,
                'w-80': showSidebar && !isNarrow
            }">
                <div :class="isNarrow ? 'flex flex-col lg:flex-row gap-4' : 'flex flex-col gap-4'">
                    <!-- 控制面板 -->
                    <div :class="isNarrow ? 'space-y-4 lg:w-1/3' : 'space-y-4'">
                        <!-- 视图切换卡片 -->
                        <ViewModeSelector v-model="viewMode" />

                        <!-- 字段显示控制卡片 -->
                        <FieldSelector v-model="visibleFields" />

                        <!-- 日志级别筛选卡片 -->
                        <LevelSelector v-model="selectedLevels" />
                    </div>

                    <!-- 日志文件列表卡片 -->
                    <div :class="isNarrow ? 'lg:w-2/3 h-64 lg:h-auto' : 'h-80'">
                        <LogFileList :logFiles="logFiles" :selectedFile="selectedFile"
                            @update:selectedFile="handleFileSelect" @refresh="loadLogFiles(true)" />
                    </div>
                </div>
            </div>

            <!-- 日志内容 - 使用 overflow-hidden 和 flex 布局 -->
            <div class="flex-1 card bg-base-200 overflow-hidden">
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