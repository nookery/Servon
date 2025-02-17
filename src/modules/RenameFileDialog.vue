<script setup lang="ts">
import { ref, watch } from 'vue'
import type { FileInfo } from '../models/FileInfo'

const props = defineProps<{
    show: boolean
    file: FileInfo | null
}>()

const emit = defineEmits<{
    (e: 'update:show', value: boolean): void
    (e: 'rename', oldPath: string, newPath: string): void
}>()

const newBaseName = ref('')
const newExtension = ref('')
const error = ref<string | null>(null)

// 当文件改变时，初始化表单数据
watch(() => props.file, (file) => {
    if (file) {
        const lastDotIndex = file.name.lastIndexOf('.')
        if (lastDotIndex > 0 && !file.isDir) {
            newBaseName.value = file.name.substring(0, lastDotIndex)
            newExtension.value = file.name.substring(lastDotIndex + 1)
        } else {
            newBaseName.value = file.name
            newExtension.value = ''
        }
    }
}, { immediate: true })

const handleRename = () => {
    if (!props.file || !newBaseName.value) return

    const newFileName = props.file.isDir || !newExtension.value
        ? newBaseName.value
        : `${newBaseName.value}.${newExtension.value}`

    emit('rename', props.file.path, `${props.file.path}/${newFileName}`.replace(/\/+/g, '/'))
}

const closeDialog = () => {
    emit('update:show', false)
    error.value = null
    newBaseName.value = ''
    newExtension.value = ''
}
</script>

<template>
    <dialog class="modal" :class="{ 'modal-open': show }">
        <div class="modal-box">
            <h3 class="font-bold text-lg mb-4">重命名{{ file?.isDir ? '目录' : '文件' }}</h3>

            <div v-if="error" class="alert alert-error shadow-lg mb-4">
                <div>
                    <i class="ri-error-warning-line"></i>
                    <span>{{ error }}</span>
                </div>
            </div>

            <div class="form-control mt-4">
                <label class="label">
                    <span class="label-text">名称</span>
                </label>
                <input type="text" v-model="newBaseName" class="input input-bordered" placeholder="输入名称"
                    :class="{ 'rounded-r-none': !file?.isDir }">
            </div>

            <div v-if="!file?.isDir" class="form-control mt-2">
                <label class="label">
                    <span class="label-text">扩展名</span>
                </label>
                <div class="join">
                    <span class="join-item flex items-center px-4 bg-base-200 border border-base-300">.</span>
                    <input type="text" v-model="newExtension" class="input input-bordered join-item"
                        placeholder="输入扩展名">
                </div>
            </div>

            <div class="mt-2 text-sm text-base-content/70">
                新文件名: {{ file?.isDir || !newExtension ? newBaseName : `${newBaseName}.${newExtension}` }}
            </div>

            <div class="modal-action">
                <button class="btn" @click="closeDialog">取消</button>
                <button class="btn btn-primary" @click="handleRename"
                    :disabled="!newBaseName || (!file?.isDir && !newExtension)">确定</button>
            </div>
        </div>
    </dialog>
</template>