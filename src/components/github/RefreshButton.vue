<script setup lang="ts">
import IconButton from '../IconButton.vue'

defineProps<{
    loading?: boolean
}>()

defineEmits<{
    (e: 'refresh'): void
}>()
</script>

<template>
    <IconButton 
    circle
        icon="ri-refresh-line" 
        variant="ghost"
        
        :loading="loading" 
        @click="$emit('refresh')" 
    />
</template>