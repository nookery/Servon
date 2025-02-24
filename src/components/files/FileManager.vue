<script setup lang="ts">
import { ref, onMounted, computed, onBeforeUnmount } from 'vue'
import FileEditor from './FileEditor.vue'
import PageContainer from '../../layouts/PageContainer.vue'
import { fileAPI } from '../../api/file_api'
import RenameFileDialog from './RenameFileDialog.vue'
import type { FileInfo, SortBy, SortOrder } from '../../types/FileInfo'
import IconButton from '../IconButton.vue'
import { useConfirm } from '../../composables/useConfirm'
import CreateFileDialog from './CreateFileDialog.vue'

const props = defineProps<{
    initialPath: string
    showToolbar?: boolean
    showBreadcrumbs?: boolean
    showPagination?: boolean
    readOnly?: boolean
    showShortcuts?: boolean
    sortBy?: SortBy
    sortOrder?: SortOrder
}>()

const emit = defineEmits<{
    'update:path': [path: string]
    'update:sort-by': [sortBy: SortBy]
    'update:sort-order': [order: SortOrder]
}>()

const files = ref<FileInfo[]>([])
const currentPath = ref(props.initialPath)
const error = ref<string | null>(null)
const itemsPerPage = 10
const currentPage = ref(1)

// 状态变量
const showEditor = ref(false)
const editingFile = ref<FileInfo | null>(null)
const showCreateDialog = ref(false)
const searchQuery = ref('')
const showRenameDialog = ref(false)
const renamingFile = ref<FileInfo | null>(null)
const currentSortBy = ref<SortBy>(props.sortBy || 'name')
const currentSortOrder = ref<SortOrder>(props.sortOrder || 'asc')

// 添加选择相关的状态
const selectedFiles = ref<Set<string>>(new Set())

// 添加自动刷新相关的状态
const autoRefresh = ref(false)
const refreshInterval = ref(5) // 默认5秒
let refreshTimer: ReturnType<typeof setInterval> | null = null

// Pagination
const totalPages = computed(() => Math.ceil(files.value.length / itemsPerPage))
const paginatedFiles = computed(() => {
    const start = (currentPage.value - 1) * itemsPerPage
    return files.value.slice(start, start + itemsPerPage)
})

const confirm = useConfirm()

async function loadFiles(path: string) {
    try {
        const res = await fileAPI.getFiles(path, currentSortBy.value, currentSortOrder.value)
        files.value = res.data
        currentPath.value = path
        emit('update:path', path)
        error.value = null
        currentPage.value = 1
    } catch (err: any) {
        error.value = `获取文件列表失败: ${err.response?.data?.error || err.message || '未知错误'}`
    }
}

async function downloadFile(file: FileInfo) {
    try {
        const response = await fileAPI.downloadFile(file.path)
        const url = window.URL.createObjectURL(new Blob([response.data]))
        const link = document.createElement('a')
        link.href = url
        link.download = file.name
        document.body.appendChild(link)
        link.click()
        document.body.removeChild(link)
    } catch (err) {
        error.value = '下载文件失败'
    }
}

async function searchFiles() {
    if (!searchQuery.value) {
        loadFiles(currentPath.value)
        return
    }
    try {
        const res = await fileAPI.searchFiles(currentPath.value, searchQuery.value)
        files.value = res.data
        error.value = null
        currentPage.value = 1
    } catch (err) {
        error.value = '搜索文件失败'
    }
}

function navigateTo(index: number) {
    const parts = currentPath.value.split('/').filter(Boolean)
    const newPath = '/' + parts.slice(0, index + 1).join('/')
    loadFiles(newPath)
}

// 其他辅助函数保持不变
function getFileIcon(file: FileInfo) {
    if (file.isDir) return 'ri-folder-fill'
    const ext = file.name.split('.').pop()?.toLowerCase()
    switch (ext) {
        case 'txt': return 'ri-file-text-fill'
        case 'pdf': return 'ri-file-pdf-fill'
        case 'doc':
        case 'docx': return 'ri-file-word-fill'
        case 'xls':
        case 'xlsx': return 'ri-file-excel-fill'
        case 'jpg':
        case 'jpeg':
        case 'png': return 'ri-image-fill'
        default: return 'ri-file-fill'
    }
}

