import { goto } from '$app/navigation';
import { convertComponentToHTML } from '$lib/utils';
import { ArkErrors, type } from 'arktype';
import artplayerPluginHlsControl from 'artplayer-plugin-hls-control';
import type Hls from 'hls.js';
import { Captions, LoaderCircle, Pause, SkipForward } from 'lucide-svelte';
import { amplifyVolumePlugin, skipPlugin, thumbnailPlugin, windowKeyBindPlugin } from './plugins';
import type { components } from '$lib/api/openapi';
import { PUBLIC_STREAMING_URL } from '$env/static/public';
import type { AppState } from '$lib/context/state.svelte';

type StreamInfo = components['schemas']['models.StreamingDataResponse'];

type Props = {
	id: string;
	appState: AppState;
	container: HTMLDivElement;
	source: StreamInfo;
	nextEpisodeUrl: string | undefined;
	updateLibrary: () => Promise<void>;
};

const artplayerSettingsSchema = type({
	times: 'Record<string, number>',
});

export const createArtPlayer = async ({
	id,
	appState,
	container,
	source,
	nextEpisodeUrl,
	updateLibrary,
}: Props) => {
	const thumbnails = source.tracks.find((track) => track.kind === 'thumbnails');
	const defaultSubtitle = source.tracks.find((track) => track.default && track.kind === 'captions');
	const isMobile = /Android|webOS|iPhone|iPad|iPod|BlackBerry|IEMobile|Opera Mini/i.test(
		navigator.userAgent,
	);

	const Hls = import('hls.js');
	const Artplayer = await import('artplayer').then((module) => module.default);

	const art = new Artplayer({
		id,
		container,
		hotkey: false,
		url: `${PUBLIC_STREAMING_URL}${source.source.proxyHls}`,
		setting: true,
		theme: 'var(--primary)',
		fullscreen: true,
		mutex: true,
		playbackRate: true,
		autoPlayback: true,
		autoOrientation: true,
		playsInline: true,
		type: 'm3u8',
		pip: !!/(chrome|edg|safari|opr)/i.exec(navigator.userAgent),
		airplay: true,
		miniProgressBar: true,
		gesture: false,
		icons: {
			loading: convertComponentToHTML(LoaderCircle, {
				size: 100,
				class: 'animate-spin text-primary',
				style: 'fill: none !important;',
			}),
		},
		subtitle: {
			url: defaultSubtitle ? `${PUBLIC_STREAMING_URL}${defaultSubtitle.url}` : '',
			type: 'vtt',
			encoding: 'utf-8',
			escape: false,
			style: {
				fontSize: isMobile ? '1rem' : '1.8rem',
			},
		},
		plugins: [
			artplayerPluginHlsControl({
				quality: {
					setting: true,
					getName: (level: { height: number }) => `${level.height}p`,
					title: 'Quality',
					auto: 'Auto',
				},
			}),
			thumbnailPlugin(thumbnails!),
			skipPlugin(source),
			windowKeyBindPlugin(),
			amplifyVolumePlugin(),
		],
		customType: {
			m3u8: (video, url, art) => {
				console.log('Custom type m3u8', url);
				Hls.then((module) => {
					const Hls = module.default;
					if (Hls.isSupported()) {
						if (art.hls) (art.hls as Hls).destroy();
						const hls = new Hls();
						hls.loadSource(url);
						hls.attachMedia(video);
						art.hls = hls;
						art.on('destroy', () => hls.destroy());
						// update art quality when hls quality changes
						hls.on(Hls.Events.LEVEL_SWITCHED, () => {
							const currentLevel = hls.levels[hls.currentLevel]?.height + 'p';
							const currentSetting = art.setting.find('hls-quality') as unknown as {
								selector: { default: boolean; html: string }[];
								tooltip: string;
							};

							if (
								currentSetting &&
								currentSetting.selector.find((item) => item.default)?.html !== currentLevel
							) {
								art.setting.update({
									...currentSetting,
									selector: currentSetting.selector.map((item) => ({
										...item,
										default: item.html === currentLevel,
									})),
									tooltip: currentLevel,
								});
							}

							art.notice.show = `Quality: ${currentLevel}`;
						});
					} else if (video.canPlayType('application/vnd.apple.mpegurl')) {
						video.src = url;
					} else {
						art.notice.show = 'Unsupported playback format: m3u8';
					}
				});
			},
		},
	});

	art.setting.add({
		icon: art.icons.setting,
		html: 'Player Settings',
		selector: [
			{
				icon: art.icons.play,
				html: 'Auto Play Episode',
				switch: appState.settings?.autoPlayEpisode,
				onSwitch: () => appState.toggleSetting('autoPlayEpisode'),
			},
			{
				icon: convertComponentToHTML(SkipForward, { size: 22 }),
				html: 'Auto Play Next Episode',
				switch: appState.settings?.autoNextEpisode,
				onSwitch: () => appState.toggleSetting('autoNextEpisode'),
			},
			{
				icon: convertComponentToHTML(Pause, { size: 22 }),
				html: 'Auto Resume Episode',
				switch: appState.settings?.autoResumeEpisode,
				onSwitch: () => appState.toggleSetting('autoResumeEpisode'),
			},
		],
	});

	art.setting.add({
		icon: convertComponentToHTML(Captions, { size: 22, style: 'fill: none !important;' }),
		html: 'Captions',
		tooltip: defaultSubtitle?.label,
		selector: [
			{
				html: 'Off',
				default: false,
				url: '',
				off: true,
			},
			...source.tracks
				.filter((track) => track.kind === 'captions')
				.map((track) => ({
					default: track.default,
					html: track.label ?? 'Unknown',
					url: `${PUBLIC_STREAMING_URL}${track.url}`,
				})),
		],
		onSelect: (item) => {
			const url = item.url as unknown;
			if (typeof url !== 'string') return;
			art.subtitle.url = url;
			art.subtitle.show = !!url;
			return item.html;
		},
	});

	art.on('ready', () => {
		if (appState.settings?.autoResumeEpisode) {
			const time = artplayerSettingsSchema(
				JSON.parse(localStorage.getItem('artplayer_settings') ?? '{}'),
			);
			if (time instanceof ArkErrors === false && time.times[id]) {
				if (time.times[id] >= art.duration) return;
				art.currentTime = time.times[id];
			}
		}

		if (appState.settings?.autoPlayEpisode) {
			art.play();
		}

		if (sessionStorage.getItem('artplayer_resume_fullscreen') === 'true') {
			art.fullscreen = true;
			sessionStorage.removeItem('artplayer_resume_fullscreen');
		}
	});

	art.on('fullscreen', (isFullScreen) => {
		const base = isMobile ? 1 : 1.8;
		const screenWidth = window.screen.width;
		const videoWidth = container.clientWidth ?? 0;
		const fontSize = isFullScreen ? `${(screenWidth / videoWidth) * base}rem` : `${base}rem`;
		art.subtitle.style('fontSize', fontSize);
	});

	art.on('video:ended', async () => {
		await updateLibrary();

		if (nextEpisodeUrl && appState.settings?.autoNextEpisode) {
			if (art.fullscreen) {
				sessionStorage.setItem('artplayer_resume_fullscreen', 'true');
			} else {
				sessionStorage.removeItem('artplayer_resume_fullscreen');
			}

			await goto(nextEpisodeUrl);
		}
	});

	return art;
};
