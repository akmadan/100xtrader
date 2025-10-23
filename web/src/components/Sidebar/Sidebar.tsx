"use client";

import { useState, useEffect } from "react";
import {
  LayoutDashboard,
  Calendar,
  FolderOpen,
  BarChart3,
  Lightbulb,
  GraduationCap,
  NotebookPen,
  BookOpen,
  Target,
  RotateCcw,
  User,
  Palette,
  Moon,
  Sun,
  Brain,
  CreditCard,
  ChevronDown,
  ChevronRight,
  Bot,
  Cpu,
  BookOpenCheck,
  Settings,
} from "lucide-react";

interface SidebarProps {
  isOpen: boolean;
  onToggle: () => void;
  onPageChange: (page: string) => void;
  activePage: string;
}

export default function Sidebar({
  isOpen,
  onToggle,
  onPageChange,
  activePage,
}: SidebarProps) {
  const [isDark, setIsDark] = useState(false);
  const [expandedSections, setExpandedSections] = useState({
    journal: true,
    intelligence: false,
  });

  useEffect(() => {
    // Check if dark mode is enabled on initial load
    const isDarkMode = document.documentElement.classList.contains("dark");
    setIsDark(isDarkMode);
  }, []);

  const toggleTheme = () => {
    const newTheme = !isDark;
    setIsDark(newTheme);

    if (newTheme) {
      document.documentElement.classList.add("dark");
      localStorage.setItem("theme", "dark");
    } else {
      document.documentElement.classList.remove("dark");
      localStorage.setItem("theme", "light");
    }
  };

  const toggleSection = (section: string) => {
    setExpandedSections((prev) => ({
      ...prev,
      [section]: !prev[section as keyof typeof prev],
    }));
  };

  const navigationSections = [
    {
      id: "journal",
      title: "Journal",
      icon: BookOpenCheck,
      items: [
        { id: "Trades", label: "Trades", icon: FolderOpen, count: null },
        { id: "Configure", label: "Configure", icon: Settings, count: null },
        { id: "Calendar", label: "Calendar", icon: Calendar, count: null },
        { id: "Reports", label: "Reports", icon: BarChart3, count: null },
        { id: "Mentor Mode", label: "Mentor Mode", icon: User, count: null },
        { id: "Knowledge", label: "Knowledge", icon: Brain, count: null },
      ],
    },
    {
      id: "intelligence",
      title: "Intelligence",
      icon: Cpu,
      items: [
        { id: "Agents", label: "Agents", icon: Bot, count: null },
        { id: "Algorithms", label: "Algorithms", icon: Cpu, count: null },
      ],
    },
  ];

  const standaloneItems = [
    { id: "Dashboard", label: "Dashboard", icon: LayoutDashboard, count: null },

    { id: "Accounts", label: "Accounts", icon: CreditCard, count: null },
  ];

  return (
    <>
      {/* Overlay for mobile */}
      {isOpen && (
        <div
          className="fixed inset-0 bg-black bg-opacity-50 z-40 lg:hidden sidebar-overlay"
          onClick={onToggle}
        />
      )}

      {/* Sidebar */}
      <div
        className={`
        fixed top-0 left-0 h-full bg-primary border-r border-primary shadow-theme-lg z-50
        sidebar-transition
        ${isOpen ? "translate-x-0" : "-translate-x-full"}
        lg:translate-x-0 lg:static lg:z-auto
        w-64
      `}
      >
        <div className="flex flex-col h-full">
          {/* Header */}
          <div className="flex items-center justify-start py-4">
            <img src="/logo.svg" alt="100x Logo" className="w-120 h-10" />
          </div>

          {/* Add Trade Button */}
          <div className="p-4">
            <button className="w-full bg-secondary hover:bg-accent text-primary font-helvetica-medium py-3 px-4 rounded-lg transition-colors duration-200">
              + Add Trade
            </button>
          </div>

          {/* Navigation */}
          <nav className="flex-1 px-4 py-2 sidebar-scroll overflow-y-auto">
            <ul className="space-y-1">
              {/* Standalone Items */}
              {standaloneItems.map((item) => {
                const IconComponent = item.icon;
                return (
                  <li key={item.id}>
                    <button
                      onClick={() => {
                        onPageChange(item.id);
                      }}
                      className={`
                        w-full flex items-center justify-between px-3 py-2 rounded-lg transition-colors duration-200
                        ${
                          activePage === item.id
                            ? "bg-tertiary text-primary"
                            : "text-secondary hover:bg-secondary hover:text-primary"
                        }
                      `}
                    >
                      <div className="flex items-center space-x-3">
                        <IconComponent className="w-5 h-5" />
                        <span className="font-helvetica">{item.label}</span>
                      </div>
                      {item.count && (
                        <span className="text-xs bg-accent text-primary px-2 py-1 rounded">
                          {item.count}
                        </span>
                      )}
                    </button>
                  </li>
                );
              })}

              {/* Grouped Sections */}
              {navigationSections.map((section) => {
                const SectionIcon = section.icon;
                const isExpanded =
                  expandedSections[section.id as keyof typeof expandedSections];

                return (
                  <li key={section.id}>
                    {/* Section Header */}
                    <button
                      onClick={() => toggleSection(section.id)}
                      className="w-full flex items-center justify-between px-3 py-2 rounded-lg transition-colors duration-200 text-secondary hover:bg-secondary hover:text-primary"
                    >
                      <div className="flex items-center space-x-3">
                        <SectionIcon className="w-5 h-5" />
                        <span className="font-helvetica-medium">
                          {section.title}
                        </span>
                      </div>
                      {isExpanded ? (
                        <ChevronDown className="w-4 h-4" />
                      ) : (
                        <ChevronRight className="w-4 h-4" />
                      )}
                    </button>

                    {/* Section Items */}
                    {isExpanded && (
                      <ul className="ml-4 mt-1 space-y-1">
                        {section.items.map((item) => {
                          const IconComponent = item.icon;
                          return (
                            <li key={item.id}>
                              <button
                                onClick={() => {
                                  onPageChange(item.id);
                                }}
                                className={`
                                  w-full flex items-center justify-between px-3 py-2 rounded-lg transition-colors duration-200
                                  ${
                                    activePage === item.id
                                      ? "bg-tertiary text-primary"
                                      : "text-tertiary hover:bg-secondary hover:text-primary"
                                  }
                                `}
                              >
                                <div className="flex items-center space-x-3">
                                  <IconComponent className="w-4 h-4" />
                                  <span className="font-helvetica text-sm">
                                    {item.label}
                                  </span>
                                </div>
                                {item.count && (
                                  <span className="text-xs bg-accent text-primary px-2 py-1 rounded">
                                    {item.count}
                                  </span>
                                )}
                              </button>
                            </li>
                          );
                        })}
                      </ul>
                    )}
                  </li>
                );
              })}
            </ul>
          </nav>

          {/* Theme Toggle */}
          <div className="p-4 border-t border-primary">
            <button
              onClick={toggleTheme}
              className="w-full flex items-center justify-center space-x-2 p-3 rounded-lg bg-secondary hover:bg-tertiary text-primary transition-colors duration-200"
            >
              {isDark ? (
                <>
                  <Sun className="w-4 h-4" />
                  <span className="font-helvetica text-sm">Light Mode</span>
                </>
              ) : (
                <>
                  <Moon className="w-4 h-4" />
                  <span className="font-helvetica text-sm">Dark Mode</span>
                </>
              )}
            </button>
          </div>

          {/* User Profile */}
          <div className="p-4 border-t border-primary">
            <div className="flex items-center space-x-3">
              <div className="w-10 h-10 bg-secondary rounded-full flex items-center justify-center">
                <span className="text-primary font-helvetica-bold text-sm">
                  UA
                </span>
              </div>
              <div className="flex-1 min-w-0">
                <p className="text-primary font-helvetica-medium text-sm truncate">
                  Umar Ashraf
                </p>
                <p className="text-tertiary font-helvetica text-xs truncate">
                  umarashraf0128@gmail.com
                </p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </>
  );
}
