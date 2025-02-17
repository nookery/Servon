<script setup lang="ts">
import { computed } from 'vue'

interface Props {
    icon: string                    // Remix icon 类名
    title?: string                  // 按钮提示文字
    variant?: 'default' | 'ghost'   // 按钮样式变体
    size?: 'sm' | 'md' | 'lg'      // 按钮大小
    circle?: boolean                // 是否为圆形按钮
}

const props = withDefaults(defineProps<Props>(), {
    variant: 'default',
    size: 'md',
    circle: false
})

// 计算按钮的 class
const buttonClass = computed(() => {
    const classes = ['btn']

    // 变体样式
    if (props.variant === 'ghost') {
        classes.push('btn-ghost')
    }

    // 尺寸
    if (props.size === 'sm') {
        classes.push('btn-sm')
    } else if (props.size === 'lg') {
        classes.push('btn-lg')
    }

    // 圆形
    if (props.circle) {
        classes.push('btn-circle')
    }

    return classes.join(' ')
})
</script>

<template>
    <button :class="buttonClass" :title="title">
        <i :class="[icon, 'text-xl text-primary']"></i>
        <slot></slot>
    </button>
</template>