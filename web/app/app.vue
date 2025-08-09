<script setup lang="ts">
useHead({
    bodyAttrs: {
        class: "dark",
    },
});
</script>

<template>
    <SidebarProvider>
        <MainSidebar />
        <main>
            <SidebarTrigger />
            <NuxtRouteAnnouncer />
            <NuxtPage />
        </main>
    </SidebarProvider>
</template>
