<template>
    <TransitionRoot appear :show="show" as="template">
        <div class="fixed bottom-4 right-4 z-50">
            <TransitionChild as="template" enter="transform ease-out duration-300 transition"
                enter-from="translate-y-2 opacity-0 sm:translate-y-0 sm:translate-x-2"
                enter-to="translate-y-0 opacity-100 sm:translate-x-0" leave="transition ease-in duration-100"
                leave-from="opacity-100" leave-to="opacity-0">
                <div :class="[
                    'max-w-sm w-full shadow-lg rounded-lg pointer-events-auto overflow-hidden',
                    type === 'success' ? 'bg-green-50' : 'bg-red-50'
                ]">
                    <div class="p-4">
                        <div class="flex items-start">
                            <div class="flex-shrink-0">
                                <CheckCircleIcon v-if="type === 'success'" class="h-6 w-6 text-green-400"
                                    aria-hidden="true" />
                                <XCircleIcon v-else class="h-6 w-6 text-red-400" aria-hidden="true" />
                            </div>
                            <div class="ml-3 w-0 flex-1 pt-0.5">
                                <p :class="[
                                    'text-sm font-medium',
                                    type === 'success' ? 'text-green-800' : 'text-red-800'
                                ]">
                                    {{ message }}
                                </p>
                            </div>
                            <div class="ml-4 flex-shrink-0 flex">
                                <button type="button" :class="[
                                    'rounded-md inline-flex text-sm font-medium focus:outline-none focus:ring-2 focus:ring-offset-2',
                                    type === 'success'
                                        ? 'text-green-500 hover:text-green-600 focus:ring-green-500'
                                        : 'text-red-500 hover:text-red-600 focus:ring-red-500'
                                ]" @click="onClose">
                                    <span class="sr-only">关闭</span>
                                    <XMarkIcon class="h-5 w-5" aria-hidden="true" />
                                </button>
                            </div>
                        </div>
                    </div>
                </div>
            </TransitionChild>
        </div>
    </TransitionRoot>
</template>

<script setup lang="ts">
import { TransitionRoot, TransitionChild } from '@headlessui/vue'
import { CheckCircleIcon, XCircleIcon, XMarkIcon } from '@heroicons/vue/24/outline'

defineProps<{
    show: boolean
    message: string
    type: 'success' | 'error'
}>()

const emit = defineEmits<{
    (e: 'close'): void
}>()

function onClose() {
    emit('close')
}
</script>