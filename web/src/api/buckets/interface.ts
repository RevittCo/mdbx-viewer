import { DataSource, KeyValuePair } from "@/types";
import { ApiResponse } from "../types";

export interface BucketsApi {
  getDataSource: () => Promise<ApiResponse<DataSource>>;
  openDBSource: (path: string) => Promise<ApiResponse<string[]>>;
  readBucket: (name: string, num: string, len: string) => Promise<ApiResponse<KeyValuePair[]>>;
}