import { createBucketsApi } from "./buckets";
import { ApiFactory } from "./types";

export const useApi = (): ApiFactory => {
  return {
    buckets: createBucketsApi(),
  }
};
