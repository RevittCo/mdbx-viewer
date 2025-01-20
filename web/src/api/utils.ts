import { AxiosResponse } from 'axios'
import { ApiResponse } from './types'

export async function handleResponse<T>(fn: () => Promise<AxiosResponse<T>>): Promise<ApiResponse<T>> {
  try {
    const response = await fn()
    const data: T = response.data
    const status: number = response.status
    const errorStatus: boolean = status < 200 || status >= 300;

    return {
      data,
      status,
      error: errorStatus ? response : null,
    }
  } catch (error: any) {
    console.error('Request failed: ', error);

    const status = error.response?.status ?? 500;
    const errorData = error.response?.data || 'An unexpected error occurred';

    return {
      data: null,
      status,
      error: {
        status,
        data: errorData,
      },
    };
  }
}
