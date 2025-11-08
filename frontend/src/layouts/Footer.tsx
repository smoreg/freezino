import { useState } from 'react';
import { Link } from 'react-router-dom';

import CookieSettings from '../components/CookieSettings';

const Footer = () => {
  const [showCookieSettings, setShowCookieSettings] = useState(false);

  return (
    <>
      <footer className="bg-gray-800 border-t border-gray-700 mt-auto">
        <div className="container mx-auto px-4 py-6">
          <div className="flex flex-col space-y-4">
            {/* Main Footer Content */}
            <div className="flex flex-col md:flex-row items-center justify-between space-y-4 md:space-y-0">
              <div className="text-gray-400 text-sm text-center md:text-left">
                <p>&copy; 2024 Freezino. Образовательный проект против игровой зависимости.</p>
              </div>

              <div className="flex items-center space-x-4 text-gray-400 text-sm">
                <span>❌ Только виртуальные деньги</span>
                <span>✅ Никаких реальных ставок</span>
              </div>
            </div>

            {/* Links Section */}
            <div className="flex flex-col md:flex-row items-center justify-between space-y-3 md:space-y-0 pt-4 border-t border-gray-700">
              {/* Legal Links */}
              <div className="flex flex-wrap items-center justify-center gap-4 text-sm">
                <Link
                  to="/about"
                  className="text-gray-400 hover:text-secondary transition-colors"
                >
                  О проекте
                </Link>
                <Link
                  to="/terms"
                  className="text-gray-400 hover:text-secondary transition-colors"
                >
                  Terms of Service
                </Link>
                <Link
                  to="/privacy"
                  className="text-gray-400 hover:text-secondary transition-colors"
                >
                  Privacy Policy
                </Link>
                <Link
                  to="/cookies"
                  className="text-gray-400 hover:text-secondary transition-colors"
                >
                  Cookie Policy
                </Link>
                <Link
                  to="/contact"
                  className="text-gray-400 hover:text-secondary transition-colors"
                >
                  Контакты
                </Link>
              </div>

              {/* Cookie Settings Button */}
              <button
                onClick={() => setShowCookieSettings(true)}
                className="flex items-center space-x-1 text-sm text-gray-400 hover:text-secondary transition-colors group"
              >
                <span>🍪</span>
                <span className="underline">Cookie Settings</span>
              </button>
            </div>
          </div>
        </div>
      </footer>

      {/* Cookie Settings Modal */}
      <CookieSettings
        isOpen={showCookieSettings}
        onClose={() => setShowCookieSettings(false)}
      />
    </>
  );
};

export default Footer;
