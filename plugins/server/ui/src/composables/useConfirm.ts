/**
 * Confirm 确认对话框 Hook
 * 
 * @example
 * ```ts
 * const confirm = useConfirm()
 * 
 * // 基础用法
 * if (await confirm.show('删除确认', '确定要删除这个文件吗？')) {
 *   // 用户点击了确认
 * }
 * 
 * // 自定义选项
 * if (await confirm.error('危险操作', '此操作不可撤销', {
 *   confirmText: '删除',
 *   cancelText: '返回'
 * })) {
 *   // 用户点击了确认
 * }
 * ```
 */

import { confirmState, type ConfirmType } from './useConfirmState'

interface ConfirmOptions {
    confirmText?: string
    cancelText?: string
}

export function useConfirm() {
    return {
        /**
         * 显示确认对话框
         */
        show(title: string, message: string, options?: ConfirmOptions) {
            return createConfirm(title, message, 'warning', options)
        },

        /**
         * 显示警告确认对话框
         */
        warning(title: string, message: string, options?: ConfirmOptions) {
            return createConfirm(title, message, 'warning', options)
        },

        /**
         * 显示错误确认对话框
         */
        error(title: string, message: string, options?: ConfirmOptions) {
            return createConfirm(title, message, 'error', options)
        },

        /**
         * 显示信息确认对话框
         */
        info(title: string, message: string, options?: ConfirmOptions) {
            return createConfirm(title, message, 'info', options)
        }
    }
}

function createConfirm(
    title: string,
    message: string,
    type: ConfirmType,
    options?: ConfirmOptions
): Promise<boolean> {
    return new Promise((resolve) => {
        confirmState.options = {
            title,
            message,
            type,
            confirmText: options?.confirmText,
            cancelText: options?.cancelText
        }
        confirmState.resolve = resolve
        confirmState.show = true
    })
} 