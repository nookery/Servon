<script setup lang="ts">
import { ref, onMounted } from 'vue'
import axios from 'axios'

const ports = ref<any[]>([])

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

async function loadPorts() {
    try {
        const res = await axios.get('/web_api/system/ports')
        ports.value = res.data
    } catch (error) {
        showToast('获取端口列表失败', 'error')
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