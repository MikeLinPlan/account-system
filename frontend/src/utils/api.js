import axios from 'axios';

// 創建 axios 實例
export const API = axios.create({
  baseURL: '/',
  withCredentials: true,
});

// 請求攔截器
API.interceptors.request.use(
  (config) => {
    // 從本地存儲獲取用戶信息
    const user = JSON.parse(localStorage.getItem('user') || '{}');
    
    // 如果用戶有訪問令牌，添加到請求頭
    if (user && user.access_token) {
      config.headers.Authorization = `Bearer ${user.access_token}`;
    }
    
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// 響應攔截器
API.interceptors.response.use(
  (response) => {
    return response;
  },
  (error) => {
    // 處理 401 未授權錯誤
    if (error.response && error.response.status === 401) {
      // 清除本地存儲並重定向到登入頁面
      localStorage.removeItem('user');
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);

// 顯示成功消息
export const showSuccess = (message) => {
  // 這裡可以使用您喜歡的通知庫，如 react-toastify
  alert(message);
};

// 顯示錯誤消息
export const showError = (message) => {
  // 這裡可以使用您喜歡的通知庫，如 react-toastify
  alert(message);
};

// 更新 API 配置
export const updateAPI = () => {
  const user = JSON.parse(localStorage.getItem('user') || '{}');
  if (user && user.access_token) {
    API.defaults.headers.common['Authorization'] = `Bearer ${user.access_token}`;
  } else {
    delete API.defaults.headers.common['Authorization'];
  }
};
