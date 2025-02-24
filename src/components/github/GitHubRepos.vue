<script setup lang="ts">
import { ref, computed } from 'vue'
import GitHubButton from './GitHubButton.vue'
import RefreshButton from './RefreshButton.vue'
import { getAuthorizedRepos } from '../../api/github_api'
import { deployRepository } from '../../api/deploy_api'
import type { GitHubRepo } from '../../types/GitHubTypes'
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

    console.log(repos.value)
}

// 筛选状态：null 表示全部，true 表示私有，false 表示公开
const filterPrivate = ref<boolean | null>(null)

// 添加搜索关键词
const searchQuery = ref('')

// 修改计算属性，增加搜索过滤
const paginatedRepos = computed(() => {
    let filtered = repos.value

    // 先按私有状态筛选
    if (filterPrivate !== null && filterPrivate.value !== null) {
        console.log(filterPrivate.value)
        filtered = filtered.filter(repo => repo.private === filterPrivate.value)
    }

    // 再按搜索关键词筛选
    if (searchQuery.value) {
        const query = searchQuery.value.toLowerCase()
        filtered = filtered.filter(repo =>
            repo.name.toLowerCase().includes(query) ||
            (repo.description?.toLowerCase() || '').includes(query)
        )
    }

    const start = (currentPage.value - 1) * pageSize.value
    const end = start + pageSize.value
    return filtered.slice(start, end)
})

// 修改总页数计算，使用相同的筛选逻辑
const totalPages = computed(() => {
    let filtered = repos.value

    // 先按私有状态筛选
    if (filterPrivate.value !== null) {
        filtered = filtered.filter(repo => repo.private === filterPrivate.value)
    }

    // 再按搜索关键词筛选
    if (searchQuery.value) {
        const query = searchQuery.value.toLowerCase()
        filtered = filtered.filter(repo =>
            repo.name.toLowerCase().includes(query) ||
            (repo.description?.toLowerCase() || '').includes(query)
        )
    }

    return Math.max(1, Math.ceil(filtered.length / pageSize.value))
})

// 页面切换
function changePage(page: number) {
    currentPage.value = page
}

function getRepoHTMLURL(repo: GitHubRepo) {
    return `https://github.com/${repo.full_name}`
}

async function handleDeploy(repo: GitHubRepo) {
    if (deployingRepo.value) return

    deployingRepo.value = repo.name
    try {
        const res = await deployRepository(getRepoHTMLURL(repo))
        toast.success(res.message)
    } catch (err: any) {
        error.value = err.response?.data?.error || err.message || '部署失败'
    } finally {
        deployingRepo.value = null
    }
}

// 计算中间要显示的页码
const middlePages = computed(() => {
    const total = totalPages.value
    const current = currentPage.value
    const delta = 2 // 当前页码左右各显示多少个页码

    let start = Math.max(2, current - delta)
    let end = Math.min(total - 1, current + delta)

    // 调整start和end，确保显示固定数量的页码
    if (end - start < delta * 2) {
        if (start === 2) {
            end = Math.min(total - 1, start + delta * 2)
        } else if (end === total - 1) {
            start = Math.max(2, end - delta * 2)
        }
    }

    const pages = []
    for (let i = start; i <= end; i++) {
        pages.push(i)
    }
    return pages
})

// 是否显示左省略号
const showLeftEllipsis = computed(() => {
    return middlePages.value[0] > 2
})

// 是否显示右省略号
const showRightEllipsis = computed(() => {
    return middlePages.value[middlePages.value.length - 1] < totalPages.value - 1
})

loadGitHubRepos()
</script>

<template>
    <div class="mb-8">
        <div class="flex justify-between mb-6">
            <div class="flex items-center gap-4">
                <GitHubButton />
                <RefreshButton :loading="loading" @refresh="loadGitHubRepos" />
                <!-- 搜索框 -->
                <div class="form-control">
                    <input type="text" v-model="searchQuery" placeholder="搜索仓库..."
                        class="input input-bordered input-sm w-64">
                </div>
                <div class="join">
                    <button class="join-item btn btn-sm" :class="{ 'btn-active': filterPrivate === null }"
                        @click="filterPrivate = null">
                        全部
                    </button>
                    <button class="join-item btn btn-sm" :class="{ 'btn-active': filterPrivate === false }"
                        @click="filterPrivate = false">
                        公开
                    </button>
                    <button class="join-item btn btn-sm" :class="{ 'btn-active': filterPrivate === true }"
                        @click="filterPrivate = true">
                        私有
                    </button>
                </div>
            </div>
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
                                    <a :href="getRepoHTMLURL(repo)" target="_blank" class="btn btn-sm btn-ghost">
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

            <!-- 分页区域 -->
            <div class="flex flex-col items-center gap-2 mt-4">
                <!-- 分页按钮 -->
                <div class="join">
                    <!-- 第一页按钮 -->
                    <button class="join-item btn btn-sm" :disabled="currentPage === 1" @click="changePage(1)"
                        title="第一页">
                        <i class="ri-skip-back-line"></i>
                    </button>

                    <!-- 上一页 -->
                    <button class="join-item btn btn-sm" :disabled="currentPage === 1"
                        @click="changePage(currentPage - 1)" title="上一页">
                        <i class="ri-arrow-left-s-line"></i>
                    </button>

                    <!-- 第一页 -->
                    <button class="join-item btn btn-sm" :class="{ 'btn-active': currentPage === 1 }"
                        @click="changePage(1)">1</button>

                    <!-- 左省略号 -->
                    <button v-if="showLeftEllipsis" class="join-item btn btn-sm btn-disabled">...</button>

                    <!-- 中间页码 -->
                    <button v-for="page in middlePages" :key="page" class="join-item btn btn-sm"
                        :class="{ 'btn-active': page === currentPage }" @click="changePage(page)">
                        {{ page }}
                    </button>

                    <!-- 右省略号 -->
                    <button v-if="showRightEllipsis" class="join-item btn btn-sm btn-disabled">...</button>

                    <!-- 最后一页 -->
                    <button v-if="totalPages > 1" class="join-item btn btn-sm"
                        :class="{ 'btn-active': currentPage === totalPages }" @click="changePage(totalPages)">{{
                            totalPages }}</button>

                    <!-- 下一页 -->
                    <button class="join-item btn btn-sm" :disabled="currentPage === totalPages"
                        @click="changePage(currentPage + 1)" title="下一页">
                        <i class="ri-arrow-right-s-line"></i>
                    </button>

                    <!-- 最后一页按钮 -->
                    <button class="join-item btn btn-sm" :disabled="currentPage === totalPages"
                        @click="changePage(totalPages)" title="最后一页">
                        <i class="ri-skip-forward-line"></i>
                    </button>
                </div>

                <!-- 分页信息 -->
                <div class="text-sm text-base-content/70">
                    共 {{ repos.length }} 个仓库，第 {{ currentPage }} / {{ totalPages }} 页
                </div>
            </div>
        </div>
    </div>
</template>