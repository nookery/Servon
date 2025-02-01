// @ts-check
import { defineConfig } from 'astro/config';
import starlight from '@astrojs/starlight';
import tailwindcss from '@tailwindcss/vite';
import vue from '@astrojs/vue';

// https://astro.build/config
export default defineConfig({
    vite: {
        // @ts-ignore
        plugins: [tailwindcss()],
    server: {
        proxy: {
            '/api': {
                target: 'http://localhost:8080',
                changeOrigin: true,
            },
        },
    },
    },
    integrations: [
        starlight({
        title: 'My Docs',
        social: {
            github: 'https://github.com/withastro/starlight',
        },
        sidebar: [
            {
                label: 'Web API',
                items: [
                    { label: '基本信息', slug: 'web_api/basic' },
                ],
            },
            {
                label: 'Reference',
                autogenerate: { directory: 'reference' },
            },
        ],
    }),
        vue()]
});