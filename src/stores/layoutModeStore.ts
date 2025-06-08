import { defineStore } from 'pinia'
import { ref } from 'vue'

export type LayoutMode = 'web' | 'os'

export const useLayoutModeStore = defineStore('layoutMode', () => {
    // 默认使用 web 布局，但会从 localStorage 中读取保存的设置
    const mode = ref<LayoutMode>(localStorage.getItem('layoutMode') as LayoutMode || 'web')

    // 切换布局模式
    function toggleMode() {
        mode.value = mode.value === 'web' ? 'os' : 'web'
        // 保存到 localStorage
        localStorage.setItem('layoutMode', mode.value)
    }

    // 设置特定布局模式
    function setMode(newMode: LayoutMode) {
        mode.value = newMode
        localStorage.setItem('layoutMode', mode.value)
    }

    return { mode, toggleMode, setMode }
}) 