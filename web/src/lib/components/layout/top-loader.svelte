<script lang="ts">
	import nProgress from 'nprogress';
	import { onMount } from 'svelte';
	import { beforeNavigate, afterNavigate } from '$app/navigation';

	let timeout: NodeJS.Timeout | null = null;

	onMount(() => {
		nProgress.configure({
			showSpinner: false,
			trickle: true,
			trickleSpeed: 200,
			minimum: 0.08,
			easing: 'ease',
			speed: 200,
			template: '<div class="bar" role="bar"><div class="peg"></div></div>',
		});
	});

	beforeNavigate(() => {
		timeout = setTimeout(() => {
			nProgress.start();
		}, 500);
	});

	afterNavigate(() => {
		if (timeout) {
			clearTimeout(timeout);
			timeout = null;
		}
		nProgress.done();
	});
</script>

<svelte:head>
	<style>
		#nprogress {
			pointer-events: none;
		}

		#nprogress .bar {
			background: var(--primary);
			position: fixed;
			z-index: 1600;
			top: 0;
			left: 0;
			width: 100%;
			height: 3px;
		}

		#nprogress .peg {
			display: block;
			position: absolute;
			right: 0;
			width: 100px;
			height: 100%;
			box-shadow:
				0 0 10px var(--primary),
				0 0 5px var(--primary);
			opacity: 1;
			-webkit-transform: rotate(3deg) translate(0px, -4px);
			-ms-transform: rotate(3deg) translate(0px, -4px);
			transform: rotate(3deg) translate(0px, -4px);
		}

		.nprogress-custom-parent {
			overflow: hidden;
			position: relative;
		}

		.nprogress-custom-parent #nprogress .bar,
		.nprogress-custom-parent #nprogress .spinner {
			position: absolute;
		}

		@-webkit-keyframes nprogress-spinner {
			0% {
				-webkit-transform: rotate(0deg);
			}
			100% {
				-webkit-transform: rotate(360deg);
			}
		}
		@keyframes nprogress-spinner {
			0% {
				transform: rotate(0deg);
			}
			100% {
				transform: rotate(360deg);
			}
		}
	</style>
</svelte:head>
