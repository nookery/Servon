<script setup lang="ts">
import { ref, computed } from 'vue'
import IconButton from '../IconButton.vue'
import GitHubButton from './GitHubButton.vue'
import { getAuthorizedRepos } from '../../api/github_api'
import { deployRepository } from '../../api/deploy_api'
import type { GitHubRepo } from '../../models/GitHubTypes'
import { useToast } from '../../composables/useToast'

const toast = useToast()
const repos = ref<GitHubRepo[]>([])
const loading = ref(false)
const error = ref<string | null>(null)

// 分页相关
const currentPage = ref(1)
const pageSize = ref(10)

// 部署相关
const deployingRepo = ref<string | null>(null)

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

// 计算当前页的数据
const paginatedRepos = computed(() => {
    const start = (currentPage.value - 1) * pageSize.value
    const end = start + pageSize.value
    return repos.value.slice(start, end)
})

// 计算总页数
const totalPages = computed(() => Math.ceil(repos.value.length / pageSize.value))

// 页面切换
function changePage(page: number) {
    currentPage.value = page
}

async function handleDeploy(repo: GitHubRepo) {
    if (deployingRepo.value) return

    deployingRepo.value = repo.name
    try {
        const res = await deployRepository(repo.full_name)
        toast.success(res.message)
    } catch (err: any) {
        error.value = err.response?.data?.error || err.message || '部署失败'
    } finally {
        deployingRepo.value = null
    }
}

loadGitHubRepos()
</script>

<template>
    <div class="mb-8">
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

        <div v-else>
            <div class="overflow-x-auto">
                <table class="table">
                    <thead>
                        <tr>
                            <th>仓库名称</th>
                            <th>描述</th>
                            <th>类型</th>
                            <th>操作</th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr v-for="repo in paginatedRepos" :key="repo.id">
                            <td class="font-medium">
                                <i class="ri-git-repository-line mr-2"></i>
                                {{ repo.name }}
                            </td>
                            <td class="text-base-content/70">{{ repo.description || '暂无描述' }}</td>
                            <td>
                                <span class="badge" :class="repo.private ? 'badge-neutral' : 'badge-ghost'">
                                    {{ repo.private ? '私有' : '公开' }}
                                </span>
                            </td>
                            <td>
                                <div class="flex gap-2">
                                    <a :href="repo.html_url" target="_blank" class="btn btn-sm btn-ghost">
                                        <i class="ri-external-link-line mr-1"></i>
                                        查看
                                    </a>
                                    <button class="btn btn-sm btn-primary"
                                        :class="{ 'loading': deployingRepo === repo.name }" @click="handleDeploy(repo)"
                                        :disabled="deployingRepo !== null">
                                        <i class="ri-rocket-line mr-1"></i>
                                        {{ deployingRepo === repo.name ? '部署中' : '部署' }}
                                    </button>
                                </div>
                            </td>
                        </tr>
                    </tbody>
                </table>
            </div>

            <!-- 分页 -->
            <div class="flex justify-center mt-4">
                <div class="join">
                    <button class="join-item btn btn-sm" :disabled="currentPage === 1"
                        @click="changePage(currentPage - 1)">
                        <i class="ri-arrow-left-s-line"></i>
                    </button>
                    <button v-for="page in totalPages" :key="page" class="join-item btn btn-sm"
                        :class="{ 'btn-active': page === currentPage }" @click="changePage(page)">
                        {{ page }}
                    </button>
                    <button class="join-item btn btn-sm" :disabled="currentPage === totalPages"
                        @click="changePage(currentPage + 1)">
                        <i class="ri-arrow-right-s-line"></i>
                    </button>
                </div>
            </div>
        </div>
    </div>
</template>