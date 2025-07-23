import axios from 'axios'
import { CreateURLRequest, URLResponse, Analytics, ApiResponse } from '@/types'

const API_BASE_URL = import.meta.env.VITE_API_URL || '/api/v1'

const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
})

// Request interceptor
api.interceptors.request.use(
  (config) => {
    // Add any auth tokens here if needed
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// Response interceptor
api.interceptors.response.use(
  (response) => response,
  (error) => {
    // Handle common errors here
    if (error.response?.status === 401) {
      // Handle unauthorized
    }
    return Promise.reject(error)
  }
)

export const urlService = {
  // Create a new short URL
  createShortURL: async (data: CreateURLRequest): Promise<URLResponse> => {
    const response = await api.post<ApiResponse<URLResponse>>('/shorten', data)
    return response.data.data
  },

  // Get all URLs with pagination
  getAllURLs: async (page = 1, limit = 10): Promise<{ data: URLResponse[], pagination: any }> => {
    const response = await api.get<ApiResponse<URLResponse[]>>(`/urls?page=${page}&limit=${limit}`)
    return {
      data: response.data.data,
      pagination: response.data.pagination,
    }
  },

  // Get URL by ID
  getURLById: async (id: string): Promise<URLResponse> => {
    const response = await api.get<ApiResponse<URLResponse>>(`/urls/${id}`)
    return response.data.data
  },

  // Delete URL
  deleteURL: async (id: string): Promise<void> => {
    await api.delete(`/urls/${id}`)
  },

  // Get URL analytics
  getURLAnalytics: async (id: string): Promise<Analytics> => {
    const response = await api.get<ApiResponse<Analytics>>(`/analytics/${id}`)
    return response.data.data
  },

  // Get all analytics
  getAllAnalytics: async (): Promise<any> => {
    const response = await api.get<ApiResponse<any>>('/analytics')
    return response.data.data
  },
}

export default api 