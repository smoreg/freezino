import { Link } from 'react-router-dom';

const Footer = () => {
  return (
    <footer className="bg-gray-800 border-t border-gray-700 mt-auto">
      <div className="container mx-auto px-4 py-6">
        <div className="flex flex-col md:flex-row items-center justify-between space-y-4 md:space-y-0">
          {/* Copyright */}
          <div className="text-gray-400 text-sm text-center md:text-left">
            <p>&copy; 2024 Freezino. Образовательный проект против игровой зависимости.</p>
          </div>

          {/* Legal Links */}
          <div className="flex items-center flex-wrap justify-center gap-3 text-gray-400 text-sm">
            <Link to="/terms" className="hover:text-red-400 transition-colors">
              Terms of Service
            </Link>
            <span className="hidden sm:inline">•</span>
            <Link to="/privacy" className="hover:text-red-400 transition-colors">
              Privacy Policy
            </Link>
            <span className="hidden sm:inline">•</span>
            <Link to="/cookies" className="hover:text-red-400 transition-colors">
              Cookie Policy
            </Link>
          </div>

          {/* Disclaimer */}
          <div className="flex items-center gap-3 text-gray-400 text-sm text-center md:text-right">
            <span className="whitespace-nowrap">❌ Только виртуальные деньги</span>
            <span className="hidden sm:inline">•</span>
            <span className="whitespace-nowrap">✅ Никаких реальных ставок</span>
          </div>
        </div>

        {/* Additional Info */}
        <div className="mt-4 text-center text-gray-500 text-xs">
          <p>18+ | Freezino is an educational platform. No real money gambling occurs here.</p>
        </div>
      </div>
    </footer>
  );
};

export default Footer;
