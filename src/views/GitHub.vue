<script setup lang="ts">
import { ref, onMounted } from 'vue'
import PageContainer from '../layouts/PageContainer.vue'
import { githubAPI } from '../api/github'

interface WebhookData {
    id: string
    type: string
    timestamp: number
    payload: any
}

const webhooks = ref<WebhookData[]>([])
const error = ref<string | null>(null)
const selectedWebhook = ref<WebhookData | null>(null)
const showPayloadDialog = ref(false)

async function loadWebhooks() {
    try {
        const res = await githubAPI.getWebhooks()
        webhooks.value = res.data.sort((a, b) => b.timestamp - a.timestamp)
        error.value = null
    } catch (err: any) {
        error.value = `获取GitHub事件列表失败: ${err.response?.data?.error || err.message || '未知错误'}`
    }
}

function viewPayload(webhook: WebhookData) {
    selectedWebhook.value = webhook
    showPayloadDialog.value = true
}

function formatTimestamp(timestamp: number) {
    return new Date(timestamp * 1000).toLocaleString()
}

onMounted(() => {
    loadWebhooks()
})
</script>

<template>
    <PageContainer title="GitHub管理">
        <template #header>
            <div v-if="error" class="alert alert-error shadow-lg mb-4">
                <div>
                    <i class="ri-error-warning-line"></i>
                    <span>{{ error }}</span>
                </div>
            </div>

            <div class="flex gap-2 mb-4">
                <button class="btn" @click="loadWebhooks">
                    <i class="ri-refresh-line mr-1"></i>刷新
                </button>
            </div>
        </template>

        <div v-if="webhooks.length === 0" class="flex flex-col items-center justify-center py-12 text-base-content/60">
            <i class="ri-github-fill text-6xl mb-4"></i>
            <p class="text-lg mb-2">暂无GitHub事件</p>
            <p class="text-sm">当收到GitHub的webhook回调时，这里会显示相关事件信息</p>
        </div>

        <div v-else class="overflow-x-auto">
            <table class="table table-zebra w-full">
                <thead>
                    <tr>
                        <th>时间</th>
                        <th>事件类型</th>
                        <th>事件ID</th>
                        <th>操作</th>
                    </tr>
                </thead>
                <tbody>
                    <tr v-for="webhook in webhooks" :key="webhook.id">
                        <td>{{ formatTimestamp(webhook.timestamp) }}</td>
                        <td>
                            <div class="badge" :class="{
                                'badge-primary': webhook.type === 'push',
                                'badge-secondary': webhook.type === 'pull_request',
                                'badge-accent': webhook.type === 'installation',
                            }">
                                {{ webhook.type }}
                            </div>
                        </td>
                        <td>{{ webhook.id }}</td>
                        <td>
                            <button class="btn btn-xs" @click="viewPayload(webhook)">
                                <i class="ri-file-text-line mr-1"></i>查看详情
                            </button>
                        </td>
                    </tr>
                </tbody>
            </table>
        </div>

        <dialog class="modal" :class="{ 'modal-open': showPayloadDialog }">
            <div class="modal-box w-11/12 max-w-5xl">
                <h3 class="font-bold text-lg mb-4">
                    Webhook详情
                    <div class="badge badge-primary ml-2">{{ selectedWebhook?.type }}</div>
                </h3>
                <div class="mockup-code">
                    <pre><code>{{ JSON.stringify(selectedWebhook?.payload, null, 2) }}</code></pre>
                </div>
                <div class="modal-action">
                    <button class="btn" @click="showPayloadDialog = false">关闭</button>
                </div>
            </div>
        </dialog>
    </PageContainer>
</template>