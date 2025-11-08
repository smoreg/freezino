const Footer = () => {
  return (
    <footer className="bg-gray-800 border-t border-gray-700 mt-auto">
      <div className="container mx-auto px-4 py-6">
        <div className="flex flex-col md:flex-row items-center justify-between space-y-4 md:space-y-0">
          <div className="text-gray-400 text-sm">
            <p>&copy; 2024 Freezino. Образовательный проект против игровой зависимости.</p>
          </div>

          <div className="flex items-center space-x-4 text-gray-400 text-sm">
            <span>❌ Только виртуальные деньги</span>
            <span>✅ Никаких реальных ставок</span>
          </div>
        </div>
      </div>
    </footer>
  );
};

export default Footer;
