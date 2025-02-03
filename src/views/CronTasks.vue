<script setup lang="ts">
import { ref, onMounted } from 'vue'
import axios from 'axios'
import ConfirmDialog from '../components/ConfirmDialog.vue'
import { useToast } from '../composables/useToast'

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

interface ValidationError {
    field: string
    message: string
}

interface ValidationErrors {
    errors: ValidationError[]
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

// 添加表单错误信息的状态
const formError = ref('')

// 添加字段错误信息的状态
const fieldErrors = ref<Record<string, string>>({})

// 添加删除确认对话框的状态
const showDeleteConfirm = ref(false)
const taskToDelete = ref<number | null>(null)

// 添加 showCronHelp 状态
const showCronHelp = ref(false)

const toast = useToast()

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
        const errorMessage = error.response?.data?.error || error.message || '未知错误'
        toast.error('获取定时任务失败: ' + errorMessage)
    }
}

// 清除错误信息
const clearErrors = () => {
    formError.value = ''
    fieldErrors.value = {}
}

// 创建或更新任务
const saveTask = async () => {
    try {
        clearErrors()
        const task = editingTask.value || newTask.value
        if (editingTask.value) {
            await axios.put(`/web_api/cron/tasks/${task.id}`, task)
            toast.success('任务更新成功')
            showModal.value = false
            await fetchTasks()
            resetForm()
        } else {
            await axios.post('/web_api/cron/tasks', task)
            toast.success('任务创建成功')
            showModal.value = false
            await fetchTasks()
            resetForm()
        }
    } catch (error: any) {
        if (error.response?.data?.errors) {
            // 处理字段验证错误
            const validationErrors = error.response.data as ValidationErrors
            validationErrors.errors.forEach(err => {
                if (err.field === 'general') {
                    // 通用错误显示在表单底部
                    formError.value = err.message
                } else {
                    // 字段错误显示在对应字段下方
                    fieldErrors.value[err.field] = err.message
                }
            })
        } else {
            // 处理其他错误
            const errorMessage = error.response?.data?.error || error.message || '未知错误'
            formError.value = errorMessage
        }
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
        toast.success('任务删除成功')
        await fetchTasks()
    } catch (error: any) {
        const errorMessage = error.response?.data?.error || error.message || '未知错误'
        toast.error('删除任务失败: ' + errorMessage)
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
        toast.error('切换任务状态失败: ' + errorMessage)
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

const handleConfirmDelete = () => {
    handleDelete()
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
                <form @submit.prevent="saveTask">
                    <div class="form-control mb-4">
                        <label class="label">
                            <span class="label-text">任务名称</span>
                        </label>
                        <input type="text" v-model="(editingTask || newTask).name" class="input input-bordered"
                            :class="{ 'input-error': fieldErrors.name }" required />
                        <label v-if="fieldErrors.name" class="label">
                            <span class="label-text-alt text-error">{{ fieldErrors.name }}</span>
                        </label>
                    </div>

                    <div class="form-control mb-4">
                        <label class="label">
                            <span class="label-text">执行命令</span>
                        </label>
                        <input type="text" v-model="(editingTask || newTask).command" class="input input-bordered"
                            :class="{ 'input-error': fieldErrors.command }" required />
                        <label v-if="fieldErrors.command" class="label">
                            <span class="label-text-alt text-error">{{ fieldErrors.command }}</span>
                        </label>
                    </div>

                    <div class="form-control mb-4">
                        <label class="label">
                            <span class="label-text">定时表达式</span>
                            <span class="label-text-alt">
                                <a href="#" @click.prevent="showCronHelp = true" class="link">
                                    帮助
                                </a>
                            </span>
                        </label>
                        <input type="text" v-model="(editingTask || newTask).schedule" class="input input-bordered"
                            :class="{ 'input-error': fieldErrors.schedule }" required placeholder="0 */5 * * * *" />
                        <div class="flex flex-col gap-1 mt-1">
                            <label class="label py-0">
                                <span class="label-text-alt text-base-content/70">格式: 秒 分 时 日 月 星期</span>
                            </label>
                            <label v-if="fieldErrors.schedule" class="label py-0">
                                <span class="label-text-alt text-error">{{ fieldErrors.schedule }}</span>
                            </label>
                        </div>
                    </div>

                    <div class="form-control mb-4">
                        <label class="label">
                            <span class="label-text">描述</span>
                        </label>
                        <textarea v-model="(editingTask || newTask).description" class="textarea textarea-bordered"
                            rows="3">
                        </textarea>
                    </div>

                    <!-- 通用错误信息显示区域 -->
                    <div v-if="formError" class="alert alert-error mb-4">
                        <i class="ri-error-warning-line"></i>
                        <span>{{ formError }}</span>
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

        <!-- Cron 帮助对话框 -->
        <dialog class="modal" :class="{ 'modal-open': showCronHelp }">
            <div class="modal-box">
                <h3 class="font-bold text-lg mb-4">Cron 表达式帮助</h3>
                <div class="space-y-4">
                    <p>格式：秒 分 时 日 月 星期</p>
                    <div>
                        <h4 class="font-semibold mb-2">常用示例：</h4>
                        <ul class="list-disc list-inside space-y-2">
                            <li><code>0 0 * * * *</code> - 每小时执行（在整点时）</li>
                            <li><code>0 */5 * * * *</code> - 每5分钟执行</li>
                            <li><code>0 0 0 * * *</code> - 每天午夜执行</li>
                            <li><code>*/10 * * * * *</code> - 每10秒执行一次</li>
                            <li><code>0 30 9 * * 1-5</code> - 工作日上午9:30执行</li>
                        </ul>
                    </div>
                </div>
                <div class="modal-action">
                    <button class="btn" @click="showCronHelp = false">关闭</button>
                </div>
            </div>
            <form method="dialog" class="modal-backdrop">
                <button @click="showCronHelp = false">关闭</button>
            </form>
        </dialog>

        <!-- 使用确认对话框组件 -->
        <ConfirmDialog v-model:show="showDeleteConfirm" title="确认删除" message="该操作无法撤销，是否确认删除此任务？" type="warning"
            confirm-text="删除" @confirm="handleConfirmDelete" />
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

.alert {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    padding: 1rem;
    border-radius: 0.5rem;
}

.alert-error {
    background-color: rgb(254, 242, 242);
    color: rgb(153, 27, 27);
    border: 1px solid rgb(252, 165, 165);
}

.input-error {
    border-color: rgb(252, 165, 165);
}

.label-text-alt.text-error {
    color: rgb(153, 27, 27);
}
</style>