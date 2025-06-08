<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { RiSaveLine, RiRefreshLine } from '@remixicon/vue'
import IconButton from '../IconButton.vue'
import { getServiceConfig, updateServiceConfig } from '../../api/services'

const props = defineProps<{
    serviceName: string
}>()

// 服务配置
const config = ref<Record<string, any>>({})
const rawConfig = ref('')
const loading = ref(false)
const saving = ref(false)
const showRawConfig = ref(false)

// 获取服务配置
const fetchServiceConfig = async () => {
    if (!props.serviceName) return

    loading.value = true
    try {
        const response = await getServiceConfig(props.serviceName)
        config.value = response
        rawConfig.value = response.raw_config || ''
    } catch (error) {
        console.error(`获取服务 ${props.serviceName} 配置失败:`, error)
    } finally {
        loading.value = false
    }
}

// 保存服务配置
const saveServiceConfig = async () => {
    if (!props.serviceName) return

    saving.value = true
    try {
        const configToSave = showRawConfig.value
            ? { raw_config: rawConfig.value }
            : config.value

        await updateServiceConfig(props.serviceName, configToSave)
        // 重新加载配置以确认更改
        await fetchServiceConfig()
    } catch (error) {
        console.error(`保存服务 ${props.serviceName} 配置失败:`, error)
    } finally {
        saving.value = false
    }
}

// 当服务名称变化时重新加载配置
watch(() => props.serviceName, (newServiceName) => {
    if (newServiceName) {
        fetchServiceConfig()
    }
})

onMounted(() => {
    if (props.serviceName) {
        fetchServiceConfig()
    }
})
</script>

<template>
    <div class="service-config-container">
        <div class="flex justify-between items-center mb-4">
            <h2 class="text-xl font-bold">{{ serviceName }} 配置</h2>
            <div class="flex gap-2">
                <div class="form-control">
                    <label class="label cursor-pointer">
                        <span class="label-text mr-2">原始配置</span>
                        <input type="checkbox" class="toggle toggle-primary" v-model="showRawConfig" />
                    </label>
                </div>

                <IconButton icon="refresh" :loading="loading" @click="fetchServiceConfig" tooltip="刷新配置">
                    <RiRefreshLine />
                </IconButton>

                <IconButton icon="save" color="primary" :loading="saving" @click="saveServiceConfig" tooltip="保存配置">
                    <RiSaveLine />
                </IconButton>
            </div>
        </div>

        <div v-if="loading" class="loading-container">
            <span class="loading loading-spinner loading-lg"></span>
            <p>加载配置中...</p>
        </div>

        <div v-else>
            <!-- 原始配置视图 -->
            <div v-if="showRawConfig" class="raw-config-editor">
                <textarea v-model="rawConfig" class="textarea textarea-bordered w-full h-96 font-mono text-sm"
                    placeholder="服务配置内容"></textarea>
            </div>

            <!-- 结构化配置视图 -->
            <div v-else class="structured-config">
                <div v-if="Object.keys(config).length === 0" class="empty-state">
                    <p>没有可用的配置项</p>
                </div>

                <div v-else class="grid grid-cols-1 gap-4">
                    <div v-for="(value, key) in config" :key="key" class="form-control">
                        <!-- 跳过原始配置字段 -->
                        <template v-if="key !== 'raw_config'">
                            <label class="label">
                                <span class="label-text font-medium">{{ key }}</span>
                            </label>

                            <!-- 根据值类型渲染不同的输入控件 -->
                            <input v-if="typeof value === 'string' || typeof value === 'number'" v-model="config[key]"
                                type="text" class="input input-bordered w-full" />

                            <select v-else-if="typeof value === 'boolean'" v-model="config[key]"
                                class="select select-bordered w-full">
                                <option :value="true">是</option>
                                <option :value="false">否</option>
                            </select>

                            <!-- 对于复杂类型，显示为只读 -->
                            <div v-else class="p-2 bg-base-200 rounded-md overflow-auto max-h-32">
                                <pre>{{ JSON.stringify(value, null, 2) }}</pre>
                            </div>
                        </template>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<style scoped>
.service-config-container {
    padding: 1rem;
}

.loading-container,
.empty-state {
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;
    height: 200px;
    color: #666;
    background-color: #f9f9f9;
    border-radius: 8px;
    margin: 20px 0;
}
</style>