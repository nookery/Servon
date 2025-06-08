<script setup lang="ts">
import { computed } from 'vue'
import { useWindowStore } from '../stores/windowStore'
import DockItem from '../components/DockItem.vue'

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
    <div class="flex justify-center items-end h-full">
        <div
            class="flex bg-base-200/60 backdrop-blur-md rounded-2xl py-2 px-6 items-center justify-center gap-2 shadow-lg mb-4 min-w-fit w-max border border-white/10 transition-all duration-300">
            <DockItem v-for="item in menuOptions" :key="item.key" :label="item.label" :icon="item.icon"
                :color="item.color" :is-active="activeItem === item.key" :onClick="() => handleMenuClick(item)" />
        </div>
    </div>
</template>