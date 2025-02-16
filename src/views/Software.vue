<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useToast } from '../composables/useToast'
import Alert from '../components/Alert.vue'
import { systemAPI, type Software } from '../api/info'

const softwares = ref<Software[]>([])
const loading = ref(false)
const currentLogs = ref<string[]>([])
const installing = ref(false)
const currentSoftware = ref<string>('')
const operationFailed = ref(false)
const error = ref<string | null>(null)

const toast = useToast()

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
    <div class="card bg-base-100 shadow-xl">
        <div class="card-body">
            <h2 class="card-title">软件管理</h2>

            <Alert v-if="error" type="error" :message="error" />


            <div class="overflow-x-auto">
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
                                    <button class="btn btn-sm"
                                        :class="item.status === 'not_installed' ? 'btn-primary' : 'btn-error'"
                                        :disabled="installing && currentSoftware === item.name"
                                        @click="handleAction(item)" v-if="item.status !== 'running'">
                                        {{ item.status === 'not_installed' ? '安装' : '卸载' }}
                                    </button>
                                    <button v-if="item.status === 'stopped'" class="btn btn-sm btn-primary"
                                        @click="handleStart(item.name)">
                                        启动
                                    </button>
                                    <button v-if="item.status === 'running'" class="btn btn-sm btn-secondary"
                                        @click="handleStop(item.name)">
                                        停止
                                    </button>
                                </div>
                            </td>
                        </tr>
                    </tbody>
                </table>
            </div>
        </div>
    </div>
</template>