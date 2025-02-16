import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useLayoutStore = defineStore('layout', () => {
    const collapsed = ref(false)

    function toggleCollapsed() {
        collapsed.value = !collapsed.value
    }

    function setCollapsed(value: boolean) {
        collapsed.value = value
    }

    return {
        collapsed,
        toggleCollapsed,
        setCollapsed
    }
}) 