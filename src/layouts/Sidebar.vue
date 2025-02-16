<script setup lang="ts">
import { useRouter, useRoute } from 'vue-router'
import 'remixicon/fonts/remixicon.css'

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
        label: 'GitHub管理',
        key: 'github',
        icon: 'ri-github-fill'
    },
    // {
    //     label: '定时任务',
    //     key: 'cron',
    //     icon: 'ri-time-line'
    // }
]

const handleMenuClick = (key: string) => {
    router.push(`/${key}`)
}
</script>

<template>
    <ul class="menu bg-base-100 w-full p-2 gap-1">
        <li v-for="item in menuOptions" :key="item.key">
            <a @click="handleMenuClick(item.key)" :class="[
                'flex items-center gap-3 px-4 py-2.5 rounded-lg transition-colors',
                route.path.includes(item.key)
                    ? 'bg-primary text-primary-content'
                    : 'hover:bg-base-200'
            ]">
                <i :class="[
                    item.icon,
                    'text-xl',
                ]"></i>
                <span>{{ item.label }}</span>
            </a>
        </li>
    </ul>
</template>

<style scoped>
:deep(.n-menu-item-content) {
    padding-left: 8px !important;
}
</style>