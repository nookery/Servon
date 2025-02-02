import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
    history: createWebHistory(),
    routes: [
        {
            path: '/',
            redirect: '/dashboard'
        },
        {
            path: '/dashboard',
            component: () => import('../views/Dashboard.vue')
        },
        {
            path: '/software',
            component: () => import('../views/Software.vue')
        },
        {
            path: '/processes',
            component: () => import('../views/Processes.vue')
        },
        {
            path: '/files',
            component: () => import('../views/Files.vue')
        },
        {
            path: '/ports',
            component: () => import('../views/Ports.vue')
        }
    ]
})

export default router 