import { Suspense, lazy } from 'react';
import { Toaster } from 'react-hot-toast';
import { BrowserRouter, Routes, Route } from 'react-router-dom';

import BankruptcyPopup from './components/BankruptcyPopup';
import CookieConsent from './components/CookieConsent';
import ErrorBoundary from './components/ErrorBoundary';
import LoadingSpinner from './components/LoadingSpinner';
import OfflineDetector from './components/OfflineDetector';
import ProtectedRoute from './components/ProtectedRoute';

// Lazy load all pages for better performance and code splitting
const MainLayout = lazy(() => import('./layouts/MainLayout'));
const Home = lazy(() => import('./pages/Home'));
const LoginPage = lazy(() => import('./pages/LoginPage'));
const NotFound = lazy(() => import('./pages/NotFound'));
const ErrorPage = lazy(() => import('./pages/ErrorPage'));
const DashboardPage = lazy(() => import('./pages/DashboardPage'));
const GameHistoryPage = lazy(() => import('./pages/GameHistoryPage'));
const ShopPage = lazy(() => import('./pages/ShopPage'));
const ProfilePage = lazy(() => import('./pages/ProfilePage'));
const CreditPage = lazy(() => import('./pages/CreditPage'));
const TermsPage = lazy(() => import('./pages/legal/TermsPage'));
const PrivacyPage = lazy(() => import('./pages/legal/PrivacyPage'));
const CookiesPage = lazy(() => import('./pages/legal/CookiesPage'));
const ContactPage = lazy(() => import('./pages/ContactPage'));
const AboutPage = lazy(() => import('./pages/AboutPage'));
const SlotsPage = lazy(() => import('./pages/SlotsPage'));
const RoulettePage = lazy(() => import('./pages/RoulettePage'));
const WorkPage = lazy(() => import('./pages/WorkPage'));

function App() {
  return (
    <ErrorBoundary>
      <BrowserRouter>
        <Suspense fallback={<LoadingSpinner />}>
          <Routes>
            <Route path="/login" element={<LoginPage />} />

            {/* Public pages */}
            <Route path="/contact" element={<ContactPage />} />
            <Route path="/about" element={<AboutPage />} />

            {/* Error pages */}
            <Route path="/error" element={<ErrorPage />} />
            <Route path="/error/:statusCode" element={<ErrorPage />} />

            {/* Legal Pages - Public */}
            <Route path="/terms" element={<TermsPage />} />
            <Route path="/privacy" element={<PrivacyPage />} />
            <Route path="/cookies" element={<CookiesPage />} />

            <Route element={<ProtectedRoute />}>
              <Route element={<MainLayout />}>
                <Route path="/" element={<Home />} />
                <Route path="/dashboard" element={<DashboardPage />} />
                <Route path="/work" element={<WorkPage />} />
                <Route path="/history" element={<GameHistoryPage />} />
                <Route path="/shop" element={<ShopPage />} />
                <Route path="/profile" element={<ProfilePage />} />
                <Route path="/credit" element={<CreditPage />} />
              </Route>
              {/* Game pages without MainLayout (full screen) */}
              <Route path="/games/slots" element={<SlotsPage />} />
              <Route path="/games/roulette" element={<RoulettePage />} />
            </Route>

            <Route path="*" element={<NotFound />} />
          </Routes>
        </Suspense>

        {/* Cookie Consent Banner - shown on first visit */}
        <CookieConsent />

        {/* Offline Detection */}
        <OfflineDetector />

        {/* Bankruptcy Popup */}
        <BankruptcyPopup />

        {/* Toast Notifications */}
        <Toaster
          position="top-right"
          toastOptions={{
            duration: 4000,
            style: {
              background: '#1F2937',
              color: '#fff',
              borderRadius: '8px',
              border: '1px solid #374151',
            },
            success: {
              iconTheme: {
                primary: '#10B981',
                secondary: '#fff',
              },
            },
            error: {
              iconTheme: {
                primary: '#DC2626',
                secondary: '#fff',
              },
            },
          }}
        />
      </BrowserRouter>
    </ErrorBoundary>
  );
}

export default App;
