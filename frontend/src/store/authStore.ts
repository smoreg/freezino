import { create } from 'zustand';
import { persist } from 'zustand/middleware';

import api from '../services/api';
import type { User, AuthResponse } from '../types';

interface AuthState {
  user: User | null;
  isAuthenticated: boolean;
  isLoading: boolean;

  // Actions
  setUser: (user: User | null) => void;
  setLoading: (isLoading: boolean) => void;
  login: (authData: AuthResponse) => void;
  logout: () => void;
  checkAuth: () => Promise<void>;
  refreshToken: () => Promise<boolean>;
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set, get) => ({
      user: null,
      isAuthenticated: false,
      isLoading: true,

      setUser: (user) => set({ user, isAuthenticated: !!user }),

      setLoading: (isLoading) => set({ isLoading }),

      login: (authData: AuthResponse) => {
        // Save tokens to localStorage
        localStorage.setItem('access_token', authData.access_token);
        localStorage.setItem('refresh_token', authData.refresh_token);

        // Update state
        set({
          user: authData.user,
          isAuthenticated: true,
          isLoading: false,
        });
      },

      logout: () => {
        // Clear tokens
        localStorage.removeItem('access_token');
        localStorage.removeItem('refresh_token');

        // Clear state
        set({
          user: null,
          isAuthenticated: false,
          isLoading: false,
        });

        // Redirect to login
        window.location.href = '/login';
      },

      checkAuth: async () => {
        const token = localStorage.getItem('access_token');

        if (!token) {
          set({ isLoading: false, isAuthenticated: false, user: null });
          return;
        }

        try {
          // Fetch current user data
          const response = await api.get<User>('/auth/me');

          set({
            user: response.data,
            isAuthenticated: true,
            isLoading: false,
          });
        } catch (error) {
          console.error('Auth check failed:', error);

          // Try to refresh token
          const refreshed = await get().refreshToken();

          if (!refreshed) {
            get().logout();
          }
        }
      },

      refreshToken: async () => {
        const refreshToken = localStorage.getItem('refresh_token');

        if (!refreshToken) {
          return false;
        }

        try {
          const response = await api.post<{ access_token: string }>('/auth/refresh', {
            refresh_token: refreshToken,
          });

          // Save new access token
          localStorage.setItem('access_token', response.data.access_token);

          // Fetch user data with new token
          const userResponse = await api.get<User>('/auth/me');

          set({
            user: userResponse.data,
            isAuthenticated: true,
            isLoading: false,
          });

          return true;
        } catch (error) {
          console.error('Token refresh failed:', error);
          return false;
        }
      },
    }),
    {
      name: 'auth-storage',
      partialize: (state) => ({
        user: state.user,
        isAuthenticated: state.isAuthenticated,
      }),
    }
  )
);
