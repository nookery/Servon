<script setup lang="ts">
import { ref, onMounted } from 'vue'
import axios from 'axios'

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

// 获取所有定时任务
const fetchTasks = async () => {
    try {
        const res = await axios.get('/web_api/cron/tasks')
        tasks.value = res.data
    } catch (error: any) {
        showToast('获取定时任务失败: ' + error.message, 'error')
    }
}

// 创建或更新任务
const saveTask = async () => {
    try {
        const task = editingTask.value || newTask.value
        if (editingTask.value) {
            await axios.put(`/web_api/cron/tasks/${task.id}`, task)
            showToast('任务更新成功', 'success')
        } else {
            await axios.post('/web_api/cron/tasks', task)
            showToast('任务创建成功', 'success')
        }
        showModal.value = false
        await fetchTasks()
        resetForm()
    } catch (error: any) {
        showToast('保存任务失败: ' + error.message, 'error')
    }
}

// 删除任务
const deleteTask = async (id: number) => {
    if (!confirm('确定要删除这个任务吗？')) return
    try {
        await axios.delete(`/web_api/cron/tasks/${id}`)
        showToast('任务删除成功', 'success')
        await fetchTasks()
    } catch (error: any) {
        showToast('删除任务失败: ' + error.message, 'error')
    }
}

// 启用/禁用任务
const toggleTask = async (id: number) => {
    try {
        await axios.post(`/web_api/cron/tasks/${id}/toggle`)
        await fetchTasks()
    } catch (error: any) {
        showToast('切换任务状态失败: ' + error.message, 'error')
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
                                <button class="btn btn-sm btn-ghost text-error" @click="deleteTask(task.id)">
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
                <form @submit.prevent="saveTask">
                    <div class="form-control mb-4">
                        <label class="label">
                            <span class="label-text">任务名称</span>
                        </label>
                        <input type="text" v-model="(editingTask || newTask).name" class="input input-bordered"
                            required />
                    </div>

                    <div class="form-control mb-4">
                        <label class="label">
                            <span class="label-text">执行命令</span>
                        </label>
                        <input type="text" v-model="(editingTask || newTask).command" class="input input-bordered"
                            required />
                    </div>

                    <div class="form-control mb-4">
                        <label class="label">
                            <span class="label-text">定时表达式</span>
                            <span class="label-text-alt">
                                <a href="https://crontab.guru/" target="_blank" class="link">
                                    帮助
                                </a>
                            </span>
                        </label>
                        <input type="text" v-model="(editingTask || newTask).schedule" class="input input-bordered"
                            required placeholder="*/5 * * * * *" />
                    </div>

                    <div class="form-control mb-4">
                        <label class="label">
                            <span class="label-text">描述</span>
                        </label>
                        <textarea v-model="(editingTask || newTask).description" class="textarea textarea-bordered"
                            rows="3">
                        </textarea>
                    </div>

                    <div class="modal-action">
                        <button type="button" class="btn" @click="showModal = false">取消</button>
                        <button type="submit" class="btn btn-primary">保存</button>
                    </div>
                </form>
            </div>
            <form method="dialog" class="modal-backdrop">
                <button @click="showModal = false">关闭</button>
            </form>
        </dialog>

        <!-- Toast -->
        <div id="toast" class="toast hidden"></div>
    </div>
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

.toast-success {
    background-color: #10B981;
    color: white;
}

.toast-error {
    background-color: #EF4444;
    color: white;
}

.hidden {
    display: none;
}
</style>