<script setup lang="ts">
import { defineProps } from 'vue'

interface DockItemProps {
    label: string
    icon: string
    color: string
    isActive: boolean
    onClick: () => void
}

const props = defineProps<DockItemProps>()
</script>

<template>
    <div class="relative flex flex-col items-center cursor-pointer" @click="props.onClick">
        <div class="tooltip tooltip-top" :data-tip="props.label">
            <!-- Icon with custom color -->
            <div class="p-2 rounded-xl flex items-center justify-center transition-all duration-200 ease-in-out border-2 border-transparent hover:-translate-y-1 w-12 h-12"
                :class="[
                    props.isActive
                        ? 'bg-opacity-90 text-white'
                        : 'bg-base-100 hover:bg-opacity-90 hover:text-white'
                ]" :style="{
                    backgroundColor: props.isActive ? props.color : '',
                    color: props.isActive ? 'white' : props.color
                }">
                <i :class="[props.icon, 'text-2xl']"></i>
            </div>
        </div>

        <!-- Indicator dot for active item -->
        <div v-if="props.isActive" class="h-1 w-1 rounded-full mt-1" :style="{ backgroundColor: props.color }">
        </div>
    </div>
</template>