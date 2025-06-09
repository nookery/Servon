<script setup lang="ts">
import { ref } from 'vue'
import GitHubRepos from './GitHubRepos.vue'
import GitHubLogs from './GitHubLogs.vue'
import GitHubConfig from './GitHubConfig.vue'

const currentTab = ref('repos')
const githubLogsRef = ref<InstanceType<typeof GitHubLogs> | null>(null)
</script>

<template>
    <div class="github-integration py-4 h-full">
        <!-- GitHub 子标签 -->
        <div class="tabs tabs-lift mb-6" role="tablist">
            <a role="tab" class="tab" :class="{ 'tab-active': currentTab === 'repos' }" @click="currentTab = 'repos'">
                <i class="ri-git-repository-line mr-2"></i>
                授权仓库
            </a>
            <a role="tab" class="tab" :class="{ 'tab-active': currentTab === 'logs' }" @click="currentTab = 'logs'">
                <i class="ri-file-list-line mr-2"></i>
                日志
            </a>
            <a role="tab" class="tab" :class="{ 'tab-active': currentTab === 'config' }" @click="currentTab = 'config'">
                <i class="ri-settings-line mr-2"></i>
                配置
            </a>
        </div>

        <!-- 授权仓库内容 -->
        <GitHubRepos v-if="currentTab === 'repos'" />

        <!-- 日志内容 -->
        <div v-else-if="currentTab === 'logs'" class="h-full">
            <GitHubLogs ref="githubLogsRef" />
        </div>

        <!-- 配置内容 -->
        <div v-else-if="currentTab === 'config'">
            <GitHubConfig ref="githubConfigRef" />
        </div>
    </div>
</template>