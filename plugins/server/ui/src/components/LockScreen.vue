<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { useLockStore } from '../stores/lockStore'
import { RiLock2Line } from '@remixicon/vue'

const lockStore = useLockStore()
const userInput = ref('')
const PASSWORD = 'admin'

const checkPassword = () => {
    if (userInput.value.includes(PASSWORD)) {
        lockStore.unlock()
        userInput.value = ''
    }
}

const handleKeyPress = (e: KeyboardEvent) => {
    userInput.value += e.key
    // 保持输入字符串在合理长度内
    if (userInput.value.length > 50) {
        userInput.value = userInput.value.slice(-50)
    }
    checkPassword()
}

onMounted(() => {
    window.addEventListener('keypress', handleKeyPress)
})

onUnmounted(() => {
    window.removeEventListener('keypress', handleKeyPress)
})
</script>

<template>
    <div class="fixed inset-0 bg-gray-900/70 backdrop-blur-md flex items-center justify-center z-50">
        <div class="text-center space-y-6">
            <div class="text-white space-y-4">
                <RiLock2Line class="w-20 h-20 mx-auto opacity-80" />
                <h2 class="text-4xl font-bold tracking-wider">应用已锁定</h2>
                <p class="text-lg text-gray-300">按任意键输入密码解锁</p>
            </div>
        </div>
    </div>
</template>