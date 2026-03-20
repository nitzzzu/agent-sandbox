const envBaseUrl = import.meta.env.VITE_API_BASE_URL?.trim()

export const API_BASE_URL = (envBaseUrl && envBaseUrl.length > 0 ? envBaseUrl : '/api/v1').replace(/\/$/, '')
