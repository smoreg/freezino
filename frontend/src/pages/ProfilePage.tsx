import { useEffect, useState } from 'react';
import { motion } from 'framer-motion';
import { useTranslation } from 'react-i18next';
import Avatar from '../components/profile/Avatar';
import { UserItem } from '../types';
import { PageTransition, ProfileSkeleton, AnimatedButton, shakeVariants } from '../components/animations';

interface ProfileStats {
  user_id: number;
  total_work_time: number;
  total_earned: number;
  total_game_sessions: number;
  total_bet: number;
  total_won: number;
  total_lost: number;
  net_profit: number;
  favorite_game: string;
  biggest_win: number;
  biggest_loss: number;
}

interface UserProfile {
  id: number;
  name: string;
  email: string;
  avatar: string;
  balance: number;
  created_at: string;
}

const ProfilePage = () => {
  const { t } = useTranslation();
  const [profile, setProfile] = useState<UserProfile | null>(null);
  const [stats, setStats] = useState<ProfileStats | null>(null);
  const [items, setItems] = useState<UserItem[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    fetchProfileData();
  }, []);

  const fetchProfileData = async () => {
    setLoading(true);
    setError(null);

    try {
      // TODO: Replace with actual user ID from auth
      const userId = 1;

      // Fetch profile, stats, and items in parallel
      const [profileRes, statsRes, itemsRes] = await Promise.all([
        fetch(`http://localhost:3000/api/user/profile?user_id=${userId}`),
        fetch(`http://localhost:3000/api/user/stats?user_id=${userId}`),
        fetch(`http://localhost:3000/api/user/items?user_id=${userId}`),
      ]);

      if (!profileRes.ok || !statsRes.ok || !itemsRes.ok) {
        throw new Error('Failed to fetch profile data');
      }

      const [profileData, statsData, itemsData] = await Promise.all([
        profileRes.json(),
        statsRes.json(),
        itemsRes.json(),
      ]);

      setProfile(profileData.data);
      setStats(statsData.data);
      setItems(itemsData.data.items || []);
    } catch (err) {
      console.error('Error fetching profile:', err);
      setError('Failed to load profile data. Please try again.');
    } finally {
      setLoading(false);
    }
  };

  const formatTime = (seconds: number): string => {
    const hours = Math.floor(seconds / 3600);
    const minutes = Math.floor((seconds % 3600) / 60);
    return `${hours}h ${minutes}m`;
  };

  const formatCurrency = (amount: number): string => {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'USD',
      minimumFractionDigits: 0,
      maximumFractionDigits: 0,
    }).format(amount);
  };

  if (loading) {
    return (
      <PageTransition>
        <div className="min-h-screen bg-gradient-to-b from-gray-900 to-black py-12 px-4">
          <div className="container mx-auto max-w-6xl">
            <ProfileSkeleton />
          </div>
        </div>
      </PageTransition>
    );
  }

  if (error) {
    return (
      <PageTransition>
        <div className="min-h-screen bg-gradient-to-b from-gray-900 to-black py-12 px-4">
          <div className="container mx-auto max-w-6xl">
            <motion.div
              variants={shakeVariants}
              initial="initial"
              animate="animate"
              className="bg-red-900/20 border border-red-500 rounded-xl p-8 text-center"
            >
              <div className="text-6xl mb-4">⚠️</div>
              <h2 className="text-2xl font-bold text-red-400 mb-2">Error</h2>
              <p className="text-gray-300">{error}</p>
              <AnimatedButton
                variant="danger"
                onClick={fetchProfileData}
                className="mt-4"
              >
                Retry
              </AnimatedButton>
            </motion.div>
          </div>
        </div>
      </PageTransition>
    );
  }

  const equippedItems = items.filter(item => item.is_equipped);

  return (
    <PageTransition>
      <div className="min-h-screen bg-gradient-to-b from-gray-900 to-black py-12 px-4">
        <div className="container mx-auto max-w-6xl">
        {/* Header */}
        <motion.div
          initial={{ opacity: 0, y: -20 }}
          animate={{ opacity: 1, y: 0 }}
          className="mb-8"
        >
          <h1 className="text-4xl font-bold text-white mb-2">
            {profile?.name}'s Profile
          </h1>
          <p className="text-gray-400">
            Member since {profile ? new Date(profile.created_at).toLocaleDateString() : ''}
          </p>
        </motion.div>

        {/* Main Content Grid */}
        <div className="grid md:grid-cols-2 gap-8">
          {/* Left Column - Avatar */}
          <motion.div
            initial={{ opacity: 0, x: -20 }}
            animate={{ opacity: 1, x: 0 }}
            transition={{ delay: 0.1 }}
          >
            <div className="bg-gray-800 rounded-2xl p-6 border border-gray-700">
              <h2 className="text-2xl font-bold text-white mb-6 flex items-center">
                <span className="mr-2">👤</span> Your Avatar
              </h2>
              <Avatar equippedItems={equippedItems} />
            </div>
          </motion.div>

          {/* Right Column - Stats */}
          <div className="space-y-6">
            {/* Balance Card */}
            <motion.div
              initial={{ opacity: 0, x: 20 }}
              animate={{ opacity: 1, x: 0 }}
              transition={{ delay: 0.2 }}
              className="bg-gradient-to-br from-secondary to-primary rounded-2xl p-6 border border-gray-700 shadow-lg"
            >
              <h3 className="text-xl font-semibold text-white mb-2 flex items-center">
                <span className="mr-2">💰</span> Current Balance
              </h3>
              <p className="text-5xl font-bold text-white">
                {formatCurrency(profile?.balance || 0)}
              </p>
            </motion.div>

            {/* Work Stats */}
            <motion.div
              initial={{ opacity: 0, x: 20 }}
              animate={{ opacity: 1, x: 0 }}
              transition={{ delay: 0.3 }}
              className="bg-gray-800 rounded-2xl p-6 border border-gray-700"
            >
              <h3 className="text-xl font-semibold text-white mb-4 flex items-center">
                <span className="mr-2">⏱️</span> Work Statistics
              </h3>
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <p className="text-gray-400 text-sm">Total Time Worked</p>
                  <p className="text-2xl font-bold text-white">
                    {formatTime(stats?.total_work_time || 0)}
                  </p>
                </div>
                <div>
                  <p className="text-gray-400 text-sm">Total Earned</p>
                  <p className="text-2xl font-bold text-green-400">
                    {formatCurrency(stats?.total_earned || 0)}
                  </p>
                </div>
              </div>
            </motion.div>

            {/* Game Stats */}
            <motion.div
              initial={{ opacity: 0, x: 20 }}
              animate={{ opacity: 1, x: 0 }}
              transition={{ delay: 0.4 }}
              className="bg-gray-800 rounded-2xl p-6 border border-gray-700"
            >
              <h3 className="text-xl font-semibold text-white mb-4 flex items-center">
                <span className="mr-2">🎮</span> Gaming Statistics
              </h3>
              <div className="space-y-4">
                <div className="grid grid-cols-2 gap-4">
                  <div>
                    <p className="text-gray-400 text-sm">Games Played</p>
                    <p className="text-2xl font-bold text-white">
                      {stats?.total_game_sessions || 0}
                    </p>
                  </div>
                  <div>
                    <p className="text-gray-400 text-sm">Favorite Game</p>
                    <p className="text-xl font-bold text-secondary capitalize">
                      {stats?.favorite_game || 'None'}
                    </p>
                  </div>
                </div>
                <div className="grid grid-cols-2 gap-4">
                  <div>
                    <p className="text-gray-400 text-sm">Total Bet</p>
                    <p className="text-xl font-bold text-yellow-400">
                      {formatCurrency(stats?.total_bet || 0)}
                    </p>
                  </div>
                  <div>
                    <p className="text-gray-400 text-sm">Total Won</p>
                    <p className="text-xl font-bold text-green-400">
                      {formatCurrency(stats?.total_won || 0)}
                    </p>
                  </div>
                </div>
                <div className="grid grid-cols-2 gap-4">
                  <div>
                    <p className="text-gray-400 text-sm">Biggest Win</p>
                    <p className="text-xl font-bold text-green-400">
                      {formatCurrency(stats?.biggest_win || 0)}
                    </p>
                  </div>
                  <div>
                    <p className="text-gray-400 text-sm">Net Profit</p>
                    <p className={`text-xl font-bold ${
                      (stats?.net_profit || 0) >= 0 ? 'text-green-400' : 'text-red-400'
                    }`}>
                      {formatCurrency(stats?.net_profit || 0)}
                    </p>
                  </div>
                </div>
              </div>
            </motion.div>

            {/* Items Count */}
            <motion.div
              initial={{ opacity: 0, x: 20 }}
              animate={{ opacity: 1, x: 0 }}
              transition={{ delay: 0.5 }}
              className="bg-gray-800 rounded-2xl p-6 border border-gray-700"
            >
              <h3 className="text-xl font-semibold text-white mb-4 flex items-center">
                <span className="mr-2">🛍️</span> Inventory
              </h3>
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <p className="text-gray-400 text-sm">Total Items</p>
                  <p className="text-2xl font-bold text-white">{items.length}</p>
                </div>
                <div>
                  <p className="text-gray-400 text-sm">Equipped Items</p>
                  <p className="text-2xl font-bold text-secondary">
                    {equippedItems.length}
                  </p>
                </div>
              </div>
              <AnimatedButton
                variant="secondary"
                fullWidth
                onClick={() => window.location.href = '/shop'}
                className="mt-4"
              >
                Visit Shop
              </AnimatedButton>
            </motion.div>
          </div>
        </div>

        {/* All Items Section */}
        {items.length > 0 && (
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: 0.6 }}
            className="mt-8 bg-gray-800 rounded-2xl p-6 border border-gray-700"
          >
            <h2 className="text-2xl font-bold text-white mb-6 flex items-center">
              <span className="mr-2">📦</span> All Purchased Items
            </h2>
            <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4">
              {items.map((userItem) => (
                <div
                  key={userItem.id}
                  className={`p-4 rounded-xl border transition-all ${
                    userItem.is_equipped
                      ? 'bg-secondary/10 border-secondary shadow-lg shadow-secondary/20'
                      : 'bg-gray-900 border-gray-700 hover:border-gray-600'
                  }`}
                >
                  <div className="flex items-center justify-center mb-3 h-20">
                    {userItem.item.image_url ? (
                      <img
                        src={userItem.item.image_url}
                        alt={userItem.item.name}
                        className="max-h-full object-contain"
                      />
                    ) : (
                      <div className="text-4xl">
                        {userItem.item.type === 'house' && '🏠'}
                        {userItem.item.type === 'clothing' && '👕'}
                        {userItem.item.type === 'car' && '🚗'}
                        {userItem.item.type === 'accessories' && '💎'}
                      </div>
                    )}
                  </div>
                  <h4 className="font-semibold text-white text-center mb-1">
                    {userItem.item.name}
                  </h4>
                  <p className="text-xs text-gray-400 text-center mb-2 capitalize">
                    {userItem.item.type}
                  </p>
                  {userItem.is_equipped && (
                    <div className="text-xs text-secondary text-center font-semibold">
                      ✓ Equipped
                    </div>
                  )}
                  <p className="text-xs text-gray-500 text-center mt-2">
                    Purchased: {new Date(userItem.purchased_at).toLocaleDateString()}
                  </p>
                </div>
              ))}
            </div>
          </motion.div>
        )}
        </div>
      </div>
    </PageTransition>
  );
};

export default ProfilePage;
