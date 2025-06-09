<script setup lang="ts">
import { ref, onMounted } from 'vue'
import axios from 'axios'
import IconButton from '../components/IconButton.vue'
const isTaskManagerVisible = ref(false)
const tasks = ref<string[]>([])

const showTaskManager = () => {
    isTaskManagerVisible.value = true
}

const closeTaskManager = () => {
    isTaskManagerVisible.value = false
}

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

<template>
    <div>
        <!-- 任务管理按钮 -->
        <IconButton @click="showTaskManager" icon="ri-task-line" variant="ghost" circle title="任务管理" />

        <!-- Task Manager Modal -->
        <div v-if="isTaskManagerVisible" class="modal modal-open">
            <div class="modal-box">
                <h2 class="font-bold text-lg">任务管理</h2>

                <!-- Task Manager Content -->
                <div class="p-4">
                    <ul class="mt-4">
                        <li v-for="task in tasks" :key="task" class="flex justify-between items-center p-2 border-b">
                            <span>{{ task }}</span>
                            <div class="flex gap-2">
                                <button @click="executeTask(task)" class="btn btn-primary btn-sm">
                                    执行
                                </button>
                                <button @click="removeTask(task)" class="btn btn-error btn-sm">
                                    删除
                                </button>
                            </div>
                        </li>
                    </ul>
                </div>

                <div class="modal-action">
                    <button @click="closeTaskManager" class="btn">关闭</button>
                </div>
            </div>
        </div>
    </div>
</template>