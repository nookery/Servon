<script setup lang="ts">
/// <reference lib="es2015" />
import { ref, onMounted } from 'vue'
import axios from 'axios'
import ThemeSwitcher from '../modules/ThemeSwitcher.vue'
import pkg from '../../package.json'
import { useLogViewerStore } from '../stores/logViewer'
import TaskManager from '../components/TaskManager.vue'

const currentUser = ref('')
const cpuUsage = ref(0)
const memoryUsage = ref(0)
const diskUsage = ref(0)
const osInfo = ref('')
const currentTheme = ref(localStorage.getItem('theme') || 'light')
const downloadSpeed = ref(0)
const uploadSpeed = ref(0)
const isTaskManagerVisible = ref(false)

const logViewerStore = useLogViewerStore()

const toggleLogViewer = () => {
    logViewerStore.toggleVisibility()
}

// 获取系统资源使用情况
const fetchSystemResources = async () => {
    try {
        const res = await axios.get('/web_api/system/resources')
        cpuUsage.value = res.data.cpu_usage
        memoryUsage.value = res.data.memory_usage
        diskUsage.value = res.data.disk_usage
    } catch (error) {
        console.error('获取系统资源信息失败:', error)
    }
}

// 获取操作系统信息
const fetchOSInfo = async () => {
    try {
        const res = await axios.get('/web_api/system/os')
        osInfo.value = res.data.os_info
    } catch (error) {
        console.error('获取操作系统信息失败:', error)
    }
}

// 获取网络资源使用情况
const fetchNetworkResources = async () => {
    try {
        const res = await axios.get('/web_api/system/network')
        downloadSpeed.value = res.data.download_speed
        uploadSpeed.value = res.data.upload_speed
    } catch (error) {
        console.error('获取网络资源信息失败:', error)
    }
}

function changeTheme(theme: string) {
    document.documentElement.setAttribute('data-theme', theme)
    localStorage.setItem('theme', theme)
    currentTheme.value = theme
}

const showTaskManager = () => {
    isTaskManagerVisible.value = true
}

const closeTaskManager = () => {
    isTaskManagerVisible.value = false
}

const startGitHubIntegration = async () => {
    try {
        // 弹出对话框让用户输入基本信息
        const appName = prompt('请输入GitHub App名称:', 'Servon App')
        if (!appName) return

        const description = prompt('请输入描述(可选):', 'Servon GitHub integration for automation')

        // 调用后端API启动GitHub App Manifest flow
        const response = await axios.post('/web_api/github/setup', {
            name: appName,
            description: description || undefined,
        }, {
            responseType: 'text'  // 接收HTML响应
        })

        // 创建一个临时div来执行返回的HTML
        const div = document.createElement('div')
        div.innerHTML = response.data
        document.body.appendChild(div)

        // 立即提交表单
        const form = div.querySelector('#github-form') as HTMLFormElement
        if (form) {
            form.submit()
        } else {
            throw new Error('表单创建失败')
        }

        // 清理临时div（延迟清理，确保表单提交完成）
        setTimeout(() => {
            document.body.removeChild(div)
        }, 2000)
    } catch (error: any) {
        console.error('启动GitHub集成失败:', error)
        alert('启动GitHub集成失败: ' + (error.response?.data?.error || error.message))
    }
}

onMounted(async () => {
    try {
        const res = await axios.get('/web_api/system/user')
        currentUser.value = res.data.username
    } catch (error) {
        console.error('获取用户信息失败:', error)
    }

    fetchSystemResources()
    fetchOSInfo()
    fetchNetworkResources()
    setInterval(() => {
        fetchSystemResources()
        fetchOSInfo()
        fetchNetworkResources()
    }, 50000)

    // 初始化主题
    const savedTheme = localStorage.getItem('theme')
    if (savedTheme) {
        changeTheme(savedTheme)
    }
})
</script>

<template>
    <div class="navbar bg-base-100 fixed top-0 left-0 z-50 h-16 px-6 shadow-sm">
        <div class="flex-1">
            <div class="flex items-center gap-2">
                <i class="ri-server-line text-2xl text-primary"></i>
                <div class="flex flex-col">
                    <span class="text-lg font-bold">{{ pkg.name.charAt(0).toUpperCase() + pkg.name.slice(1) }}</span>
                    <span class="text-xs text-base-content/60">{{ osInfo }}</span>
                </div>
            </div>
        </div>

        <div class="flex-none gap-6">
            <div class="flex items-center gap-6">
                <!-- Add GitHub integration button before the log viewer button -->
                <button @click="startGitHubIntegration" class="btn btn-ghost btn-circle" title="GitHub集成">
                    <i class="ri-github-line text-xl"></i>
                </button>

                <!-- CPU Usage -->
                <div class="w-36">
                    <div class="text-xs text-base-content/70 mb-1">
                        CPU: {{ cpuUsage.toFixed(1) }}%
                    </div>
                    <progress class="progress progress-primary h-2" :value="cpuUsage" max="100"></progress>
                </div>

                <!-- Memory Usage -->
                <div class="w-36">
                    <div class="text-xs text-base-content/70 mb-1">
                        内存: {{ memoryUsage.toFixed(1) }}%
                    </div>
                    <progress class="progress progress-primary h-2" :value="memoryUsage" max="100"></progress>
                </div>

                <!-- Disk Usage -->
                <div class="w-36">
                    <div class="text-xs text-base-content/70 mb-1">
                        磁盘: {{ diskUsage.toFixed(1) }}%
                    </div>
                    <progress class="progress progress-primary h-2" :value="diskUsage" max="100"></progress>
                </div>

                <!-- Network Usage -->
                <div class="flex gap-4">
                    <div class="flex items-center gap-1">
                        <i class="ri-download-line text-xs text-base-content/70"></i>
                        <span class="text-xs text-base-content/70">
                            {{ (downloadSpeed / 1024 / 1024).toFixed(1) }} MB/s
                        </span>
                    </div>
                    <div class="flex items-center gap-1">
                        <i class="ri-upload-line text-xs text-base-content/70"></i>
                        <span class="text-xs text-base-content/70">
                            {{ (uploadSpeed / 1024 / 1024).toFixed(1) }} MB/s
                        </span>
                    </div>
                </div>

                <!-- Log Viewer Button -->
                <button @click="toggleLogViewer" class="btn btn-ghost btn-circle">
                    <i class="ri-file-list-line text-xl"></i>
                </button>

                <!-- 任务管理按钮 -->
                <button @click="showTaskManager" class="btn btn-ghost btn-circle">
                    <i class="ri-task-line text-xl"></i>
                </button>

                <!-- Theme Switcher Component -->
                <ThemeSwitcher />

                <!-- User Avatar -->
                <div class="flex items-center gap-2">
                    <div class="avatar placeholder">
                        <div class="bg-primary text-primary-content rounded-full w-8">
                            <span class="text-xs">{{ currentUser.charAt(0).toUpperCase() }}</span>
                        </div>
                    </div>
                    <span class="text-sm">{{ currentUser }}</span>
                </div>
            </div>
        </div>
    </div>

    <!-- Task Manager Modal -->
    <div v-if="isTaskManagerVisible" class="modal modal-open">
        <div class="modal-box">
            <h2 class="font-bold text-lg">任务管理</h2>
            <TaskManager />
            <div class="modal-action">
                <button @click="closeTaskManager" class="btn">关闭</button>
            </div>
        </div>
    </div>
</template>

<style scoped>
.dropdown-content {
    max-height: 300px;
}
</style>