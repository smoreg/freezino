import { motion, AnimatePresence } from 'framer-motion';

interface CountryComparison {
  code: string;
  name: string;
  flag: string;
  avgSalaryPerHour: number;
  timeNeeded: string;
}

interface StatsModalProps {
  isOpen: boolean;
  onClose: () => void;
  earned: number;
  totalWorkTime: number; // in seconds
}

const StatsModal = ({ isOpen, onClose, earned, totalWorkTime }: StatsModalProps) => {
  // Country data with average hourly wages in USD
  const countries: CountryComparison[] = [
    {
      code: 'US',
      name: 'США',
      flag: '🇺🇸',
      avgSalaryPerHour: 30,
      timeNeeded: '',
    },
    {
      code: 'GB',
      name: 'Великобритания',
      flag: '🇬🇧',
      avgSalaryPerHour: 25,
      timeNeeded: '',
    },
    {
      code: 'DE',
      name: 'Германия',
      flag: '🇩🇪',
      avgSalaryPerHour: 22,
      timeNeeded: '',
    },
    {
      code: 'FR',
      name: 'Франция',
      flag: '🇫🇷',
      avgSalaryPerHour: 20,
      timeNeeded: '',
    },
    {
      code: 'JP',
      name: 'Япония',
      flag: '🇯🇵',
      avgSalaryPerHour: 18,
      timeNeeded: '',
    },
    {
      code: 'CN',
      name: 'Китай',
      flag: '🇨🇳',
      avgSalaryPerHour: 8,
      timeNeeded: '',
    },
    {
      code: 'RU',
      name: 'Россия',
      flag: '🇷🇺',
      avgSalaryPerHour: 5,
      timeNeeded: '',
    },
    {
      code: 'IN',
      name: 'Индия',
      flag: '🇮🇳',
      avgSalaryPerHour: 3,
      timeNeeded: '',
    },
  ];

  // Calculate time needed to earn the same amount in each country
  const countriesWithTime = countries.map((country) => {
    const hoursNeeded = earned / country.avgSalaryPerHour;
    const minutesNeeded = hoursNeeded * 60;

    let timeNeeded = '';
    if (minutesNeeded < 60) {
      timeNeeded = `${minutesNeeded.toFixed(1)} минут`;
    } else if (minutesNeeded < 1440) {
      // less than 24 hours
      const hours = Math.floor(minutesNeeded / 60);
      const minutes = Math.floor(minutesNeeded % 60);
      timeNeeded = minutes > 0
        ? `${hours} ч ${minutes} мин`
        : `${hours} ч`;
    } else {
      const days = Math.floor(minutesNeeded / 1440);
      const hours = Math.floor((minutesNeeded % 1440) / 60);
      timeNeeded = hours > 0
        ? `${days} дн ${hours} ч`
        : `${days} дн`;
    }

    return { ...country, timeNeeded };
  });

  // Format total work time
  const formatWorkTime = (seconds: number): string => {
    const hours = Math.floor(seconds / 3600);
    const minutes = Math.floor((seconds % 3600) / 60);
    const secs = seconds % 60;

    if (hours > 0) {
      return `${hours} ч ${minutes} мин`;
    } else if (minutes > 0) {
      return `${minutes} мин ${secs} сек`;
    } else {
      return `${secs} сек`;
    }
  };

  return (
    <AnimatePresence>
      {isOpen && (
        <>
          {/* Backdrop */}
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            onClick={onClose}
            className="fixed inset-0 bg-black/70 backdrop-blur-sm z-40"
          />

          {/* Modal */}
          <motion.div
            initial={{ opacity: 0, scale: 0.9, y: 20 }}
            animate={{ opacity: 1, scale: 1, y: 0 }}
            exit={{ opacity: 0, scale: 0.9, y: 20 }}
            transition={{ type: 'spring', duration: 0.5 }}
            className="fixed inset-0 z-50 flex items-center justify-center p-4"
          >
            <div className="bg-gray-900 border-2 border-primary/30 rounded-2xl shadow-2xl max-w-2xl w-full max-h-[90vh] overflow-y-auto">
              {/* Header */}
              <div className="sticky top-0 bg-gray-900 border-b border-gray-800 p-6 flex items-center justify-between">
                <div>
                  <h2 className="text-2xl font-bold text-white flex items-center gap-2">
                    <span className="text-3xl">📊</span>
                    Статистика работы
                  </h2>
                  <p className="text-gray-400 text-sm mt-1">
                    Сравнение с мировыми зарплатами
                  </p>
                </div>
                <button
                  onClick={onClose}
                  className="text-gray-400 hover:text-white transition-colors p-2 hover:bg-gray-800 rounded-lg"
                  aria-label="Закрыть"
                >
                  <svg
                    width="24"
                    height="24"
                    viewBox="0 0 24 24"
                    fill="none"
                    stroke="currentColor"
                    strokeWidth="2"
                    strokeLinecap="round"
                    strokeLinejoin="round"
                  >
                    <line x1="18" y1="6" x2="6" y2="18"></line>
                    <line x1="6" y1="6" x2="18" y2="18"></line>
                  </svg>
                </button>
              </div>

              {/* Content */}
              <div className="p-6 space-y-6">
                {/* Earned Amount */}
                <motion.div
                  initial={{ opacity: 0, y: 10 }}
                  animate={{ opacity: 1, y: 0 }}
                  transition={{ delay: 0.1 }}
                  className="bg-gradient-to-r from-primary/20 to-secondary/20 border border-primary/30 rounded-xl p-6 text-center"
                >
                  <p className="text-gray-300 text-sm mb-2">Вы заработали</p>
                  <p className="text-5xl font-bold text-secondary">
                    ${earned.toFixed(2)}
                  </p>
                  <p className="text-gray-400 text-sm mt-2">
                    За {formatWorkTime(totalWorkTime)}
                  </p>
                </motion.div>

                {/* Country Comparison Table */}
                <motion.div
                  initial={{ opacity: 0, y: 10 }}
                  animate={{ opacity: 1, y: 0 }}
                  transition={{ delay: 0.2 }}
                >
                  <h3 className="text-lg font-semibold text-white mb-4">
                    Сравнение с другими странами
                  </h3>
                  <p className="text-gray-400 text-sm mb-4">
                    Чтобы заработать <span className="text-secondary font-semibold">${earned}</span>,
                    в других странах вам потребовалось бы:
                  </p>

                  <div className="space-y-3">
                    {countriesWithTime.map((country, index) => (
                      <motion.div
                        key={country.code}
                        initial={{ opacity: 0, x: -20 }}
                        animate={{ opacity: 1, x: 0 }}
                        transition={{ delay: 0.3 + index * 0.05 }}
                        className="bg-gray-800 border border-gray-700 rounded-lg p-4 hover:border-primary/50 transition-colors"
                      >
                        <div className="flex items-center justify-between">
                          <div className="flex items-center gap-3">
                            <span className="text-3xl">{country.flag}</span>
                            <div>
                              <p className="text-white font-semibold">
                                {country.name}
                              </p>
                              <p className="text-gray-400 text-xs">
                                ~${country.avgSalaryPerHour}/час
                              </p>
                            </div>
                          </div>
                          <div className="text-right">
                            <p className="text-secondary font-bold text-lg">
                              {country.timeNeeded}
                            </p>
                            <p className="text-gray-500 text-xs">работы</p>
                          </div>
                        </div>
                      </motion.div>
                    ))}
                  </div>
                </motion.div>

                {/* Educational Note */}
                <motion.div
                  initial={{ opacity: 0, y: 10 }}
                  animate={{ opacity: 1, y: 0 }}
                  transition={{ delay: 0.8 }}
                  className="bg-blue-500/10 border border-blue-500/30 rounded-lg p-4"
                >
                  <p className="text-blue-300 text-sm">
                    <span className="font-semibold">💡 Обратите внимание:</span> Данные о зарплатах усреднены
                    и приведены для сравнения. Реальные значения могут отличаться в зависимости от профессии и региона.
                  </p>
                </motion.div>
              </div>

              {/* Footer */}
              <div className="sticky bottom-0 bg-gray-900 border-t border-gray-800 p-6">
                <button
                  onClick={onClose}
                  className="w-full bg-primary hover:bg-primary/80 text-white font-semibold py-3 px-6 rounded-lg transition-colors"
                >
                  Закрыть
                </button>
              </div>
            </div>
          </motion.div>
        </>
      )}
    </AnimatePresence>
  );
};

export default StatsModal;
