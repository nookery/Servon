<script setup lang="ts">
import { ref } from 'vue'
import { fileAPI } from '../../api/file_api'

const props = defineProps<{
    show: boolean
    currentPath: string
}>()

const emit = defineEmits<{
    'update:show': [value: boolean]
    'created': []
}>()

const error = ref<string | null>(null)
const newFileName = ref('')
const newFileType = ref<'file' | 'directory'>('file')

async function createFile() {
    if (!newFileName.value) return
    try {
        const newPath = `${props.currentPath}/${newFileName.value}`.replace(/\/+/g, '/')
        await fileAPI.createFile(newPath, newFileType.value)
        emit('update:show', false)
        emit('created')
        resetForm()
    } catch (err: any) {
        const errorMessage = err.response?.data?.error || err.message || '创建文件失败'
        error.value = errorMessage
        setTimeout(() => {
            error.value = null
        }, 5000)
    }
}

function resetForm() {
    error.value = null
    newFileName.value = ''
    newFileType.value = 'file'
}

function handleClose() {
    emit('update:show', false)
    resetForm()
}
</script>

<template>
    <dialog class="modal" :class="{ 'modal-open': show }">
        <div class="modal-box">
            <h3 class="font-bold text-lg mb-4">新建文件/目录</h3>

            <div v-if="error" class="alert alert-error shadow-lg mb-4">
                <div>
                    <i class="ri-error-warning-line"></i>
                    <span>{{ error }}</span>
                </div>
            </div>

            <div class="form-control">
                <label class="label">
                    <span class="label-text">类型</span>
                </label>
                <select v-model="newFileType" class="select select-bordered">
                    <option value="file">文件</option>
                    <option value="directory">目录</option>
                </select>
            </div>
            <div class="form-control mt-4">
                <label class="label">
                    <span class="label-text">名称</span>
                </label>
                <input type="text" v-model="newFileName" class="input input-bordered" placeholder="输入名称">
            </div>
            <div class="modal-action">
                <button class="btn" @click="handleClose">取消</button>
                <button class="btn btn-primary" @click="createFile" :disabled="!newFileName">创建</button>
            </div>
        </div>
    </dialog>
</template>