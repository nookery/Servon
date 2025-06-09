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
    tooltipPosition?: 'top' | 'bottom' | 'left' | 'right' | 'none'  // tooltip 位置配置，none表示不显示
}

const props = withDefaults(defineProps<Props>(), {
    variant: 'default',
    size: 'md',
    circle: false,
    active: false,
    disabled: false,
    loading: false,
    tooltipPosition: 'none'  // 默认不显示
})

// Compute button classes
const buttonClass = computed(() => {
    const classes = ['btn', 'transition-all', 'duration-200']

    // Only add tooltip class if tooltipPosition is not 'none'
    if (props.tooltipPosition !== 'none') {
        classes.push('tooltip')

        // Add tooltip position class
        if (props.tooltipPosition !== 'top') {
            classes.push(`tooltip-${props.tooltipPosition}`)
        }
    }

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
    <button :class="buttonClass" :data-tip="tooltipPosition !== 'none' ? title : undefined"
        :title="tooltipPosition !== 'none' ? title : undefined" :disabled="disabled" v-bind="$attrs"
        class="flex items-center justify-center">
        <i v-if="icon" :class="[icon, loading ? 'hidden' : '']"></i>
        <slot v-if="$slots.default" :class="icon ? 'ml-2' : ''"></slot>
    </button>
</template>