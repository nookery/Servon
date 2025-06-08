<script setup lang="ts">
import { ref, onMounted } from 'vue'
import CronTaskForm from '../modules/CronTaskForm.vue'
import PageContainer from '../layouts/PageContainer.vue'
import { type CronTask, getTasks, createTask, updateTask, deleteTask, toggleTask } from '../api/cronTasks'
import { useConfirm } from '../composables/useConfirm'
import { useToast } from '../composables/useToast'

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
const error = ref<string | null>(null)
const confirm = useConfirm()
const toast = useToast()

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
            toast.success('更新任务成功')
        } else {
            await createTask(task)
            toast.success('创建任务成功')
        }
        showModal.value = false
        await fetchTasks()
        resetForm()
    } catch (err: any) {
        throw new Error(err.response?.data?.error || err.message || '未知错误')
    }
}

// 删除任务
const handleDelete = async (task: CronTask) => {
    if (await confirm.error('删除任务', `确定要删除任务 "${task.name}" 吗？此操作不可撤销。`, {
        confirmText: '删除'
    })) {
        try {
            await deleteTask(task.id)
            await fetchTasks()
            toast.success('删除任务成功')
        } catch (err: any) {
            error.value = '删除任务失败: ' + (err.response?.data?.error || err.message || '未知错误')
        }
    }
}

// 启用/禁用任务
const handleToggleTask = async (id: number) => {
    try {
        await toggleTask(id)
        await fetchTasks()
        toast.success('切换任务状态成功')
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
    <PageContainer title="定时任务管理" :error="error">
        <template #header>
            <div class="flex justify-between items-center mb-6">
                <button class="btn btn-primary btn-md" @click="showModal = true">
                    <i class="ri-add-line"></i>新建任务
                </button>
            </div>
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
                                <button class="btn btn-ghost btn-sm" @click="handleDelete(task)">
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
    </PageContainer>
</template>