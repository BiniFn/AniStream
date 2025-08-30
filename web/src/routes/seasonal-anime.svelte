<script lang="ts">
	import type { components } from '$lib/api/openapi';

	type Props = {
		seasonalAnime: components['schemas']['models.SeasonalAnimeResponse'];
	};

	let { seasonalAnime }: Props = $props();
</script>
