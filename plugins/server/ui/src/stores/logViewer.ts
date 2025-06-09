import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useLogViewerStore = defineStore('logViewer', () => {
    // Initialize from localStorage, default to false if not found
    const isVisible = ref(localStorage.getItem('logViewerVisible') === 'true')

    function toggleVisibility() {
        isVisible.value = !isVisible.value
        // Save to localStorage when value changes
        localStorage.setItem('logViewerVisible', isVisible.value.toString())
    }

    function setVisibility(value: boolean) {
        isVisible.value = value
        // Save to localStorage when value changes
        localStorage.setItem('logViewerVisible', value.toString())
    }

    return {
        isVisible,
        toggleVisibility,
        setVisibility
    }
}) 