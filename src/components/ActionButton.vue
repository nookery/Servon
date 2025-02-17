<template>
    <button :type="type" :class="[
        'btn',
        variant && `btn-${variant}`,
        size && `btn-${size}`,
        { 'btn-disabled': disabled }
    ]" @click="$emit('click')">
        <i v-if="icon" :class="['text-primary text-xl', icon]"></i>
        <slot></slot>
    </button>
</template>

<script setup lang="ts">
interface Props {
    type?: 'button' | 'submit' | 'reset'
    variant?: 'primary' | 'secondary' | 'accent' | 'ghost' | 'link'
    size?: 'lg' | 'md' | 'sm' | 'xs'
    disabled?: boolean
    icon?: string
}

defineProps<Props>()
defineEmits<{
    (e: 'click'): void
}>()
</script>
