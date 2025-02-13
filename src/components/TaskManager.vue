<template>
    <div class="p-4">
        <h1 class="text-2xl font-bold">任务管理</h1>
        <ul class="mt-4">
            <li v-for="task in tasks" :key="task" class="flex justify-between items-center p-2 border-b">
                <span>{{ task }}</span>
                <div>
                    <button @click="executeTask(task)" class="btn btn-primary">执行</button>
                    <button @click="removeTask(task)" class="btn btn-danger">删除</button>
                </div>
            </li>
        </ul>
    </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import axios from 'axios'

const tasks = ref<string[]>([])

const fetchTasks = async () => {
    try {
        const response = await axios.get('/web_api/tasks')
        tasks.value = response.data
    } catch (error) {
        console.error('获取任务列表失败:', error)
    }
}

const removeTask = async (task: string) => {
    try {
        await axios.delete(`/web_api/tasks/${task}`)
        await fetchTasks()
    } catch (error) {
        console.error('删除任务失败:', error)
    }
}

const executeTask = async (task: string) => {
    try {
        await axios.post(`/web_api/tasks/${task}/execute`)
        await fetchTasks()
    } catch (error) {
        console.error('执行任务失败:', error)
    }
}

onMounted(() => {
    fetchTasks()
})
</script>