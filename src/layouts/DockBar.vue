<script setup lang="ts">
import { useRouter, useRoute } from 'vue-router'
import { computed } from 'vue'

const router = useRouter()
const route = useRoute()

const menuOptions = [
    {
        label: '系统概览',
        key: 'dashboard',
        icon: 'ri-server-line'
    },
    {
        label: '软件管理',
        key: 'software',
        icon: 'ri-apps-line'
    },
    {
        label: '进程管理',
        key: 'processes',
        icon: 'ri-terminal-line'
    },
    {
        label: '文件管理',
        key: 'files',
        icon: 'ri-folder-line'
    },
    {
        label: '端口管理',
        key: 'ports',
        icon: 'ri-swap-line'
    },
    {
        label: '用户管理',
        key: 'users',
        icon: 'ri-user-settings-line'
    },
    {
        label: '集成管理',
        key: 'integrations',
        icon: 'ri-settings-3-line'
    },
    {
        label: '数据中心',
        key: 'data',
        icon: 'ri-database-line'
    },
    {
        label: '日志管理',
        key: 'logs',
        icon: 'ri-file-list-line'
    },
    {
        label: '项目管理',
        key: 'projects',
        icon: 'ri-folder-open-line'
    }
]

// 当前激活的菜单项
const activeItem = computed(() => {
    return menuOptions.find(item => route.path.includes(item.key))?.key || ''
})

const handleMenuClick = (key: string) => {
    router.push(`/${key}`)
}
</script>

<template>
    <!-- macOS风格的Dock -->
    <div class="dock-container flex justify-center items-end h-full w-full">
        <div class="dock bg-base-200/80 backdrop-blur-md rounded-2xl px-4 py-2 flex items-center gap-3 shadow-lg">
            <div v-for="item in menuOptions" :key="item.key" class="dock-item relative flex flex-col items-center"
                @click="handleMenuClick(item.key)">

                <!-- Icon -->
                <div class="icon-container p-3 rounded-xl flex items-center justify-center" :class="[
                    activeItem === item.key
                        ? 'bg-primary text-primary-content'
                        : 'bg-base-100 text-base-content hover:bg-base-300'
                ]">
                    <i :class="[item.icon, 'text-2xl']"></i>
                </div>

                <!-- Tooltip -->
                <div class="tooltip tooltip-top" :data-tip="item.label"></div>

                <!-- Indicator dot for active item -->
                <div v-if="activeItem === item.key" class="indicator-dot bg-primary h-1 w-1 rounded-full mt-1"></div>
            </div>
        </div>
    </div>
</template>

<style scoped>
.dock {
    transition: all 0.3s ease;
}

.dock-item {
    cursor: pointer;
}
</style>