<script setup lang="ts">
import { ref, onMounted } from 'vue'
import ConfirmDialog from '../components/ConfirmDialog.vue'
import Alert from '../components/Alert.vue'
import CronTaskForm from '../components/CronTaskForm.vue'
import PageContainer from '../layouts/PageContainer.vue'
import { type CronTask, getTasks, createTask, updateTask, deleteTask, toggleTask } from '../api/cronTasks'

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
        tasks.value = await getTasks()
        error.value = null
    } catch (err: any) {
        error.value = '获取定时任务失败: ' + (err.response?.data?.error || err.message || '未知错误')
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
            await updateTask(task)
        } else {
            await createTask(task)
        }
        showModal.value = false
        await fetchTasks()
        resetForm()
    } catch (err: any) {
        throw new Error(err.response?.data?.error || err.message || '未知错误')
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
        await deleteTask(taskToDelete.value)
        await fetchTasks()
    } catch (err: any) {
        error.value = '删除任务失败: ' + (err.response?.data?.error || err.message || '未知错误')
    } finally {
        showDeleteConfirm.value = false
        taskToDelete.value = null
    }
}

// 启用/禁用任务
const handleToggleTask = async (id: number) => {
    try {
        await toggleTask(id)
        await fetchTasks()
    } catch (err: any) {
        error.value = '切换任务状态失败: ' + (err.response?.data?.error || err.message || '未知错误')
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
    <PageContainer title="定时任务管理">
        <template #header>
            <div class="flex justify-between items-center mb-6">
                <button class="btn btn-primary btn-md">
                    <i class="ri-add-line"></i>新建任务
                </button>
            </div>

            <Alert v-if="error" type="error" :message="error" />
        </template>

        <!-- 任务列表 -->
        <div class="overflow-x-auto bg-base-100 rounded-lg shadow">
            <table class="table table-zebra">
                <thead>
                    <tr class="bg-base-200">
                        <th class="font-bold">名称</th>
                        <th class="font-bold">命令</th>
                        <th class="font-bold">定时表达式</th>
                        <th class="font-bold">描述</th>
                        <th class="font-bold">状态</th>
                        <th class="font-bold">上次执行</th>
                        <th class="font-bold">下次执行</th>
                        <th class="font-bold">操作</th>
                    </tr>
                </thead>
                <tbody>
                    <tr v-for="task in tasks" :key="task.id" class="hover">
                        <td>{{ task.name }}</td>
                        <td class="max-w-xs truncate">
                            <div class="tooltip" :data-tip="task.command">
                                {{ task.command }}
                            </div>
                        </td>
                        <td>{{ task.schedule }}</td>
                        <td class="max-w-xs truncate">
                            <div class="tooltip" :data-tip="task.description">
                                {{ task.description }}
                            </div>
                        </td>
                        <td>
                            <div class="form-control">
                                <input type="checkbox" class="toggle toggle-primary toggle-sm" :checked="task.enabled"
                                    @change="handleToggleTask(task.id)" />
                            </div>
                        </td>
                        <td>{{ formatTime(task.last_run) }}</td>
                        <td>{{ formatTime(task.next_run) }}</td>
                        <td>
                            <div class="flex gap-2">
                                <button class="btn btn-ghost btn-sm" @click="editTask(task)">
                                    <i class="ri-edit-line text-primary"></i>
                                </button>
                                <button class="btn btn-ghost btn-sm" @click="confirmDelete(task.id)">
                                    <i class="ri-delete-bin-line text-error"></i>
                                </button>
                            </div>
                        </td>
                    </tr>
                </tbody>
            </table>
        </div>

        <!-- 创建/编辑任务模态框 -->
        <dialog class="modal" :class="{ 'modal-open': showModal }">
            <div class="modal-box max-w-2xl">
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
    </PageContainer>
</template>