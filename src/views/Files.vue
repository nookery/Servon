<script setup lang="ts">
import { ref, onMounted } from 'vue'
import axios from 'axios'

const files = ref<any[]>([])
const currentPath = ref('/')

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

async function loadFiles(path: string) {
    try {
        const res = await axios.get(`/web_api/system/files?path=${path}`)
        files.value = res.data
    } catch (error) {
        showToast('获取文件列表失败', 'error')
    }
}

onMounted(() => {
    loadFiles(currentPath.value)
})
</script>

<template>
    <div class="card bg-base-100 shadow-xl">
        <div class="card-body">
            <h2 class="card-title">文件管理</h2>

            <div class="breadcrumbs text-sm mb-4">
                <ul>
                    <li><a>根目录</a></li>
                    <li v-for="(part, index) in currentPath.split('/').filter(Boolean)" :key="index">
                        {{ part }}
                    </li>
                </ul>
            </div>

            <div class="overflow-x-auto">
                <table class="table table-zebra w-full">
                    <thead>
                        <tr>
                            <th>名称</th>
                            <th>类型</th>
                            <th>大小</th>
                            <th>修改时间</th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr v-for="file in files" :key="file.name">
                            <td>{{ file.name }}</td>
                            <td>{{ file.type }}</td>
                            <td>{{ file.size }}</td>
                            <td>{{ file.modTime }}</td>
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