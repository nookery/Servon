<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { systemAPI, type Process } from '../api/info'
import PageContainer from '../components/PageContainer.vue'

const processes = ref<Process[]>([])
const filteredProcesses = ref<Process[]>([])
const error = ref<string | null>(null)
const sortKey = ref<string>('pid')
const sortOrder = ref<string>('asc')
const searchText = ref('')

// 预定义的常用软件过滤器
const commonFilters = [
    { name: 'Nginx', pattern: /nginx/i },
    { name: 'MySQL', pattern: /(mysql|mysqld)/i },
    { name: 'PHP', pattern: /(php|php-fpm)/i },
    { name: 'Redis', pattern: /redis/i },
    { name: 'Node.js', pattern: /node/i }
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

// 过滤进程
function filterProcesses(pattern: RegExp | null) {
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
                    <button class="btn btn-sm bg-base-100 hover:bg-base-300"
                        :class="{ 'btn-primary': !searchText && filteredProcesses.length !== processes.length }"
                        @click="filterProcesses(null)">
                        全部
                    </button>
                    <button v-for="filter in commonFilters" :key="filter.name"
                        class="btn btn-sm bg-base-100 hover:bg-base-300"
                        :class="{ 'btn-primary': filteredProcesses.length !== processes.length }"
                        @click="filterProcesses(filter.pattern)">
                        {{ filter.name }}
                    </button>
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
                    <tr v-for="process in filteredProcesses" :key="process.pid">
                        <td>{{ process.pid }}</td>
                        <td>{{ process.command }}</td>
                        <td>{{ process.cpu.toFixed(1) }}%</td>
                        <td>{{ process.memory.toFixed(1) }}%</td>
                        <td>
                            <button class="btn btn-error btn-sm" @click="killProcess(process.pid)">结束</button>
                        </td>
                    </tr>
                </tbody>
            </table>
        </div>
    </PageContainer>
</template>