import { Outlet } from 'react-router-dom';
import Header from './Header';
import Sidebar from './Sidebar';
import Footer from './Footer';

const MainLayout = () => {
  return (
    <div className="min-h-screen flex flex-col bg-dark">
      <Header />

      <div className="flex flex-1">
        <Sidebar />

        <main className="flex-1 p-6">
          <Outlet />
        </main>
      </div>

      <Footer />
    </div>
  );
};

export default MainLayout;