function formatFileSize(size: number) {
    const units = ['B', 'KB', 'MB', 'GB', 'TB']
    let index = 0
    while (size >= 1024 && index < units.length - 1) {
        size /= 1024
        index++
    }
    return `${size.toFixed(1)} ${units[index]}`
}

async function openFile(file: FileInfo) {
    if (file.isDir) {
        loadFiles(file.path)
    } else {
        editingFile.value = file
        showEditor.value = true
    }
}

async function renameFile(file: FileInfo) {
    renamingFile.value = file
    showRenameDialog.value = true
}

async function handleRename(oldPath: string, newPath: string) {
    try {
        await fileAPI.renameFile(oldPath, newPath)
        showRenameDialog.value = false
        renamingFile.value = null
        loadFiles(currentPath.value)
    } catch (err: any) {
        const errorMessage = err.response?.data?.error || err.message || '重命名失败'
        error.value = errorMessage
        setTimeout(() => {
            error.value = null
        }, 5000)
    }
}

// 处理排序点击
function handleSortClick(field: SortBy) {
    if (currentSortBy.value === field) {
        // 如果点击相同字段，切换排序顺序
        currentSortOrder.value = currentSortOrder.value === 'asc' ? 'desc' : 'asc'
    } else {
        // 如果点击不同字段，设置新字段并默认升序
        currentSortBy.value = field
        currentSortOrder.value = 'asc'
    }
    emit('update:sort-by', currentSortBy.value)
    emit('update:sort-order', currentSortOrder.value)
    loadFiles(currentPath.value)
}

// 获取排序图标
function getSortIcon(field: SortBy) {
    if (currentSortBy.value !== field) return 'ri-arrow-up-down-line'
    return currentSortOrder.value === 'asc' ? 'ri-sort-asc' : 'ri-sort-desc'
}

// 删除处理函数
async function handleDelete(file: FileInfo) {
    if (await confirm.error('删除文件', `确定要删除 ${file.name} 吗？`, {
        confirmText: '删除'
    })) {
        try {
            await fileAPI.deleteFile(file.path)
            loadFilesWithClear(currentPath.value)
        } catch (err: any) {
            error.value = err.response?.data?.error || err.message || '删除失败'
        }
    }
}

// 批量删除处理函数
async function handleBatchDelete() {
    const selectedFilesList = paginatedFiles.value.filter(f =>
        selectedFiles.value.has(f.path)
    )

    if (await confirm.error('批量删除',
        `确定要删除选中的 ${selectedFilesList.length} 个文件吗？`, {
        confirmText: '删除'
    })) {
        try {
            await fileAPI.batchDeleteFiles(selectedFilesList.map(f => f.path))
            selectedFiles.value.clear()
            loadFilesWithClear(currentPath.value)
        } catch (err: any) {
            error.value = err.response?.data?.error || err.message || '删除失败'
        }
    }
}

// 选择处理函数
function toggleSelect(file: FileInfo) {
    if (selectedFiles.value.has(file.path)) {
        selectedFiles.value.delete(file.path)
    } else {
        selectedFiles.value.add(file.path)
    }
}

// 全选/取消全选
function toggleSelectAll() {
    if (selectedFiles.value.size === paginatedFiles.value.length) {
        selectedFiles.value.clear()
    } else {
        selectedFiles.value = new Set(paginatedFiles.value.map(f => f.path))
    }
}

// 替换原来的 loadFiles 重新赋值代码
async function loadFilesWithClear(path: string) {
    selectedFiles.value.clear()
    await loadFiles(path)
}

// 处理自动刷新
function toggleAutoRefresh() {
    autoRefresh.value = !autoRefresh.value
    if (autoRefresh.value) {
        startAutoRefresh()
    } else {
        stopAutoRefresh()
    }
}

function startAutoRefresh() {
    if (refreshTimer) clearInterval(refreshTimer)
    refreshTimer = setInterval(() => {
        loadFiles(currentPath.value)
    }, refreshInterval.value * 1000)
}

function stopAutoRefresh() {
    if (refreshTimer) {
        clearInterval(refreshTimer)
        refreshTimer = null
    }
}

// 在组件卸载时清理定时器
onBeforeUnmount(() => {
    stopAutoRefresh()
})

onMounted(() => {
    loadFilesWithClear(props.initialPath)
})

