<script setup lang="ts">
interface WebhookData {
    id: string
    type: string
    timestamp: number
    payload: any
}

const props = defineProps<{
    webhook: WebhookData | null
}>()

function formatTimestamp(timestamp: number) {
    return new Date(timestamp).toLocaleString()
}
</script>

<template>
    <div class="p-4">
        <div v-if="webhook?.payload" class="space-y-4">
            <!-- Basic Info -->
            <div class="card bg-base-200">
                <div class="card-body">
                    <h3 class="card-title text-base">Basic Information</h3>
                    <div class="grid grid-cols-2 gap-2">
                        <div>Event Type: {{ webhook.type }}</div>
                        <div>Time: {{ formatTimestamp(webhook.timestamp) }}</div>
                        <div v-if="webhook.payload.action">Action: {{ webhook.payload.action }}</div>
                    </div>
                </div>
            </div>

            <!-- Repository Info -->
            <div v-if="webhook.payload.repository" class="card bg-base-200">
                <div class="card-body">
                    <h3 class="card-title text-base">Repository</h3>
                    <div class="grid grid-cols-2 gap-2">
                        <div>Name: {{ webhook.payload.repository.full_name }}</div>
                        <div>Visibility: {{ webhook.payload.repository.visibility }}</div>
                        <div>
                            <a :href="webhook.payload.repository.html_url" target="_blank" class="link">
                                View on GitHub
                            </a>
                        </div>
                    </div>
                </div>
            </div>

            <!-- Sender Info -->
            <div v-if="webhook.payload.sender" class="card bg-base-200">
                <div class="card-body">
                    <h3 class="card-title text-base">Sender</h3>
                    <div class="grid grid-cols-2 gap-2">
                        <div class="flex items-center gap-2">
                            <img :src="webhook.payload.sender.avatar_url" class="w-6 h-6 rounded-full"
                                :alt="webhook.payload.sender.login" />
                            {{ webhook.payload.sender.login }}
                        </div>
                        <div>
                            <a :href="webhook.payload.sender.html_url" target="_blank" class="link">
                                View Profile
                            </a>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>