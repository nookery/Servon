<script setup lang="ts">
import { ref } from 'vue'
import IconButton from '../IconButton.vue'
import GitHubButton from './GitHubButton.vue'
import GitHubLogs from './GitHubLogs.vue'
import { type GitHubRepo, getAuthorizedRepos } from '../../api/github_api'

const currentTab = ref('repos')
const repos = ref<GitHubRepo[]>([])
const loading = ref(false)
const error = ref<string | null>(null)
const githubLogsRef = ref<InstanceType<typeof GitHubLogs> | null>(null)
const webhookLogsRef = ref<InstanceType<typeof GitHubLogs> | null>(null)

async function loadGitHubRepos() {
    loading.value = true
    error.value = null
    try {
        repos.value = await getAuthorizedRepos()
    } catch (err: any) {
        error.value = err.response?.data?.error || err.message || '加载失败'
    } finally {
        loading.value = false
    }
}
</script>

<template>
    <div class="github-integration py-4">
        <!-- GitHub 子标签 -->
        <div class="tabs tabs-boxed mb-6">
            <a class="tab" :class="{ 'tab-active': currentTab === 'repos' }" @click="currentTab = 'repos'">
                <i class="ri-git-repository-line mr-2"></i>
                授权仓库
            </a>
            <a class="tab" :class="{ 'tab-active': currentTab === 'logs' }" @click="currentTab = 'logs'">
                <i class="ri-file-list-line mr-2"></i>
                日志
            </a>
            <a class="tab" :class="{ 'tab-active': currentTab === 'webhook' }" @click="currentTab = 'webhook'">
                <i class="ri-webhook-line mr-2"></i>
                Webhook日志
            </a>
        </div>

        <!-- 授权仓库内容 -->
        <div v-if="currentTab === 'repos'" class="mb-8">
            <div class="flex justify-between mb-6">
                <div class="flex items-center gap-4">
                    <IconButton icon="ri-refresh-line" :loading="loading" @click="loadGitHubRepos" />
                </div>
                <GitHubButton />
            </div>

            <div v-if="error" class="alert alert-error mb-4">
                {{ error }}
            </div>

            <div v-if="loading" class="flex justify-center py-8">
                <span class="loading loading-spinner loading-lg"></span>
            </div>

            <div v-else-if="repos.length === 0"
                class="flex flex-col items-center justify-center py-12 text-base-content/60">
                <i class="ri-git-repository-line text-6xl mb-4"></i>
                <p class="text-lg mb-2">暂无授权仓库</p>
                <p class="text-sm">请先完成 GitHub 授权并选择要集成的仓库</p>
            </div>

            <div v-else class="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div v-for="repo in repos" :key="repo.id" class="card bg-base-200 hover:bg-base-300 transition-colors">
                    <div class="card-body">
                        <h3 class="card-title">
                            <i class="ri-git-repository-line mr-2"></i>
                            {{ repo.name }}
                        </h3>
                        <p class="text-base-content/70">{{ repo.description || '暂无描述' }}</p>
                        <div class="card-actions justify-end mt-4">
                            <a :href="repo.htmlUrl" target="_blank" class="btn btn-sm btn-ghost">
                                <i class="ri-external-link-line mr-1"></i>
                                查看仓库
                            </a>
                        </div>
                    </div>
                </div>
            </div>
        </div>

        <!-- 集成日志内容 -->
        <div v-else-if="currentTab === 'logs'">
            <GitHubLogs ref="githubLogsRef" />
        </div>

        <!-- Webhook日志内容 -->
        <div v-else-if="currentTab === 'webhook'">
            <GitHubLogs ref="webhookLogsRef" initial-path="/data/github/webhook" />
        </div>
    </div>
</template>