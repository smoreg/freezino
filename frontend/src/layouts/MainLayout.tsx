import { Outlet } from 'react-router-dom';
import Header from '../components/layout/Header';
import Sidebar from '../components/layout/Sidebar';
import MobileMenu from '../components/layout/MobileMenu';
import Footer from './Footer';
import { useMobileMenu } from '../hooks/useMobileMenu';

const MainLayout = () => {
  const { isOpen, toggleMenu, closeMenu } = useMobileMenu();

  return (
    <div className="min-h-screen flex flex-col bg-dark">
      <Header onMenuClick={toggleMenu} />

      <div className="flex flex-1">
        <Sidebar />

        <main className="flex-1 p-3 md:p-4 lg:p-6 overflow-y-auto">
          <Outlet />
        </main>
      </div>

      <Footer />

      {/* Mobile Menu */}
      <MobileMenu isOpen={isOpen} onClose={closeMenu} />
    </div>
  );
};

export default MainLayout;
