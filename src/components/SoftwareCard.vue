<template>
    <div class="bg-white shadow rounded-lg p-6" :class="{
        'border-2 border-red-300': installFailed,
        'border-2 border-blue-300': installing,
        'border-2 border-green-300': installSuccess
    }">
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
                <span class="px-2 py-1 text-sm rounded-full" :class="getStatusClass(software.status)">
                    {{ getStatusText(software.status) }}
                </span>
                <button class="px-3 py-1 text-sm rounded-md text-white" :class="[
                    software.status === 'not_installed' ? 'bg-blue-500 hover:bg-blue-600' : 'bg-red-500 hover:bg-red-600',
                    installing ? 'opacity-50 cursor-not-allowed' : ''
                ]" :disabled="installing" @click="handleAction">
                    {{
                        installing ? '安装中...' :
                            software.status === 'not_installed' ? '安装' : '卸载'
                    }}
                </button>
            </div>
        </div>
        <p v-if="software.path" class="text-sm text-gray-500 mt-2">
            路径: {{ software.path }}
        </p>
        <!-- 安装状态提示 -->
        <div v-if="installFailed" class="mt-2 p-2 bg-red-50 border border-red-200 rounded-md">
            <div class="flex justify-between items-center">
                <p class="text-sm text-red-600">安装过程中断</p>
                <button @click="clearError" class="text-sm text-red-500 hover:text-red-700" title="关闭错误提示">
                    <XMarkIcon class="h-5 w-5" />
                </button>
            </div>
        </div>
        <!-- 安装日志 -->
        <TransitionRoot v-if="logs.length > 0" as="template" show>
            <div class="mt-4">
                <div class="flex justify-between items-center mb-2">
                    <h3 class="text-sm font-medium text-gray-900">安装日志</h3>
                    <div class="flex gap-2">
                        <button @click="copyLogs" class="text-sm text-gray-500 hover:text-gray-700">
                            复制日志
                        </button>
                        <button @click="clearLogs" class="text-sm text-gray-500 hover:text-gray-700">
                            关闭日志
                        </button>
                    </div>
                </div>
                <div class="bg-gray-50 p-4 rounded-md">
                    <div class="max-h-60 overflow-auto space-y-1">
                        <div v-for="(log, index) in logs" :key="index" class="text-sm text-gray-600">
                            {{ log }}
                        </div>
                    </div>
                </div>
            </div>
        </TransitionRoot>
        <!-- 卸载确认对话框 -->
        <TransitionRoot appear :show="showUninstallDialog" as="template">
            <Dialog as="div" @close="showUninstallDialog = false" class="relative z-10">
                <TransitionChild as="template" enter="duration-300 ease-out" enter-from="opacity-0"
                    enter-to="opacity-100" leave="duration-200 ease-in" leave-from="opacity-100" leave-to="opacity-0">
                    <div class="fixed inset-0 bg-black bg-opacity-25" />
                </TransitionChild>

                <div class="fixed inset-0 overflow-y-auto">
                    <div class="flex min-h-full items-center justify-center p-4 text-center">
                        <TransitionChild as="template" enter="duration-300 ease-out" enter-from="opacity-0 scale-95"
                            enter-to="opacity-100 scale-100" leave="duration-200 ease-in"
                            leave-from="opacity-100 scale-100" leave-to="opacity-0 scale-95">
                            <DialogPanel
                                class="w-full max-w-md transform overflow-hidden rounded-2xl bg-white p-6 text-left align-middle shadow-xl transition-all">
                                <DialogTitle as="h3" class="text-lg font-medium leading-6 text-gray-900">
                                    确认卸载
                                </DialogTitle>
                                <div class="mt-2">
                                    <p class="text-sm text-gray-500">
                                        确定要卸载 {{ software.name }} 吗？这可能会删除相关的数据和配置。
                                    </p>
                                </div>

                                <div class="mt-4 flex justify-end space-x-3">
                                    <button type="button"
                                        class="inline-flex justify-center rounded-md border border-transparent bg-red-100 px-4 py-2 text-sm font-medium text-red-900 hover:bg-red-200"
                                        @click="confirmUninstall">
                                        确定卸载
                                    </button>
                                    <button type="button"
                                        class="inline-flex justify-center rounded-md border border-transparent bg-gray-100 px-4 py-2 text-sm font-medium text-gray-900 hover:bg-gray-200"
                                        @click="showUninstallDialog = false">
                                        取消
                                    </button>
                                </div>
                            </DialogPanel>
                        </TransitionChild>
                    </div>
                </div>
            </Dialog>
        </TransitionRoot>
        <!-- Toast 消息 -->
        <Toast :show="showToast" :message="toastMessage" :type="toastType" @close="showToast = false" />
    </div>