async function copyFile(file: FileInfo) {
    try {
        const baseName = file.name.split('.')
        const ext = baseName.length > 1 ? baseName.pop() : ''
        const newName = `${baseName.join('.')} - 副本${ext ? '.' + ext : ''}`
        const newPath = `${currentPath.value}/${newName}`

        await fileAPI.copyFile(file.path, newPath)
        loadFiles(currentPath.value)
    } catch (err: any) {
        error.value = err.response?.data?.error || err.message || '复制失败'
        setTimeout(() => {
            error.value = null
        }, 5000)
    }
}
</script>

<template>
    <PageContainer title="文件管理">
        <template #header>
            <div v-if="error" class="alert alert-error shadow-lg mb-4">
                <div>
                    <i class="ri-error-warning-line"></i>
                    <span>{{ error }}</span>
                </div>
            </div>

            <div class="bg-base-200 border border-base-300 rounded-lg p-1 mb-4">
                <div class="flex items-center gap-2">
                    <div class="breadcrumbs text-sm">
                        <ul>
                            <li>
                                <a @click="loadFiles('/')" class="hover:bg-base-300 px-2 py-0.5 rounded no-underline">
                                    <i class="ri-home-line mr-1"></i>根目录
                                </a>
                            </li>
                            <li v-for="(part, index) in currentPath.split('/').filter(Boolean)" :key="index">
                                <a @click="navigateTo(index)"
                                    class="hover:bg-base-300 px-2 py-0.5 rounded no-underline">
                                    {{ part }}
                                </a>
                            </li>
                        </ul>
                    </div>
                </div>
            </div>

            <div v-if="showToolbar" class="flex flex-wrap gap-2 mb-4">
                <div class="flex gap-2">
                    <IconButton v-if="!readOnly" icon="ri-add-line" @click="showCreateDialog = true">
                        新建
                    </IconButton>
                    <IconButton icon="ri-refresh-line" @click="loadFiles(currentPath)">
                        刷新
                    </IconButton>
                    <IconButton v-if="!readOnly && selectedFiles.size > 0" icon="ri-delete-bin-2-line" variant="error"
                        @click="handleBatchDelete">
                        删除选中 ({{ selectedFiles.size }})
                    </IconButton>
                    <IconButton :icon="autoRefresh ? 'ri-time-fill' : 'ri-time-line'"
                        :variant="autoRefresh ? 'primary' : 'default'" @click="toggleAutoRefresh"
                        :title="`自动刷新 (${refreshInterval}秒)`">
                        {{ autoRefresh ? '停止刷新' : '自动刷新' }}
                    </IconButton>
                    <div v-if="autoRefresh" class="flex items-center gap-2">
                        <select v-model="refreshInterval" class="select select-bordered select-sm"
                            @change="startAutoRefresh">
                            <option value="3">3秒</option>
                            <option value="5">5秒</option>
                            <option value="10">10秒</option>
                            <option value="30">30秒</option>
                        </select>
                    </div>
                </div>

                <template v-if="showShortcuts">
                    <div class="divider divider-horizontal"></div>
                    <div class="flex flex-wrap gap-2 items-center">
                        <IconButton icon="ri-settings-3-line" variant="ghost" size="sm" title="数据目录"
                            @click="loadFiles('/data')">
                            data
                        </IconButton>
                        <IconButton icon="ri-settings-3-line" variant="ghost" size="sm" title="系统配置目录"
                            @click="loadFiles('/etc')">
                            etc
                        </IconButton>
                        <IconButton icon="ri-home-4-line" variant="ghost" size="sm" title="用户目录"
                            @click="loadFiles('/home')">
                            home
                        </IconButton>
                        <IconButton icon="ri-file-list-3-line" variant="ghost" size="sm" title="系统日志目录"
                            @click="loadFiles('/var/log')">
                            logs
                        </IconButton>
                        <IconButton icon="ri-apps-2-line" variant="ghost" size="sm" title="本地安装的软件"
                            @click="loadFiles('/usr/local')">
                            local
                        </IconButton>
                        <IconButton icon="ri-time-line" variant="ghost" size="sm" title="临时文件目录"
                            @click="loadFiles('/tmp')">
                            tmp
                        </IconButton>
                    </div>
                </template>

                <div class="flex-1 flex justify-end items-center">
                    <div class="form-control">
                        <div class="input-group">
                            <input type="text" placeholder="搜索文件..." class="input input-bordered input-sm w-64"
                                v-model="searchQuery" @keyup.enter="searchFiles">
                            <button class="btn btn-square btn-sm" @click="searchFiles">
                                <i class="ri-search-line"></i>
                            </button>
                        </div>
                    </div>
                </div>
            </div>
        </template>

        <div class="overflow-x-auto">
            <table class="table table-zebra w-full">
                <thead>
                    <tr>
                        <th v-if="!readOnly" class="w-10">
                            <input type="checkbox" class="checkbox checkbox-sm"
                                :checked="selectedFiles.size === paginatedFiles.length"
                                :indeterminate="selectedFiles.size > 0 && selectedFiles.size < paginatedFiles.length"
                                @change="toggleSelectAll" />
                        </th>
                        <th @click="handleSortClick('name')" class="cursor-pointer hover:bg-base-200">
                            <div class="flex items-center gap-2">
                                名称
                                <i :class="getSortIcon('name')"></i>
                            </div>
                        </th>
                        <th>权限</th>
                        <th>所有者</th>
                        <th @click="handleSortClick('size')" class="cursor-pointer hover:bg-base-200">
                            <div class="flex items-center gap-2">
                                大小
                                <i :class="getSortIcon('size')"></i>
                            </div>
                        </th>
                        <th @click="handleSortClick('modTime')" class="cursor-pointer hover:bg-base-200">
                            <div class="flex items-center gap-2">
                                修改时间
                                <i :class="getSortIcon('modTime')"></i>
                            </div>
                        </th>
                        <th v-if="!readOnly">操作</th>
                    </tr>
                </thead>
                <tbody>
                    <tr v-for="file in paginatedFiles" :key="file.path">
                        <td v-if="!readOnly">
                            <input type="checkbox" class="checkbox checkbox-sm" :checked="selectedFiles.has(file.path)"
                                @change="toggleSelect(file)" />
                        </td>
                        <td class="flex items-center gap-2">
                            <i :class="getFileIcon(file)" class="text-lg"></i>
                            <span class="cursor-pointer" @click="openFile(file)">
                                {{ file.name }}
                            </span>
                        </td>
                        <td>{{ file.mode }}</td>
                        <td>{{ file.owner }}:{{ file.group }}</td>
                        <td>{{ formatFileSize(file.size) }}</td>
                        <td>{{ new Date(file.modTime).toLocaleString() }}</td>
                        <td v-if="!readOnly" class="space-x-2">
                            <div class="join gap-0">
                                <button class="btn btn-xs btn-error join-item" @click="handleDelete(file)">
                                    <i class="ri-delete-bin-line"></i>
                                </button>
                                <button class="btn btn-xs btn-warning join-item" @click="renameFile(file)">
                                    <i class="ri-edit-line"></i>
                                </button>
                                <button v-if="!file.isDir" class="btn btn-xs join-item" @click="downloadFile(file)">
                                    <i class="ri-download-line"></i>
                                </button>
                                <button v-if="!file.isDir" class="btn btn-xs join-item" @click="copyFile(file)">
                                    <i class="ri-file-copy-line"></i>
                                </button>
                            </div>
                        </td>
                    </tr>
                </tbody>
            </table>
        </div>

        <div class="flex justify-center mt-4 space-x-2">
            <button class="btn btn-sm" :disabled="currentPage === 1" @click="currentPage--">
                上一页
            </button>
            <span class="flex items-center">{{ currentPage }} / {{ totalPages }}</span>
            <button class="btn btn-sm" :disabled="currentPage === totalPages" @click="currentPage++">
                下一页
            </button>
        </div>

        <FileEditor v-model:show="showEditor" :file="editingFile" @saved="loadFiles(currentPath)" />

        <CreateFileDialog v-model:show="showCreateDialog" :current-path="currentPath"
            @created="loadFiles(currentPath)" />

        <RenameFileDialog v-model:show="showRenameDialog" :file="renamingFile" @rename="handleRename" />
    </PageContainer>
</template>

<style>
.breadcrumbs a {
    text-decoration: none !important;
}

.breadcrumbs a:hover {
    text-decoration: none !important;
}
</style>