import React, { useState, useEffect } from 'react';
import { API, showError, showSuccess } from '../utils/api';

const Tokens = () => {
  const [tokens, setTokens] = useState([]);
  const [loading, setLoading] = useState(true);
  const [modalOpen, setModalOpen] = useState(false);
  const [tokenInput, setTokenInput] = useState({
    name: '',
    remainQuota: 0,
    unlimitedQuota: false,
  });

  // 加載令牌列表
  const loadTokens = async () => {
    setLoading(true);
    try {
      const res = await API.get('/api/token');
      if (res.data.success) {
        setTokens(res.data.data);
      } else {
        showError(res.data.message);
      }
    } catch (error) {
      showError('加載令牌失敗');
      console.error(error);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadTokens();
  }, []);

  // 處理輸入變化
  const handleInputChange = (name, value) => {
    setTokenInput((prev) => ({ ...prev, [name]: value }));
  };

  // 創建新令牌
  const createToken = async () => {
    try {
      const res = await API.post('/api/token', tokenInput);
      if (res.data.success) {
        showSuccess('創建成功');
        setModalOpen(false);
        setTokenInput({
          name: '',
          remainQuota: 0,
          unlimitedQuota: false,
        });
        loadTokens();
      } else {
        showError(res.data.message);
      }
    } catch (error) {
      showError('創建失敗');
      console.error(error);
    }
  };

  // 刪除令牌
  const deleteToken = async (id) => {
    if (!window.confirm('確定要刪除此令牌嗎？')) {
      return;
    }
    
    try {
      const res = await API.delete(`/api/token/${id}`);
      if (res.data.success) {
        showSuccess('刪除成功');
        loadTokens();
      } else {
        showError(res.data.message);
      }
    } catch (error) {
      showError('刪除失敗');
      console.error(error);
    }
  };

  // 更新令牌狀態
  const updateTokenStatus = async (id, status) => {
    try {
      const token = tokens.find((t) => t.id === id);
      if (!token) return;
      
      const res = await API.put('/api/token', {
        ...token,
        status,
      });
      
      if (res.data.success) {
        showSuccess('更新成功');
        loadTokens();
      } else {
        showError(res.data.message);
      }
    } catch (error) {
      showError('更新失敗');
      console.error(error);
    }
  };

  return (
    <div className="max-w-6xl mx-auto py-8 px-4 sm:px-6 lg:px-8">
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-2xl font-bold text-gray-900">API 令牌管理</h1>
        <button
          onClick={() => setModalOpen(true)}
          className="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
        >
          創建新令牌
        </button>
      </div>

      {loading ? (
        <div className="text-center py-4">載入中...</div>
      ) : tokens.length === 0 ? (
        <div className="text-center py-4 text-gray-500">暫無令牌</div>
      ) : (
        <div className="overflow-x-auto">
          <table className="min-w-full divide-y divide-gray-200">
            <thead className="bg-gray-50">
              <tr>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  名稱
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  令牌
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  狀態
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  配額
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  創建時間
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  操作
                </th>
              </tr>
            </thead>
            <tbody className="bg-white divide-y divide-gray-200">
              {tokens.map((token) => (
                <tr key={token.id}>
                  <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                    {token.name}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                    <div className="flex items-center">
                      <span className="truncate max-w-xs">{token.key}</span>
                      <button
                        onClick={() => {
                          navigator.clipboard.writeText(token.key);
                          showSuccess('已複製到剪貼板');
                        }}
                        className="ml-2 text-blue-600 hover:text-blue-800"
                      >
                        複製
                      </button>
                    </div>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                    <span
                      className={`px-2 inline-flex text-xs leading-5 font-semibold rounded-full ${
                        token.status === 1
                          ? 'bg-green-100 text-green-800'
                          : 'bg-red-100 text-red-800'
                      }`}
                    >
                      {token.status === 1 ? '啟用' : '禁用'}
                    </span>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                    {token.unlimited_quota ? '無限制' : token.remain_quota}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                    {new Date(token.created_time).toLocaleString()}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm font-medium">
                    <div className="flex space-x-2">
                      {token.status === 1 ? (
                        <button
                          onClick={() => updateTokenStatus(token.id, 2)}
                          className="text-yellow-600 hover:text-yellow-900"
                        >
                          禁用
                        </button>
                      ) : (
                        <button
                          onClick={() => updateTokenStatus(token.id, 1)}
                          className="text-green-600 hover:text-green-900"
                        >
                          啟用
                        </button>
                      )}
                      <button
                        onClick={() => deleteToken(token.id)}
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

      {/* 創建令牌模態框 */}
      {modalOpen && (
        <div className="fixed inset-0 bg-gray-500 bg-opacity-75 flex items-center justify-center p-4">
          <div className="bg-white rounded-lg overflow-hidden shadow-xl max-w-md w-full">
            <div className="px-6 py-4">
              <h3 className="text-lg font-medium text-gray-900">創建新令牌</h3>
              <div className="mt-4">
                <label className="block text-sm font-medium text-gray-700">
                  名稱
                </label>
                <input
                  type="text"
                  value={tokenInput.name}
                  onChange={(e) => handleInputChange('name', e.target.value)}
                  className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
                />
              </div>
              <div className="mt-4">
                <label className="block text-sm font-medium text-gray-700">
                  配額
                </label>
                <input
                  type="number"
                  value={tokenInput.remainQuota}
                  onChange={(e) => handleInputChange('remainQuota', parseInt(e.target.value))}
                  disabled={tokenInput.unlimitedQuota}
                  className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
                />
              </div>
              <div className="mt-4 flex items-center">
                <input
                  type="checkbox"
                  id="unlimitedQuota"
                  checked={tokenInput.unlimitedQuota}
                  onChange={(e) => handleInputChange('unlimitedQuota', e.target.checked)}
                  className="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded"
                />
                <label htmlFor="unlimitedQuota" className="ml-2 block text-sm text-gray-900">
                  無限制配額
                </label>
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
                onClick={createToken}
                className="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
              >
                創建
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default Tokens;
