<script setup lang="ts">
import ThemeSwitcher from '../modules/ThemeSwitcher.vue'
import NetworkStatus from '../components/NetworkStatus.vue'
import SystemResources from '../components/SystemResources.vue'
import OSInfo from '../components/OSInfo.vue'
import UserInfo from '../components/UserInfo.vue'
import TaskManagerButton from '../components/TaskManagerButton.vue'
import IPInfoButton from '../components/IPInfoButton.vue'
import LogViewerButton from '../modules/LogViewerButton.vue'
import GitHubButton from '../components/GitHubButton.vue'
</script>

<template>
    <div class="navbar bg-base-100 fixed top-0 left-0 z-50 h-16 px-6 shadow-sm">
        <div class="flex-1">
            <OSInfo />
        </div>

        <div class="flex-none gap-6">
            <div class="flex items-center gap-6">
                <IPInfoButton />
                <GitHubButton />
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