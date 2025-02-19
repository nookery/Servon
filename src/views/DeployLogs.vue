<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import PageContainer from '../layouts/PageContainer.vue'
import ConfirmDialog from '../components/ConfirmDialog.vue'
import IconButton from '../components/IconButton.vue'
import DeployLogDetailDialog from '../components/DeployLogDetailDialog.vue'
import { type DeployLog, getDeployLogs, deleteDeployLog } from '../api/deploy_api'

const logs = ref<DeployLog[]>([])
const error = ref<string | null>(null)
const showDeleteConfirm = ref(false)
const logToDelete = ref<string | null>(null)
const selectedLog = ref<DeployLog | null>(null)
const showLogDetailDialog = ref(false)

const isEmpty = computed(() => logs.value.length === 0)

// 加载日志列表
async function loadLogs() {
    try {
        logs.value = await getDeployLogs()
        error.value = null
    } catch (err: any) {
        error.value = '获取部署日志失败: ' +
            (err.response?.data?.error || err.message || '未知错误')
    }
}

// 查看日志详情
function viewLogDetail(log: DeployLog) {
    selectedLog.value = log
    showLogDetailDialog.value = true
}

// 格式化时间
function formatDate(dateStr: string): string {
    if (!dateStr) return '未知'
    return new Date(dateStr).toLocaleString()
}

// 删除日志
const confirmDelete = (id: string) => {
    logToDelete.value = id
    showDeleteConfirm.value = true
}

const handleDelete = async () => {
    if (!logToDelete.value) return

    try {
        await deleteDeployLog(logToDelete.value)
        await loadLogs()
    } catch (err: any) {
        error.value = '删除日志失败: ' +
            (err.response?.data?.error || err.message || '未知错误')
    } finally {
        showDeleteConfirm.value = false
        logToDelete.value = null
    }
}

onMounted(() => {
    loadLogs()
})
</script>

<template>
    <PageContainer title="部署日志">
        <template #header>
            <div v-if="error" class="alert alert-error mb-4">
                {{ error }}
            </div>
        </template>

        <!-- 空状态显示 -->
        <div v-if="isEmpty" class="flex flex-col items-center justify-center py-16 text-base-content/60">
            <i class="ri-file-list-3-line text-4xl mb-4"></i>
            <p class="text-lg">暂无部署日志</p>
            <p class="text-sm mt-2">当有新的部署任务完成后，相关日志将会显示在这里</p>
        </div>

        <!-- 日志列表 -->
        <div v-else class="overflow-x-auto">
            <table class="table table-zebra w-full">
                <thead>
                    <tr>
                        <th>ID</th>
                        <th>状态</th>
                        <th>消息</th>
                        <th>时间</th>
                        <th>操作</th>
                    </tr>
                </thead>
                <tbody>
                    <tr v-for="log in logs" :key="log.id">
                        <td>{{ log.id }}</td>
                        <td>
                            <span class="badge" :class="{
                                'badge-success': log.status === 'success',
                                'badge-error': log.status === 'error',
                                'badge-warning': log.status === 'running'
                            }">
                                {{ log.status }}
                            </span>
                        </td>
                        <td class="whitespace-pre-wrap">{{ log.message }}</td>
                        <td>{{ formatDate(log.timestamp) }}</td>
                        <td>
                            <div class="flex gap-2">
                                <IconButton icon="ri-file-text-line" type="primary" size="sm"
                                    @click="viewLogDetail(log)" />
                                <IconButton icon="ri-delete-bin-line" type="error" size="sm"
                                    @click="confirmDelete(log.id)" />
                            </div>
                        </td>
                    </tr>
                </tbody>
            </table>
        </div>

        <!-- 确认删除对话框 -->
        <ConfirmDialog v-model:show="showDeleteConfirm" title="确认删除" message="该操作无法撤销，是否确认删除此日志？" type="warning"
            confirm-text="删除" @confirm="handleDelete" />

        <!-- 日志详情对话框 -->
        <DeployLogDetailDialog v-model:show="showLogDetailDialog" :log="selectedLog" />
    </PageContainer>
</template>