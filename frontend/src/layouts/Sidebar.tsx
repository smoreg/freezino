import { Link, useLocation } from 'react-router-dom';

const Sidebar = () => {
  const location = useLocation();

  const menuItems = [
    { path: '/', label: 'Главная', icon: '🏠' },
    { path: '/games', label: 'Игры', icon: '🎮' },
    { path: '/work', label: 'Работать', icon: '⏰' },
    { path: '/shop', label: 'Магазин', icon: '🛍️' },
    { path: '/profile', label: 'Профиль', icon: '👤' },
    { path: '/stats', label: 'Статистика', icon: '📊' },
  ];

  const isActive = (path: string) => location.pathname === path;

  return (
    <aside className="bg-gray-800 border-r border-gray-700 w-64 min-h-screen">
      <nav className="p-4">
        <ul className="space-y-2">
          {menuItems.map((item) => (
            <li key={item.path}>
              <Link
                to={item.path}
                className={`flex items-center space-x-3 px-4 py-3 rounded-lg transition-colors ${
                  isActive(item.path)
                    ? 'bg-primary text-white'
                    : 'text-gray-300 hover:bg-gray-700 hover:text-white'
                }`}
              >
                <span className="text-xl">{item.icon}</span>
                <span className="font-medium">{item.label}</span>
              </Link>
            </li>
          ))}
        </ul>
      </nav>
    </aside>
  );
};

export default Sidebar;