</template>

<script setup lang="ts">
import { ref, onUnmounted } from 'vue'
import {
    TransitionRoot,
    TransitionChild,
    Dialog,
    DialogPanel,
    DialogTitle
} from '@headlessui/vue'
import Toast from './Toast.vue'
import { XMarkIcon } from '@heroicons/vue/24/outline'

interface Software {
    name: string
    version?: string
    status: string
    path?: string
    description?: string
}

const props = defineProps<{
    software: Software
}>()

const installing = ref(false)
const installFailed = ref(false)
const installSuccess = ref(false)
const logs = ref<string[]>([])
const showUninstallDialog = ref(false)
const showToast = ref(false)
const toastMessage = ref('')
const toastType = ref<'success' | 'error'>('success')

function getStatusClass(status: string): string {
    switch (status) {
        case 'running':
            return 'bg-green-100 text-green-800'
        case 'stopped':
            return 'bg-red-100 text-red-800'
        default:
            return 'bg-gray-100 text-gray-800'
    }
}

function getStatusText(status: string): string {
    switch (status) {
        case 'running':
            return '运行中'
        case 'stopped':
            return '已停止'
        default:
            return '未安装'
    }
}

function install() {
    installing.value = true
    installFailed.value = false // 重置失败状态
    const eventSource = new EventSource(`/api/system/software/${props.software.name}/install`)

    eventSource.onmessage = (event) => {
        logs.value.push(event.data)

        if (event.data === '安装完成') {
            eventSource.close()
            installing.value = false
            installSuccess.value = true
            setTimeout(() => {
                installSuccess.value = false
            }, 3000)
        }
    }

    eventSource.onerror = () => {
        eventSource.close()
        installing.value = false
        const lastLog = logs.value[logs.value.length - 1]
        if (lastLog === '安装完成') {
            installSuccess.value = true
            setTimeout(() => {
                installSuccess.value = false
            }, 3000)
        } else {
            installFailed.value = true
            logs.value.push('连接已关闭，安装可能未完成')
            // 移除失败状态的自动清除
        }
    }
}

function clearLogs() {
    logs.value = []
}

async function confirmUninstall() {
    showUninstallDialog.value = false
    installing.value = true
    installFailed.value = false // 重置失败状态
    logs.value.push(`开始卸载 ${props.software.name}...`)

    try {
        // 先停止服务
        logs.value.push('正在停止服务...')
        const stopResponse = await fetch(`/api/system/software/${props.software.name}/stop`, {
            method: 'POST'
        })
        if (!stopResponse.ok) {
            const data = await stopResponse.json()
            throw new Error(data.error || '停止服务失败')
        }
        logs.value.push('服务已停止')

        // 执行卸载
        logs.value.push('正在卸载软件...')
        const response = await fetch(`/api/system/software/${props.software.name}/uninstall`, {
            method: 'POST'
        })
        if (!response.ok) {
            const data = await response.json()
            throw new Error(data.error || '卸载失败')
        }
        logs.value.push('卸载完成')
        installSuccess.value = true
        setTimeout(() => {
            installSuccess.value = false
        }, 3000)
    } catch (err) {
        installFailed.value = true
        logs.value.push(`错误: ${err instanceof Error ? err.message : '未知错误'}`)
        // 移除失败状态的自动清除
    } finally {
        installing.value = false
    }
}

function handleAction() {
    if (props.software.status === 'not_installed') {
        install()
    } else {
        showUninstallDialog.value = true
    }
}

function showSuccessToast(message: string) {
    toastMessage.value = message
    toastType.value = 'success'
    showToast.value = true
    setTimeout(() => {
        showToast.value = false
    }, 3000)
}

function showErrorToast(message: string) {
    toastMessage.value = message
    toastType.value = 'error'
    showToast.value = true
    setTimeout(() => {
        showToast.value = false
    }, 3000)
}

function copyLogs() {
    const logText = logs.value.join('\n')
    navigator.clipboard.writeText(logText).then(() => {
        showSuccessToast('日志已复制到剪贴板')
    }).catch(err => {
        console.error('复制失败:', err)
        showErrorToast('复制失败，请手动复制')
    })
}

// 添加清除错误状态的函数
function clearError() {
    installFailed.value = false
}

// 组件卸载时清理
onUnmounted(() => {
    clearLogs()
})
</script>