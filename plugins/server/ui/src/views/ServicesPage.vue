<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { RiPlayCircleLine, RiSettings4Line, RiFileListLine } from '@remixicon/vue'
import PageContainer from '../layouts/PageContainer.vue'
import ServiceList from '../components/services/ServiceList.vue'
import ServiceConfig from '../components/services/ServiceConfig.vue'
import ServiceLogs from '../components/services/ServiceLogs.vue'

const route = useRoute()
const router = useRouter()

// 本地存储键名
const STORAGE_KEY = 'servicesLastService'

// 当前选中的服务
const currentService = ref(route.query.service as string || '')

// 当前激活的 Tab
const activeTab = ref(route.query.tab as string || 'list')

// 定义标签页
const tabs = [
    { key: 'list', title: '服务列表', icon: RiPlayCircleLine },
    { key: 'config', title: '服务配置', icon: RiSettings4Line },
    { key: 'logs', title: '服务日志', icon: RiFileListLine }
]

// 监听服务变化，更新 URL 和本地存储
watch(() => currentService.value, (newService) => {
    // 更新 URL
    router.replace({
        query: { ...route.query, service: newService }
    })

    // 保存到本地存储
    if (newService) {
        localStorage.setItem(STORAGE_KEY, newService)
    }
})

// 监听 Tab 变化，更新 URL
watch(() => activeTab.value, (newTab) => {
    router.replace({
        query: { ...route.query, tab: newTab }
    })
})

// 初始化时加载服务和 Tab
onMounted(() => {
    // 加载 Tab
    if (route.query.tab) {
        activeTab.value = route.query.tab as string
    }

    // 优先级：URL 参数 > 本地存储
    if (route.query.service) {
        // 如果 URL 中有服务参数，使用它
        currentService.value = route.query.service as string
    } else {
        // 尝试从本地存储获取上次访问的服务
        const savedService = localStorage.getItem(STORAGE_KEY)
        if (savedService) {
            currentService.value = savedService
        }
    }
})

// 处理服务选择
const handleServiceSelect = (serviceName: string) => {
    currentService.value = serviceName
    // 如果当前不在列表页，且选择了新服务，自动切换到配置页
    if (activeTab.value !== 'list' && serviceName) {
        activeTab.value = 'config'
    }
}
</script>

<template>
    <PageContainer title="服务管理" :tabs="tabs" v-model="activeTab">
        <!-- 服务列表 Tab -->
        <template #list>
            <ServiceList :current-service="currentService" @select-service="handleServiceSelect" />
        </template>

        <!-- 服务配置 Tab -->
        <template #config>
            <ServiceConfig :service-name="currentService" v-if="currentService" />
            <div v-else class="empty-state">
                <p>请先从服务列表中选择一个服务</p>
            </div>
        </template>

        <!-- 服务日志 Tab -->
        <template #logs>
            <ServiceLogs :service-name="currentService" v-if="currentService" />
            <div v-else class="empty-state">
                <p>请先从服务列表中选择一个服务</p>
            </div>
        </template>
    </PageContainer>
</template>

<style scoped>
.empty-state {
    display: flex;
    justify-content: center;
    align-items: center;
    height: 200px;
    color: #666;
    background-color: #f9f9f9;
    border-radius: 8px;
    margin: 20px;
}
</style>