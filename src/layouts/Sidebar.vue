<script setup lang="ts">
import { NMenu, NIcon } from 'naive-ui'
import { h } from 'vue'
import { useRouter } from 'vue-router'
import {
    ServerOutline,
    AppsOutline,
    TerminalOutline,
    FolderOutline,
    SwapHorizontalOutline
} from '@vicons/ionicons5'

const router = useRouter()
const renderIcon = (icon: any) => () => h(NIcon, null, { default: () => h(icon) })

const menuOptions = [
    {
        label: '系统概览',
        key: 'dashboard',
        icon: renderIcon(ServerOutline)
    },
    {
        label: '软件管理',
        key: 'software',
        icon: renderIcon(AppsOutline)
    },
    {
        label: '进程管理',
        key: 'processes',
        icon: renderIcon(TerminalOutline)
    },
    {
        label: '文件管理',
        key: 'files',
        icon: renderIcon(FolderOutline)
    },
    {
        label: '端口管理',
        key: 'ports',
        icon: renderIcon(SwapHorizontalOutline)
    }
]

const handleMenuClick = (key: string) => {
    router.push(`/${key}`)
}

defineProps<{
    collapsed?: boolean
}>()
</script>

<template>
    <n-menu :options="menuOptions" :collapsed-icon-size="24" :collapsed="collapsed" :indent="0"
        @update:value="handleMenuClick" />
</template>

<style scoped>
:deep(.n-menu-item-content) {
    padding-left: 8px !important;
}
</style>