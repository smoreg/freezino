import { useEffect, useState, useCallback } from 'react';
import { useNavigate, useSearchParams, useLocation } from 'react-router-dom';
import { motion } from 'framer-motion';
import { useAuthStore } from '../store/authStore';
import api from '../services/api';
import type { AuthResponse } from '../types';
import { PageTransition, rotateVariants, scaleFadeVariants } from '../components/animations';

const LoginPage = () => {
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();
  const location = useLocation();
  const { login, isAuthenticated } = useAuthStore();
  const [error, setError] = useState<string | null>(null);
  const [isProcessing, setIsProcessing] = useState(false);

  // Get redirect location from state (for protected route redirects)
  const from = (location.state as { from?: { pathname?: string } })?.from?.pathname || '/';

  const handleOAuthCallback = useCallback(async (code: string) => {
    setIsProcessing(true);
    setError(null);

    try {
      // Exchange code for tokens
      const response = await api.get<AuthResponse>(`/auth/google/callback?code=${code}`);

      // Save auth data to store
      login(response.data);

      // Redirect to home or previous location
      navigate(from, { replace: true });
    } catch (err: unknown) {
      console.error('OAuth callback failed:', err);
      setError((err as { response?: { data?: { message?: string } } }).response?.data?.message || 'Ошибка входа. Попробуйте снова.');
      setIsProcessing(false);
    }
  }, [login, navigate, from]);

  useEffect(() => {
    // If already authenticated, redirect to home or previous location
    if (isAuthenticated) {
      navigate(from, { replace: true });
      return;
    }

    // Handle OAuth callback
    const code = searchParams.get('code');
    const errorParam = searchParams.get('error');

    if (errorParam) {
      setError('Ошибка авторизации. Попробуйте снова.');
      return;
    }

    if (code && !isProcessing) {
      handleOAuthCallback(code);
    }
  }, [searchParams, isAuthenticated, navigate, from, isProcessing, handleOAuthCallback]);

  const handleGoogleLogin = () => {
    // Redirect to backend OAuth endpoint
    const apiUrl = import.meta.env.VITE_API_URL || 'http://localhost:8080';
    window.location.href = `${apiUrl}/api/auth/google`;
  };

  // Show loading state during OAuth processing
  if (isProcessing) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-dark">
        <PageTransition>
          <div className="text-center">
            <motion.div
              className="inline-block rounded-full h-12 w-12 border-t-2 border-b-2 border-primary"
              variants={rotateVariants}
              initial="initial"
              animate="animate"
            />
            <p className="mt-4 text-gray-400">Вход в систему...</p>
          </div>
        </PageTransition>
      </div>
    );
  }

  return (
    <div className="min-h-screen flex items-center justify-center bg-dark">
      <PageTransition>
        <motion.div
          variants={scaleFadeVariants}
          initial="initial"
          animate="animate"
          className="bg-gray-800 p-8 rounded-lg border border-gray-700 max-w-md w-full"
        >
          <div className="text-center mb-8">
            <h1 className="text-4xl font-bold text-primary mb-2">🎰 FREEZINO</h1>
            <p className="text-gray-400">Казино-симулятор против игровой зависимости</p>
          </div>

          {error && (
            <motion.div
              initial={{ opacity: 0, x: -10 }}
              animate={{ opacity: 1, x: 0 }}
              className="mb-6 p-4 bg-red-900/50 border border-red-700 rounded-lg text-red-200 text-sm"
            >
              {error}
            </motion.div>
          )}

          <motion.button
            whileHover={{ scale: 1.02 }}
            whileTap={{ scale: 0.98 }}
            onClick={handleGoogleLogin}
            disabled={isProcessing}
            className="w-full bg-white text-gray-800 font-semibold py-3 px-6 rounded-lg flex items-center justify-center space-x-3 hover:bg-gray-100 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
          >
          <svg className="w-6 h-6" viewBox="0 0 24 24">
            <path
              fill="currentColor"
              d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z"
            />
            <path
              fill="currentColor"
              d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"
            />
            <path
              fill="currentColor"
              d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z"
            />
            <path
              fill="currentColor"
              d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"
            />
          </svg>
          <span>Войти через Google</span>
        </motion.button>

        <div className="mt-6 text-center text-gray-400 text-sm">
          <p>Входя, вы соглашаетесь с тем, что:</p>
          <ul className="mt-2 space-y-1 text-xs">
            <li>❌ Это не настоящее казино</li>
            <li>✅ Используются только виртуальные деньги</li>
            <li>🎓 Цель - образовательная</li>
          </ul>
        </div>
      </motion.div>
      </PageTransition>
    </div>
  );
};

export default LoginPage;
