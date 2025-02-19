<script setup lang="ts">
import ThemeSwitcher from '../modules/ThemeSwitcher.vue'
import NetworkStatus from '../modules/NetworkStatus.vue'
import SystemResources from '../modules/SystemResources.vue'
import OSInfo from '../modules/OSInfo.vue'
import UserInfo from '../modules/UserInfo.vue'
import TaskManagerButton from '../modules/TaskManagerButton.vue'
import IPInfoButton from '../modules/IPInfoButton.vue'
import LogViewerButton from '../modules/LogViewerButton.vue'
</script>

<template>
    <div class="navbar bg-base-100 fixed top-0 left-0 z-50 h-16 px-6 shadow-sm">
        <div class="flex-1">
            <OSInfo />
        </div>

        <div class="flex-none gap-6">
            <div class="flex items-center gap-6">
                <IPInfoButton />
                <LogViewerButton />
                <TaskManagerButton />
                <ThemeSwitcher />
                <NetworkStatus />
                <SystemResources />
                <UserInfo />
            </div>
        </div>
    </div>
</template>