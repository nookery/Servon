<template>
    <div class="p-8">
        <h1 class="text-2xl font-bold mb-8">软件信息</h1>

        <!-- 搜索框 -->
        <div class="mb-6">
            <UInput v-model="searchQuery" icon="i-mdi-magnify" placeholder="搜索软件..." class="max-w-md" />
        </div>

        <!-- 软件列表 -->
        <UCard>
            <template #header>
                <div class="flex items-center justify-between">
                    <div class="flex items-center">
                        <UIcon name="i-mdi-apps" class="text-2xl mr-2" />
                        <h2 class="text-lg font-semibold">已安装的软件</h2>
                    </div>
                    <UBadge :label="filteredSoftware.length.toString()" />
                </div>
            </template>

            <UTable :columns="columns" :rows="filteredSoftware" :loading="loading">
                <template #name-data="{ row }">
                    <div class="font-medium">{{ row.name }}</div>
                </template>

                <template #version-data="{ row }">
                    <UBadge :label="row.version" variant="soft" />
                </template>
            </UTable>
        </UCard>
    </div>
</template>

<script setup lang="ts">
const searchQuery = ref('')
const loading = ref(true)
const software = ref<any[]>([])

// 表格列定义
const columns = [
    {
        key: 'name',
        label: '名称'
    },
    {
        key: 'version',
        label: '版本'
    },
    {
        key: 'description',
        label: '描述'
    }
]

// 获取软件列表
async function fetchSoftware() {
    try {
        loading.value = true
        software.value = await $fetch('/api/software')
    } catch (error) {
        console.error('Failed to fetch software list:', error)
    } finally {
        loading.value = false
    }
}

// 过滤后的软件列表
const filteredSoftware = computed(() => {
    if (!searchQuery.value) return software.value
    const query = searchQuery.value.toLowerCase()
    return software.value.filter(item =>
        item.name.toLowerCase().includes(query) ||
        item.description?.toLowerCase().includes(query)
    )
})

// 初始加载
onMounted(() => {
    fetchSoftware()
})
</script>