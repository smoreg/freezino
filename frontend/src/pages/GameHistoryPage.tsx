import { useState, useEffect } from 'react';
import { motion } from 'framer-motion';
import {
  BarChart,
  Bar,
  PieChart,
  Pie,
  Cell,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
  ResponsiveContainer,
} from 'recharts';

interface GameHistoryItem {
  id: number;
  game_type: string;
  bet: number;
  win: number;
  profit: number;
  created_at: string;
}

interface GameStats {
  total_games: number;
  total_wins: number;
  total_losses: number;
  total_bet: number;
  total_won: number;
  net_profit: number;
  favorite_game: string;
  win_rate: number;
  biggest_win: number;
  biggest_loss: number;
  game_breakdown: {
    game_type: string;
    games_played: number;
    total_bet: number;
    total_won: number;
    net_profit: number;
  }[];
}

const GAME_NAMES: Record<string, string> = {
  roulette: 'Рулетка',
  slots: 'Слоты',
  blackjack: 'Блэкджек',
  craps: 'Кости',
  baccara: 'Баккара',
  wheel: 'Колесо фортуны',
  keno: 'Кено',
  poker: 'Покер',
  hilo: 'Hi-Lo',
  crash: 'Crash',
  bingo: 'Бинго',
  plinko: 'Plinko',
};

const COLORS = ['#DC2626', '#FBBF24', '#10B981', '#3B82F6', '#8B5CF6', '#EC4899'];

