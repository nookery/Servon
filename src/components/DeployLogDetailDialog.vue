<script setup lang="ts">
import { ref, onMounted, watch, onBeforeUnmount } from 'vue'
import * as monaco from 'monaco-editor'
import { type DeployLog } from '../api/deploy_api'

const props = defineProps<{
    show: boolean
    log: DeployLog | null
}>()

const emit = defineEmits<{
    'update:show': [value: boolean]
}>()

const editorContainer = ref<HTMLElement | null>(null)
let editor: any = null

function formatDate(dateStr: string): string {
    if (!dateStr) return '未知'
    return new Date(dateStr).toLocaleString()
}

onMounted(() => {
    if (editorContainer.value) {
        editor = monaco.editor.create(editorContainer.value, {
            value: props.log?.message || '',
            language: 'plaintext',
            theme: 'vs-dark',
            readOnly: true,
            minimap: { enabled: false },
            automaticLayout: true,
            wordWrap: 'on',
        })
    }
})

watch(
    () => props.log?.message,
    (newMessage) => {
        if (editor) {
            editor.setValue(newMessage || '')
        }
    }
)

watch(
    () => props.show,
    (newShow) => {
        if (newShow && editor) {
            // 当对话框显示时，触发编辑器布局更新
            setTimeout(() => {
                editor.layout()
            }, 0)
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
    <dialog :class="{ 'modal': true, 'modal-open': show }">
        <div class="modal-box w-11/12 max-w-5xl h-[80vh]">
            <h3 class="font-bold text-lg mb-4">部署日志详情</h3>
            <div class="mb-4 h-[calc(100%-8rem)]">
                <div class="grid grid-cols-2 gap-4 mb-4">
                    <div>
                        <span class="font-bold">ID:</span> {{ log?.id }}
                    </div>
                    <div>
                        <span class="font-bold">时间:</span> {{ formatDate(log?.timestamp || '') }}
                    </div>
                </div>
                <div class="h-[calc(100%-4rem)]" ref="editorContainer"></div>
            </div>
            <div class="modal-action">
                <button class="btn" @click="$emit('update:show', false)">关闭</button>
            </div>
        </div>
        <form method="dialog" class="modal-backdrop">
            <button @click="$emit('update:show', false)">关闭</button>
        </form>
    </dialog>
</template>