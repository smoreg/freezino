import { create } from 'zustand';
import { persist } from 'zustand/middleware';

// import api from '../services/api';
import type { WorkSession } from '../types';

const WORK_DURATION_SECONDS = 180; // 3 minutes
const WORK_REWARD = 500; // $500

interface WorkStats {
  total_work_time: number; // in seconds
  total_earned: number;
  sessions_count: number;
  last_session?: string;
}

interface WorkState {
  isWorking: boolean;
  isPaused: boolean;
  timeRemaining: number; // in seconds
  startTime: number | null;
  pausedTime: number;
  workSessions: WorkSession[];
  stats: WorkStats;
  showStatsModal: boolean;
  lastCompletedSession: WorkSession | null;

  // Actions
  startWork: () => void;
  pauseWork: () => void;
  resumeWork: () => void;
  completeWork: () => void;
  cancelWork: () => void;
  tick: () => void;
  closeStatsModal: () => void;
  loadStats: () => Promise<void>;
}

export const useWorkStore = create<WorkState>()(
  persist(
    (set, get) => ({
      isWorking: false,
      isPaused: false,
      timeRemaining: WORK_DURATION_SECONDS,
      startTime: null,
      pausedTime: 0,
      workSessions: [],
      stats: {
        total_work_time: 0,
        total_earned: 0,
        sessions_count: 0,
      },
      showStatsModal: false,
      lastCompletedSession: null,

      startWork: () => {
        const now = Date.now();
        set({
          isWorking: true,
          isPaused: false,
          timeRemaining: WORK_DURATION_SECONDS,
          startTime: now,
          pausedTime: 0,
        });
      },

      pauseWork: () => {
        set({ isPaused: true });
      },

      resumeWork: () => {
        set({ isPaused: false });
      },

      completeWork: () => {
        const state = get();
        const completedSession: WorkSession = {
          id: String(Date.now()),
          user_id: 'current_user', // Will be replaced with actual user_id from auth
          duration_seconds: WORK_DURATION_SECONDS,
          earned: WORK_REWARD,
          completed_at: new Date().toISOString(),
        };

        // In production, this would be an API call
        // await api.post('/work/complete', { duration: WORK_DURATION_SECONDS });

        set({
          isWorking: false,
          isPaused: false,
          timeRemaining: WORK_DURATION_SECONDS,
          startTime: null,
          pausedTime: 0,
          workSessions: [completedSession, ...state.workSessions],
          stats: {
            total_work_time: state.stats.total_work_time + WORK_DURATION_SECONDS,
            total_earned: state.stats.total_earned + WORK_REWARD,
            sessions_count: state.stats.sessions_count + 1,
            last_session: completedSession.completed_at,
          },
          showStatsModal: true,
          lastCompletedSession: completedSession,
        });
      },

      cancelWork: () => {
        set({
          isWorking: false,
          isPaused: false,
          timeRemaining: WORK_DURATION_SECONDS,
          startTime: null,
          pausedTime: 0,
        });
      },

      tick: () => {
        const state = get();

        if (!state.isWorking || state.isPaused) return;

        const newTimeRemaining = state.timeRemaining - 1;

        if (newTimeRemaining <= 0) {
          get().completeWork();
        } else {
          set({ timeRemaining: newTimeRemaining });
        }
      },

      closeStatsModal: () => {
        set({ showStatsModal: false });
      },

      loadStats: async () => {
        // In production, this would be an API call
        // const response = await api.get('/work/stats');
        // set({ stats: response.data });

        // For now, stats are persisted locally
        console.error('Stats loaded from local storage');
      },
    }),
    {
      name: 'work-storage',
      partialize: (state) => ({
        workSessions: state.workSessions,
        stats: state.stats,
      }),
    }
  )
);
