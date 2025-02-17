<script setup lang="ts">
import { ref } from 'vue'
import type { IPInfo } from '../models/IpInfo'

defineProps<{
    ipInfo: IPInfo | null
    isVisible: boolean
}>()
</script>

<template>
    <div v-if="isVisible" class="absolute right-0 mt-2 w-96 card bg-base-100 shadow-xl z-50">
        <div class="card-body p-4">
            <h3 class="card-title text-sm mb-2">IP 信息</h3>

            <!-- Public IPs -->
            <div class="mb-4">
                <div class="text-sm font-semibold mb-1">公网IP</div>
                <div class="text-sm">IPv4: {{ ipInfo?.public_ip || '未获取' }}</div>
                <div class="text-sm">IPv6: {{ ipInfo?.public_ipv6 || '未获取' }}</div>
            </div>

            <!-- Local IPs -->
            <div class="mb-4">
                <div class="text-sm font-semibold mb-1">本地IP</div>
                <div v-for="ip in ipInfo?.local_ips" :key="ip.ip" class="text-sm mb-1">
                    <div>{{ ip.interface }}: {{ ip.ip }}</div>
                    <div class="text-xs opacity-70">掩码: {{ ip.netmask }}</div>
                </div>
            </div>

            <!-- DNS Servers -->
            <div class="mb-4">
                <div class="text-sm font-semibold mb-1">DNS 服务器</div>
                <div v-for="dns in ipInfo?.dns_servers" :key="dns" class="text-sm">
                    {{ dns }}
                </div>
            </div>
        </div>
    </div>
</template>
