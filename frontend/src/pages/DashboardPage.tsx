import { useState, useEffect } from 'react';
import { motion } from 'framer-motion';
import GameCard from '../components/GameCard';

interface Game {
  id: string;
  title: string;
  icon: string;
  description: string;
  minBet: number;
  isComingSoon: boolean;
}

const DashboardPage = () => {
  const [games, setGames] = useState<Game[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    // Simulate API call to fetch games
    const fetchGames = async () => {
      try {
        // TODO: Replace with actual API call
        await new Promise(resolve => setTimeout(resolve, 800));

        const gamesData: Game[] = [
          {
            id: '1',
            title: 'Рулетка',
            icon: '🎡',
            description: 'Классическая европейская рулетка',
            minBet: 10,
            isComingSoon: false,
          },
          {
            id: '2',
            title: 'Слоты',
            icon: '🎰',
            description: 'Однорукий бандит с 5 барабанами',
            minBet: 5,
            isComingSoon: false,
          },
          {
            id: '3',
            title: 'Блэкджек',
            icon: '🃏',
            description: 'Карточная игра против дилера',
            minBet: 20,
            isComingSoon: true,
          },
          {
            id: '4',
            title: 'Кости',
            icon: '🎲',
            description: 'Классическая игра Craps',
            minBet: 15,
            isComingSoon: true,
          },
          {
            id: '5',
            title: 'Crash',
            icon: '📈',
            description: 'График с растущим множителем',
            minBet: 10,
            isComingSoon: true,
          },
          {
            id: '6',
            title: 'Hi-Lo',
            icon: '🔼',
            description: 'Угадай выше или ниже',
            minBet: 5,
            isComingSoon: true,
          },
          {
            id: '7',
            title: 'Колесо Фортуны',
            icon: '🎪',
            description: 'Крути колесо и выиграй приз',
            minBet: 10,
            isComingSoon: true,
          },
          {
            id: '8',
            title: 'Покер',
            icon: '♠️',
            description: 'Video Poker - 5 карт',
            minBet: 25,
            isComingSoon: true,
          },
        ];

        setGames(gamesData);
      } catch (error) {
        console.error('Failed to fetch games:', error);
      } finally {
        setLoading(false);
      }
    };

    fetchGames();
  }, []);

  const handleGameClick = (gameId: string) => {
    console.log('Game clicked:', gameId);
    // TODO: Navigate to game page
  };

  return (
    <div className="min-h-screen">
      {/* Welcome Section */}
      <motion.div
        initial={{ opacity: 0, y: -20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.5 }}
        className="mb-8"
      >
        <h1 className="text-4xl font-bold text-white mb-2">
          Добро пожаловать в <span className="text-primary">Freezino</span>
        </h1>
        <p className="text-gray-400">
          Выберите игру и испытайте удачу! Играйте на виртуальные деньги безопасно.
        </p>
      </motion.div>

      {/* Stats Section */}
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.5, delay: 0.1 }}
        className="grid grid-cols-1 md:grid-cols-3 gap-4 mb-8"
      >
        <div className="bg-gray-800 border border-gray-700 rounded-xl p-6">
          <div className="flex items-center space-x-3">
            <span className="text-4xl">🎮</span>
            <div>
              <p className="text-gray-400 text-sm">Доступно игр</p>
              <p className="text-2xl font-bold text-white">
                {loading ? '...' : games.filter(g => !g.isComingSoon).length}
              </p>
            </div>
          </div>
        </div>

        <div className="bg-gray-800 border border-gray-700 rounded-xl p-6">
          <div className="flex items-center space-x-3">
            <span className="text-4xl">⏰</span>
            <div>
              <p className="text-gray-400 text-sm">Времени играно</p>
              <p className="text-2xl font-bold text-white">0 ч</p>
            </div>
          </div>
        </div>

        <div className="bg-gray-800 border border-gray-700 rounded-xl p-6">
          <div className="flex items-center space-x-3">
            <span className="text-4xl">🏆</span>
            <div>
              <p className="text-gray-400 text-sm">Всего выиграно</p>
              <p className="text-2xl font-bold text-secondary">$0</p>
            </div>
          </div>
        </div>
      </motion.div>

      {/* Games Section */}
      <div className="mb-8">
        <h2 className="text-2xl font-bold text-white mb-4">Казино Игры</h2>

        {loading ? (
          // Loading Skeleton
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
            {[...Array(8)].map((_, index) => (
              <div
                key={index}
                className="bg-gray-800 border border-gray-700 rounded-xl p-6 animate-pulse"
              >
                <div className="flex items-center justify-center mb-4">
                  <div className="w-24 h-24 bg-gray-700 rounded-full"></div>
                </div>
                <div className="h-6 bg-gray-700 rounded mb-2"></div>
                <div className="h-4 bg-gray-700 rounded mb-4"></div>
                <div className="h-8 bg-gray-700 rounded"></div>
              </div>
            ))}
          </div>
        ) : (
          // Games Grid
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            transition={{ duration: 0.5, delay: 0.2 }}
            className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6"
          >
            {games.map((game, index) => (
              <motion.div
                key={game.id}
                initial={{ opacity: 0, y: 20 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ duration: 0.3, delay: index * 0.05 }}
              >
                <GameCard
                  title={game.title}
                  icon={game.icon}
                  description={game.description}
                  minBet={game.minBet}
                  isComingSoon={game.isComingSoon}
                  onClick={() => handleGameClick(game.id)}
                />
              </motion.div>
            ))}
          </motion.div>
        )}
      </div>

      {/* Info Banner */}
      <motion.div
        initial={{ opacity: 0 }}
        animate={{ opacity: 1 }}
        transition={{ duration: 0.5, delay: 0.3 }}
        className="bg-gradient-to-r from-primary/20 to-secondary/20 border border-primary/30 rounded-xl p-6 text-center"
      >
        <p className="text-white font-semibold mb-2">
          💡 Помните: это образовательный проект
        </p>
        <p className="text-gray-300 text-sm">
          Вы играете на виртуальные деньги. Цель - показать, как быстро можно потерять деньги в казино.
        </p>
      </motion.div>
    </div>
  );
};

export default DashboardPage;
