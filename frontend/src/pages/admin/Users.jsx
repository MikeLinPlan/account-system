import React, { useState, useEffect } from 'react';
import { API, showError, showSuccess } from '../../utils/api';

const Users = () => {
  const [users, setUsers] = useState([]);
  const [loading, setLoading] = useState(true);
  const [modalOpen, setModalOpen] = useState(false);
  const [editingUser, setEditingUser] = useState(null);
  const [userInput, setUserInput] = useState({
    username: '',
    displayName: '',
    email: '',
    password: '',
    role: 1,
    status: 1,
  });
  const [page, setPage] = useState(1);
  const [totalPages, setTotalPages] = useState(1);
  const [searchKeyword, setSearchKeyword] = useState('');

  // 加載用戶列表
  const loadUsers = async (pageNum = 1, keyword = '') => {
    setLoading(true);
    try {
      let url = `/api/user?page=${pageNum}&page_size=10`;
      if (keyword) {
        url = `/api/user/search?keyword=${encodeURIComponent(keyword)}&page=${pageNum}&page_size=10`;
      }
      
      const res = await API.get(url);
      if (res.data.success) {
        setUsers(res.data.data);
        setTotalPages(Math.ceil(res.data.total / 10));
      } else {
        showError(res.data.message);
      }
    } catch (error) {
      showError('加載用戶失敗');
      console.error(error);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadUsers(page, searchKeyword);
  }, [page, searchKeyword]);

  // 處理輸入變化
  const handleInputChange = (name, value) => {
    setUserInput((prev) => ({ ...prev, [name]: value }));
  };

  // 打開創建用戶模態框
  const openCreateModal = () => {
    setEditingUser(null);
    setUserInput({
      username: '',
      displayName: '',
      email: '',
      password: '',
      role: 1,
      status: 1,
    });
    setModalOpen(true);
  };

  // 打開編輯用戶模態框
  const openEditModal = (user) => {
    setEditingUser(user);
    setUserInput({
      username: user.username,
      displayName: user.display_name || '',
      email: user.email || '',
      password: '',
      role: user.role,
      status: user.status,
    });
    setModalOpen(true);
  };

  // 創建或更新用戶
  const saveUser = async () => {
    try {
      const userData = {
        username: userInput.username,
        display_name: userInput.displayName,
        email: userInput.email,
        role: parseInt(userInput.role),
        status: parseInt(userInput.status),
      };
      
      if (userInput.password) {
        userData.password = userInput.password;
      }
      
      let res;
      if (editingUser) {
        // 更新用戶
        userData.id = editingUser.id;
        res = await API.put('/api/user', userData);
      } else {
        // 創建用戶
        res = await API.post('/api/user', userData);
      }
      
      if (res.data.success) {
        showSuccess(editingUser ? '更新成功' : '創建成功');
        setModalOpen(false);
        loadUsers(page, searchKeyword);
      } else {
        showError(res.data.message);
      }
    } catch (error) {
      showError(editingUser ? '更新失敗' : '創建失敗');
      console.error(error);
    }
  };

  // 刪除用戶
  const deleteUser = async (id) => {
    if (!window.confirm('確定要刪除此用戶嗎？')) {
      return;
    }
    
    try {
      const res = await API.delete(`/api/user/${id}`);
      if (res.data.success) {
        showSuccess('刪除成功');
        loadUsers(page, searchKeyword);
      } else {
        showError(res.data.message);
      }
    } catch (error) {
      showError('刪除失敗');
      console.error(error);
    }
  };

  // 處理搜索
  const handleSearch = (e) => {
    e.preventDefault();
    setPage(1);
    loadUsers(1, searchKeyword);
  };

  // 處理分頁
  const handlePageChange = (newPage) => {
    if (newPage < 1 || newPage > totalPages) return;
    setPage(newPage);
  };

  return (
    <div className="max-w-6xl mx-auto py-8 px-4 sm:px-6 lg:px-8">
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-2xl font-bold text-gray-900">用戶管理</h1>
        <button
          onClick={openCreateModal}
          className="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
        >
          創建用戶
        </button>
      </div>

      <div className="mb-6">
        <form onSubmit={handleSearch} className="flex">
          <input
            type="text"
            value={searchKeyword}
            onChange={(e) => setSearchKeyword(e.target.value)}
            placeholder="搜索用戶名、顯示名稱或郵箱"
            className="flex-1 border border-gray-300 rounded-l-md px-4 py-2 focus:outline-none focus:ring-blue-500 focus:border-blue-500"
          />
          <button
            type="submit"
            className="bg-blue-600 text-white px-4 py-2 rounded-r-md hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
          >
            搜索
          </button>
        </form>
      </div>

      {loading ? (
        <div className="text-center py-4">載入中...</div>
      ) : users.length === 0 ? (
        <div className="text-center py-4 text-gray-500">暫無用戶</div>
      ) : (
        <div className="overflow-x-auto">
          <table className="min-w-full divide-y divide-gray-200">
            <thead className="bg-gray-50">
              <tr>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  ID
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  用戶名
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  顯示名稱
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  郵箱
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  角色
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  狀態
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  操作
                </th>
              </tr>
            </thead>
            <tbody className="bg-white divide-y divide-gray-200">
              {users.map((user) => (
                <tr key={user.id}>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                    {user.id}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                    {user.username}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                    {user.display_name || '-'}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                    {user.email || '-'}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                    {user.role === 100 ? (
                      <span className="px-2 inline-flex text-xs leading-5 font-semibold rounded-full bg-purple-100 text-purple-800">
                        超級管理員
                      </span>
                    ) : user.role === 10 ? (
                      <span className="px-2 inline-flex text-xs leading-5 font-semibold rounded-full bg-blue-100 text-blue-800">
                        管理員
                      </span>
                    ) : (
                      <span className="px-2 inline-flex text-xs leading-5 font-semibold rounded-full bg-gray-100 text-gray-800">
                        普通用戶
                      </span>
                    )}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                    <span
                      className={`px-2 inline-flex text-xs leading-5 font-semibold rounded-full ${
                        user.status === 1
                          ? 'bg-green-100 text-green-800'
                          : 'bg-red-100 text-red-800'
                      }`}
                    >
                      {user.status === 1 ? '啟用' : '禁用'}
                    </span>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm font-medium">
                    <div className="flex space-x-2">
                      <button
                        onClick={() => openEditModal(user)}
                        className="text-blue-600 hover:text-blue-900"
                      >
                        編輯
                      </button>
                      <button
                        onClick={() => deleteUser(user.id)}
                        className="text-red-600 hover:text-red-900"
                      >
                        刪除
                      </button>
                    </div>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}

      {/* 分頁 */}
      {totalPages > 1 && (
        <div className="flex justify-center mt-6">
          <nav className="relative z-0 inline-flex rounded-md shadow-sm -space-x-px">
            <button
              onClick={() => handlePageChange(page - 1)}
              disabled={page === 1}
              className={`relative inline-flex items-center px-2 py-2 rounded-l-md border border-gray-300 bg-white text-sm font-medium ${
                page === 1
                  ? 'text-gray-300 cursor-not-allowed'
                  : 'text-gray-500 hover:bg-gray-50'
              }`}
            >
              上一頁
            </button>
            {Array.from({ length: totalPages }, (_, i) => i + 1).map((pageNum) => (
              <button
                key={pageNum}
                onClick={() => handlePageChange(pageNum)}
                className={`relative inline-flex items-center px-4 py-2 border border-gray-300 bg-white text-sm font-medium ${
                  page === pageNum
                    ? 'z-10 bg-blue-50 border-blue-500 text-blue-600'
                    : 'text-gray-500 hover:bg-gray-50'
                }`}
              >
                {pageNum}
              </button>
            ))}
            <button
              onClick={() => handlePageChange(page + 1)}
              disabled={page === totalPages}
              className={`relative inline-flex items-center px-2 py-2 rounded-r-md border border-gray-300 bg-white text-sm font-medium ${
                page === totalPages
                  ? 'text-gray-300 cursor-not-allowed'
                  : 'text-gray-500 hover:bg-gray-50'
              }`}
            >
              下一頁
            </button>
          </nav>
        </div>
      )}

      {/* 創建/編輯用戶模態框 */}
      {modalOpen && (
        <div className="fixed inset-0 bg-gray-500 bg-opacity-75 flex items-center justify-center p-4">
          <div className="bg-white rounded-lg overflow-hidden shadow-xl max-w-md w-full">
            <div className="px-6 py-4">
              <h3 className="text-lg font-medium text-gray-900">
                {editingUser ? '編輯用戶' : '創建用戶'}
              </h3>
              <div className="mt-4 space-y-4">
                <div>
                  <label className="block text-sm font-medium text-gray-700">
                    用戶名
                  </label>
                  <input
                    type="text"
                    value={userInput.username}
                    onChange={(e) => handleInputChange('username', e.target.value)}
                    className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700">
                    顯示名稱
                  </label>
                  <input
                    type="text"
                    value={userInput.displayName}
                    onChange={(e) => handleInputChange('displayName', e.target.value)}
                    className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700">
                    郵箱
                  </label>
                  <input
                    type="email"
                    value={userInput.email}
                    onChange={(e) => handleInputChange('email', e.target.value)}
                    className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700">
                    密碼 {editingUser && '(留空表示不修改)'}
                  </label>
                  <input
                    type="password"
                    value={userInput.password}
                    onChange={(e) => handleInputChange('password', e.target.value)}
                    className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700">
                    角色
                  </label>
                  <select
                    value={userInput.role}
                    onChange={(e) => handleInputChange('role', e.target.value)}
                    className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
                  >
                    <option value={1}>普通用戶</option>
                    <option value={10}>管理員</option>
                    <option value={100}>超級管理員</option>
                  </select>
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700">
                    狀態
                  </label>
                  <select
                    value={userInput.status}
                    onChange={(e) => handleInputChange('status', e.target.value)}
                    className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
                  >
                    <option value={1}>啟用</option>
                    <option value={2}>禁用</option>
                  </select>
                </div>
              </div>
            </div>
            <div className="px-6 py-4 bg-gray-50 flex justify-end">
              <button
                onClick={() => setModalOpen(false)}
                className="px-4 py-2 bg-gray-200 text-gray-800 rounded-md mr-2 hover:bg-gray-300 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-gray-500"
              >
                取消
              </button>
              <button
                onClick={saveUser}
                className="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
              >
                保存
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default Users;
