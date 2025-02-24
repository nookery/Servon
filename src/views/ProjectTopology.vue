<script setup lang="ts">
import { ref } from 'vue'
import PageContainer from '../layouts/PageContainer.vue'
import GatewayConfigModal from '../components/GatewayConfigModal.vue'
import { RiGlobalLine, RiRouteLine, RiServerLine } from '@remixicon/vue'
import { topologyAPI } from '../api/topology'
import { useToast } from '../composables/useToast'

// 模拟数据
const projects = ref([
    { id: 1, name: 'Blog System', status: 'running', type: 'nodejs' },
    { id: 2, name: 'API Gateway', status: 'running', type: 'nginx' },
    { id: 3, name: 'User Service', status: 'stopped', type: 'php' },
    { id: 4, name: 'Data Analytics', status: 'running', type: 'python' },
])

const showGatewayConfig = ref(false)
const gatewayConfig = ref('')
const error = ref('')
const currentGateway = ref('')

const toast = useToast()

// 获取状态对应的样式
function getStatusStyle(status: string) {
    return {
        'running': 'bg-success/10 text-success border-success',
        'stopped': 'bg-error/10 text-error border-error',
    }[status] || 'bg-base-200 text-base-content border-base-300'
}

async function handleGatewayClick() {
    try {
        const gateways = await topologyAPI.getGateways()
        if (gateways.data.length === 0) {
            error.value = '当前系统中没有可用的网关软件'
            return
        }

        currentGateway.value = gateways.data[0]
        const config = await topologyAPI.getGatewayConfig(currentGateway.value)
        gatewayConfig.value = config.data.config
        showGatewayConfig.value = true
        error.value = ''
    } catch (err: any) {
        error.value = err.response?.data?.error || '获取网关配置失败'
        showGatewayConfig.value = false
    }
}

async function handleConfigSaved() {
    toast.success('配置已更新')
}
</script>

<template>
    <PageContainer title="项目拓扑">
        <div class="flex flex-col items-center gap-8 p-8">
            <!-- 互联网层 -->
            <div class="flex flex-col items-center gap-2">
                <div
                    class="w-24 h-24 rounded-full flex items-center justify-center bg-base-200 border-4 border-base-300">
                    <RiGlobalLine class="w-12 h-12" />
                </div>
                <div class="text-sm font-medium">Internet</div>
                <div class="h-8 w-0.5 bg-base-300"></div>
            </div>

            <!-- 网关层 -->
            <div class="flex flex-col items-center gap-2">
                <div class="w-32 h-16 rounded-lg flex items-center justify-center bg-primary/10 text-primary 
                    border-2 border-primary cursor-pointer hover:bg-primary/20" @click="handleGatewayClick">
                    <RiRouteLine class="w-6 h-6 mr-2" />
                    <span class="text-sm font-medium">Gateway</span>
                </div>
                <div class="h-8 w-0.5 bg-base-300"></div>
            </div>

            <!-- 项目层 -->
            <div class="grid grid-cols-2 gap-8 w-full max-w-2xl">
                <div v-for="project in projects" :key="project.id" class="flex flex-col items-center">
                    <!-- 连接线 -->
                    <div class="h-8 w-0.5 bg-base-300"></div>
                    <!-- 项目卡片 -->
                    <div class="w-full p-4 rounded-lg border-2 flex items-center gap-3"
                        :class="getStatusStyle(project.status)">
                        <RiServerLine class="w-5 h-5" />
                        <div class="flex-1">
                            <div class="font-medium">{{ project.name }}</div>
                            <div class="text-xs opacity-70">{{ project.type }}</div>
                        </div>
                        <!-- 状态指示器 -->
                        <div class="w-2 h-2 rounded-full"
                            :class="project.status === 'running' ? 'bg-success animate-pulse' : 'bg-error'">
                        </div>
                    </div>
                </div>
            </div>
        </div>

        <!-- 错误提示 -->
        <div v-if="error" class="alert alert-error shadow-lg mb-4">
            <span>{{ error }}</span>
        </div>

        <!-- 使用网关配置模态框组件 -->
        <GatewayConfigModal v-model="showGatewayConfig" :config="gatewayConfig" :gateway="currentGateway"
            @saved="handleConfigSaved" />
    </PageContainer>
</template>

<style scoped>
.animate-pulse {
    animation: pulse 2s cubic-bezier(0.4, 0, 0.6, 1) infinite;
}

@keyframes pulse {

    0%,
    100% {
        opacity: 1;
    }

    50% {
        opacity: 0.5;
    }
}
</style>