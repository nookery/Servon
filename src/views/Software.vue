<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useToast } from '../composables/useToast'
import PageContainer from '../layouts/PageContainer.vue'
import IconButton from '../components/IconButton.vue'
import LogView from '../components/logs/LogView.vue'
import FileManager from '../components/files/FileManager.vue'
import { systemAPI, type Software } from '../api/info'
import { RiAppsLine, RiListSettingsLine, RiFileListLine, RiFolderOpenLine } from '@remixicon/vue'

const activeTab = ref('software')
const softwares = ref<Software[]>([])
const loading = ref(false)
const currentLogs = ref<string[]>([])
const installing = ref(false)
const currentSoftware = ref<string>('')
const operationFailed = ref(false)
const error = ref<string | null>(null)
const isRefreshing = ref(false)

const toast = useToast()

// 定义标签页配置
const tabs = [
    {
        key: 'software',
        title: '软件列表',
        icon: RiListSettingsLine
    },
    {
        key: 'files',
        title: '文件浏览',
        icon: RiFolderOpenLine
    },
    {
        key: 'logs',
        title: '日志查看',
        icon: RiFileListLine
    }
]

// 刷新软件列表
async function refreshSoftwareList() {
    if (isRefreshing.value) return
    isRefreshing.value = true
    try {
        await loadSoftwareList()
        toast.success('软件列表已刷新')
    } catch (err) {
        // 错误已在 loadSoftwareList 中处理
    } finally {
        isRefreshing.value = false
    }
}

async function handleAction(software: Software) {
    currentSoftware.value = software.name
    currentLogs.value = []
    installing.value = true
    operationFailed.value = false
    error.value = null

    try {
        const action = software.status === 'not_installed' ? 'install' : 'uninstall'
        await systemAPI.manageSoftware(software.name, action)
        installing.value = false
        toast.success(`${action === 'install' ? '安装任务' : '卸载任务'}已启动`)
    } catch (err) {
        console.error('操作失败:', err)
        installing.value = false
        operationFailed.value = true
        error.value = '操作失败'
    }
}

async function handleStop(name: string) {
    try {
        await systemAPI.stopSoftware(name)
        toast.success('服务已停止')
        error.value = null
    } catch (error: any) {
        const errorMessage = error.response?.data?.error || error.message
        error.value = `停止服务失败: ${errorMessage}`
    }
}

async function handleStart(name: string) {
    try {
        await systemAPI.startSoftware(name)
        toast.success('服务已启动')
        error.value = null
    } catch (err: any) {
        const errorMessage = err.response?.data?.error || err.message
        error.value = `启动服务失败: ${errorMessage}`
    }
}

async function handleStatus(name: string) {
    try {
        const response = await systemAPI.getSoftwareStatus(name)
        const software = softwares.value.find(s => s.name === name)
        if (software) {
            software.status = response.data.status
        }
    } catch (error) {
        console.error('获取软件状态失败:', error)
    }
}

async function loadSoftwareList() {
    loading.value = true
    try {
        const response = await systemAPI.getSoftwareList()
        softwares.value = response.data.map((item: string) => ({
            name: item,
            status: 'not_installed'
        }))
        error.value = null

        for (const software of softwares.value) {
            await handleStatus(software.name)
        }
    } catch (err) {
        error.value = '获取软件列表失败'
    } finally {
        loading.value = false
    }
}

onMounted(() => {
    loadSoftwareList()
})
</script>

<template>
    <PageContainer title="软件管理" :error="error" :empty="!loading && softwares.length === 0" emptyText="还没有可用的软件包"
        emptyDescription="系统正在等待您添加第一个软件包，让我们开始吧！" :emptyIcon="RiAppsLine" :tabs="tabs" v-model="activeTab">

        <!-- 软件列表标签页 -->
        <template #software>
            <div class="flex justify-between items-center mb-4">
                <h2 class="text-xl font-semibold">可用软件</h2>
                <IconButton icon="ri-refresh-line" variant="primary" size="sm" :loading="isRefreshing"
                    @click="refreshSoftwareList">
                    刷新列表
                </IconButton>
            </div>

            <div v-if="softwares.length > 0" class="overflow-x-auto">
                <table class="table table-zebra w-full">
                    <thead>
                        <tr>
                            <th>软件名称</th>
                            <th>状态</th>
                            <th>操作</th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr v-for="item in softwares" :key="item.name">
                            <td>{{ item.name }}</td>
                            <td :class="{
                                'text-info': item.status === 'not_installed',
                                'text-warning': item.status === 'stopped',
                                'text-success': item.status === 'running'
                            }">{{ item.status }}</td>
                            <td>
                                <div class="flex gap-2">
                                    <IconButton v-if="item.status !== 'running'" size="sm" tooltip-position="none"
                                        :icon="item.status === 'not_installed' ? 'ri-download-line' : 'ri-delete-bin-line'"
                                        :variant="item.status === 'not_installed' ? 'primary' : 'error'"
                                        :disabled="installing && currentSoftware === item.name"
                                        @click="handleAction(item)">
                                        {{ item.status === 'not_installed' ? '安装' : '卸载' }}
                                    </IconButton>
                                    <IconButton v-if="item.status === 'stopped'" tooltip-position="none"
                                        icon="ri-play-line" variant="primary" size="sm" @click="handleStart(item.name)">
                                        启动
                                    </IconButton>
                                    <IconButton v-if="item.status === 'running'" size="sm" tooltip-position="none"
                                        icon="ri-stop-line" variant="secondary" @click="handleStop(item.name)">
                                        停止
                                    </IconButton>
                                </div>
                            </td>
                        </tr>
                    </tbody>
                </table>
            </div>
        </template>

        <!-- 文件浏览标签页 -->
        <template #files>
            <FileManager :initial-path="'/data/softwares'" :show-toolbar="true" :show-breadcrumbs="true"
                :show-pagination="true" :read-only="false" :show-shortcuts="false" />
        </template>

        <!-- 日志查看标签页 -->
        <template #logs>
            <LogView current-dir="" />
        </template>

    </PageContainer>
</template>