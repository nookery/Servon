<template>
    <div class="min-h-screen bg-gray-50">
        <!-- 侧边栏 -->
        <div class="fixed inset-y-0 left-0 w-64 bg-white border-r border-gray-200">
            <div class="p-4">
                <h1 class="text-xl font-bold text-gray-800">Servon</h1>
            </div>
            <nav class="mt-4">
                <NuxtLink to="/" class="flex items-center px-4 py-2 text-gray-700 hover:bg-gray-100"
                    :class="{ 'bg-gray-100': $route.path === '/' }">
                    <UIcon name="i-mdi-monitor-dashboard" class="w-5 h-5 mr-2" />
                    系统监控
                </NuxtLink>
                <NuxtLink to="/software" class="flex items-center px-4 py-2 text-gray-700 hover:bg-gray-100"
                    :class="{ 'bg-gray-100': $route.path === '/software' }">
                    <UIcon name="i-mdi-apps" class="w-5 h-5 mr-2" />
                    软件信息
                </NuxtLink>
            </nav>
        </div>

        <!-- 主内容区 -->
        <div class="ml-64">
            <slot />
        </div>
    </div>
</template>