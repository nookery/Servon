<script setup lang="ts">
import { defineProps, ref, onMounted, watch, onBeforeUnmount } from 'vue'
import IconButton from '../../components/IconButton.vue'
import * as monaco from 'monaco-editor'

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

// 添加 editor 实例引用
const editorContainer = ref<HTMLElement | null>(null)
let editor: any = null

onMounted(() => {
    if (editorContainer.value) {
        editor = monaco.editor.create(editorContainer.value, {
            value: props.webhook?.payload ? JSON.stringify(props.webhook.payload, null, 2) : '',
            language: 'json',
            theme: 'vs-dark',
            readOnly: true,
            minimap: { enabled: false },
            automaticLayout: true,
        })
    }
})

// 添加 watch 来监听 webhook 变化
watch(
    () => props.webhook?.payload,
    (newPayload) => {
        if (editor) {
            editor.setValue(newPayload ? JSON.stringify(newPayload, null, 2) : '')
        }
    }
)

// 组件卸载时清理 editor
onBeforeUnmount(() => {
    if (editor) {
        editor.dispose()
    }
})

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
                        <IconButton icon="ri-file-copy-line" @click="copyPayload">复制</IconButton>
                        <IconButton icon="ri-close-line" circle @click="$emit('update:show', false)" />
                    </div>
                </div>
            </div>
            <div class="h-full overflow-y-auto">
                <div ref="editorContainer" style="height: 100%;"></div>
            </div>
        </div>
    </dialog>
</template>