const GameHistoryPage = () => {
  const [history, setHistory] = useState<GameHistoryItem[]>([]);
  const [stats, setStats] = useState<GameStats | null>(null);
  const [loading, setLoading] = useState(true);
  const [selectedGame, setSelectedGame] = useState<string>('');
  const [currentPage, setCurrentPage] = useState(1);
  // const itemsPerPage = 10; // для будущей пагинации

  useEffect(() => {
    fetchHistory();
    fetchStats();
  }, [selectedGame, currentPage]);

  const fetchHistory = async () => {
    try {
      setLoading(true);

      // В реальном приложении использовать axios или fetch
      // const userId = 1; // получить из auth контекста
      // const offset = (currentPage - 1) * itemsPerPage;
      // const gameParam = selectedGame ? `&game=${selectedGame}` : '';
      // const response = await fetch(`/api/games/history?user_id=${userId}&limit=${itemsPerPage}&offset=${offset}${gameParam}`);
      // const data = await response.json();
      // setHistory(data.data.games);

      // Демо данные
      const mockData: GameHistoryItem[] = [
        { id: 1, game_type: 'roulette', bet: 100, win: 200, profit: 100, created_at: '2025-11-08T10:30:00Z' },
        { id: 2, game_type: 'slots', bet: 50, win: 0, profit: -50, created_at: '2025-11-08T10:25:00Z' },
        { id: 3, game_type: 'blackjack', bet: 200, win: 400, profit: 200, created_at: '2025-11-08T10:20:00Z' },
        { id: 4, game_type: 'roulette', bet: 150, win: 0, profit: -150, created_at: '2025-11-08T10:15:00Z' },
        { id: 5, game_type: 'slots', bet: 75, win: 150, profit: 75, created_at: '2025-11-08T10:10:00Z' },
        { id: 6, game_type: 'crash', bet: 100, win: 250, profit: 150, created_at: '2025-11-08T10:05:00Z' },
        { id: 7, game_type: 'poker', bet: 300, win: 0, profit: -300, created_at: '2025-11-08T10:00:00Z' },
        { id: 8, game_type: 'roulette', bet: 50, win: 100, profit: 50, created_at: '2025-11-08T09:55:00Z' },
      ];
      setHistory(mockData);
    } catch (error) {
      console.error('Failed to fetch game history:', error);
    } finally {
      setLoading(false);
    }
  };

  const fetchStats = async () => {
    try {
      // В реальном приложении:
      // const response = await fetch(`/api/games/stats?user_id=${userId}`);
      // const data = await response.json();
      // setStats(data.data);

      // Демо данные
      const mockStats: GameStats = {
        total_games: 42,
        total_wins: 18,
        total_losses: 24,
        total_bet: 5000,
        total_won: 3500,
        net_profit: -1500,
        favorite_game: 'roulette',
        win_rate: 42.86,
        biggest_win: 500,
        biggest_loss: 300,
        game_breakdown: [
          { game_type: 'roulette', games_played: 15, total_bet: 2000, total_won: 1500, net_profit: -500 },
          { game_type: 'slots', games_played: 12, total_bet: 1200, total_won: 800, net_profit: -400 },
          { game_type: 'blackjack', games_played: 8, total_bet: 1000, total_won: 900, net_profit: -100 },
          { game_type: 'crash', games_played: 5, total_bet: 500, total_won: 200, net_profit: -300 },
          { game_type: 'poker', games_played: 2, total_bet: 300, total_won: 100, net_profit: -200 },
        ],
      };
      setStats(mockStats);
    } catch (error) {
      console.error('Failed to fetch game stats:', error);
    }
  };

  const formatDate = (dateString: string) => {
    const date = new Date(dateString);
    return date.toLocaleString('ru-RU', {
      day: '2-digit',
      month: '2-digit',
      year: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    });
  };

  if (loading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-center">
          <div className="animate-spin rounded-full h-16 w-16 border-t-4 border-b-4 border-primary mx-auto"></div>
          <p className="text-gray-400 mt-4">Загрузка истории игр...</p>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen pb-8">
      {/* Header */}
      <motion.div
        initial={{ opacity: 0, y: -20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.5 }}
        className="mb-8"
      >
        <h1 className="text-3xl md:text-4xl font-bold text-white mb-2 flex items-center gap-2 md:gap-3">
          <span className="text-4xl md:text-5xl">🎮</span>
          История игр
        </h1>
        <p className="text-sm md:text-base text-gray-400">
          Подробная статистика и история всех ваших игр
        </p>
      </motion.div>

      {/* Stats Cards */}
      {stats && (
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.5, delay: 0.1 }}
          className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8"
        >
          <div className="bg-gray-800 border border-gray-700 rounded-xl p-6">
            <div className="flex items-center justify-between mb-4">
              <span className="text-4xl">🎲</span>
            </div>
            <p className="text-gray-400 text-sm mb-1">Всего игр</p>
            <p className="text-2xl font-bold text-white">{stats.total_games}</p>
          </div>

          <div className="bg-gray-800 border border-gray-700 rounded-xl p-6">
            <div className="flex items-center justify-between mb-4">
              <span className="text-4xl">🏆</span>
            </div>
            <p className="text-gray-400 text-sm mb-1">Процент побед</p>
            <p className="text-2xl font-bold text-secondary">{stats.win_rate.toFixed(1)}%</p>
          </div>

          <div className="bg-gray-800 border border-gray-700 rounded-xl p-6">
            <div className="flex items-center justify-between mb-4">
              <span className="text-4xl">💰</span>
            </div>
            <p className="text-gray-400 text-sm mb-1">Чистая прибыль</p>
            <p className={`text-2xl font-bold ${stats.net_profit >= 0 ? 'text-green-400' : 'text-red-400'}`}>
              ${stats.net_profit.toFixed(2)}
            </p>
          </div>

          <div className="bg-gray-800 border border-gray-700 rounded-xl p-6">
            <div className="flex items-center justify-between mb-4">
              <span className="text-4xl">❤️</span>
            </div>
            <p className="text-gray-400 text-sm mb-1">Любимая игра</p>
            <p className="text-2xl font-bold text-white">{GAME_NAMES[stats.favorite_game] || stats.favorite_game}</p>
          </div>
        </motion.div>
      )}

      {/* Charts */}
      {stats && (
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6 mb-8">
          {/* Games Breakdown Pie Chart */}
          <motion.div
            initial={{ opacity: 0, x: -20 }}
            animate={{ opacity: 1, x: 0 }}
            transition={{ duration: 0.5, delay: 0.2 }}
            className="bg-gray-800 border border-gray-700 rounded-xl p-6"
          >
            <h2 className="text-xl font-bold text-white mb-4">Распределение игр</h2>
            <ResponsiveContainer width="100%" height={300}>
              <PieChart>
                <Pie
                  data={stats.game_breakdown}
                  dataKey="games_played"
                  nameKey="game_type"
                  cx="50%"
                  cy="50%"
                  outerRadius={100}
                  label={(entry: { game_type: string }) => GAME_NAMES[entry.game_type] || entry.game_type}
                >
                  {stats.game_breakdown.map((_, index) => (
                    <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
                  ))}
                </Pie>
                <Tooltip
                  contentStyle={{ backgroundColor: '#1F2937', border: '1px solid #374151' }}
                  labelStyle={{ color: '#F3F4F6' }}
                />
              </PieChart>
            </ResponsiveContainer>
          </motion.div>

          {/* Profit/Loss Bar Chart */}
          <motion.div
            initial={{ opacity: 0, x: 20 }}
            animate={{ opacity: 1, x: 0 }}
            transition={{ duration: 0.5, delay: 0.2 }}
            className="bg-gray-800 border border-gray-700 rounded-xl p-6"
          >
            <h2 className="text-xl font-bold text-white mb-4">Прибыль/Убыток по играм</h2>
            <ResponsiveContainer width="100%" height={300}>
              <BarChart data={stats.game_breakdown}>
                <CartesianGrid strokeDasharray="3 3" stroke="#374151" />
                <XAxis
                  dataKey="game_type"
                  stroke="#9CA3AF"
                  tickFormatter={(value) => GAME_NAMES[value] || value}
                />
                <YAxis stroke="#9CA3AF" />
                <Tooltip
                  contentStyle={{ backgroundColor: '#1F2937', border: '1px solid #374151' }}
                  labelStyle={{ color: '#F3F4F6' }}
                  labelFormatter={(value) => GAME_NAMES[value] || value}
                />
                <Legend />
                <Bar dataKey="net_profit" name="Чистая прибыль" fill="#FBBF24" />
              </BarChart>
            </ResponsiveContainer>
          </motion.div>

          {/* Win/Loss Line Chart */}
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.5, delay: 0.3 }}
            className="bg-gray-800 border border-gray-700 rounded-xl p-6 lg:col-span-2"
          >
            <h2 className="text-xl font-bold text-white mb-4">Статистика выигрышей и проигрышей</h2>
            <ResponsiveContainer width="100%" height={300}>
              <BarChart
                data={[
                  { name: 'Победы', value: stats.total_wins },
                  { name: 'Поражения', value: stats.total_losses },
                ]}
              >
                <CartesianGrid strokeDasharray="3 3" stroke="#374151" />
                <XAxis dataKey="name" stroke="#9CA3AF" />
                <YAxis stroke="#9CA3AF" />
                <Tooltip
                  contentStyle={{ backgroundColor: '#1F2937', border: '1px solid #374151' }}
                  labelStyle={{ color: '#F3F4F6' }}
                />
                <Bar dataKey="value" fill="#DC2626">
                  {[
                    <Cell key="cell-0" fill="#10B981" />,
                    <Cell key="cell-1" fill="#DC2626" />,
                  ]}
                </Bar>
              </BarChart>
            </ResponsiveContainer>
          </motion.div>
        </div>
      )}

      {/* Filters */}
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.5, delay: 0.3 }}
        className="bg-gray-800 border border-gray-700 rounded-xl p-6 mb-6"
      >
        <div className="flex flex-wrap gap-4 items-center">
          <label className="text-white font-semibold">Фильтр по игре:</label>
          <select
            value={selectedGame}
            onChange={(e) => {
              setSelectedGame(e.target.value);
              setCurrentPage(1);
            }}
            className="bg-gray-700 border border-gray-600 text-white rounded-lg px-4 py-2 focus:outline-none focus:ring-2 focus:ring-primary"
          >
            <option value="">Все игры</option>
            {Object.entries(GAME_NAMES).map(([key, name]) => (
              <option key={key} value={key}>
                {name}
              </option>
            ))}
          </select>
        </div>
      </motion.div>

      {/* History Table */}
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.5, delay: 0.4 }}
        className="bg-gray-800 border border-gray-700 rounded-xl overflow-hidden"
      >
        <div className="overflow-x-auto">
          <table className="w-full">
            <thead className="bg-gray-900">
              <tr>
                <th className="px-3 md:px-6 py-3 md:py-4 text-left text-xs md:text-sm font-semibold text-gray-300">Дата</th>
                <th className="px-3 md:px-6 py-3 md:py-4 text-left text-xs md:text-sm font-semibold text-gray-300">Игра</th>
                <th className="hidden sm:table-cell px-3 md:px-6 py-3 md:py-4 text-right text-xs md:text-sm font-semibold text-gray-300">Ставка</th>
                <th className="hidden sm:table-cell px-3 md:px-6 py-3 md:py-4 text-right text-xs md:text-sm font-semibold text-gray-300">Выигрыш</th>
                <th className="px-3 md:px-6 py-3 md:py-4 text-right text-xs md:text-sm font-semibold text-gray-300">Прибыль</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-gray-700">
              {history.length === 0 ? (
                <tr>
                  <td colSpan={5} className="px-6 py-8 text-center text-gray-400">
                    Нет истории игр
                  </td>
                </tr>
              ) : (
                history.map((item) => (
                  <tr key={item.id} className="hover:bg-gray-700/50 transition-colors">
                    <td className="px-3 md:px-6 py-3 md:py-4 text-xs md:text-sm text-gray-300">
                      <div className="sm:hidden">{new Date(item.created_at).toLocaleDateString()}</div>
                      <div className="hidden sm:block">{formatDate(item.created_at)}</div>
                    </td>
                    <td className="px-3 md:px-6 py-3 md:py-4 text-xs md:text-sm text-white font-medium">
                      {GAME_NAMES[item.game_type] || item.game_type}
                    </td>
                    <td className="hidden sm:table-cell px-3 md:px-6 py-3 md:py-4 text-xs md:text-sm text-right text-gray-300">${item.bet.toFixed(2)}</td>
                    <td className="hidden sm:table-cell px-3 md:px-6 py-3 md:py-4 text-xs md:text-sm text-right text-secondary">${item.win.toFixed(2)}</td>
                    <td
                      className={`px-3 md:px-6 py-3 md:py-4 text-xs md:text-sm text-right font-semibold ${
                        item.profit >= 0 ? 'text-green-400' : 'text-red-400'
                      }`}
                    >
                      {item.profit >= 0 ? '+' : ''}${item.profit.toFixed(2)}
                    </td>
                  </tr>
                ))
              )}
            </tbody>
          </table>
        </div>

        {/* Pagination */}
        {history.length > 0 && (
          <div className="bg-gray-900 px-3 md:px-6 py-4 flex items-center justify-between gap-2">
            <button
              onClick={() => setCurrentPage((prev) => Math.max(1, prev - 1))}
              disabled={currentPage === 1}
              className="px-3 md:px-4 py-2.5 md:py-2 text-sm md:text-base bg-gray-700 text-white rounded-lg hover:bg-gray-600 disabled:opacity-50 disabled:cursor-not-allowed transition-colors touch-manipulation min-h-[44px]"
            >
              <span className="hidden sm:inline">Предыдущая</span>
              <span className="sm:hidden">←</span>
            </button>
            <span className="text-xs md:text-sm text-gray-400">Стр. {currentPage}</span>
            <button
              onClick={() => setCurrentPage((prev) => prev + 1)}
              className="px-3 md:px-4 py-2.5 md:py-2 text-sm md:text-base bg-gray-700 text-white rounded-lg hover:bg-gray-600 transition-colors touch-manipulation min-h-[44px]"
            >
              <span className="hidden sm:inline">Следующая</span>
              <span className="sm:hidden">→</span>
            </button>
          </div>
        )}
      </motion.div>
    </div>
  );
};

export default GameHistoryPage;
