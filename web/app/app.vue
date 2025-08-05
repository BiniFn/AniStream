<script setup lang="ts">
useHead({
    bodyAttrs: {
        class: "dark",
    },
});
</script>

<template>
    <div>
        <NuxtRouteAnnouncer />
        <NuxtPage />
    </div>
</template>
