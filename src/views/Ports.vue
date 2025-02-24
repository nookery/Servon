<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getPorts } from '../api/ports'
import PageContainer from '../layouts/PageContainer.vue'
import IconButton from '../components/IconButton.vue'

const ports = ref<any[]>([])
const filteredPorts = ref<any[]>([])
const error = ref<string | null>(null)
const activeFilter = ref<string | null>(null)

const commonFilters = [
    { name: 'Nginx', pattern: /nginx/i, icon: 'ri-server-line' },
    { name: 'MySQL', pattern: /(mysql|mysqld)/i, icon: 'ri-database-2-line' },
    { name: 'PHP', pattern: /(php|php-fpm)/i, icon: 'ri-code-s-slash-line' },
    { name: 'Redis', pattern: /redis/i, icon: 'ri-database-line' },
    { name: 'Node.js', pattern: /node/i, icon: 'ri-nodejs-line' }
]

async function loadPorts() {
    try {
        const data = await getPorts()
        ports.value = data
        filteredPorts.value = data
        error.value = null
    } catch (err: any) {
        error.value = `获取端口列表失败: ${err.response?.data?.message || err.message || '未知错误'}`
    }
}

function filterPorts(pattern: RegExp | null, filterName: string | null = null) {
    activeFilter.value = filterName

    if (!pattern) {
        filteredPorts.value = ports.value
        return
    }

    filteredPorts.value = ports.value.filter(port => {
        const process = port.process.toLowerCase()
        return pattern.test(process)
    })
}

onMounted(() => {
    loadPorts()
})
</script>

<template>
    <PageContainer title="端口列表" :error="error">
        <template #header>
            <div class="flex flex-wrap gap-2 mb-4 p-4 bg-base-200 rounded-lg border border-base-300">
                <IconButton icon="ri-list-unordered" size="sm" :variant="activeFilter === null ? 'primary' : 'default'"
                    @click="filterPorts(null, null)">
                    全部
                </IconButton>
                <IconButton v-for="filter in commonFilters" :key="filter.name" :icon="filter.icon" size="sm"
                    :variant="activeFilter === filter.name ? 'primary' : 'default'"
                    @click="filterPorts(filter.pattern, filter.name)">
                    {{ filter.name }}
                </IconButton>
            </div>
        </template>

        <div class="overflow-x-auto">
            <table class="table table-zebra w-full">
                <thead>
                    <tr>
                        <th>端口</th>
                        <th>协议</th>
                        <th>进程</th>
                        <th>状态</th>
                    </tr>
                </thead>
                <tbody>
                    <tr v-for="port in filteredPorts" :key="port.port">
                        <td>{{ port.port }}</td>
                        <td>{{ port.protocol }}</td>
                        <td>{{ port.process }}</td>
                        <td>{{ port.state }}</td>
                    </tr>
                </tbody>
            </table>
        </div>
    </PageContainer>
</template>