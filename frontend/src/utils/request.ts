const BASE_URL = import.meta.env.VITE_APP_API_BASE_URL || 'https://occont.asia:18081';

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

export const request = <T = any>(options: RequestOptions): Promise<T> => {
  const token = uni.getStorageSync('token');

  return new Promise((resolve, reject) => {
    const headers: any = { ...options.header };
    // Only include Authorization for non-public routes
    if (token && !isPublicRoute(options.url)) {
      headers['Authorization'] = `Bearer ${token}`;
    }

    uni.request({
      url: options.url.startsWith('http') ? options.url : `${BASE_URL}${options.url}`,
      method: options.method || 'GET',
      data: options.data,
      header: headers,
      success: (res) => {
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
        reject(err);
      }
    });
  });
};
