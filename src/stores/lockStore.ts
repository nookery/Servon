import { defineStore } from 'pinia'

export const useLockStore = defineStore('lock', {
    state: () => ({
        isLocked: true
    }),
    actions: {
        lock() {
            this.isLocked = true
        },
        unlock() {
            this.isLocked = false
        }
    }
})