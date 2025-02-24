<script setup lang="ts">
import { ref } from 'vue'
import FileEditor from './files/FileEditor.vue'
import type { FileInfo } from '../types/FileInfo'
import { topologyAPI } from '../api/topology'
import { useToast } from '../composables/useToast'

const props = defineProps<{
    modelValue: boolean
    config: string
    gateway: string
}>()

const emit = defineEmits<{
    'update:modelValue': [value: boolean]
    'saved': []
}>()

const toast = useToast()

// 创建一个虚拟的 FileInfo 对象
const fileInfo = ref<FileInfo>({
    name: 'Caddyfile',
    path: `gateways/${props.gateway}/config`,
    size: props.config.length,
    isDir: false,
    modTime: new Date().toISOString(),
    mode: '0644',
    owner: 'root',
    group: 'root'
})

// 处理保存事件
async function handleSaved() {
    try {
        await topologyAPI.setGatewayConfig(props.gateway, props.config)
        toast.success('配置保存成功')
        emit('saved')
        emit('update:modelValue', false)
    } catch (err: any) {
        toast.error(err.response?.data?.error || '保存配置失败')
    }
}
</script>

<template>
    <FileEditor :show="modelValue" @update:show="(val) => emit('update:modelValue', val)" :file="fileInfo"
        :initial-content="config" @saved="handleSaved" />
</template>