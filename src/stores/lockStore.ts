import { defineStore } from 'pinia'

export const useLockStore = defineStore('lock', {
    state: () => ({
        isLocked: localStorage.getItem('isLocked') !== 'false'
    }),
    actions: {
        lock() {
            this.isLocked = true
            localStorage.setItem('isLocked', 'true')
        },
        unlock() {
            this.isLocked = false
            localStorage.setItem('isLocked', 'false')
        }
    }
})