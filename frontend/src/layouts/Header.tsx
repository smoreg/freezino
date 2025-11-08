import { Link } from 'react-router-dom';

const Header = () => {
  return (
    <header className="bg-gray-800 border-b border-gray-700">
      <div className="container mx-auto px-4 py-4">
        <div className="flex items-center justify-between">
          <Link to="/" className="flex items-center space-x-2">
            <span className="text-2xl font-bold text-primary">🎰 FREEZINO</span>
          </Link>

          <div className="flex items-center space-x-4">
            <div className="flex items-center space-x-2 bg-gray-700 px-4 py-2 rounded-lg">
              <span className="text-secondary text-xl font-bold">💰</span>
              <span className="text-white font-semibold">$0</span>
            </div>

            <div className="flex items-center space-x-2">
              <div className="w-10 h-10 bg-gray-600 rounded-full flex items-center justify-center">
                <span className="text-xl">👤</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </header>
  );
};

export default Header;
