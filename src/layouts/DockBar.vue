<script setup lang="ts">
import { computed } from 'vue'
import { useWindowStore } from '../stores/windowStore'

const windowStore = useWindowStore()

const menuOptions = [
    {
        label: '系统概览',
        key: 'dashboard',
        icon: 'ri-server-line',
        component: 'DashboardView',
        color: '#FF6B6B' // 红色
    },
    {
        label: '软件管理',
        key: 'software',
        icon: 'ri-apps-line',
        component: 'SoftwareView',
        color: '#4ECDC4' // 蓝绿色
    },
    {
        label: '进程管理',
        key: 'processes',
        icon: 'ri-terminal-line',
        component: 'ProcessesView',
        color: '#FFD166' // 黄色
    },
    {
        label: '文件管理',
        key: 'files',
        icon: 'ri-folder-line',
        component: 'FilesView',
        color: '#118AB2' // 蓝色
    },
    {
        label: '端口管理',
        key: 'ports',
        icon: 'ri-swap-line',
        component: 'PortsView',
        color: '#06D6A0' // 绿色
    },
    {
        label: '用户管理',
        key: 'users',
        icon: 'ri-user-settings-line',
        component: 'UsersView',
        color: '#9381FF' // 紫色
    },
    {
        label: '集成管理',
        key: 'integrations',
        icon: 'ri-settings-3-line',
        component: 'IntegrationsView',
        color: '#F78C6B' // 橙色
    },
    {
        label: '数据中心',
        key: 'data',
        icon: 'ri-database-line',
        component: 'DataView',
        color: '#3A86FF' // 蓝色
    },
    {
        label: '日志管理',
        key: 'logs',
        icon: 'ri-file-list-line',
        component: 'LogsView',
        color: '#8338EC' // 深紫色
    },
    {
        label: '项目管理',
        key: 'projects',
        icon: 'ri-folder-open-line',
        component: 'ProjectsView',
        color: '#FB5607' // 橙红色
    }
]

// 当前激活的菜单项
const activeItem = computed(() => {
    return menuOptions.find(item => {
        const windowWithComponent = windowStore.windows.find(w => w.component === item.component)
        return windowWithComponent && windowWithComponent.isActive
    })?.key || ''
})

const handleMenuClick = (menuItem: typeof menuOptions[0]) => {
    windowStore.openWindow({
        title: menuItem.label,
        icon: menuItem.icon,
        component: menuItem.component,
        props: {}
    })
}
</script>

<template>
    <!-- macOS风格的Dock -->
    <div class="dock-container flex justify-center items-end h-full px-8 border border-red-200">
        <div
            class="dock bg-info/50 backdrop-blur-md rounded-2xl w-full px-24 py-2 flex items-center gap-1 shadow-lg mb-4">
            <div v-for="item in menuOptions" :key="item.key" class="dock-item relative flex flex-col items-center"
                @click="handleMenuClick(item)">

                <!-- Icon with custom color -->
                <div class="icon-container p-2 rounded-xl flex items-center justify-center" :class="[
                    activeItem === item.key
                        ? 'bg-opacity-90 text-white'
                        : 'bg-base-100 hover:bg-opacity-90 hover:text-white'
                ]" :style="{
                    backgroundColor: activeItem === item.key ? item.color : '',
                    color: activeItem === item.key ? 'white' : item.color
                }">
                    <i :class="[item.icon, 'text-2xl']"></i>
                </div>

                <!-- Tooltip -->
                <div class="tooltip tooltip-top" :data-tip="item.label"></div>

                <!-- Indicator dot for active item -->
                <div v-if="activeItem === item.key" class="indicator-dot h-1 w-1 rounded-full mt-1"
                    :style="{ backgroundColor: item.color }">
                </div>
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

.icon-container {
    transition: all 0.2s ease;
    border: 2px solid transparent;
    border-radius: 12px;
}

.icon-container:hover {
    transform: translateY(-2px);
}
</style>