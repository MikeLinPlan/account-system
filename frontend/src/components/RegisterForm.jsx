import React, { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { API, showError, showSuccess } from '../utils/api';

const RegisterForm = () => {
  const [inputs, setInputs] = useState({
    username: '',
    password: '',
    confirmPassword: '',
    email: '',
  });
  const [loading, setLoading] = useState(false);
  const navigate = useNavigate();

  const handleChange = (name, value) => {
    setInputs((inputs) => ({ ...inputs, [name]: value }));
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    const { username, password, confirmPassword, email } = inputs;
    
    if (!username || !password) {
      showError('請輸入用戶名和密碼！');
      return;
    }
    
    if (password.length < 8) {
      showError('密碼長度不得小於 8 位！');
      return;
    }
    
    if (password !== confirmPassword) {
      showError('兩次輸入的密碼不一致！');
      return;
    }
    
    setLoading(true);
    try {
      const res = await API.post('/api/user/register', {
        username,
        password,
        email,
      });
      
      const { success, message } = res.data;
      if (success) {
        showSuccess('註冊成功！請登入');
        navigate('/login');
      } else {
        showError(message);
      }
    } catch (error) {
      showError('註冊失敗，請稍後重試');
      console.error(error);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
      <div className="max-w-md w-full space-y-8">
        <div>
          <h2 className="mt-6 text-center text-3xl font-extrabold text-gray-900">
            註冊新帳號
          </h2>
          <p className="mt-2 text-center text-sm text-gray-600">
            或{' '}
            <Link to="/login" className="font-medium text-blue-600 hover:text-blue-500">
              登入現有帳號
            </Link>
          </p>
        </div>
        <form className="mt-8 space-y-6" onSubmit={handleSubmit}>
          <div className="rounded-md shadow-sm -space-y-px">
            <div>
              <label htmlFor="username" className="sr-only">
                用戶名
              </label>
              <input
                id="username"
                name="username"
                type="text"
                required
                className="appearance-none rounded-none relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 rounded-t-md focus:outline-none focus:ring-blue-500 focus:border-blue-500 focus:z-10 sm:text-sm"
                placeholder="用戶名"
                value={inputs.username}
                onChange={(e) => handleChange('username', e.target.value)}
              />
            </div>
            <div>
              <label htmlFor="email" className="sr-only">
                郵箱
              </label>
              <input
                id="email"
                name="email"
                type="email"
                className="appearance-none rounded-none relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 focus:outline-none focus:ring-blue-500 focus:border-blue-500 focus:z-10 sm:text-sm"
                placeholder="郵箱（可選）"
                value={inputs.email}
                onChange={(e) => handleChange('email', e.target.value)}
              />
            </div>
            <div>
              <label htmlFor="password" className="sr-only">
                密碼
              </label>
              <input
                id="password"
                name="password"
                type="password"
                required
                className="appearance-none rounded-none relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 focus:outline-none focus:ring-blue-500 focus:border-blue-500 focus:z-10 sm:text-sm"
                placeholder="密碼（至少 8 位）"
                value={inputs.password}
                onChange={(e) => handleChange('password', e.target.value)}
              />
            </div>
            <div>
              <label htmlFor="confirmPassword" className="sr-only">
                確認密碼
              </label>
              <input
                id="confirmPassword"
                name="confirmPassword"
                type="password"
                required
                className="appearance-none rounded-none relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 rounded-b-md focus:outline-none focus:ring-blue-500 focus:border-blue-500 focus:z-10 sm:text-sm"
                placeholder="確認密碼"
                value={inputs.confirmPassword}
                onChange={(e) => handleChange('confirmPassword', e.target.value)}
              />
            </div>
          </div>

          <div>
            <button
              type="submit"
              disabled={loading}
              className="group relative w-full flex justify-center py-2 px-4 border border-transparent text-sm font-medium rounded-md text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
            >
              {loading ? '註冊中...' : '註冊'}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};

export default RegisterForm;
