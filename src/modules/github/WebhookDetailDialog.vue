<script setup lang="ts">
import { ref } from 'vue'
import IconButton from '../../components/IconButton.vue'
import ParsedWebhookView from './ParsedWebhookView.vue'
import RawWebhookView from './RawWebhookView.vue'

interface WebhookData {
    id: string
    type: string
    timestamp: number
    payload: any
}

const props = defineProps<{
    webhook: WebhookData | null
    show: boolean
}>()

const viewMode = ref<'raw' | 'parsed'>('parsed')

async function copyPayload() {
    if (props.webhook?.payload) {
        try {
            await navigator.clipboard.writeText(JSON.stringify(props.webhook.payload, null, 2))
        } catch (err) {
            console.error('复制失败:', err)
        }
    }
}
</script>

<template>
    <dialog class="modal" :class="{ 'modal-open': show }">
        <div class="modal-box w-11/12 max-w-5xl p-0 h-4/5">
            <div class="sticky top-0 bg-base-100 z-10 px-4 py-1">
                <div class="flex items-center justify-between">
                    <h3 class="font-bold text-lg">
                        Webhook详情
                        <div class="badge badge-primary ml-2">{{ webhook?.type }}</div>
                    </h3>
                    <div class="flex gap-2">
                        <IconButton :icon="viewMode === 'raw' ? 'ri-code-line' : 'ri-file-list-line'"
                            @click="viewMode = viewMode === 'raw' ? 'parsed' : 'raw'">
                            {{ viewMode === 'raw' ? '原始数据' : '解析视图' }}
                        </IconButton>
                        <IconButton icon="ri-file-copy-line" @click="copyPayload">复制</IconButton>
                        <IconButton icon="ri-close-line" circle @click="$emit('update:show', false)" />
                    </div>
                </div>
            </div>
            <div class="h-full overflow-y-auto">
                <template v-if="viewMode === 'raw'">
                    <RawWebhookView :webhook="webhook" />
                </template>
                <template v-else>
                    <ParsedWebhookView :webhook="webhook" />
                </template>
            </div>
        </div>
    </dialog>
</template>