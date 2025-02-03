<script setup lang="ts">
import { ref, onMounted } from 'vue'
import axios from 'axios'

const processes = ref<any[]>([])
const showRawData = ref(false)
const rawProcessesData = ref<any>(null)

function showToast(message: string, type: 'success' | 'error') {
    const toast = document.getElementById('toast') as HTMLDivElement
    if (toast) {
        toast.textContent = message
        toast.className = `toast toast-${type}`
        setTimeout(() => {
            toast.className = 'toast hidden'
        }, 3000)
    }
}

async function loadProcesses() {
    try {
        const res = await axios.get('/web_api/system/processes')
        processes.value = res.data
        rawProcessesData.value = res.data
    } catch (error) {
        showToast('获取进程列表失败', 'error')
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
        </div>
    </div>

    <!-- Toast -->
    <div id="toast" class="toast hidden"></div>
</template>

<style scoped>
.toast {
    position: fixed;
    bottom: 1rem;
    right: 1rem;
    padding: 1rem;
    border-radius: 0.5rem;
    z-index: 1000;
}
</style>