import { BucketsApi } from "./buckets/interface";

export type ApiResponse<T> = {
	data: T | null
	status?: number
	error: T | null | any
};

export interface ApiFactory {
	buckets: BucketsApi;
};
