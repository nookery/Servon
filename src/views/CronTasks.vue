<script setup lang="ts">
import { ref, onMounted } from 'vue'
import axios from 'axios'
import ConfirmDialog from '../components/ConfirmDialog.vue'
import Alert from '../components/Alert.vue'
import CronTaskForm from '../components/CronTaskForm.vue'

interface CronTask {
    id: number
    name: string
    command: string
    schedule: string
    description: string
    enabled: boolean
    last_run?: string
    next_run?: string
}

const tasks = ref<CronTask[]>([])
const showModal = ref(false)
const editingTask = ref<CronTask | null>(null)
const newTask = ref<CronTask>({
    id: 0,
    name: '',
    command: '',
    schedule: '',
    description: '',
    enabled: true
})

const formError = ref('')
const fieldErrors = ref<Record<string, string>>({})
const showDeleteConfirm = ref(false)
const taskToDelete = ref<number | null>(null)
const error = ref<string | null>(null)

// 获取所有定时任务
const fetchTasks = async () => {
    try {
        const res = await axios.get('/web_api/cron/tasks')
        tasks.value = res.data
        error.value = null
    } catch (error: any) {
        const errorMessage = error.response?.data?.error || error.message || '未知错误'
        error.value = '获取定时任务失败: ' + errorMessage
    }
}

// 清除错误信息
const clearErrors = () => {
    formError.value = ''
    fieldErrors.value = {}
}

// 创建或更新任务
const saveTask = async (task: CronTask) => {
    try {
        if (editingTask.value) {
            await axios.put(`/web_api/cron/tasks/${task.id}`, task)
        } else {
            await axios.post('/web_api/cron/tasks', task)
        }
        showModal.value = false
        await fetchTasks()
        resetForm()
    } catch (error: any) {
        const errorMessage = error.response?.data?.error || error.message || '未知错误'
        throw new Error(errorMessage)
    }
}

// 修改删除任务的处理逻辑
const confirmDelete = (id: number) => {
    taskToDelete.value = id
    showDeleteConfirm.value = true
}

const handleDelete = async () => {
    if (!taskToDelete.value) return

    try {
        await axios.delete(`/web_api/cron/tasks/${taskToDelete.value}`)
        await fetchTasks()
    } catch (error: any) {
        const errorMessage = error.response?.data?.error || error.message || '未知错误'
        error.value = '删除任务失败: ' + errorMessage
    } finally {
        showDeleteConfirm.value = false
        taskToDelete.value = null
    }
}

// 启用/禁用任务
const toggleTask = async (id: number) => {
    try {
        await axios.post(`/web_api/cron/tasks/${id}/toggle`)
        await fetchTasks()
    } catch (error: any) {
        const errorMessage = error.response?.data?.error || error.message || '未知错误'
        error.value = '切换任务状态失败: ' + errorMessage
    }
}

// 编辑任务
const editTask = (task: CronTask) => {
    editingTask.value = { ...task }
    showModal.value = true
}

// 重置表单
const resetForm = () => {
    editingTask.value = null
    clearErrors()
    newTask.value = {
        id: 0,
        name: '',
        command: '',
        schedule: '',
        description: '',
        enabled: true
    }
}

// 格式化时间
const formatTime = (time: string | undefined) => {
    if (!time) return '-'
    return new Date(time).toLocaleString()
}

onMounted(fetchTasks)
</script>

<template>
    <div class="p-6">
        <div class="flex justify-between items-center mb-6">
            <h1 class="text-2xl font-bold">定时任务管理</h1>
            <button class="btn btn-primary" @click="showModal = true">
                <i class="ri-add-line mr-1"></i>新建任务
            </button>
        </div>

        <Alert v-if="error" type="error" :message="error" />

        <!-- 任务列表 -->
        <div class="overflow-x-auto">
            <table class="table w-full">
                <thead>
                    <tr>
                        <th>名称</th>
                        <th>命令</th>
                        <th>定时表达式</th>
                        <th>描述</th>
                        <th>状态</th>
                        <th>上次执行</th>
                        <th>下次执行</th>
                        <th>操作</th>
                    </tr>
                </thead>
                <tbody>
                    <tr v-for="task in tasks" :key="task.id">
                        <td>{{ task.name }}</td>
                        <td class="max-w-xs truncate">{{ task.command }}</td>
                        <td>{{ task.schedule }}</td>
                        <td class="max-w-xs truncate">{{ task.description }}</td>
                        <td>
                            <div class="form-control">
                                <input type="checkbox" class="toggle toggle-primary" :checked="task.enabled"
                                    @change="toggleTask(task.id)" />
                            </div>
                        </td>
                        <td>{{ formatTime(task.last_run) }}</td>
                        <td>{{ formatTime(task.next_run) }}</td>
                        <td>
                            <div class="flex gap-2">
                                <button class="btn btn-sm btn-ghost" @click="editTask(task)">
                                    <i class="ri-edit-line"></i>
                                </button>
                                <button class="btn btn-sm btn-ghost text-error" @click="confirmDelete(task.id)">
                                    <i class="ri-delete-bin-line"></i>
                                </button>
                            </div>
                        </td>
                    </tr>
                </tbody>
            </table>
        </div>

        <!-- 创建/编辑任务模态框 -->
        <dialog class="modal" :class="{ 'modal-open': showModal }">
            <div class="modal-box">
                <h3 class="font-bold text-lg mb-4">
                    {{ editingTask ? '编辑任务' : '新建任务' }}
                </h3>
                <CronTaskForm :task="editingTask || newTask" :is-editing="!!editingTask" @submit="saveTask"
                    @cancel="showModal = false" />
            </div>
            <form method="dialog" class="modal-backdrop">
                <button @click="showModal = false">关闭</button>
            </form>
        </dialog>

        <!-- 使用确认对话框组件 -->
        <ConfirmDialog v-model:show="showDeleteConfirm" title="确认删除" message="该操作无法撤销，是否确认删除此任务？" type="warning"
            confirm-text="删除" @confirm="handleDelete" />
    </div>
</template>

<style scoped>
.input-error {
    border-color: rgb(252, 165, 165);
}

.label-text-alt.text-error {
    color: rgb(153, 27, 27);
}
</style>