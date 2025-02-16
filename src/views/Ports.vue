<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getPorts } from '../api/ports'

const ports = ref<any[]>([])
const error = ref<string | null>(null)

async function loadPorts() {
    try {
        ports.value = await getPorts()
        error.value = null
    } catch (err: any) {
        error.value = `获取端口列表失败: ${err.response?.data?.message || err.message || '未知错误'}`
    }
}

onMounted(() => {
    loadPorts()
})
</script>

<template>
    <div class="card bg-base-100 shadow-xl">
        <div class="card-body">
            <h2 class="card-title">端口列表</h2>

            <div v-if="error" class="alert alert-error">
                {{ error }}
            </div>

            <div v-else class="overflow-x-auto">
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
                        <tr v-for="port in ports" :key="port.port">
                            <td>{{ port.port }}</td>
                            <td>{{ port.protocol }}</td>
                            <td>{{ port.process }}</td>
                            <td>{{ port.state }}</td>
                        </tr>
                    </tbody>
                </table>
            </div>
        </div>
    </div>
</template>