<script setup lang="ts">
import { ref, onMounted } from 'vue'
import IconButton from '../components/IconButton.vue'
const themes = [
    'light',
    'dark',
]

const currentTheme = ref(localStorage.getItem('theme') || 'light')

function changeTheme(theme: string) {
    document.documentElement.setAttribute('data-theme', theme)
    localStorage.setItem('theme', theme)
    currentTheme.value = theme
}

onMounted(() => {
    const savedTheme = localStorage.getItem('theme')
    if (savedTheme) {
        changeTheme(savedTheme)
    }
})
</script>

<template>
    <div class="dropdown dropdown-end m-0 p-0">
        <label tabindex="0" class="btn btn-ghost gap-2">
            <IconButton icon="ri-palette-line" variant="ghost" circle title="切换主题" />
        </label>
        <ul tabindex="0"
            class="dropdown-content z-[1] menu p-2 shadow bg-base-100 rounded-box w-52 max-h-96 overflow-y-auto">
            <li v-for="theme in themes" :key="theme">
                <button class="flex items-center gap-2"
                    :class="{ 'bg-primary text-primary-content': currentTheme === theme }" @click="changeTheme(theme)">
                    <i v-if="currentTheme === theme" class="ri-check-line"></i>
                    <span>{{ theme }}</span>
                </button>
            </li>
        </ul>
    </div>
</template>

<style scoped>
.dropdown-content {
    max-height: 300px;
}

/* 自定义滚动条样式 */
.dropdown-content::-webkit-scrollbar {
    width: 6px;
}

.dropdown-content::-webkit-scrollbar-track {
    background: transparent;
}

.dropdown-content::-webkit-scrollbar-thumb {
    background-color: var(--tw-primary);
    border-radius: 3px;
}
</style>