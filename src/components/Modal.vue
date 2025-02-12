<template>
    <dialog class="modal" :open="show">
        <div class="modal-box p-0 w-11/12 max-w-5xl" :class="[boxClass, {
            'border-2 border-info': !error,
            'border-2 border-error': error
        }]">
            <div class="flex justify-between items-center px-4 py-2 bg-base-200">
                <h3 class="font-bold text-lg" :class="[titleClass, { 'text-error': error }]">{{ title }}</h3>
                <button class="btn btn-sm btn-circle border border-info" @click="close" :disabled="loading">
                    ✕
                </button>
            </div>
            <div class="py-0">
                <slot></slot>
            </div>
            <div class="modal-action bg-base-200 px-4 p-0 m-0 -translate-y-4">
                <slot name="actions"></slot>
            </div>
        </div>
    </dialog>
</template>

<script setup lang="ts">
import { defineEmits } from 'vue'

defineProps<{
    show: boolean
    title: string
    loading?: boolean
    boxClass?: string
    titleClass?: string
    error?: boolean
}>()

const emit = defineEmits<{
    (e: 'update:show', value: boolean): void
}>()

const close = () => {
    emit('update:show', false)
}
</script>