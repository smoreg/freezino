import { Link } from 'react-router-dom';
import { useState, useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import LanguageSwitcher from '../LanguageSwitcher';
import { useSoundStore } from '../../store/soundStore';
import { soundManager } from '../../utils/sounds';

interface User {
  name: string;
  avatar: string;
  balance: number;
}

const Header = () => {
  const { t } = useTranslation();
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);
  const { isMusicEnabled, isSfxEnabled, musicVolume, toggleMusic, toggleSfx } = useSoundStore();

  useEffect(() => {
    // Simulate API call to fetch user data
    const fetchUser = async () => {
      try {
        // TODO: Replace with actual API call
        await new Promise(resolve => setTimeout(resolve, 500));

        setUser({
          name: 'Player',
          avatar: '👤',
          balance: 1000,
        });
      } catch (error) {
        console.error('Failed to fetch user data:', error);
      } finally {
        setLoading(false);
      }
    };

    fetchUser();
  }, []);

  // Initialize sound manager
  useEffect(() => {
    soundManager.init();
  }, []);

  // Control background music
  useEffect(() => {
    if (isMusicEnabled) {
      soundManager.playMusic(musicVolume);
    } else {
      soundManager.stopMusic();
    }
  }, [isMusicEnabled, musicVolume]);

  const handleMusicToggle = () => {
    toggleMusic();
    if (isSfxEnabled) {
      soundManager.play('click', 0.3);
    }
  };

  const handleSfxToggle = () => {
    soundManager.play('click', 0.3);
    toggleSfx();
  };

  return (
    <header className="bg-gray-800 border-b border-gray-700 sticky top-0 z-50">
      <div className="container mx-auto px-4 py-4">
        <div className="flex items-center justify-between">
          {/* Logo */}
          <Link to="/dashboard" className="flex items-center space-x-2 hover:opacity-80 transition-opacity">
            <span className="text-2xl font-bold text-primary">🎰 FREEZINO</span>
          </Link>

          {/* User Info */}
          <div className="flex items-center space-x-4">
            {/* Sound Controls */}
            <div className="flex items-center space-x-2">
              <button
                onClick={handleMusicToggle}
                className="p-2 hover:bg-gray-700 rounded-lg transition-colors"
                title={isMusicEnabled ? t('sound.musicOn') || 'Music On' : t('sound.musicOff') || 'Music Off'}
              >
                <span className="text-xl">{isMusicEnabled ? '🎵' : '🔇'}</span>
              </button>
              <button
                onClick={handleSfxToggle}
                className="p-2 hover:bg-gray-700 rounded-lg transition-colors"
                title={isSfxEnabled ? t('sound.sfxOn') || 'Sound Effects On' : t('sound.sfxOff') || 'Sound Effects Off'}
              >
                <span className="text-xl">{isSfxEnabled ? '🔊' : '🔈'}</span>
              </button>
            </div>

            {/* Language Switcher */}
            <LanguageSwitcher />
            {/* Balance */}
            <div className="flex items-center space-x-2 bg-gray-700 px-4 py-2 rounded-lg hover:bg-gray-600 transition-colors">
              <span className="text-secondary text-xl font-bold">💰</span>
              {loading ? (
                <div className="w-16 h-5 bg-gray-600 animate-pulse rounded"></div>
              ) : (
                <span className="text-white font-semibold">
                  {t('common.currency', { amount: user?.balance || 0 })}
                </span>
              )}
            </div>

            {/* Avatar */}
            <Link to="/profile" className="flex items-center space-x-2 hover:opacity-80 transition-opacity">
              {loading ? (
                <div className="w-10 h-10 bg-gray-600 animate-pulse rounded-full"></div>
              ) : (
                <div className="w-10 h-10 bg-gradient-to-br from-primary to-secondary rounded-full flex items-center justify-center border-2 border-gray-600 hover:border-secondary transition-colors">
                  <span className="text-xl">{user?.avatar || '👤'}</span>
                </div>
              )}
              {!loading && (
                <span className="hidden md:block text-white font-medium">{user?.name}</span>
              )}
            </Link>
          </div>
        </div>
      </div>
    </header>
  );
};

export default Header;
