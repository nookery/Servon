import { createRouter, createWebHistory } from 'vue-router'
import CronTasks from '../views/CronTasks.vue'

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
        },
        {
            path: '/deploy',
            component: () => import('../views/Deploy.vue')
        },
        {
            path: '/cron',
            name: 'cron',
            component: CronTasks,
            meta: {
                title: '定时任务'
            }
        }
    ]
})

export default router 