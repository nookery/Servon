<script setup lang="ts">
import { ref, onMounted } from 'vue'
import PageContainer from '../layouts/PageContainer.vue'
import { type User, type NewUser, getUsers, createUser, deleteUser } from '../api/users'
import { useConfirm } from '../composables/useConfirm'
import { useToast } from '../composables/useToast'

const users = ref<User[]>([])
const error = ref<string | null>(null)
const newUser = ref<NewUser>({
    username: '',
    password: ''
})
const showCreateModal = ref(false)
const confirm = useConfirm()
const toast = useToast()

// 加载用户列表
async function loadUsers() {
    try {
        users.value = await getUsers()
        error.value = null
    } catch (err: any) {
        error.value = '获取用户列表失败: ' + err.message
    }
}

// 格式化时间
function formatDate(dateStr: string): string {
    if (!dateStr) return '未知'
    return new Date(dateStr).toLocaleString()
}

// 创建用户
async function handleCreateUser() {
    try {
        await createUser(newUser.value)
        showCreateModal.value = false
        newUser.value = { username: '', password: '' }
        await loadUsers()
        toast.success('创建用户成功')
    } catch (err: any) {
        error.value = '创建用户失败: ' + err.message
    }
}

// 删除用户
async function handleDelete(username: string) {
    if (await confirm.error('删除用户', `确定要删除用户 ${username} 吗？此操作不可撤销。`, {
        confirmText: '删除'
    })) {
        try {
            await deleteUser(username)
            await loadUsers()
            toast.success('删除用户成功')
        } catch (err: any) {
            error.value = '删除用户失败: ' + err.message
        }
    }
}

onMounted(() => {
    loadUsers()
})
</script>

<template>
    <PageContainer title="用户管理" :error="error">
        <template #header>
            <div class="flex justify-end mb-4">
                <button class="btn btn-primary" @click="showCreateModal = true">
                    创建用户
                </button>
            </div>
        </template>

        <div class="overflow-x-auto">
            <table class="table table-zebra w-full">
                <thead>
                    <tr>
                        <th>用户名</th>
                        <th>用户组</th>
                        <th>Shell</th>
                        <th>主目录</th>
                        <th>创建时间</th>
                        <th>最后登录</th>
                        <th>权限</th>
                        <th>操作</th>
                    </tr>
                </thead>
                <tbody>
                    <tr v-for="user in users" :key="user.username">
                        <td>{{ user.username }}</td>
                        <td>
                            <div class="flex flex-wrap gap-1">
                                <span v-for="group in user.groups" :key="group" class="badge badge-sm">
                                    {{ group }}
                                </span>
                            </div>
                        </td>
                        <td>{{ user.shell }}</td>
                        <td class="text-xs">{{ user.home_dir }}</td>
                        <td>{{ formatDate(user.create_time) }}</td>
                        <td>{{ formatDate(user.last_login) }}</td>
                        <td>
                            <span class="badge" :class="user.sudo ? 'badge-warning' : 'badge-ghost'">
                                {{ user.sudo ? 'sudo' : 'normal' }}
                            </span>
                        </td>
                        <td>
                            <button class="btn btn-error btn-sm" @click="handleDelete(user.username)"
                                :disabled="user.sudo">
                                <i class="ri-delete-bin-line"></i>
                            </button>
                        </td>
                    </tr>
                </tbody>
            </table>
        </div>

        <!-- 创建用户模态框 -->
        <dialog class="modal" :class="{ 'modal-open': showCreateModal }">
            <div class="modal-box">
                <h3 class="font-bold text-lg mb-4">创建新用户</h3>
                <form @submit.prevent="handleCreateUser">
                    <div class="form-control">
                        <label class="label">
                            <span class="label-text">用户名</span>
                        </label>
                        <input type="text" v-model="newUser.username" class="input input-bordered" required />
                    </div>
                    <div class="form-control mt-4">
                        <label class="label">
                            <span class="label-text">密码</span>
                        </label>
                        <input type="password" v-model="newUser.password" class="input input-bordered" required />
                    </div>
                    <div class="modal-action">
                        <button type="submit" class="btn btn-primary">创建</button>
                        <button type="button" class="btn" @click="showCreateModal = false">
                            取消
                        </button>
                    </div>
                </form>
            </div>
        </dialog>
    </PageContainer>
</template>