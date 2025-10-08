import { buttonVariants } from '$lib/components/ui/button';
import type Artplayer from 'artplayer';
import type { components } from '$lib/api/openapi';
import { PUBLIC_STREAMING_URL } from '$env/static/public';

type StreamInfo = components['schemas']['models.StreamingDataResponse'];

export const thumbnailPlugin = (thumbnails: { raw: string; url: string }) => {
	return async (art: Artplayer) => {
		const {
			template: { $progress },
		} = art;

		// eslint-disable-next-line @typescript-eslint/no-explicit-any
		let timer: any = null;

		const url = `${PUBLIC_STREAMING_URL}${thumbnails.url}`;
		if (!url) return;

		const tns = await fetch(url)
			.then((res) => res.text())
			.then((res) =>
				res
					.split('\n')
					.filter((line) => line.trim())
					.slice(1),
			);

		const data: {
			start: number;
			end: number;
			url: string;
			x: number;
			y: number;
			w: number;
			h: number;
		}[] = [];

		tns.forEach((_, index) => {
			if (index % 3 !== 0) return;
			const time = tns[index + 1];
			const url = tns[index + 2];
			if (!time || !url) return;
			const start = time.split(' --> ')[0]!;
			const end = time.split(' --> ')[1]!;

			const startSeconds = start.split(':').reduce((acc, time, i) => {
				return acc + Number(time) * Math.pow(60, 2 - i);
			}, 0);

			const endSeconds = end.split(':').reduce((acc, time, i) => {
				return acc + Number(time) * Math.pow(60, 2 - i);
			}, 0);

			const [x, y, w, h] = url.split('#xywh=')[1]!.split(',').map(Number);

			data.push({
				start: startSeconds,
				end: endSeconds,
				url: `${PUBLIC_STREAMING_URL}/${url.split('#xywh=')[0]}`,
				x: x!,
				y: y!,
				w: w!,
				h: h!,
			});
		});

		art.controls.add({
			name: 'vtt-thumbnail',
			position: 'top',
			mounted($control) {
				$control.classList.add('art-control-thumbnails');
				art.on('setBar', async (type, percentage, event) => {
					const isMobile = /Android|webOS|iPhone|iPad|iPod|BlackBerry|IEMobile|Opera Mini/i.test(
						navigator.userAgent,
					);

					const isMobileDragging = type === 'played' && event && isMobile;

					if (type === 'hover' || isMobileDragging) {
						const width = $progress.clientWidth * percentage;
						const second = percentage * art.duration;
						$control.style.display = 'flex';

						const find = data.find((item) => item.start <= second && item.end >= second);

						if (!find) {
							$control.style.display = 'none';
							return;
						}

						if (width > 0 && width < $progress.clientWidth) {
							$control.style.backgroundImage = `url(${find.url})`;
							$control.style.height = `${find.h}px`;
							$control.style.width = `${find.w}px`;
							$control.style.backgroundPosition = `-${find.x}px -${find.y}px`;
							if (width <= find.w / 2) {
								$control.style.left = '0px';
							} else if (width > $progress.clientWidth - find.w / 2) {
								$control.style.left = `${$progress.clientWidth - find.w}px`;
							} else {
								$control.style.left = `${width - find.w / 2}px`;
							}
						} else {
							if (!isMobile) {
								$control.style.display = 'none';
							}
						}

						if (isMobileDragging) {
							if (timer) clearTimeout(timer);
							timer = setTimeout(() => {
								$control.style.display = 'none';
							}, 1000);
						}
					}
				});
			},
		});
	};
};

export const skipPlugin = (source: StreamInfo) => {
	return async (art: Artplayer) => {
		art.on('ready', () => {
			function addElement(title: string, start: number, end: number) {
				const startPercentage = (start / art.duration) * 100;
				const endPercentage = ((end - start) / art.duration) * 100;
				const highlightElement = art.template.$progress.querySelector('.art-progress-highlight');

				highlightElement?.insertAdjacentHTML(
					'beforeend',
					`<span data-time="${start}" data-text="${title}" style="left: ${startPercentage}%; width: ${endPercentage}% !important"></span>`,
				);
			}

			if (source.intro) {
				addElement('Opening', source.intro.start, source.intro.end);
			}

			if (source.outro) {
				addElement('Ending', source.outro.start, source.outro.end);
			}
		});

		art.on('video:timeupdate', () => {
			if (
				source.intro &&
				art.currentTime >= source.intro.start &&
				art.currentTime <= source.intro.end &&
				!art.controls['opening']
			) {
				art.controls.add({
					name: 'opening',
					position: 'top',
					html: `<button class="${buttonVariants({ class: 'absolute bottom-6 right-0' })}">Skip Opening</button>`,
					click: (_, e) => {
						e.preventDefault();
						e.stopPropagation();

						art.seek = source.intro!.end;
						art.notice.show = 'Skipped Opening';
					},
				});
			}

			if (
				source.outro &&
				art.currentTime >= source.outro.start &&
				art.currentTime <= source.outro.end &&
				!art.controls['ending']
			) {
				art.controls.add({
					name: 'ending',
					position: 'top',
					html: `<button class="${buttonVariants({ class: 'absolute bottom-6 right-0' })}">Skip Ending</button>`,
					click: (_, e) => {
						e.preventDefault();
						e.stopPropagation();

						art.seek = source.outro!.end;
						art.notice.show = 'Skipped Ending';
					},
				});
			}

			if (
				source.intro &&
				(art.currentTime < source.intro.start || art.currentTime > source.intro.end) &&
				art.controls['opening']
			) {
				art.controls.remove('opening');
			}

			if (
				source.outro &&
				(art.currentTime < source.outro.start || art.currentTime > source.outro.end) &&
				art.controls['ending']
			) {
				art.controls.remove('ending');
			}
		});
	};
};

export const windowKeyBindPlugin = () => {
	return async (art: Artplayer) => {
		art.events.proxy(window, 'keydown', (e) => {
			if (e instanceof KeyboardEvent === false) return;
			if (e.target instanceof HTMLInputElement || e.target instanceof HTMLTextAreaElement) return;

			if (Object.keys(art.hotkey.keys).includes(e.code)) {
				e.preventDefault();
				art.hotkey.keys[e.code]?.forEach((fn) => fn?.(e));
			}
		});

		art.on('ready', () => {
			art.hotkey.add('KeyF', () => {
				art.fullscreen = !art.fullscreen;
			});

			art.hotkey.add('KeyM', () => {
				art.muted = !art.muted;
			});

			art.hotkey.add('Space', () => {
				art.toggle();
			});

			art.hotkey.add('ArrowLeft', () => {
				art.backward = 10;
			});

			art.hotkey.add('ArrowRight', () => {
				art.forward = 10;
			});

			art.hotkey.add('ArrowUp', () => {
				art.volume += 0.1;
			});

			art.hotkey.add('ArrowDown', () => {
				art.volume -= 0.1;
			});
		});
	};
};

export const amplifyVolumePlugin = () => {
	return async (art: Artplayer) => {
		art.on('ready', () => {
			art.volume = 100;

			const context = new AudioContext();
			const source = context.createMediaElementSource(art.video);
			const gainNode = context.createGain();
			source.connect(gainNode);
			gainNode.connect(context.destination);
			gainNode.gain.value = 2;

			art.on('video:play', () => {
				context.resume();
			});

			art.on('video:pause', () => {
				context.suspend();
			});

			art.on('destroy', () => {
				context.close();
			});
		});
	};
};
