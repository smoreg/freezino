import WorkTimer from '../components/WorkTimer';
import api from '../services/api';
import { useAuthStore } from '../store/authStore';
import type { User } from '../types';

const Home = () => {
  const { user, setUser } = useAuthStore();

  const handleWorkComplete = async (earned: number) => {
    try {
      // Refetch user data to get updated balance
      const response = await api.get<{ user: User }>('/auth/me');
      setUser(response.data.user);
    } catch (error) {
      console.error('Failed to update user balance:', error);
      // Fallback: update balance optimistically
      if (user) {
        setUser({
          ...user,
          balance: user.balance + earned,
        });
      }
    }
  };

  return (
    <div className="space-y-6">
      <h1 className="text-4xl font-bold text-white mb-2">
        Добро пожаловать в Freezino! 🎰
      </h1>
      <p className="text-gray-400 text-lg">
        Образовательный казино-симулятор для борьбы с игровой зависимостью
      </p>

      {/* Work Timer Component */}
      <WorkTimer
        userBalance={user?.balance || 0}
        onWorkComplete={handleWorkComplete}
      />

      <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mt-8">
        <div className="bg-gray-800 p-6 rounded-lg border border-gray-700">
          <div className="text-4xl mb-4">🎮</div>
          <h3 className="text-xl font-semibold text-white mb-2">Играйте</h3>
          <p className="text-gray-400">
            10+ различных казино-игр с виртуальными деньгами
          </p>
        </div>

        <div className="bg-gray-800 p-6 rounded-lg border border-gray-700">
          <div className="text-4xl mb-4">⏰</div>
          <h3 className="text-xl font-semibold text-white mb-2">Работайте</h3>
          <p className="text-gray-400">
            Зарабатывайте виртуальные деньги, работая 3 минуты = $500
          </p>
        </div>

        <div className="bg-gray-800 p-6 rounded-lg border border-gray-700">
          <div className="text-4xl mb-4">📊</div>
          <h3 className="text-xl font-semibold text-white mb-2">Анализируйте</h3>
          <p className="text-gray-400">
            Смотрите статистику и сравнивайте с реальными зарплатами
          </p>
        </div>
      </div>

      <div className="bg-primary bg-opacity-20 border border-primary p-6 rounded-lg mt-8">
        <h3 className="text-xl font-semibold text-white mb-2">⚠️ Важно!</h3>
        <ul className="text-gray-300 space-y-2">
          <li>❌ НЕЛЬЗЯ вводить реальные деньги</li>
          <li>✅ Только виртуальная валюта (псевдодоллары)</li>
          <li>🎓 Образовательная цель - показать проблему азартных игр</li>
        </ul>
      </div>
    </div>
  );
};

export default Home;
