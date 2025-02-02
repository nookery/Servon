<!-- SoftwareList.vue -->
<template>
    <div>
        <!-- 加载状态 -->
        <div v-if="loading" class="text-gray-600">加载中...</div>
        <div v-if="error" class="text-red-600">{{ error }}</div>

        <!-- 软件列表 -->
        <div v-if="!loading && !error">
            <div class="space-y-4">
                <div v-if="softwareList.length === 0" class="text-center py-8 text-gray-500">
                    暂无可用软件
                </div>
                <SoftwareCard v-for="software in softwareList" :key="software.name" :software="software" />
            </div>

            <!-- 原始响应显示 -->
            <div class="mt-4">
                <button class="px-4 py-2 bg-gray-100 hover:bg-gray-200 rounded-md text-sm" @click="showRaw = !showRaw">
                    {{ showRaw ? '隐藏原始响应' : '显示原始响应' }}
                </button>
                <pre v-if="showRaw" class="mt-4 p-4 bg-gray-50 rounded-md overflow-auto">
                    {{ JSON.stringify(softwareList, null, 2) }}
                </pre>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import SoftwareCard from './SoftwareCard.vue'

interface Software {
    name: string
    version?: string
    status: string
    path?: string
    description?: string
}

const loading = ref(true)
const error = ref<string | null>(null)
const softwareList = ref<Software[]>([])
const showRaw = ref(false)

async function fetchSoftwareStatus(name: string): Promise<Software> {
    const response = await fetch(`/api/system/software/${name}/status`)
    if (!response.ok) {
        throw new Error(`获取软件 ${name} 状态失败`)
    }
    return await response.json()
}

async function fetchSoftwareList() {
    try {
        loading.value = true
        error.value = null

        // 获取软件名称列表
        const response = await fetch('/api/system/software')
        if (!response.ok) {
            throw new Error('获取软件列表失败')
        }
        const names: string[] = await response.json()

        // 获取每个软件的详细信息
        const softwareDetails = await Promise.all(
            names.map(name => fetchSoftwareStatus(name))
        )

        softwareList.value = softwareDetails
    } catch (err) {
        error.value = err instanceof Error ? err.message : '未知错误'
    } finally {
        loading.value = false
    }
}

onMounted(fetchSoftwareList)
</script>