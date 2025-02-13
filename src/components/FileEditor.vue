<script setup lang="ts">
import { ref, watch } from 'vue'
import axios from 'axios'

interface FileInfo {
    name: string
    path: string
    size: number
    isDir: boolean
    mode: string
    modTime: string
    owner: string
    group: string
}

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

watch(() => props.show, async (newVal) => {
    if (newVal && props.file) {
        await loadFileContent()
    }
})

async function loadFileContent() {
    if (!props.file) return
    try {
        const res = await axios.get(`/web_api/system/files/content?path=${props.file.path}`)
        content.value = res.data.content
        error.value = null
    } catch (err: any) {
        error.value = `打开文件失败: ${err.response?.data?.error || err.message || '未知错误'}`
    }
}

async function saveFile() {
    if (!props.file) return
    try {
        await axios.post('/web_api/system/files/save', {
            path: props.file.path,
            content: content.value
        })
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

            <textarea 
                v-model="content" 
                class="textarea textarea-bordered w-full h-96 font-mono"
            ></textarea>

            <div class="modal-action">
                <button class="btn" @click="$emit('update:show', false)">取消</button>
                <button class="btn btn-primary" @click="saveFile">保存</button>
            </div>
        </div>
    </dialog>
</template> 