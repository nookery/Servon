<script setup lang="ts">
import { computed } from 'vue'

interface Props {
    icon: string                    // Remix icon 类名
    title?: string                  // 按钮提示文字
    variant?: 'default' | 'ghost' | 'primary' | 'error' | 'warning'   // 按钮样式变体
    size?: 'xs' | 'sm' | 'md' | 'lg'      // 按钮大小
    circle?: boolean                // 是否为圆形按钮
    active?: boolean                // 是否激活状态
}

const props = withDefaults(defineProps<Props>(), {
    variant: 'default',
    size: 'md',
    circle: false,
    active: false
})

// 计算按钮的 class
const buttonClass = computed(() => {
    const classes = ['btn']

    // 变体样式
    if (props.variant === 'ghost') {
        classes.push('btn-ghost')
    } else if (props.variant === 'primary') {
        classes.push('btn-primary')
    } else if (props.variant === 'error') {
        classes.push('btn-error')
    } else if (props.variant === 'warning') {
        classes.push('btn-warning')
    }

    // 尺寸
    if (props.size === 'xs') {
        classes.push('btn-xs')
    } else if (props.size === 'sm') {
        classes.push('btn-sm')
    } else if (props.size === 'lg') {
        classes.push('btn-lg')
    }

    // 圆形
    if (props.circle) {
        classes.push('btn-circle')
    }

    // 激活状态
    if (props.active) {
        classes.push('btn-active')
    }

    return classes.join(' ')
})
</script>

<template>
    <button :class="buttonClass" :title="title">
        <i :class="[icon, 'mr-1']"></i>
        <slot></slot>
    </button>
</template>