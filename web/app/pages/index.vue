<script setup lang="ts">
const apiUrls = useApiUrls();
const { data, error } = await useFetch(apiUrls.anime.listings.recentlyUpdated, {
    transform: (res) => PaginatedAnimeSchema.parse(res),
});

const errorToaster = useErrorToaster();
watch(
    error,
    (err) => {
        if (!err) return;
        errorToaster(err);
    },
    { immediate: true },
);
</script>

<template>
    <div v-if="!error && data" class="grid grid-cols-6 w-full gap-4">
        <NuxtLink
            v-for="anime in data.items"
            :key="anime.id"
            :to="`/anime/${anime.id}`"
            class="group relative overflow-hidden rounded-md border transition hover:scale-105"
        >
            <div
                class="bg-card aspect-[3/4] w-full min-w-full overflow-hidden rounded-md"
            >
                <NuxtImg
                    :src="anime.imageUrl"
                    :alt="`${anime.ename} image`"
                    class="aspect-[3/4] w-full scale-105 rounded-md object-cover transition group-hover:scale-100"
                    loading="lazy"
                />
            </div>
            <div
                class="from-background absolute bottom-0 left-0 top-0 flex w-full flex-col justify-end bg-gradient-to-t to-transparent p-3"
            >
                <h3
                    class="font-sora line-clamp-2 text-sm font-bold md:text-base"
                >
                    {{ anime.jname }}
                </h3>
                <p
                    class="text-muted-foreground mt-1 hidden text-xs md:block md:text-sm"
                >
                    {{ anime.genre }}
                </p>
                <p class="text-muted-foreground mt-1 text-xs md:text-sm">
                    {{ anime.lastEpisode }} episodes
                </p>
            </div>
        </NuxtLink>
    </div>
</template>
