/**
 * Toast 提示 Hook
 * 
 * @example
 * ```ts
 * const toast = useToast()
 * 
 * toast.success('操作成功')
 * toast.error('操作失败')
 * toast.show('提示信息')
 * ```
 */

import { toastState } from './useToastState'

export function useToast() {
    return {
        /**
         * 显示成功提示
         */
        success(message: string) {
            toastState.add(message, 'success')
        },

        /**
         * 显示错误提示
         */
        error(message: string) {
            toastState.add(message, 'error')
        },

        /**
         * 显示普通提示
         */
        info(message: string) {
            toastState.add(message, 'info')
        },

        show(message: string) {
            toastState.add(message, 'info')
        },

        /**
         * 显示警告提示
         */
        warning(message: string) {
            toastState.add(message, 'warning')
        }
    }
} 