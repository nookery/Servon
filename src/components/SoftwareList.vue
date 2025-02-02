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

<script setup>
import { ref, onMounted } from 'vue'
import SoftwareCard from './SoftwareCard.vue'

const loading = ref(true)
const error = ref(null)
const softwareList = ref([])
const showRaw = ref(false)

async function fetchSoftwareList() {
    try {
        loading.value = true
        error.value = null
        const response = await fetch('/api/system/software')
        if (!response.ok) {
            throw new Error('获取软件列表失败')
        }
        softwareList.value = await response.json()
    } catch (err) {
        error.value = err.message
    } finally {
        loading.value = false
    }
}

onMounted(fetchSoftwareList)
</script>