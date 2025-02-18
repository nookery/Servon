<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { systemAPI, type Process } from '../api/info'
import PageContainer from '../layouts/PageContainer.vue'
import IconButton from '../components/IconButton.vue'

const processes = ref<Process[]>([])
const filteredProcesses = ref<Process[]>([])
const error = ref<string | null>(null)
const sortKey = ref<string>('pid')
const sortOrder = ref<string>('asc')
const searchText = ref('')

// 添加当前激活的过滤器状态
const activeFilter = ref<string | null>(null)

// 预定义的常用软件过滤器
const commonFilters = [
    { name: 'Nginx', pattern: /nginx/i, icon: 'ri-server-line' },
    { name: 'MySQL', pattern: /(mysql|mysqld)/i, icon: 'ri-database-2-line' },
    { name: 'PHP', pattern: /(php|php-fpm)/i, icon: 'ri-code-s-slash-line' },
    { name: 'Redis', pattern: /redis/i, icon: 'ri-database-line' },
    { name: 'Node.js', pattern: /node/i, icon: 'ri-nodejs-line' }
]

// 加载进程列表
async function loadProcesses() {
    try {
        const res = await systemAPI.getProcesses()
        processes.value = res.data
        filteredProcesses.value = res.data
        error.value = null
    } catch (err) {
        error.value = '获取进程列表失败'
    }
}

// 修改过滤进程函数
function filterProcesses(pattern: RegExp | null, filterName: string | null = null) {
    activeFilter.value = filterName // 设置当前激活的过滤器

    if (!pattern && !searchText.value) {
        filteredProcesses.value = processes.value
        return
    }

    filteredProcesses.value = processes.value.filter(process => {
        const command = process.command.toLowerCase()
        const textMatch = searchText.value ?
            command.includes(searchText.value.toLowerCase()) :
            true
        const patternMatch = pattern ?
            pattern.test(command) :
            true
        return textMatch && patternMatch
    })
}

// 监听搜索文本变化
watch(searchText, () => {
    filterProcesses(null)
})

function sortProcesses(key: string) {
    if (sortKey.value === key) {
        sortOrder.value = sortOrder.value === 'asc' ? 'desc' : 'asc'
    } else {
        sortKey.value = key
        sortOrder.value = 'asc'
    }

    processes.value.sort((a, b) => {
        if (sortOrder.value === 'asc') {
            return a[key as keyof Process] > b[key as keyof Process] ? 1 : -1
        } else {
            return a[key as keyof Process] < b[key as keyof Process] ? 1 : -1
        }
    })
}

async function killProcess(pid: number) {
    try {
        await systemAPI.killProcess(pid)
        loadProcesses()
    } catch (err) {
        error.value = `结束进程 ${pid} 失败`
    }
}

onMounted(() => {
    loadProcesses()
})
</script>

<template>
    <PageContainer title="进程列表">
        <template #header>
            <div v-if="error" class="alert alert-error">
                {{ error }}
            </div>

            <!-- 操作栏 -->
            <div class="flex flex-wrap gap-2 mb-4 p-4 bg-base-200 rounded-lg border border-base-300">
                <!-- 搜索框组 -->
                <div class="flex-1 min-w-[240px]">
                    <input type="text" v-model="searchText" placeholder="搜索进程..."
                        class="input input-bordered w-full max-w-xs bg-base-100" />
                </div>

                <!-- 快速过滤按钮组 -->
                <div class="flex flex-wrap gap-2 flex-1">
                    <IconButton icon="ri-list-unordered" size="sm"
                        :variant="activeFilter === null ? 'primary' : 'default'" @click="filterProcesses(null, null)">
                        全部
                    </IconButton>
                    <IconButton v-for="filter in commonFilters" :key="filter.name" :icon="filter.icon" size="sm"
                        :variant="activeFilter === filter.name ? 'primary' : 'default'"
                        @click="filterProcesses(filter.pattern, filter.name)">
                        {{ filter.name }}
                    </IconButton>
                </div>
            </div>
        </template>

        <!-- 默认插槽用于主要内容 -->
        <div class="overflow-x-auto">
            <table class="table table-zebra w-full">
                <thead>
                    <tr>
                        <th @click="sortProcesses('pid')">PID</th>
                        <th @click="sortProcesses('command')">Command</th>
                        <th @click="sortProcesses('cpu')">CPU</th>
                        <th @click="sortProcesses('memory')">内存</th>
                        <th>操作</th>
                    </tr>
                </thead>
                <tbody>
                    <template v-if="filteredProcesses.length > 0">
                        <tr v-for="process in filteredProcesses" :key="process.pid">
                            <td>{{ process.pid }}</td>
                            <td>{{ process.command }}</td>
                            <td>{{ process.cpu.toFixed(1) }}%</td>
                            <td>{{ process.memory.toFixed(1) }}%</td>
                            <td>
                                <IconButton icon="ri-close-line" size="sm" variant="error"
                                    :title="`结束进程 ${process.pid}`" @click="killProcess(process.pid)">
                                    结束
                                </IconButton>
                            </td>
                        </tr>
                    </template>
                    <tr v-else>
                        <td colspan="5" class="text-center py-8">
                            <div class="flex flex-col items-center gap-2 text-base-content/60">
                                <i class="ri-inbox-line text-4xl"></i>
                                <p>没有找到匹配的进程</p>
                                <IconButton v-if="activeFilter !== null || searchText" icon="ri-refresh-line"
                                    variant="ghost" size="sm" @click="filterProcesses(null, null); searchText = ''">
                                    清除筛选
                                </IconButton>
                            </div>
                        </td>
                    </tr>
                </tbody>
            </table>
        </div>
    </PageContainer>
</template>