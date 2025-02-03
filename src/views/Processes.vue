<script setup lang="ts">
import { ref, onMounted } from 'vue'
import axios from 'axios'

const processes = ref<any[]>([])
const showRawData = ref(false)
const rawProcessesData = ref<any>(null)
const error = ref<string | null>(null)

async function loadProcesses() {
    try {
        const res = await axios.get('/web_api/system/processes')
        processes.value = res.data
        rawProcessesData.value = res.data
        error.value = null
    } catch (err) {
        error.value = '获取进程列表失败'
    }
}

onMounted(() => {
    loadProcesses()
})
</script>

<template>
    <div class="card bg-base-100 shadow-xl">
        <div class="card-body">
            <h2 class="card-title">进程列表</h2>

            <div v-if="error" class="alert alert-error">
                {{ error }}
            </div>

            <template v-else>
                <div class="flex justify-end mb-4">
                    <button class="btn btn-primary" @click="showRawData = !showRawData">
                        {{ showRawData ? '显示进程列表' : '显示原始数据' }}
                    </button>
                </div>

                <div v-if="!showRawData" class="overflow-x-auto">
                    <table class="table table-zebra w-full">
                        <thead>
                            <tr>
                                <th>PID</th>
                                <th>Command</th>
                                <th>CPU</th>
                                <th>内存</th>
                            </tr>
                        </thead>
                        <tbody>
                            <tr v-for="process in processes" :key="process.pid">
                                <td>{{ process.pid }}</td>
                                <td>{{ process.command }}</td>
                                <td>{{ process.cpu }}</td>
                                <td>{{ process.memory }}</td>
                            </tr>
                        </tbody>
                    </table>
                </div>

                <pre v-else
                    class="bg-base-200 p-4 rounded-lg overflow-auto font-mono text-left whitespace-pre">{{ JSON.stringify(rawProcessesData, null, 2) }}</pre>
            </template>
        </div>
    </div>
</template>