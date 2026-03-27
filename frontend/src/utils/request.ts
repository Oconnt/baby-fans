export const BASE_URL = import.meta.env.VITE_APP_API_BASE_URL || 'https://occont.asia';

export interface RequestOptions {
  url: string;
  method?: 'GET' | 'POST' | 'PUT' | 'DELETE';
  data?: any;
  header?: any;
}

// Public routes that should not include Authorization header
const PUBLIC_ROUTES = ['/register', '/login/face', '/login/code'];

const isPublicRoute = (url: string): boolean => {
  return PUBLIC_ROUTES.some(route => url.includes(route));
};

// In-flight request deduplication to prevent rapid repeated calls
const pendingRequests = new Map<string, Promise<any>>();

const getRequestKey = (options: RequestOptions): string => {
  return `${options.method || 'GET'}:${options.url}:${JSON.stringify(options.data || {})}`;
};

// Generate a unique request ID for backend idempotency
const generateRequestId = (): string => {
  return `${Date.now()}-${Math.random().toString(36).substring(2, 11)}`;
};

export const request = <T = any>(options: RequestOptions): Promise<T> => {
  const token = uni.getStorageSync('token');
  const requestKey = getRequestKey(options);

  // For POST/PUT/DELETE, check if there's already an in-flight request with same key
  if (pendingRequests.has(requestKey)) {
    return pendingRequests.get(requestKey)!;
  }

  const requestPromise = new Promise<T>((resolve, reject) => {
    const headers: any = { ...options.header };
    // Only include Authorization for non-public routes
    if (token && !isPublicRoute(options.url)) {
      headers['Authorization'] = `Bearer ${token}`;
    }
    // Add request ID for backend idempotency
    headers['X-Request-ID'] = generateRequestId();

    uni.request({
      url: options.url.startsWith('http') ? options.url : `${BASE_URL}${options.url}`,
      method: options.method || 'GET',
      data: options.data,
      header: headers,
      success: (res) => {
        // Remove from pending requests on completion
        pendingRequests.delete(requestKey);
        if (res.statusCode >= 200 && res.statusCode < 300) {
          resolve(res.data as T);
        } else if (res.statusCode === 401) {
          uni.removeStorageSync('token');
          uni.navigateTo({ url: '/pages/index/index' });
          reject(res.data);
        } else {
          reject(res.data);
        }
      },
      fail: (err) => {
        pendingRequests.delete(requestKey);
        reject(err);
      }
    });
  });

  // Only track POST/PUT/DELETE requests (mutating operations)
  if (options.method && ['POST', 'PUT', 'DELETE'].includes(options.method)) {
    pendingRequests.set(requestKey, requestPromise);
  }

  return requestPromise;
};
