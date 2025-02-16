<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import FileEditor from '../components/FileEditor.vue'
import PageContainer from '../components/PageContainer.vue'
import { fileAPI, type FileInfo } from '../api/file'

const files = ref<FileInfo[]>([])
const currentPath = ref('/')
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

// Pagination
const totalPages = computed(() => Math.ceil(files.value.length / itemsPerPage))
const paginatedFiles = computed(() => {
    const start = (currentPage.value - 1) * itemsPerPage
    return files.value.slice(start, start + itemsPerPage)
})

async function loadFiles(path: string) {
    try {
        const res = await fileAPI.getFiles(path)
        files.value = res.data
        currentPath.value = path
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

onMounted(() => {
    loadFiles(currentPath.value)
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

            <div class="bg-base-200 border border-base-300 rounded-lg p-3 mb-4">
                <div class="flex items-center gap-2">
                    <i class="ri-folder-line text-lg"></i>
                    <div class="breadcrumbs text-sm">
                        <ul>
                            <li>
                                <a @click="loadFiles('/')" class="hover:bg-base-300 px-2 py-1 rounded">
                                    <i class="ri-home-line mr-1"></i>根目录
                                </a>
                            </li>
                            <li v-for="(part, index) in currentPath.split('/').filter(Boolean)" :key="index">
                                <a @click="navigateTo(index)" class="hover:bg-base-300 px-2 py-1 rounded">
                                    {{ part }}
                                </a>
                            </li>
                        </ul>
                    </div>
                </div>
            </div>

            <div class="flex flex-wrap gap-2 mb-4">
                <div class="flex gap-2">
                    <button class="btn btn-primary" @click="showCreateDialog = true">
                        <i class="ri-add-line mr-1"></i>新建
                    </button>
                    <button class="btn" @click="loadFiles(currentPath)">
                        <i class="ri-refresh-line mr-1"></i>刷新
                    </button>
                </div>

                <div class="divider divider-horizontal"></div>
                <div class="flex flex-wrap gap-2">
                    <button class="btn btn-sm btn-ghost" @click="loadFiles('/data')" title="数据目录">
                        <i class="ri-settings-3-line mr-1"></i>data
                    </button>
                    <button class="btn btn-sm btn-ghost" @click="loadFiles('/etc')" title="系统配置目录">
                        <i class="ri-settings-3-line mr-1"></i>etc
                    </button>
                    <button class="btn btn-sm btn-ghost" @click="loadFiles('/home')" title="用户目录">
                        <i class="ri-home-4-line mr-1"></i>home
                    </button>
                    <button class="btn btn-sm btn-ghost" @click="loadFiles('/var/log')" title="系统日志目录">
                        <i class="ri-file-list-3-line mr-1"></i>logs
                    </button>
                    <button class="btn btn-sm btn-ghost" @click="loadFiles('/usr/local')" title="本地安装的软件">
                        <i class="ri-apps-2-line mr-1"></i>local
                    </button>
                    <button class="btn btn-sm btn-ghost" @click="loadFiles('/tmp')" title="临时文件目录">
                        <i class="ri-time-line mr-1"></i>tmp
                    </button>
                </div>

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
                        <th>名称</th>
                        <th>权限</th>
                        <th>所有者</th>
                        <th>大小</th>
                        <th>修改时间</th>
                        <th>操作</th>
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
                        <td class="space-x-2">
                            <button class="btn btn-xs" @click="downloadFile(file)" v-if="!file.isDir">
                                <i class="ri-download-line"></i>
                            </button>
                            <button class="btn btn-xs btn-error" @click="deleteFile(file)">
                                <i class="ri-delete-bin-line"></i>
                            </button>
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
    </PageContainer>
</template>

<style>
@import 'remixicon/fonts/remixicon.css';
</style>