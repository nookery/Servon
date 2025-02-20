<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import FileEditor from './FileEditor.vue'
import PageContainer from '../../layouts/PageContainer.vue'
import { fileAPI } from '../../api/file_api'
import RenameFileDialog from './RenameFileDialog.vue'
import type { FileInfo, SortBy, SortOrder } from '../../models/FileInfo'
import IconButton from '../IconButton.vue'

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
const newFileName = ref('')
const newFileType = ref<'file' | 'directory'>('file')
const searchQuery = ref('')
const showRenameDialog = ref(false)
const renamingFile = ref<FileInfo | null>(null)
const currentSortBy = ref<SortBy>(props.sortBy || 'name')
const currentSortOrder = ref<SortOrder>(props.sortOrder || 'asc')

// Pagination
const totalPages = computed(() => Math.ceil(files.value.length / itemsPerPage))
const paginatedFiles = computed(() => {
    const start = (currentPage.value - 1) * itemsPerPage
    return files.value.slice(start, start + itemsPerPage)
})

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

async function deleteFile(file: FileInfo) {
    if (!confirm(`确定要删除 ${file.name} 吗？`)) return
    try {
        await fileAPI.deleteFile(file.path)
        loadFiles(currentPath.value)
    } catch (err) {
        error.value = '删除文件失败'
    }
}

async function createFile() {
    if (!newFileName.value) return
    try {
        const newPath = `${currentPath.value}/${newFileName.value}`.replace(/\/+/g, '/')
        await fileAPI.createFile(newPath, newFileType.value)
        showCreateDialog.value = false
        newFileName.value = ''
        loadFiles(currentPath.value)
    } catch (err: any) {
        const errorMessage = err.response?.data?.error || err.message || '创建文件失败'
        error.value = errorMessage
        setTimeout(() => {
            error.value = null
        }, 5000)
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

onMounted(() => {
    loadFiles(props.initialPath)
})
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
                </div>

                <template v-if="showShortcuts">
                    <div class="divider divider-horizontal"></div>
                    <div class="flex flex-wrap gap-2">
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

                <div class="flex-1 flex justify-end">
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
                            <IconButton v-if="!file.isDir" icon="ri-download-line" size="xs"
                                @click="downloadFile(file)" />
                            <IconButton icon="ri-edit-line" size="xs" variant="warning" @click="renameFile(file)" />
                            <IconButton icon="ri-delete-bin-line" size="xs" variant="error" @click="deleteFile(file)" />
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

        <dialog class="modal" :class="{ 'modal-open': showCreateDialog }">
            <div class="modal-box">
                <h3 class="font-bold text-lg mb-4">新建文件/目录</h3>

                <div v-if="error" class="alert alert-error shadow-lg mb-4">
                    <div>
                        <i class="ri-error-warning-line"></i>
                        <span>{{ error }}</span>
                    </div>
                </div>

                <div class="form-control">
                    <label class="label">
                        <span class="label-text">类型</span>
                    </label>
                    <select v-model="newFileType" class="select select-bordered">
                        <option value="file">文件</option>
                        <option value="directory">目录</option>
                    </select>
                </div>
                <div class="form-control mt-4">
                    <label class="label">
                        <span class="label-text">名称</span>
                    </label>
                    <input type="text" v-model="newFileName" class="input input-bordered" placeholder="输入名称">
                </div>
                <div class="modal-action">
                    <button class="btn" @click="() => {
                        showCreateDialog = false;
                        if (error) error = null;
                        newFileName = '';
                    }">取消</button>
                    <button class="btn btn-primary" @click="createFile" :disabled="!newFileName">创建</button>
                </div>
            </div>
        </dialog>

        <RenameFileDialog v-model:show="showRenameDialog" :file="renamingFile" @rename="handleRename" />
    </PageContainer>
</template>

<style>
@import 'remixicon/fonts/remixicon.css';

.breadcrumbs a {
    text-decoration: none !important;
}

.breadcrumbs a:hover {
    text-decoration: none !important;
}
</style>