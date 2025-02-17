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
            language: 'json',
            theme: 'vs-dark',
            automaticLayout: true,
            minimap: { enabled: true },
            scrollBeyondLastLine: false,
            formatOnPaste: true,
        })

        editor.addCommand(monaco.KeyMod.Shift | monaco.KeyMod.Alt | monaco.KeyCode.KeyF, () => {
            formatJson()
        })

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

const formatJson = () => {
    if (!editor) return

    try {
        const content = editor.getValue()
        const formatted = JSON.stringify(JSON.parse(content), null, 2)
        editor.setValue(formatted)
    } catch (e) {
        error.value = '无效的 JSON 格式'
    }
}
</script>

<template>
    <dialog class="modal" :class="{ 'modal-open': show }">
        <div class="modal-box w-11/12 max-w-5xl">
            <h3 class="font-bold text-lg">编辑文件: {{ file?.name }}</h3>

            <!-- Toolbar -->
            <div class="flex justify-between items-center py-2 mb-4 border-b">
                <div class="flex gap-2">
                    <button class="btn btn-sm" @click="formatJson" title="格式化 JSON (Shift+Alt+F)">
                        格式化
                    </button>
                </div>
                <div class="flex gap-2">
                    <button class="btn btn-sm" @click="$emit('update:show', false)">取消</button>
                    <button class="btn btn-sm btn-primary" @click="saveFile">保存</button>
                </div>
            </div>

            <div v-if="error" class="alert alert-error shadow-lg mb-4">
                <div>
                    <i class="ri-error-warning-line"></i>
                    <span>{{ error }}</span>
                </div>
            </div>

            <div ref="editorContainer" class="h-96"></div>
        </div>
    </dialog>
</template>