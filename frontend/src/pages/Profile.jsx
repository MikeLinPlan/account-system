import React, { useState, useContext, useEffect } from 'react';
import { AuthContext } from '../context/AuthContext';
import { API, showError, showSuccess } from '../utils/api';

const Profile = () => {
  const { user, updateUser } = useContext(AuthContext);
  const [inputs, setInputs] = useState({
    username: '',
    displayName: '',
    email: '',
    password: '',
    confirmPassword: '',
  });
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    if (user) {
      setInputs({
        username: user.username || '',
        displayName: user.display_name || '',
        email: user.email || '',
        password: '',
        confirmPassword: '',
      });
    }
  }, [user]);

  const handleChange = (name, value) => {
    setInputs((inputs) => ({ ...inputs, [name]: value }));
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    const { username, displayName, email, password, confirmPassword } = inputs;
    
    if (password && password !== confirmPassword) {
      showError('兩次輸入的密碼不一致！');
      return;
    }
    
    if (password && password.length < 8) {
      showError('密碼長度不得小於 8 位！');
      return;
    }
    
    setLoading(true);
    try {
      const updateData = {
        username,
        display_name: displayName,
        email,
      };
      
      if (password) {
        updateData.password = password;
      }
      
      const res = await API.put('/api/user/self', updateData);
      
      const { success, message } = res.data;
      if (success) {
        showSuccess('更新成功！');
        
        // 更新本地用戶信息
        const updatedUser = {
          ...user,
          username,
          display_name: displayName,
          email,
        };
        updateUser(updatedUser);
        
        // 清空密碼字段
        setInputs((inputs) => ({
          ...inputs,
          password: '',
          confirmPassword: '',
        }));
      } else {
        showError(message);
      }
    } catch (error) {
      showError('更新失敗，請稍後重試');
      console.error(error);
    } finally {
      setLoading(false);
    }
  };

  const generateAccessToken = async () => {
    try {
      const res = await API.get('/api/user/token');
      
      const { success, message, data } = res.data;
      if (success) {
        showSuccess('生成成功！');
        
        // 更新本地用戶信息
        const updatedUser = {
          ...user,
          access_token: data,
        };
        updateUser(updatedUser);
        
        // 顯示令牌
        alert(`您的訪問令牌：${data}`);
      } else {
        showError(message);
      }
    } catch (error) {
      showError('生成失敗，請稍後重試');
      console.error(error);
    }
  };

  return (
    <div className="max-w-4xl mx-auto py-8 px-4 sm:px-6 lg:px-8">
      <div className="bg-white shadow overflow-hidden sm:rounded-lg">
        <div className="px-4 py-5 sm:px-6">
          <h3 className="text-lg leading-6 font-medium text-gray-900">個人資料</h3>
          <p className="mt-1 max-w-2xl text-sm text-gray-500">更新您的個人信息和密碼</p>
        </div>
        <div className="border-t border-gray-200">
          <form onSubmit={handleSubmit} className="px-4 py-5 sm:p-6">
            <div className="grid grid-cols-1 gap-6 sm:grid-cols-2">
              <div>
                <label htmlFor="username" className="block text-sm font-medium text-gray-700">
                  用戶名
                </label>
                <input
                  type="text"
                  name="username"
                  id="username"
                  value={inputs.username}
                  onChange={(e) => handleChange('username', e.target.value)}
                  className="mt-1 focus:ring-blue-500 focus:border-blue-500 block w-full shadow-sm sm:text-sm border-gray-300 rounded-md"
                />
              </div>
              <div>
                <label htmlFor="displayName" className="block text-sm font-medium text-gray-700">
                  顯示名稱
                </label>
                <input
                  type="text"
                  name="displayName"
                  id="displayName"
                  value={inputs.displayName}
                  onChange={(e) => handleChange('displayName', e.target.value)}
                  className="mt-1 focus:ring-blue-500 focus:border-blue-500 block w-full shadow-sm sm:text-sm border-gray-300 rounded-md"
                />
              </div>
              <div>
                <label htmlFor="email" className="block text-sm font-medium text-gray-700">
                  郵箱
                </label>
                <input
                  type="email"
                  name="email"
                  id="email"
                  value={inputs.email}
                  onChange={(e) => handleChange('email', e.target.value)}
                  className="mt-1 focus:ring-blue-500 focus:border-blue-500 block w-full shadow-sm sm:text-sm border-gray-300 rounded-md"
                />
              </div>
              <div>
                <label htmlFor="password" className="block text-sm font-medium text-gray-700">
                  新密碼
                </label>
                <input
                  type="password"
                  name="password"
                  id="password"
                  value={inputs.password}
                  onChange={(e) => handleChange('password', e.target.value)}
                  className="mt-1 focus:ring-blue-500 focus:border-blue-500 block w-full shadow-sm sm:text-sm border-gray-300 rounded-md"
                  placeholder="留空表示不修改"
                />
              </div>
              <div>
                <label htmlFor="confirmPassword" className="block text-sm font-medium text-gray-700">
                  確認新密碼
                </label>
                <input
                  type="password"
                  name="confirmPassword"
                  id="confirmPassword"
                  value={inputs.confirmPassword}
                  onChange={(e) => handleChange('confirmPassword', e.target.value)}
                  className="mt-1 focus:ring-blue-500 focus:border-blue-500 block w-full shadow-sm sm:text-sm border-gray-300 rounded-md"
                  placeholder="留空表示不修改"
                />
              </div>
            </div>
            <div className="mt-6 flex justify-between">
              <button
                type="submit"
                disabled={loading}
                className="inline-flex justify-center py-2 px-4 border border-transparent shadow-sm text-sm font-medium rounded-md text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
              >
                {loading ? '更新中...' : '更新資料'}
              </button>
              <button
                type="button"
                onClick={generateAccessToken}
                className="inline-flex justify-center py-2 px-4 border border-transparent shadow-sm text-sm font-medium rounded-md text-white bg-green-600 hover:bg-green-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-green-500"
              >
                生成訪問令牌
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
  );
};

export default Profile;
