<script setup lang="ts">
import { ref, watch } from 'vue'
import { fileAPI } from '@/api/file_api'
import type { FileInfo } from '@/models/FileInfo'
import * as monaco from 'monaco-editor'
import { onMounted, onBeforeUnmount } from 'vue'

const props = defineProps<{
    show: boolean
    file: FileInfo | null
}>()

const emit = defineEmits<{
    'update:show': [value: boolean]
    'saved': []
}>()

const content = ref('')
const error = ref<string | null>(null)
const editorContainer = ref<HTMLElement | null>(null)
let editor: monaco.editor.IStandaloneCodeEditor | null = null

onMounted(() => {
    if (editorContainer.value) {
        editor = monaco.editor.create(editorContainer.value, {
            value: content.value,
            language: 'plaintext',
            theme: 'vs-dark',
            automaticLayout: true,
            minimap: { enabled: true },
            scrollBeyondLastLine: false,
        })

        // Sync editor content with content ref
        editor.onDidChangeModelContent(() => {
            content.value = editor?.getValue() || ''
        })
    }
})

onBeforeUnmount(() => {
    editor?.dispose()
})

watch(() => props.show, async (newVal) => {
    if (newVal && props.file) {
        await loadFileContent()
        editor?.setValue(content.value)
    }
})

async function loadFileContent() {
    if (!props.file) return
    try {
        const res = await fileAPI.getFileContent(props.file.path)
        content.value = res.data.content
        error.value = null
    } catch (err: any) {
        error.value = `打开文件失败: ${err.response?.data?.error || err.message || '未知错误'}`
    }
}

async function saveFile() {
    if (!props.file) return
    try {
        await fileAPI.saveFileContent(props.file.path, content.value)
        error.value = null
        emit('saved')
        emit('update:show', false)
    } catch (err: any) {
        error.value = `保存文件失败: ${err.response?.data?.error || err.message || '未知错误'}`
    }
}
</script>

<template>
    <dialog class="modal" :class="{ 'modal-open': show }">
        <div class="modal-box w-11/12 max-w-5xl">
            <h3 class="font-bold text-lg mb-4">编辑文件: {{ file?.name }}</h3>

            <div v-if="error" class="alert alert-error shadow-lg mb-4">
                <div>
                    <i class="ri-error-warning-line"></i>
                    <span>{{ error }}</span>
                </div>
            </div>

            <div ref="editorContainer" class="h-96"></div>

            <div class="modal-action">
                <button class="btn" @click="$emit('update:show', false)">取消</button>
                <button class="btn btn-primary" @click="saveFile">保存</button>
            </div>
        </div>
    </dialog>
</template>