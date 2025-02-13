import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useLogViewerStore = defineStore('logViewer', () => {
    const isVisible = ref(false)

    function toggleVisibility() {
        isVisible.value = !isVisible.value
    }

    function setVisibility(value: boolean) {
        isVisible.value = value
    }

    return {
        isVisible,
        toggleVisibility,
        setVisibility
    }
}) 