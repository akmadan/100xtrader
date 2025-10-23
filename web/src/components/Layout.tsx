'use client';

import { useState } from 'react';
import Sidebar from './Sidebar/Sidebar';
import { ColorDemo, ConfigurePage } from './index';
import TradesPage from './Trade/TradesPage';
import { CalendarPage } from './Calendar';

interface LayoutProps {
  children: React.ReactNode;
}

export default function Layout({ children }: LayoutProps) {
  const [sidebarOpen, setSidebarOpen] = useState(false);
  const [activePage, setActivePage] = useState('Dashboard');

  const toggleSidebar = () => {
    setSidebarOpen(!sidebarOpen);
  };

  const handlePageChange = (page: string) => {
    setActivePage(page);
  };

  const renderContent = () => {
    switch (activePage) {
      case 'Color Demo':
        return <ColorDemo />;
      case 'Configure':
        return <ConfigurePage />;
      case 'Trades':
        return <TradesPage />;
      case 'Calendar':
        return <CalendarPage />;
      default:
        return children;
    }
  };

  return (
    <div className="flex h-screen bg-primary">
      <Sidebar 
        isOpen={sidebarOpen} 
        onToggle={toggleSidebar} 
        onPageChange={handlePageChange}
        activePage={activePage}
      />
      
      {/* Main Content */}
      <div className="flex-1 flex flex-col overflow-hidden">
        {/* Top Bar */}
        <div className="bg-primary border-b border-primary px-4 py-3 flex items-center justify-between">
          <button
            onClick={toggleSidebar}
            className="lg:hidden p-2 rounded-lg bg-secondary hover:bg-tertiary text-primary transition-colors"
          >
            <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 6h16M4 12h16M4 18h16" />
            </svg>
          </button>
          
          <div className="flex items-center space-x-4">
            <h2 className="text-lg font-helvetica-medium text-primary">{activePage}</h2>
          </div>
        </div>

        {/* Content Area */}
        <main className="flex-1 overflow-auto bg-primary">
          {renderContent()}
        </main>
      </div>
    </div>
  );
}
