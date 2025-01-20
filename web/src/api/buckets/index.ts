import { BucketsApi } from './interface'
import { http } from '../http'

const Buckets_Url = '/v1/buckets'

export const createBucketsApi = (): BucketsApi => ({
  getDataSource: async () => {
    return await http.get(`${Buckets_Url}/data-source`)
  },
  openDBSource: async (path: string) => {
    return await http.post(`${Buckets_Url}/open`, { path: path })
  },
  readBucket: async (name: string, num: string, len: string) => {
    return await http.get(`${Buckets_Url}/${name}/read/${num}/${len}`)
  }
})