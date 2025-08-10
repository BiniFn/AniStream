<script setup lang="ts">
import { Search, Menu } from 'lucide-vue-next';

const items = [
    { label: 'Home', to: '/' },
    { label: 'Catalog', to: '/catalog' },
    { label: 'Season', to: '/season' },
    { label: 'Top', to: '/top' },
];

const route = useRoute();
const mobileOpen = ref(false);
const isActive = (to: string) => to === route.path;
const openSearch = () => {
    mobileOpen.value = true;
};

watch(
    () => route.path,
    () => {
        mobileOpen.value = false;
    },
);
</script>
<template>
    <nav class="sticky top-0 z-50 border-b bg-transparent">
        <div
            class="mx-auto flex w-full max-w-[1640px] items-center gap-4 px-4 py-3 lg:gap-8 lg:px-12 lg:py-4"
        >
            <Sheet v-model:open="mobileOpen">
                <SheetTrigger as-child>
                    <Button variant="outline" size="icon" class="lg:hidden">
                        <Menu class="size-5" />
                        <span class="sr-only">Open menu</span>
                    </Button>
                </SheetTrigger>

                <NuxtLink
                    class="font-bungee text-2xl font-bold lg:text-3xl"
                    to="/"
                >
                    ANIWAYS
                </NuxtLink>

                <ul class="mx-1 hidden gap-4 lg:flex">
                    <li v-for="item in items" :key="item.label">
                        <NuxtLink
                            :to="item.to"
                            :data-active="isActive(item.to)"
                            class="text-muted-foreground hover:text-foreground data-[active=true]:text-foreground"
                        >
                            {{ item.label }}
                        </NuxtLink>
                    </li>
                </ul>

                <div class="relative ml-auto hidden w-full lg:block">
                    <label for="search" class="sr-only">Search</label>
                    <Search
                        class="text-muted-foreground pointer-events-none absolute top-1/2 left-3 size-4 -translate-y-1/2"
                    />
                    <Input
                        id="search"
                        placeholder="Search"
                        class="bg-muted/50 w-full pl-10"
                    />
                </div>

                <div class="hidden items-center gap-3 lg:flex">
                    <Button variant="secondary" as-child>
                        <NuxtLink to="/login">Login</NuxtLink>
                    </Button>
                    <Button as-child>
                        <NuxtLink to="/signup">Get Started</NuxtLink>
                    </Button>
                </div>

                <Button
                    variant="ghost"
                    size="icon"
                    class="ml-auto lg:hidden"
                    aria-label="Search"
                    @click="openSearch"
                >
                    <Search class="size-5" />
                </Button>

                <SheetContent side="left" class="w-80 p-6">
                    <SheetHeader class="mb-4">
                        <SheetTitle class="font-bungee text-2xl"
                            >ANIWAYS</SheetTitle
                        >
                    </SheetHeader>

                    <div class="relative mb-4">
                        <label for="m-search" class="sr-only">Search</label>
                        <Search
                            class="text-muted-foreground pointer-events-none absolute top-1/2 left-3 size-4 -translate-y-1/2"
                        />
                        <Input
                            id="m-search"
                            placeholder="Search"
                            class="bg-muted/50 w-full pl-10"
                        />
                    </div>

                    <nav class="grid gap-1">
                        <SheetClose
                            v-for="item in items"
                            :key="item.label"
                            as-child
                        >
                            <NuxtLink
                                :to="item.to"
                                :data-active="isActive(item.to)"
                                class="text-muted-foreground hover:bg-muted/60 hover:text-foreground data-[active=true]:bg-muted data-[active=true]:text-foreground block rounded-md px-3 py-2 text-sm"
                            >
                                {{ item.label }}
                            </NuxtLink>
                        </SheetClose>
                    </nav>

                    <div class="mt-6 flex gap-2">
                        <SheetClose as-child>
                            <Button variant="secondary" class="flex-1" as-child>
                                <NuxtLink to="/login">Login</NuxtLink>
                            </Button>
                        </SheetClose>
                        <SheetClose as-child>
                            <Button class="flex-1" as-child>
                                <NuxtLink to="/signup">Get Started</NuxtLink>
                            </Button>
                        </SheetClose>
                    </div>
                </SheetContent>
            </Sheet>
        </div>
    </nav>
</template>
