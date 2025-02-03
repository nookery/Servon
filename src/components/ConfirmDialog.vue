<script setup lang="ts">
/**
 * 通用确认对话框组件
 * 
 * @example
 * ```vue
 * <template>
 *   <!-- 基础用法 -->
 *   <ConfirmDialog
 *     v-model:show="showConfirm"
 *     title="确认删除"
 *     message="是否确认删除此项？"
 *     @confirm="handleConfirm"
 *   />
 * 
 *   <!-- 自定义类型和按钮文本 -->
 *   <ConfirmDialog
 *     v-model:show="showConfirm"
 *     title="危险操作"
 *     message="此操作不可撤销，是否继续？"
 *     type="error"
 *     confirm-text="删除"
 *     cancel-text="返回"
 *     @confirm="handleConfirm"
 *     @cancel="handleCancel"
 *   />
 * </template>
 * ```
 */

import { defineProps, defineEmits } from 'vue'

interface Props {
    /** 控制对话框显示/隐藏 */
    show: boolean
    /** 对话框标题 */
    title: string
    /** 对话框消息内容 */
    message: string
    /** 
     * 对话框类型，影响图标和确认按钮的样式
     * - warning: 警告（默认）
     * - error: 错误/危险
     * - info: 信息
     */
    type?: 'warning' | 'error' | 'info'
    /** 确认按钮文本 */
    confirmText?: string
    /** 取消按钮文本 */
    cancelText?: string
}

/**
 * 组件属性默认值
 */
const props = withDefaults(defineProps<Props>(), {
    type: 'warning',
    confirmText: '确认',
    cancelText: '取消'
})

/**
 * 组件事件
 * - update:show: 更新显示状态
 * - confirm: 点击确认按钮
 * - cancel: 点击取消按钮
 */
const emit = defineEmits<{
    (e: 'update:show', value: boolean): void
    (e: 'confirm'): void
    (e: 'cancel'): void
}>()

const close = () => {
    emit('update:show', false)
}

const handleConfirm = () => {
    emit('confirm')
}

const handleCancel = () => {
    emit('cancel')
    close()
}

// 根据类型获取图标路径
const getIconPath = () => {
    switch (props.type) {
        case 'error':
            return 'M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z'
        case 'info':
            return 'M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z'
        default: // warning
            return 'M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z'
    }
}

// 根据类型获取图标颜色
const getIconColor = () => {
    switch (props.type) {
        case 'error':
            return 'stroke-error'
        case 'info':
            return 'stroke-info'
        default: // warning
            return 'stroke-warning'
    }
}
</script>

<template>
    <dialog class="modal" :class="{ 'modal-open': show }">
        <div class="modal-box p-0">
            <div role="alert" class="alert">
                <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" class="h-6 w-6 shrink-0"
                    :class="getIconColor()">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" :d="getIconPath()">
                    </path>
                </svg>
                <div class="flex-1">
                    <h3 class="font-bold">{{ title }}</h3>
                    <div class="text-sm">{{ message }}</div>
                </div>
                <div class="flex gap-2">
                    <button class="btn btn-sm" @click="handleCancel">{{ cancelText }}</button>
                    <button class="btn btn-sm" :class="{
                        'btn-error': type === 'error',
                        'btn-warning': type === 'warning',
                        'btn-info': type === 'info'
                    }" @click="handleConfirm">
                        {{ confirmText }}
                    </button>
                </div>
            </div>
        </div>
        <form method="dialog" class="modal-backdrop">
            <button @click="handleCancel">关闭</button>
        </form>
    </dialog>
</template>

<style scoped>
/* 如果需要添加组件特定的样式 */
</style>