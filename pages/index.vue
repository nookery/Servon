<template>
    <div class="min-h-screen bg-gray-50 p-8">
        <UContainer>
            <h1 class="text-2xl font-bold mb-8">系统监控</h1>

            <!-- 系统信息 -->
            <div class="grid grid-cols-1 md:grid-cols-2 gap-6 mb-8">
                <UCard>
                    <template #header>
                        <div class="flex items-center">
                            <div class="mr-2">
                                <UIcon name="i-mdi-desktop-tower-monitor" class="text-2xl" />
                            </div>
                            <h2 class="text-lg font-semibold">系统信息</h2>
                        </div>
                    </template>
                    <div class="space-y-2">
                        <p><span class="font-medium">主机名：</span>{{ systemInfo?.os?.hostname }}</p>
                        <p><span class="font-medium">系统：</span>{{ systemInfo?.os?.distro }}</p>
                        <p><span class="font-medium">版本：</span>{{ systemInfo?.os?.release }}</p>
                    </div>
                </UCard>

                <!-- CPU信息 -->
                <UCard>
                    <template #header>
                        <div class="flex items-center">
                            <div class="mr-2">
                                <UIcon name="i-mdi-cpu-64-bit" class="text-2xl" />
                            </div>
                            <h2 class="text-lg font-semibold">CPU</h2>
                        </div>
                    </template>
                    <div class="space-y-4">
                        <p><span class="font-medium">核心数：</span>{{ systemInfo?.cpu?.cores }}</p>
                        <div>
                            <span class="font-medium">使用率：</span>
                            <UProgress :value="systemInfo?.cpu?.load || 0" :color="getCpuColor(systemInfo?.cpu?.load)"
                                class="mt-2" />
                        </div>
                    </div>
                </UCard>

                <!-- 内存信息 -->
                <UCard>
                    <template #header>
                        <div class="flex items-center">
                            <div class="mr-2">
                                <UIcon name="i-mdi-memory" class="text-2xl" />
                            </div>
                            <h2 class="text-lg font-semibold">内存</h2>
                        </div>
                    </template>
                    <div class="space-y-4">
                        <p>
                            <span class="font-medium">总内存：</span>
                            {{ formatBytes(systemInfo?.memory?.total) }}
                        </p>
                        <div>
                            <span class="font-medium">使用率：</span>
                            <UProgress :value="getMemoryUsage(systemInfo?.memory)"
                                :color="getMemoryColor(getMemoryUsage(systemInfo?.memory))" class="mt-2" />
                        </div>
                    </div>
                </UCard>

                <!-- 磁盘信息 -->
                <UCard>
                    <template #header>
                        <div class="flex items-center">
                            <div class="mr-2">
                                <UIcon name="i-mdi-harddisk" class="text-2xl" />
                            </div>
                            <h2 class="text-lg font-semibold">磁盘</h2>
                        </div>
                    </template>
                    <div class="space-y-4">
                        <template v-for="disk in systemInfo?.disk" :key="disk.mount">
                            <div>
                                <p class="font-medium">{{ disk.mount }}</p>
                                <p class="text-sm text-gray-500 mb-2">
                                    {{ formatBytes(disk.used) }} / {{ formatBytes(disk.size) }}
                                </p>
                                <UProgress :value="(disk.used / disk.size) * 100"
                                    :color="getDiskColor((disk.used / disk.size) * 100)" />
                            </div>
                        </template>
                    </div>
                </UCard>
            </div>
        </UContainer>
    </div>
</template>

<script setup lang="ts">
const systemInfo = ref<any>(null);

// 获取系统信息
async function fetchSystemInfo() {
    try {
        systemInfo.value = await $fetch('/api/system');
    } catch (error) {
        console.error('Failed to fetch system info:', error);
    }
}

// 格式化字节
function formatBytes(bytes?: number) {
    if (!bytes) return '0 B';
    const k = 1024;
    const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return `${(bytes / Math.pow(k, i)).toFixed(2)} ${sizes[i]}`;
}

// 获取内存使用率
function getMemoryUsage(memory?: { total: number; used: number }) {
    if (!memory?.total || !memory?.used) return 0;
    return (memory.used / memory.total) * 100;
}

// 获取CPU颜色
function getCpuColor(load?: number) {
    if (!load) return 'primary';
    if (load > 90) return 'red';
    if (load > 70) return 'orange';
    return 'green';
}

// 获取内存颜色
function getMemoryColor(usage: number) {
    if (usage > 90) return 'red';
    if (usage > 70) return 'orange';
    return 'green';
}

// 获取磁盘颜色
function getDiskColor(usage: number) {
    if (usage > 90) return 'red';
    if (usage > 70) return 'orange';
    return 'green';
}

// 定时更新系统信息
onMounted(() => {
    fetchSystemInfo();
    const timer = setInterval(fetchSystemInfo, 5000);
    onUnmounted(() => clearInterval(timer));
});
</script>