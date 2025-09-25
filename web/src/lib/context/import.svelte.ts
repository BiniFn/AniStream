type ImportJob = {
	id: string | null;
};

export const importjob = $state<ImportJob>({
	id: null,
});
