<script setup lang="ts">
import { PaginatedAnimeSchema } from "~/types/anime";

const config = useRuntimeConfig();
const { data, error } = await useFetch(
    `${config.public.apiBase}/anime/listings/recently-updated`,
    {
        transform: (res) => PaginatedAnimeSchema.parse(res),
    },
);
</script>

<template>
    <div>
        <div v-if="error">Error: {{ error.data.error }}</div>
        <pre v-else-if="data">{{ JSON.stringify(data, null, 2) }}</pre>
    </div>
</template>
