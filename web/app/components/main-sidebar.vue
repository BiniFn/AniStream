<script lang="ts" setup>
import {
    Home,
    TrendingUp,
    Bookmark,
    Settings,
    Play,
    Calendar,
    Award,
    Heart,
    Zap,
    Crown,
    Film,
} from "lucide-vue-next";

const mainNavItems = [
    {
        title: "Dashboard",
        url: "#",
        icon: Home,
        isActive: true,
    },
    {
        title: "Trending",
        url: "#",
        icon: TrendingUp,
    },
    {
        title: "Top Rated",
        url: "#",
        icon: Crown,
    },
    {
        title: "New Releases",
        url: "#",
        icon: Zap,
    },
    {
        title: "My Library",
        url: "#",
        icon: Bookmark,
    },
];

const genreItems = [
    {
        title: "Action",
        url: "#",
        icon: Play,
    },
    {
        title: "Romance",
        url: "#",
        icon: Heart,
    },
    {
        title: "Comedy",
        url: "#",
        icon: Award,
    },
    {
        title: "Drama",
        url: "#",
        icon: Calendar,
    },
    {
        title: "Movies",
        url: "#",
        icon: Film,
    },
];
</script>

<template>
    <Sidebar collapsible="icon" class="border-r border-border/50">
        <SidebarHeader class="border-b border-border/50 pb-4">
            <div
                class="flex items-center gap-3 px-4 py-2 group-data-[collapsible=icon]:px-0 group-data-[collapsible=icon]:justify-center"
            >
                <div
                    class="flex h-10 w-10 items-center justify-center rounded-xl bg-gradient-to-br from-primary to-primary/80 shadow-lg flex-shrink-0"
                >
                    <Play class="h-5 w-5 text-primary-foreground" />
                </div>
                <div
                    class="group-data-[collapsible=icon]:hidden min-w-0 font-bold text-xl bg-gradient-to-r from-primary to-primary/80 bg-clip-text text-transparent truncate"
                >
                    AniWays
                </div>
            </div>
        </SidebarHeader>

        <SidebarContent class="px-4 group-data-[collapsible=icon]:px-0">
            <SidebarGroup>
                <SidebarGroupLabel
                    class="text-xs font-semibold text-muted-foreground/80 uppercase tracking-wider mb-2 px-0 group-data-[collapsible=icon]:sr-only"
                >
                    Browse
                </SidebarGroupLabel>
                <SidebarGroupContent>
                    <SidebarMenu class="space-y-1">
                        <SidebarMenuItem
                            v-for="item in mainNavItems"
                            :key="item.title"
                        >
                            <SidebarMenuButton
                                as-child
                                :is-active="item.isActive"
                                class="h-10 rounded-lg hover:bg-accent/50 data-[active=true]:bg-primary/10 data-[active=true]:text-primary data-[active=true]:border data-[active=true]:border-primary/20 px-3 group-data-[collapsible=icon]:px-0 group-data-[collapsible=icon]:justify-center"
                            >
                                <a
                                    :href="item.url"
                                    class="flex items-center gap-3 w-full"
                                >
                                    <component
                                        :is="item.icon"
                                        class="h-4 w-4 flex-shrink-0"
                                    />
                                    <span
                                        class="font-medium group-data-[collapsible=icon]:sr-only"
                                        >{{ item.title }}</span
                                    >
                                </a>
                            </SidebarMenuButton>
                        </SidebarMenuItem>
                    </SidebarMenu>
                </SidebarGroupContent>
            </SidebarGroup>

            <SidebarGroup>
                <SidebarGroupLabel
                    class="text-xs font-semibold text-muted-foreground/80 uppercase tracking-wider mb-2 px-0 group-data-[collapsible=icon]:sr-only"
                >
                    Genres
                </SidebarGroupLabel>
                <SidebarGroupContent>
                    <SidebarMenu class="space-y-1">
                        <SidebarMenuItem
                            v-for="item in genreItems"
                            :key="item.title"
                        >
                            <SidebarMenuButton
                                as-child
                                class="h-9 rounded-lg hover:bg-accent/50 px-3 group-data-[collapsible=icon]:px-0 group-data-[collapsible=icon]:justify-center"
                            >
                                <a
                                    :href="item.url"
                                    class="flex items-center gap-3 w-full"
                                >
                                    <component
                                        :is="item.icon"
                                        class="h-4 w-4 flex-shrink-0"
                                    />
                                    <span
                                        class="font-medium group-data-[collapsible=icon]:sr-only"
                                        >{{ item.title }}</span
                                    >
                                </a>
                            </SidebarMenuButton>
                        </SidebarMenuItem>
                    </SidebarMenu>
                </SidebarGroupContent>
            </SidebarGroup>
        </SidebarContent>

        <SidebarFooter
            class="border-t border-border/50 pt-4 px-4 group-data-[collapsible=icon]:px-2"
        >
            <SidebarMenu>
                <SidebarMenuItem>
                    <SidebarMenuButton
                        as-child
                        class="h-12 rounded-lg hover:bg-accent/50 px-3 group-data-[collapsible=icon]:px-0 group-data-[collapsible=icon]:justify-center"
                    >
                        <a href="#" class="flex items-center gap-3 w-full">
                            <div
                                class="h-8 w-8 rounded-full bg-primary/10 flex items-center justify-center flex-shrink-0"
                            >
                                <span class="text-primary font-semibold text-sm"
                                    >JD</span
                                >
                            </div>
                            <div
                                class="group-data-[collapsible=icon]:sr-only min-w-0 flex-1"
                            >
                                <p class="font-medium text-sm truncate">
                                    John Doe
                                </p>
                            </div>
                        </a>
                    </SidebarMenuButton>
                </SidebarMenuItem>
                <SidebarMenuItem>
                    <SidebarMenuButton
                        as-child
                        class="h-9 rounded-lg hover:bg-accent/50 px-3 group-data-[collapsible=icon]:px-0 group-data-[collapsible=icon]:justify-center"
                    >
                        <a href="#" class="flex items-center gap-3 w-full">
                            <Settings class="h-4 w-4 flex-shrink-0" />
                            <span
                                class="font-medium group-data-[collapsible=icon]:sr-only"
                                >Settings</span
                            >
                        </a>
                    </SidebarMenuButton>
                </SidebarMenuItem>
            </SidebarMenu>
        </SidebarFooter>
        <SidebarRail />
    </Sidebar>
</template>
