<script lang="ts">
	import { apiClient } from '$lib/api/client';
	import { importjob } from '$lib/context/import.svelte';
	import { toast } from 'svelte-sonner';

	let pollingInterval: NodeJS.Timeout | null = null;

	async function checkImportStatus(jobId: string) {
		try {
			const response = await apiClient.GET('/library/import/{id}', {
				params: {
					path: { id: jobId },
				},
			});

			if (response.data) {
				const status = response.data.status;

				if (status === 'completed') {
					clearInterval(pollingInterval!);
					pollingInterval = null;
					localStorage.removeItem('import_job_id');
					toast.success('Library import completed successfully!');
				} else if (status === 'failed') {
					clearInterval(pollingInterval!);
					pollingInterval = null;
					localStorage.removeItem('import_job_id');
					toast.error('Library import failed. Please try again.');
				}
			}
		} catch (error) {
			console.error('Failed to check import status:', error);
		}
	}

	$effect(() => {
		const jobId = importjob.id ?? localStorage.getItem('import_job_id');

		if (jobId && !pollingInterval) {
			pollingInterval = setInterval(() => {
				checkImportStatus(jobId);
			}, 5000);
		} else if (!jobId && pollingInterval) {
			clearInterval(pollingInterval);
			pollingInterval = null;
		}

		return () => {
			if (pollingInterval) {
				clearInterval(pollingInterval);
				pollingInterval = null;
			}
		};
	});
</script>
