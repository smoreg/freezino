import { useEffect, useCallback } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { useWorkStore } from '../store/workStore';

interface CountryComparison {
  name: string;
  flag: string;
  avgHourlyWage: number;
  minutesToEarn500: number;
}

// Country wage data for comparison
const countryComparisons: CountryComparison[] = [
  { name: 'США', flag: '🇺🇸', avgHourlyWage: 28.16, minutesToEarn500: 1065 }, // ~17.75 hours
  { name: 'Германия', flag: '🇩🇪', avgHourlyWage: 24.5, minutesToEarn500: 1224 }, // ~20.4 hours
  { name: 'Россия', flag: '🇷🇺', avgHourlyWage: 6.5, minutesToEarn500: 4615 }, // ~76.9 hours
  { name: 'Индия', flag: '🇮🇳', avgHourlyWage: 2.8, minutesToEarn500: 10714 }, // ~178.6 hours
  { name: 'Китай', flag: '🇨🇳', avgHourlyWage: 5.2, minutesToEarn500: 5769 }, // ~96.15 hours
];

interface WorkTimerProps {
  userBalance?: number;
  onWorkComplete?: (earned: number) => void;
}

const WorkTimer = ({ userBalance = 0, onWorkComplete }: WorkTimerProps) => {
  const {
    isWorking,
    isPaused,
    timeRemaining,
    showStatsModal,
    lastCompletedSession,
    stats,
    startWork,
    pauseWork,
    resumeWork,
    completeWork,
    cancelWork,
    tick,
    closeStatsModal,
  } = useWorkStore();

  // Timer tick effect
  useEffect(() => {
    if (!isWorking || isPaused) return;

    const interval = setInterval(() => {
      tick();
    }, 1000);

    return () => clearInterval(interval);
  }, [isWorking, isPaused, tick]);

  // Format time as MM:SS
  const formatTime = (seconds: number): string => {
    const mins = Math.floor(seconds / 60);
    const secs = seconds % 60;
    return `${mins.toString().padStart(2, '0')}:${secs.toString().padStart(2, '0')}`;
  };

  // Calculate progress percentage
  const progressPercentage = ((180 - timeRemaining) / 180) * 100;

  // Format duration for display
  const formatDuration = (seconds: number): string => {
    const hours = Math.floor(seconds / 3600);
    const minutes = Math.floor((seconds % 3600) / 60);

    if (hours > 0) {
      return `${hours} ч ${minutes} мин`;
    }
    return `${minutes} мин`;
  };

  // Handle modal close attempt
  const handleCloseAttempt = useCallback(() => {
    if (isWorking && !isPaused) {
      // Show warning or do nothing when timer is active
      return;
    }
    cancelWork();
  }, [isWorking, isPaused, cancelWork]);

  // Prevent closing modal with Escape key when working
  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      if (e.key === 'Escape' && isWorking && !isPaused) {
        e.preventDefault();
        e.stopPropagation();
      }
    };

    if (isWorking) {
      window.addEventListener('keydown', handleKeyDown, true);
      return () => window.removeEventListener('keydown', handleKeyDown, true);
    }
  }, [isWorking, isPaused]);

  return (
    <>
      {/* Work Button (shown when balance is 0) */}
      {userBalance === 0 && !isWorking && (
        <motion.button
          whileHover={{ scale: 1.05 }}
          whileTap={{ scale: 0.95 }}
          onClick={startWork}
          className="fixed bottom-8 right-8 bg-gradient-to-r from-primary to-secondary text-white font-bold py-4 px-8 rounded-full shadow-2xl hover:shadow-primary/50 transition-all duration-300 z-50 flex items-center space-x-3"
        >
          <span className="text-2xl">💼</span>
          <span className="text-lg">Работать</span>
        </motion.button>
      )}

      {/* Work Timer Modal */}
      <AnimatePresence>
        {isWorking && (
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            className="fixed inset-0 bg-black/80 backdrop-blur-sm flex items-center justify-center z-50 p-4"
            onClick={(e) => {
              if (e.target === e.currentTarget) {
                handleCloseAttempt();
              }
            }}
          >
            <motion.div
              initial={{ scale: 0.9, opacity: 0 }}
              animate={{ scale: 1, opacity: 1 }}
              exit={{ scale: 0.9, opacity: 0 }}
              className="bg-gradient-to-br from-gray-800 to-gray-900 rounded-2xl p-8 max-w-md w-full border-2 border-gray-700 shadow-2xl"
              onClick={(e) => e.stopPropagation()}
            >
              {/* Header */}
              <div className="text-center mb-6">
                <div className="text-6xl mb-4">💼</div>
                <h2 className="text-3xl font-bold text-white mb-2">
                  Рабочий процесс
                </h2>
                <p className="text-gray-400">
                  Завершите работу, чтобы заработать
                </p>
              </div>

              {/* Timer Display */}
              <div className="mb-6">
                <div className="text-center mb-4">
                  <div className="text-6xl font-bold text-transparent bg-clip-text bg-gradient-to-r from-primary to-secondary">
                    {formatTime(timeRemaining)}
                  </div>
                  <p className="text-gray-400 mt-2">
                    до завершения
                  </p>
                </div>

                {/* Progress Bar */}
                <div className="relative h-4 bg-gray-700 rounded-full overflow-hidden">
                  <motion.div
                    initial={{ width: 0 }}
                    animate={{ width: `${progressPercentage}%` }}
                    transition={{ duration: 0.3 }}
                    className="absolute inset-y-0 left-0 bg-gradient-to-r from-primary to-secondary rounded-full"
                  />
                  {/* Animated shine effect */}
                  <motion.div
                    animate={{
                      x: ['-100%', '200%'],
                    }}
                    transition={{
                      duration: 2,
                      repeat: Infinity,
                      ease: 'linear',
                    }}
                    className="absolute inset-0 bg-gradient-to-r from-transparent via-white/20 to-transparent"
                  />
                </div>

                {/* Progress percentage */}
                <div className="text-center mt-2">
                  <span className="text-secondary font-semibold">
                    {progressPercentage.toFixed(1)}%
                  </span>
                </div>
              </div>

              {/* Reward Info */}
              <div className="bg-gray-700/50 rounded-xl p-4 mb-6 border border-gray-600">
                <div className="flex items-center justify-between">
                  <span className="text-gray-300">Награда:</span>
                  <span className="text-secondary font-bold text-xl">
                    +$500
                  </span>
                </div>
              </div>

              {/* Action Buttons */}
              <div className="flex space-x-3">
                {!isPaused ? (
                  <motion.button
                    whileHover={{ scale: 1.02 }}
                    whileTap={{ scale: 0.98 }}
                    onClick={pauseWork}
                    className="flex-1 bg-gray-700 hover:bg-gray-600 text-white font-semibold py-3 rounded-lg transition-colors"
                  >
                    ⏸ Пауза
                  </motion.button>
                ) : (
                  <motion.button
                    whileHover={{ scale: 1.02 }}
                    whileTap={{ scale: 0.98 }}
                    onClick={resumeWork}
                    className="flex-1 bg-green-600 hover:bg-green-700 text-white font-semibold py-3 rounded-lg transition-colors"
                  >
                    ▶️ Продолжить
                  </motion.button>
                )}

                <motion.button
                  whileHover={{ scale: 1.02 }}
                  whileTap={{ scale: 0.98 }}
                  onClick={cancelWork}
                  disabled={!isPaused}
                  className={`flex-1 font-semibold py-3 rounded-lg transition-colors ${
                    isPaused
                      ? 'bg-red-600 hover:bg-red-700 text-white'
                      : 'bg-gray-800 text-gray-500 cursor-not-allowed'
                  }`}
                >
                  ❌ Отменить
                </motion.button>
              </div>

              {/* Warning */}
              {!isPaused && (
                <p className="text-center text-gray-500 text-sm mt-4">
                  ⚠️ Поставьте на паузу, чтобы отменить
                </p>
              )}
            </motion.div>
          </motion.div>
        )}
      </AnimatePresence>

      {/* Stats Modal */}
      <AnimatePresence>
        {showStatsModal && lastCompletedSession && (
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            className="fixed inset-0 bg-black/80 backdrop-blur-sm flex items-center justify-center z-50 p-4"
            onClick={closeStatsModal}
          >
            <motion.div
              initial={{ scale: 0.9, opacity: 0, y: 20 }}
              animate={{ scale: 1, opacity: 1, y: 0 }}
              exit={{ scale: 0.9, opacity: 0, y: 20 }}
              className="bg-gradient-to-br from-gray-800 to-gray-900 rounded-2xl p-8 max-w-2xl w-full border-2 border-secondary shadow-2xl max-h-[90vh] overflow-y-auto"
              onClick={(e) => e.stopPropagation()}
            >
              {/* Success Header */}
              <div className="text-center mb-6">
                <motion.div
                  initial={{ scale: 0 }}
                  animate={{ scale: 1 }}
                  transition={{ delay: 0.2, type: 'spring' }}
                  className="text-8xl mb-4"
                >
                  ✅
                </motion.div>
                <h2 className="text-4xl font-bold text-transparent bg-clip-text bg-gradient-to-r from-primary to-secondary mb-2">
                  Работа завершена!
                </h2>
                <p className="text-gray-400">
                  Вы заработали <span className="text-secondary font-bold">${lastCompletedSession.earned}</span>
                </p>
              </div>

              {/* Stats Grid */}
              <div className="grid grid-cols-2 gap-4 mb-6">
                <div className="bg-gray-700/50 rounded-xl p-4 border border-gray-600">
                  <div className="text-gray-400 text-sm mb-1">Всего сессий</div>
                  <div className="text-2xl font-bold text-white">
                    {stats.sessions_count}
                  </div>
                </div>
                <div className="bg-gray-700/50 rounded-xl p-4 border border-gray-600">
                  <div className="text-gray-400 text-sm mb-1">Всего заработано</div>
                  <div className="text-2xl font-bold text-secondary">
                    ${stats.total_earned.toLocaleString()}
                  </div>
                </div>
                <div className="bg-gray-700/50 rounded-xl p-4 border border-gray-600 col-span-2">
                  <div className="text-gray-400 text-sm mb-1">Общее время работы</div>
                  <div className="text-2xl font-bold text-white">
                    {formatDuration(stats.total_work_time)}
                  </div>
                </div>
              </div>

              {/* Country Comparison */}
              <div className="mb-6">
                <h3 className="text-xl font-bold text-white mb-4">
                  📊 Сравнение с реальными зарплатами
                </h3>
                <div className="space-y-3">
                  {countryComparisons.map((country) => (
                    <div
                      key={country.name}
                      className="bg-gray-700/30 rounded-lg p-4 border border-gray-600"
                    >
                      <div className="flex items-center justify-between mb-2">
                        <div className="flex items-center space-x-2">
                          <span className="text-2xl">{country.flag}</span>
                          <span className="text-white font-semibold">
                            {country.name}
                          </span>
                        </div>
                        <span className="text-gray-400 text-sm">
                          ${country.avgHourlyWage}/час
                        </span>
                      </div>
                      <div className="text-gray-300 text-sm">
                        Время для заработка $500:{' '}
                        <span className="text-secondary font-semibold">
                          {formatDuration(country.minutesToEarn500 * 60)}
                        </span>
                      </div>
                      <div className="mt-2 text-primary text-sm font-semibold">
                        В игре: 3 минуты (в {Math.floor(country.minutesToEarn500 / 3)}x быстрее!)
                      </div>
                    </div>
                  ))}
                </div>
              </div>

              {/* Educational Message */}
              <div className="bg-primary/10 border border-primary/50 rounded-xl p-4 mb-6">
                <p className="text-gray-300 text-sm leading-relaxed">
                  <span className="font-bold text-primary">⚠️ Помните:</span>{' '}
                  В реальной жизни заработок требует намного больше времени и усилий.
                  Азартные игры — это развлечение с высоким риском потерь.
                </p>
              </div>

              {/* Close Button */}
              <motion.button
                whileHover={{ scale: 1.02 }}
                whileTap={{ scale: 0.98 }}
                onClick={() => {
                  closeStatsModal();
                  if (onWorkComplete) {
                    onWorkComplete(lastCompletedSession.earned);
                  }
                }}
                className="w-full bg-gradient-to-r from-primary to-secondary text-white font-bold py-4 rounded-xl hover:shadow-lg hover:shadow-primary/50 transition-all"
              >
                Продолжить
              </motion.button>
            </motion.div>
          </motion.div>
        )}
      </AnimatePresence>
    </>
  );
};

export default WorkTimer;
