<script setup lang="ts">
import { ref, onMounted } from 'vue'
import axios from 'axios'

const files = ref<any[]>([])
const currentPath = ref('/')
const error = ref<string | null>(null)

async function loadFiles(path: string) {
    try {
        const res = await axios.get(`/web_api/system/files?path=${path}`)
        files.value = res.data
        error.value = null
    } catch (err) {
        error.value = '获取文件列表失败'
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

            <div v-if="error" class="alert alert-error">
                {{ error }}
            </div>

            <template v-else>
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
            </template>
        </div>
    </div>
</template>