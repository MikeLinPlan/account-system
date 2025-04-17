import React, { createContext, useState, useEffect } from 'react';
import { API } from '../utils/api';

export const AuthContext = createContext();

export const AuthProvider = ({ children }) => {
  const [user, setUser] = useState(null);
  const [loading, setLoading] = useState(true);
  const [initialized, setInitialized] = useState(false);

  // 從本地存儲加載用戶
  useEffect(() => {
    const loadUser = async () => {
      const storedUser = localStorage.getItem('user');
      if (storedUser) {
        try {
          const userData = JSON.parse(storedUser);
          setUser(userData);
          
          // 驗證用戶會話是否有效
          try {
            const res = await API.get('/api/user/self');
            if (res.data.success) {
              setUser(res.data.data);
            } else {
              // 會話無效，清除本地存儲
              localStorage.removeItem('user');
              setUser(null);
            }
          } catch (error) {
            console.error('Failed to validate session:', error);
            localStorage.removeItem('user');
            setUser(null);
          }
        } catch (error) {
          console.error('Failed to parse stored user:', error);
          localStorage.removeItem('user');
        }
      }
      setLoading(false);
      setInitialized(true);
    };

    loadUser();
  }, []);

  // 登入
  const login = (userData) => {
    setUser(userData);
    localStorage.setItem('user', JSON.stringify(userData));
  };

  // 登出
  const logout = async () => {
    try {
      await API.get('/api/user/logout');
    } catch (error) {
      console.error('Logout error:', error);
    } finally {
      setUser(null);
      localStorage.removeItem('user');
    }
  };

  // 更新用戶信息
  const updateUser = (userData) => {
    setUser(userData);
    localStorage.setItem('user', JSON.stringify(userData));
  };

  return (
    <AuthContext.Provider
      value={{
        user,
        loading,
        initialized,
        login,
        logout,
        updateUser,
        isAuthenticated: !!user,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
};
