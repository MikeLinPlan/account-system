import React from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import { AuthProvider } from './context/AuthContext';
import ProtectedRoute from './components/ProtectedRoute';
import Navbar from './components/Navbar';
import Home from './pages/Home';
import LoginForm from './components/LoginForm';
import RegisterForm from './components/RegisterForm';
import Profile from './pages/Profile';
import Tokens from './pages/Tokens';
import AdminUsers from './pages/admin/Users';
import NotFound from './pages/NotFound';

function App() {
  return (
    <AuthProvider>
      <Router>
        <Navbar />
        <Routes>
          {/* 公共路由 */}
          <Route path="/" element={<Home />} />
          <Route path="/login" element={<LoginForm />} />
          <Route path="/register" element={<RegisterForm />} />
          
          {/* 需要用戶認證的路由 */}
          <Route element={<ProtectedRoute />}>
            <Route path="/profile" element={<Profile />} />
            <Route path="/tokens" element={<Tokens />} />
          </Route>
          
          {/* 需要管理員認證的路由 */}
          <Route element={<ProtectedRoute requiredRole={10} />}>
            <Route path="/admin/users" element={<AdminUsers />} />
          </Route>
          
          {/* 404 頁面 */}
          <Route path="*" element={<NotFound />} />
        </Routes>
      </Router>
    </AuthProvider>
  );
}

export default App;
