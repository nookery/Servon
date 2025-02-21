<script setup lang="ts">
import { computed } from 'vue'

interface Props {
    icon?: string                    // Remix icon class name
    title?: string                   // Button tooltip text
    variant?: 'default' | 'ghost' | 'primary' | 'secondary' | 'accent' | 'info' | 'success' | 'error' | 'warning' | 'link'   // Button style variant
    size?: 'xs' | 'sm' | 'md' | 'lg' // Button size
    circle?: boolean                 // Whether it's a circular button
    active?: boolean                 // Whether in active state
    disabled?: boolean              // Whether button is disabled
    loading?: boolean               // Loading state
    customClass?: string           // Additional custom classes
}

const props = withDefaults(defineProps<Props>(), {
    variant: 'default',
    size: 'md',
    circle: false,
    active: false,
    disabled: false,
    loading: false
})

// Compute button classes
const buttonClass = computed(() => {
    const classes = ['btn', 'transition-all', 'duration-200']

    // Add variant classes
    if (props.variant !== 'default') {
        classes.push(`btn-${props.variant}`)
    }

    // Add size class
    if (props.size !== 'md') {
        classes.push(`btn-${props.size}`)
    }

    // Add states
    if (props.circle) classes.push('btn-circle')
    if (props.active) classes.push('btn-active')
    if (props.disabled) classes.push('btn-disabled')
    if (props.loading) classes.push('loading')

    // Add hover effects (only when not disabled or loading)
    if (!props.disabled && !props.loading) {
        classes.push('hover:scale-105', 'hover:shadow-md')
    }

    // Add custom classes
    if (props.customClass) {
        classes.push(props.customClass)
    }

    return classes.join(' ')
})
</script>

<template>
    <button :class="buttonClass" :title="title" :disabled="disabled" v-bind="$attrs">
        <i v-if="icon" :class="[icon, loading ? 'hidden' : 'mr-1']"></i>
        <slot></slot>
    </button>
</template>