<script setup lang="ts">
import "vue-sonner/style.css";
useHead({
    bodyAttrs: {
        class: "dark",
    },
    link: [
        {
            rel: "preconnect",
            href: "https://fonts.gstatic.com",
        },
        {
            rel: "preconnect",
            href: "https://fonts.googleapis.com",
        },
    ],
});
</script>

<template>
    <main class="p-4">
        <NuxtRouteAnnouncer />
        <NuxtPage />
    </main>
    <Toaster rich-colors theme="dark" />
</template>
