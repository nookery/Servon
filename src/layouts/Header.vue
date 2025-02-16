<script setup lang="ts">
/// <reference lib="es2015" />
import { ref, onMounted } from 'vue'
import ThemeSwitcher from '../modules/ThemeSwitcher.vue'
import pkg from '../../package.json'
import { useLogViewerStore } from '../stores/logViewer'
import TaskManager from '../components/TaskManager.vue'
import { systemAPI } from '../api/info'
import GitHubAppForm from '../components/GitHubAppForm.vue'
import type { IPInfo } from '../api/info'

const currentUser = ref('')
const cpuUsage = ref(0)
const memoryUsage = ref(0)
const diskUsage = ref(0)
const osInfo = ref('')
const currentTheme = ref(localStorage.getItem('theme') || 'light')
const downloadSpeed = ref(0)
const uploadSpeed = ref(0)
const isTaskManagerVisible = ref(false)
const ipInfo = ref<IPInfo | null>(null)
const isIPInfoVisible = ref(false)

const logViewerStore = useLogViewerStore()
const githubFormRef = ref<InstanceType<typeof GitHubAppForm> | null>(null)

const toggleLogViewer = () => {
    logViewerStore.toggleVisibility()
}

// 获取系统资源使用情况
const fetchSystemResources = async () => {
    try {
        const res = await systemAPI.getResources()
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
        const res = await systemAPI.getOSInfo()
        osInfo.value = res.data.os_info
    } catch (error) {
        console.error('获取操作系统信息失败:', error)
    }
}

// 获取网络资源使用情况
const fetchNetworkResources = async () => {
    try {
        const res = await systemAPI.getNetworkResources()
        downloadSpeed.value = res.data.download_speed
        uploadSpeed.value = res.data.upload_speed
    } catch (error) {
        console.error('获取网络资源信息失败:', error)
    }
}

// 获取IP信息
const fetchIPInfo = async () => {
    try {
        const res = await systemAPI.getIPInfo()
        ipInfo.value = res.data
    } catch (error) {
        console.error('获取IP信息失败:', error)
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

const startGitHubIntegration = () => {
    githubFormRef.value?.showModal()
}

const handleGitHubSuccess = () => {
    // 可以添加成功提示
}

const handleGitHubError = (error: string) => {
    alert('启动GitHub集成失败: ' + error)
}

const toggleIPInfo = () => {
    isIPInfoVisible.value = !isIPInfoVisible.value
}

onMounted(async () => {
    try {
        const res = await systemAPI.getCurrentUser()
        currentUser.value = res.data.username
    } catch (error) {
        console.error('获取用户信息失败:', error)
    }

    fetchSystemResources()
    fetchOSInfo()
    fetchNetworkResources()
    fetchIPInfo()
    setInterval(() => {
        fetchSystemResources()
        fetchOSInfo()
        fetchNetworkResources()
        fetchIPInfo()
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
                <!-- Add IP info button before GitHub integration button -->
                <div class="relative">
                    <button @click="toggleIPInfo" class="btn btn-ghost btn-circle" title="IP信息">
                        <i class="ri-global-line text-xl"></i>
                    </button>

                    <!-- IP Info Dropdown -->
                    <div v-if="isIPInfoVisible" class="absolute right-0 mt-2 w-96 card bg-base-100 shadow-xl z-50">
                        <div class="card-body p-4">
                            <h3 class="card-title text-sm mb-2">IP 信息</h3>

                            <!-- Public IPs -->
                            <div class="mb-4">
                                <div class="text-sm font-semibold mb-1">公网IP</div>
                                <div class="text-sm">IPv4: {{ ipInfo?.public_ip || '未获取' }}</div>
                                <div class="text-sm">IPv6: {{ ipInfo?.public_ipv6 || '未获取' }}</div>
                            </div>

                            <!-- Local IPs -->
                            <div class="mb-4">
                                <div class="text-sm font-semibold mb-1">本地IP</div>
                                <div v-for="ip in ipInfo?.local_ips" :key="ip.ip" class="text-sm mb-1">
                                    <div>{{ ip.interface }}: {{ ip.ip }}</div>
                                    <div class="text-xs opacity-70">掩码: {{ ip.netmask }}</div>
                                </div>
                            </div>

                            <!-- DNS Servers -->
                            <div class="mb-4">
                                <div class="text-sm font-semibold mb-1">DNS 服务器</div>
                                <div v-for="dns in ipInfo?.dns_servers" :key="dns" class="text-sm">
                                    {{ dns }}
                                </div>
                            </div>
                        </div>
                    </div>
                </div>

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

    <GitHubAppForm ref="githubFormRef" @success="handleGitHubSuccess" @error="handleGitHubError" />
</template>

<style scoped>
.dropdown-content {
    max-height: 300px;
}
</style>