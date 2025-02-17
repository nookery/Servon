<script setup lang="ts">
import { ref, onMounted, watch, onBeforeUnmount } from 'vue'
import * as monaco from 'monaco-editor'

interface WebhookData {
    id: string
    type: string
    timestamp: number
    payload: any
}

const props = defineProps<{
    webhook: WebhookData | null
}>()

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

watch(
    () => props.webhook?.payload,
    (newPayload) => {
        if (editor) {
            editor.setValue(newPayload ? JSON.stringify(newPayload, null, 2) : '')
        }
    }
)

onBeforeUnmount(() => {
    if (editor) {
        editor.dispose()
    }
})
</script>

<template>
    <div ref="editorContainer" style="height: 100%;"></div>
</template>