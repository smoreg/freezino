import { Link } from 'react-router-dom';
import { useState, useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import LanguageSwitcher from '../LanguageSwitcher';

interface User {
  name: string;
  avatar: string;
  balance: number;
}

interface HeaderProps {
  onMenuClick?: () => void;
}

const Header = ({ onMenuClick }: HeaderProps) => {
  const { t } = useTranslation();
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);

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

  return (
    <header className="bg-gray-800 border-b border-gray-700 sticky top-0 z-50">
      <div className="container mx-auto px-3 md:px-4 py-3 md:py-4">
        <div className="flex items-center justify-between">
          {/* Left: Hamburger + Logo */}
          <div className="flex items-center space-x-2 md:space-x-3">
            {/* Hamburger button for mobile */}
            {onMenuClick && (
              <button
                onClick={onMenuClick}
                className="lg:hidden text-gray-300 hover:text-white transition-colors p-2 -ml-2 touch-manipulation"
                aria-label="Toggle menu"
              >
                <svg
                  className="w-6 h-6"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M4 6h16M4 12h16M4 18h16"
                  />
                </svg>
              </button>
            )}

            {/* Logo */}
            <Link to="/dashboard" className="flex items-center space-x-2 hover:opacity-80 transition-opacity">
              <span className="text-xl md:text-2xl font-bold text-primary">🎰 FREEZINO</span>
            </Link>
          </div>

          {/* Right: Language + Balance + Avatar */}
          <div className="flex items-center space-x-2 md:space-x-4">
            {/* Language Switcher - hide on small screens */}
            <div className="hidden sm:block">
              <LanguageSwitcher />
            </div>

            {/* Balance - responsive sizing */}
            <div className="flex items-center space-x-1 md:space-x-2 bg-gray-700 px-2 md:px-4 py-1.5 md:py-2 rounded-lg hover:bg-gray-600 transition-colors">
              <span className="text-secondary text-lg md:text-xl font-bold">💰</span>
              {loading ? (
                <div className="w-12 md:w-16 h-4 md:h-5 bg-gray-600 animate-pulse rounded"></div>
              ) : (
                <span className="text-white text-sm md:text-base font-semibold">
                  ${user?.balance?.toFixed(0) || 0}
                </span>
              )}
            </div>

            {/* Avatar */}
            <Link to="/profile" className="flex items-center space-x-2 hover:opacity-80 transition-opacity touch-manipulation">
              {loading ? (
                <div className="w-8 h-8 md:w-10 md:h-10 bg-gray-600 animate-pulse rounded-full"></div>
              ) : (
                <div className="w-8 h-8 md:w-10 md:h-10 bg-gradient-to-br from-primary to-secondary rounded-full flex items-center justify-center border-2 border-gray-600 hover:border-secondary transition-colors">
                  <span className="text-lg md:text-xl">{user?.avatar || '👤'}</span>
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
