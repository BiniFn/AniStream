<script setup lang="ts">
import 'vue-sonner/style.css';
useHead({
    bodyAttrs: {
        class: 'dark',
    },
    link: [
        {
            rel: 'preconnect',
            href: 'https://fonts.gstatic.com',
        },
        {
            rel: 'preconnect',
            href: 'https://fonts.googleapis.com',
        },
    ],
});
</script>

<template>
    <NuxtLayout>
        <NuxtRouteAnnouncer />
        <NuxtPage />
    </NuxtLayout>
    <Toaster rich-colors theme="dark" />
</template>
