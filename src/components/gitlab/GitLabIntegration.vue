<script setup lang="ts">
import { ref, onMounted } from 'vue'
import IconButton from '../IconButton.vue'
import { RiGitlabFill } from '@remixicon/vue'
import { useToast } from '../../composables/useToast'
import { useConfirm } from '../../composables/useConfirm'

interface GitLabToken {
    id: number
    name: string
    created_at: string
    last_used_at: string | null
}

const tokens = ref<GitLabToken[]>([])
const loading = ref(false)
const showAddDialog = ref(false)
const toast = useToast()
const confirm = useConfirm()

const newToken = ref({
    name: '',
    token: ''
})

async function loadTokens() {
    try {
        loading.value = true
        // TODO: 替换为实际的 API 调用
        tokens.value = [
            {
                id: 1,
                name: "CI Token",
                created_at: "2024-02-24T12:00:00Z",
                last_used_at: "2024-02-24T13:00:00Z"
            }
        ]
    } catch (err: any) {
        toast.error('加载 GitLab 令牌失败')
    } finally {
        loading.value = false
    }
}

async function addToken() {
    if (!newToken.value.name || !newToken.value.token) {
        toast.error('请填写完整的令牌信息')
        return
    }

    try {
        loading.value = true
        // TODO: 替换为实际的 API 调用
        await new Promise(resolve => setTimeout(resolve, 1000))
        showAddDialog.value = false
        newToken.value = { name: '', token: '' }
        await loadTokens()
        toast.success('添加 GitLab 令牌成功')
    } catch (err: any) {
        toast.error('添加 GitLab 令牌失败')
    } finally {
        loading.value = false
    }
}

async function deleteToken(token: GitLabToken) {
    if (await confirm.error('删除令牌', `确定要删除令牌 "${token.name}" 吗？`, {
        confirmText: '删除'
    })) {
        try {
            loading.value = true
            // TODO: 替换为实际的 API 调用
            await new Promise(resolve => setTimeout(resolve, 1000))
            await loadTokens()
            toast.success('删除 GitLab 令牌成功')
        } catch (err: any) {
            toast.error('删除 GitLab 令牌失败')
        } finally {
            loading.value = false
        }
    }
}

onMounted(() => {
    loadTokens()
})
</script>

<template>
    <div>
        <!-- 头部操作栏 -->
        <div class="flex justify-between items-center mb-4">
            <div class="flex items-center gap-2">
                <RiGitlabFill class="text-2xl text-[#FC6D26]" />
                <h2 class="text-lg font-bold">GitLab 集成</h2>
            </div>
            <div class="flex gap-2">
                <IconButton icon="ri-refresh-line" :loading="loading" @click="loadTokens">
                    刷新
                </IconButton>
                <IconButton icon="ri-add-line" variant="primary" @click="showAddDialog = true">
                    添加令牌
                </IconButton>
            </div>
        </div>

        <!-- 令牌列表 -->
        <div class="overflow-x-auto">
            <table class="table">
                <thead>
                    <tr>
                        <th>名称</th>
                        <th>创建时间</th>
                        <th>最后使用</th>
                        <th>操作</th>
                    </tr>
                </thead>
                <tbody>
                    <tr v-for="token in tokens" :key="token.id">
                        <td>{{ token.name }}</td>
                        <td>{{ new Date(token.created_at).toLocaleString() }}</td>
                        <td>{{ token.last_used_at ? new Date(token.last_used_at).toLocaleString() : '从未使用' }}</td>
                        <td>
                            <IconButton icon="ri-delete-bin-line" variant="error" size="sm" @click="deleteToken(token)">
                                删除
                            </IconButton>
                        </td>
                    </tr>
                    <tr v-if="tokens.length === 0">
                        <td colspan="4" class="text-center">
                            暂无 GitLab 令牌
                        </td>
                    </tr>
                </tbody>
            </table>
        </div>

        <!-- 添加令牌对话框 -->
        <dialog class="modal" :class="{ 'modal-open': showAddDialog }">
            <div class="modal-box">
                <h3 class="font-bold text-lg mb-4">添加 GitLab 令牌</h3>
                <div class="form-control w-full">
                    <label class="label">
                        <span class="label-text">令牌名称</span>
                    </label>
                    <input type="text" v-model="newToken.name" placeholder="例如：CI Token"
                        class="input input-bordered w-full" />
                </div>
                <div class="form-control w-full">
                    <label class="label">
                        <span class="label-text">访问令牌</span>
                    </label>
                    <input type="password" v-model="newToken.token" placeholder="GitLab 访问令牌"
                        class="input input-bordered w-full" />
                </div>
                <div class="modal-action">
                    <button class="btn" @click="showAddDialog = false">取消</button>
                    <button class="btn btn-primary" :disabled="loading" @click="addToken">
                        <span class="loading loading-spinner" v-if="loading"></span>
                        添加
                    </button>
                </div>
            </div>
            <form method="dialog" class="modal-backdrop">
                <button @click="showAddDialog = false">关闭</button>
            </form>
        </dialog>
    </div>
</template>