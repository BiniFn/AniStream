import { goto } from '$app/navigation';
import { convertComponentToHTML } from '$lib/utils';
import { ArkErrors, type } from 'arktype';
import artplayerPluginHlsControl from 'artplayer-plugin-hls-control';
import Hls from 'hls.js';
import { Captions, LoaderCircle, Pause, SkipForward } from 'lucide-svelte';
import { amplifyVolumePlugin, skipPlugin, thumbnailPlugin, windowKeyBindPlugin } from './plugins';
import type { components } from '$lib/api/openapi';
import { PUBLIC_STREAMING_URL } from '$env/static/public';
import type { AppState } from '$lib/context/state.svelte';
import Artplayer from 'artplayer';
import { isElectron } from '$lib/hooks/is-electron';

type StreamInfo = components['schemas']['models.StreamingDataResponse'];

type Props = {
	id: string;
	appState: AppState;
	container: HTMLDivElement;
	source: StreamInfo;
	nextEpisodeUrl: string | null;
	updateLibrary: () => Promise<void>;
};

const artplayerSettingsSchema = type({
	times: 'Record<string, number>',
});

export const createArtPlayer = ({
	id,
	appState,
	container,
	source,
	nextEpisodeUrl,
	updateLibrary,
}: Props) => {
	const abortController = new AbortController();
	const thumbnails = source.tracks.find((track) => track.kind === 'thumbnails');
	const defaultSubtitle = source.tracks.find((track) => track.default && track.kind === 'captions');
	const isMobile = /Android|webOS|iPhone|iPad|iPod|BlackBerry|IEMobile|Opera Mini/i.test(
		navigator.userAgent,
	);
	const inElectron = isElectron();

	// In Electron, use raw URLs directly (no CORS issues). In browser, use proxy.
	const getStreamUrl = () => {
		if (inElectron && source.source.hls) {
			return source.source.hls;
		}
		return `${PUBLIC_STREAMING_URL}${source.source.proxyHls}`;
	};

	const getTrackUrl = (track: { raw: string; url: string }) => {
		if (inElectron) {
			return track.raw;
		}
		return `${PUBLIC_STREAMING_URL}${track.url}`;
	};

	const plugins = [
		artplayerPluginHlsControl({
			quality: {
				setting: true,
				getName: (level: { height: number }) => `${level.height}p`,
				title: 'Quality',
				auto: 'Auto',
			},
		}),
		skipPlugin(source),
		windowKeyBindPlugin(),
		amplifyVolumePlugin(),
	];

	if (thumbnails) {
		plugins.push(thumbnailPlugin(getTrackUrl(thumbnails)));
	}

	const art = new Artplayer({
		id,
		container,
		hotkey: false,
		url: getStreamUrl(),
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
			url: defaultSubtitle ? getTrackUrl(defaultSubtitle) : '',
			type: 'vtt',
			encoding: 'utf-8',
			escape: false,
			style: {
				fontSize: isMobile ? '1rem' : '1.8rem',
			},
		},
		plugins,
		customType: {
			m3u8: (video, url, art) => {
				console.log('Custom type m3u8', url);

				if (Hls.isSupported()) {
					if (art.hls) {
						const existingHls = art.hls as Hls;
						existingHls.detachMedia();
						existingHls.stopLoad();
						existingHls.destroy();
						art.hls = null;
					}

					const hls = new Hls({
						fetchSetup: (context, initParams) => {
							return new Request(context.url, {
								...initParams,
								signal: abortController.signal,
							});
						},
					});
					hls.loadSource(url);
					hls.attachMedia(video);
					art.hls = hls;

					const cleanup = () => {
						abortController.abort();
						if (hls) {
							hls.detachMedia();
							hls.stopLoad();
							hls.destroy();
							art.hls = null;
						}
					};

					art.once('destroy', cleanup);
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
					url: getTrackUrl(track),
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
				const savedTime = Math.floor(time.times[id]);
				const duration = Math.floor(art.duration);
				if (savedTime >= duration - 5) return;
				art.currentTime = time.times[id];
			}
		}

		if (appState.settings?.autoPlayEpisode) {
			art.play();
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
			await goto(nextEpisodeUrl);
		}
	});

	return art;
};
