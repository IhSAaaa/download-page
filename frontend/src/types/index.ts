export interface URL {
  id: string
  original_url: string
  short_code: string
  custom_code?: string
  title?: string
  description?: string
  user_id?: string
  is_active: boolean
  expires_at?: string
  click_count: number
  created_at: string
  updated_at: string
}

export interface URLResponse {
  id: string
  original_url: string
  short_url: string
  custom_code?: string
  title?: string
  description?: string
  qr_code?: string
  click_count: number
  expires_at?: string
  created_at: string
}

export interface CreateURLRequest {
  original_url: string
  custom_code?: string
  title?: string
  description?: string
  expires_at?: string
}

export interface Click {
  id: string
  url_id: string
  ip_address: string
  user_agent: string
  referer?: string
  country?: string
  city?: string
  device?: string
  browser?: string
  os?: string
  clicked_at: string
}

export interface Analytics {
  url_id: string
  total_clicks: number
  unique_clicks: number
  top_countries: Country[]
  top_devices: Device[]
  top_browsers: Browser[]
  click_timeline: Timeline[]
  last_clicked_at?: string
}

export interface Country {
  country: string
  clicks: number
}

export interface Device {
  device: string
  clicks: number
}

export interface Browser {
  browser: string
  clicks: number
}

export interface Timeline {
  date: string
  clicks: number
}

export interface Pagination {
  page: number
  limit: number
  total: number
  pages: number
}

export interface ApiResponse<T> {
  message?: string
  data: T
  pagination?: Pagination
} 