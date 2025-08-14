<script lang="ts">
	import Button from '$lib/components/ui/button/button.svelte';

	let number = $state(0);
</script>

<Button onclick={() => number++}>
	Increment {number}
</Button>
