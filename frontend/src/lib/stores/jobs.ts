/**
 * Jobs store for tracking active background jobs
 * Requirements: 4.4
 */

import { writable, derived, get } from 'svelte/store';
import {
	listJobs,
	getJob,
	createJob,
	cancelJob,
	isJobActive,
	isJobTerminal,
	type Job,
	type JobState,
	type CreateJobRequest
} from '$lib/api/jobs';

/**
 * Job update from WebSocket
 */
export interface JobUpdate {
	jobId: string;
	state: JobState;
	progress: number;
	error?: string;
}

/**
 * Jobs state
 */
export interface JobsState {
	jobs: Map<string, Job>;
	isLoading: boolean;
	error: string | null;
}

/**
 * Initial jobs state
 */
const initialState: JobsState = {
	jobs: new Map(),
	isLoading: false,
	error: null
};

/**
 * Create the jobs store
 */
function createJobsStore() {
	const { subscribe, set, update } = writable<JobsState>(initialState);

	/**
	 * Load all jobs from the API
	 */
	async function loadJobs(): Promise<void> {
		update((state) => ({ ...state, isLoading: true, error: null }));

		try {
			const response = await listJobs();
			const jobsMap = new Map<string, Job>();
			for (const job of response.jobs) {
				jobsMap.set(job.id, job);
			}
			update((state) => ({
				...state,
				jobs: jobsMap,
				isLoading: false
			}));
		} catch (err) {
			const message = err instanceof Error ? err.message : 'Failed to load jobs';
			update((state) => ({
				...state,
				isLoading: false,
				error: message
			}));
		}
	}

	/**
	 * Add or update a job in the store
	 */
	function upsertJob(job: Job): void {
		update((state) => {
			const newJobs = new Map(state.jobs);
			newJobs.set(job.id, job);
			return { ...state, jobs: newJobs };
		});
	}

	/**
	 * Update a job from a WebSocket update
	 */
	function updateFromWebSocket(jobUpdate: JobUpdate): void {
		update((state) => {
			const existingJob = state.jobs.get(jobUpdate.jobId);
			if (!existingJob) return state;

			const updatedJob: Job = {
				...existingJob,
				state: jobUpdate.state,
				progress: jobUpdate.progress,
				error: jobUpdate.error
			};

			// Set completedAt if job is now terminal
			if (isJobTerminal(updatedJob) && !updatedJob.completedAt) {
				updatedJob.completedAt = new Date().toISOString();
			}

			const newJobs = new Map(state.jobs);
			newJobs.set(jobUpdate.jobId, updatedJob);
			return { ...state, jobs: newJobs };
		});
	}

	/**
	 * Remove a job from the store
	 */
	function removeJob(jobId: string): void {
		update((state) => {
			const newJobs = new Map(state.jobs);
			newJobs.delete(jobId);
			return { ...state, jobs: newJobs };
		});
	}

	/**
	 * Clear all completed/failed/cancelled jobs
	 */
	function clearTerminalJobs(): void {
		update((state) => {
			const newJobs = new Map<string, Job>();
			for (const [id, job] of state.jobs) {
				if (isJobActive(job)) {
					newJobs.set(id, job);
				}
			}
			return { ...state, jobs: newJobs };
		});
	}

	/**
	 * Get a job by ID
	 */
	function getJobById(jobId: string): Job | undefined {
		return get({ subscribe }).jobs.get(jobId);
	}

	/**
	 * Get all jobs as an array
	 */
	function getAllJobs(): Job[] {
		return Array.from(get({ subscribe }).jobs.values());
	}

	/**
	 * Get active jobs
	 */
	function getActiveJobs(): Job[] {
		return getAllJobs().filter(isJobActive);
	}

	/**
	 * Clear error
	 */
	function clearError(): void {
		update((state) => ({ ...state, error: null }));
	}

	/**
	 * Reset store to initial state
	 */
	function reset(): void {
		set(initialState);
	}

	return {
		subscribe,
		loadJobs,
		upsertJob,
		updateFromWebSocket,
		removeJob,
		clearTerminalJobs,
		getJobById,
		getAllJobs,
		getActiveJobs,
		clearError,
		reset
	};
}

/**
 * Jobs store singleton
 */
export const jobsStore = createJobsStore();

/**
 * Derived store for jobs as array (sorted by creation time, newest first)
 */
export const jobsList = derived(jobsStore, ($jobs) =>
	Array.from($jobs.jobs.values()).sort(
		(a, b) => new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime()
	)
);

/**
 * Derived store for active jobs only
 */
export const activeJobs = derived(jobsStore, ($jobs) =>
	Array.from($jobs.jobs.values())
		.filter(isJobActive)
		.sort((a, b) => new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime())
);

/**
 * Derived store for active jobs count
 */
export const activeJobsCount = derived(
	jobsStore,
	($jobs) => Array.from($jobs.jobs.values()).filter(isJobActive).length
);

/**
 * Derived store for whether there are any active jobs
 */
export const hasActiveJobs = derived(jobsStore, ($jobs) =>
	Array.from($jobs.jobs.values()).some(isJobActive)
);

/**
 * Derived store for completed jobs
 */
export const completedJobs = derived(jobsStore, ($jobs) =>
	Array.from($jobs.jobs.values())
		.filter((job) => job.state === 'completed')
		.sort(
			(a, b) =>
				new Date(b.completedAt || b.createdAt).getTime() -
				new Date(a.completedAt || a.createdAt).getTime()
		)
);

/**
 * Derived store for failed jobs
 */
export const failedJobs = derived(jobsStore, ($jobs) =>
	Array.from($jobs.jobs.values())
		.filter((job) => job.state === 'failed')
		.sort(
			(a, b) =>
				new Date(b.completedAt || b.createdAt).getTime() -
				new Date(a.completedAt || a.createdAt).getTime()
		)
);

/**
 * Query key factory for jobs
 */
export const jobQueryKeys = {
	all: ['jobs'] as const,
	list: () => [...jobQueryKeys.all, 'list'] as const,
	detail: (id: string) => [...jobQueryKeys.all, 'detail', id] as const
};

/**
 * Query options factory for listing all jobs
 */
export function jobsQueryOptions() {
	return {
		queryKey: jobQueryKeys.list(),
		queryFn: () => listJobs(),
		refetchInterval: 5000 // Refetch every 5 seconds for active jobs
	};
}

/**
 * Query options factory for a specific job
 */
export function jobQueryOptions(jobId: string) {
	return {
		queryKey: jobQueryKeys.detail(jobId),
		queryFn: () => getJob(jobId),
		enabled: !!jobId
	};
}

/**
 * Mutation options for creating a new job
 */
export function createJobMutationOptions() {
	return {
		mutationFn: (request: CreateJobRequest) => createJob(request),
		onSuccess: (job: Job) => {
			jobsStore.upsertJob(job);
		}
	};
}

/**
 * Mutation options for cancelling a job
 */
export function cancelJobMutationOptions() {
	return {
		mutationFn: (jobId: string) => cancelJob(jobId),
		onSuccess: (_: unknown, jobId: string) => {
			// Update the job state locally
			const job = jobsStore.getJobById(jobId);
			if (job) {
				jobsStore.upsertJob({
					...job,
					state: 'cancelled',
					completedAt: new Date().toISOString()
				});
			}
		}
	};
}

// Re-export utility functions
export { isJobActive, isJobTerminal } from '$lib/api/jobs';
