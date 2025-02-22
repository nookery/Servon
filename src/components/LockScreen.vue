<script setup lang="ts">
import { ref } from 'vue'

const password = ref('')
const isLocked = ref(true)
const errorMessage = ref('')

const unlock = () => {
    // 这里应该与后端 API 进行验证
    // 临时使用硬编码的密码进行演示
    if (password.value === 'admin') {
        isLocked.value = false
        errorMessage.value = ''
    } else {
        errorMessage.value = '密码错误'
        password.value = ''
    }
}
</script>

<template>
    <div v-if="isLocked" class="fixed inset-0 bg-gray-900/50 backdrop-blur-sm flex items-center justify-center z-50">
        <div class="bg-white dark:bg-gray-800 p-8 rounded-lg shadow-xl w-96">
            <h2 class="text-2xl font-bold mb-6 text-center dark:text-white">应用已锁定</h2>
            <form @submit.prevent="unlock" class="space-y-4">
                <div>
                    <input
                        type="password"
                        v-model="password"
                        placeholder="请输入密码"
                        class="w-full px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 dark:bg-gray-700 dark:border-gray-600 dark:text-white"
                        autocomplete="current-password"
                    />
                </div>
                <p v-if="errorMessage" class="text-red-500 text-sm">{{ errorMessage }}</p>
                <button
                    type="submit"
                    class="w-full bg-blue-500 text-white py-2 rounded-lg hover:bg-blue-600 transition-colors"
                >
                    解锁
                </button>
            </form>
        </div>
    </div>
</template>