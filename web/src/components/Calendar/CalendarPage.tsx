'use client';

import { useState, useEffect } from 'react';
import { 
  ChevronLeft, 
  ChevronRight, 
  Calendar as CalendarIcon, 
  Plus, 
  TrendingUp, 
  TrendingDown, 
  Target,
  DollarSign,
  BarChart3,
  Filter,
  Search
} from 'lucide-react';

interface CalendarEvent {
  id: string;
  title: string;
  type: 'trade' | 'strategy' | 'meeting' | 'reminder';
  date: Date;
  time?: string;
  description?: string;
  pnl?: number;
  status?: 'completed' | 'pending' | 'cancelled';
}

interface CalendarDay {
  date: Date;
  isCurrentMonth: boolean;
  isToday: boolean;
  isSelected: boolean;
  events: CalendarEvent[];
}

export default function CalendarPage() {
  const [currentDate, setCurrentDate] = useState(new Date());
  const [selectedDate, setSelectedDate] = useState<Date | null>(null);
  const [viewMode, setViewMode] = useState<'month' | 'week' | 'day'>('month');
  const [events, setEvents] = useState<CalendarEvent[]>([]);

  // Sample events for demonstration
  useEffect(() => {
    const sampleEvents: CalendarEvent[] = [
      {
        id: '1',
        title: 'RELIANCE Trade',
        type: 'trade',
        date: new Date(2024, 0, 15),
        time: '09:30',
        pnl: 2500,
        status: 'completed'
      },
      {
        id: '2',
        title: 'NIFTY Strategy Review',
        type: 'strategy',
        date: new Date(2024, 0, 18),
        time: '14:00',
        status: 'pending'
      },
      {
        id: '3',
        title: 'Market Analysis Meeting',
        type: 'meeting',
        date: new Date(2024, 0, 22),
        time: '10:00',
        status: 'pending'
      },
      {
        id: '4',
        title: 'Stop Loss Reminder',
        type: 'reminder',
        date: new Date(2024, 0, 25),
        time: '15:30',
        status: 'pending'
      }
    ];
    setEvents(sampleEvents);
  }, []);

  const getDaysInMonth = (date: Date) => {
    const year = date.getFullYear();
    const month = date.getMonth();
    const firstDay = new Date(year, month, 1);
    const lastDay = new Date(year, month + 1, 0);
    const daysInMonth = lastDay.getDate();
    const startingDayOfWeek = firstDay.getDay();

    const days: CalendarDay[] = [];

    // Previous month's trailing days
    const prevMonth = new Date(year, month - 1, 0);
    for (let i = startingDayOfWeek - 1; i >= 0; i--) {
      const dayDate = new Date(year, month - 1, prevMonth.getDate() - i);
      days.push({
        date: dayDate,
        isCurrentMonth: false,
        isToday: false,
        isSelected: false,
        events: events.filter(event => 
          event.date.toDateString() === dayDate.toDateString()
        )
      });
    }

    // Current month's days
    for (let day = 1; day <= daysInMonth; day++) {
      const dayDate = new Date(year, month, day);
      const isToday = dayDate.toDateString() === new Date().toDateString();
      const isSelected = selectedDate?.toDateString() === dayDate.toDateString();
      
      days.push({
        date: dayDate,
        isCurrentMonth: true,
        isToday,
        isSelected,
        events: events.filter(event => 
          event.date.toDateString() === dayDate.toDateString()
        )
      });
    }

    // Next month's leading days
    const totalCells = 42; // 6 weeks * 7 days
    const remainingCells = totalCells - days.length;
    for (let day = 1; day <= remainingCells; day++) {
      const dayDate = new Date(year, month + 1, day);
      days.push({
        date: dayDate,
        isCurrentMonth: false,
        isToday: false,
        isSelected: false,
        events: events.filter(event => 
          event.date.toDateString() === dayDate.toDateString()
        )
      });
    }

    return days;
  };

  const navigateMonth = (direction: 'prev' | 'next') => {
    setCurrentDate(prev => {
      const newDate = new Date(prev);
      if (direction === 'prev') {
        newDate.setMonth(prev.getMonth() - 1);
      } else {
        newDate.setMonth(prev.getMonth() + 1);
      }
      return newDate;
    });
  };

  const goToToday = () => {
    setCurrentDate(new Date());
    setSelectedDate(new Date());
  };

  const handleDateClick = (date: Date) => {
    setSelectedDate(date);
  };

  const getEventTypeColor = (type: string) => {
    switch (type) {
      case 'trade':
        return 'bg-green-500/20 text-green-400 border-green-500/30';
      case 'strategy':
        return 'bg-blue-500/20 text-blue-400 border-blue-500/30';
      case 'meeting':
        return 'bg-purple-500/20 text-purple-400 border-purple-500/30';
      case 'reminder':
        return 'bg-orange-500/20 text-orange-400 border-orange-500/30';
      default:
        return 'bg-gray-500/20 text-gray-400 border-gray-500/30';
    }
  };

  const getEventIcon = (type: string) => {
    switch (type) {
      case 'trade':
        return <TrendingUp className="w-3 h-3" />;
      case 'strategy':
        return <Target className="w-3 h-3" />;
      case 'meeting':
        return <BarChart3 className="w-3 h-3" />;
      case 'reminder':
        return <CalendarIcon className="w-3 h-3" />;
      default:
        return <CalendarIcon className="w-3 h-3" />;
    }
  };

  const formatDate = (date: Date) => {
    return date.toLocaleDateString('en-IN', {
      year: 'numeric',
      month: 'long',
      day: 'numeric'
    });
  };

  const days = getDaysInMonth(currentDate);
  const weekDays = ['Sun', 'Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat'];

  return (
    <div className="p-6">
      

      {/* Calendar Navigation */}
      <div className="bg-primary border border-primary rounded-lg p-4 mb-6">
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-4">
            <button
              onClick={() => navigateMonth('prev')}
              className="p-2 rounded-lg bg-primary hover:bg-tertiary text-secondary transition-colors duration-200"
            >
              <ChevronLeft className="w-4 h-4" />
            </button>
            
            <h2 className="text-xl font-helvetica-bold text-primary">
              {currentDate.toLocaleDateString('en-IN', { 
                month: 'long', 
                year: 'numeric' 
              })}
            </h2>
            
            <button
              onClick={() => navigateMonth('next')}
              className="p-2 rounded-lg bg-primary hover:bg-tertiary text-secondary transition-colors duration-200"
            >
              <ChevronRight className="w-4 h-4" />
            </button>
          </div>

          <div className="flex items-center gap-2">
            <button
              onClick={goToToday}
              className="bg-accent hover:bg-primary text-primary font-helvetica-medium px-3 py-2 rounded-lg transition-colors duration-200"
            >
              Today
            </button>
            
            <div className="flex bg-primary rounded-lg p-1">
              {(['month', 'week', 'day'] as const).map((mode) => (
                <button
                  key={mode}
                  onClick={() => setViewMode(mode)}
                  className={`px-3 py-1 rounded text-sm font-helvetica-medium transition-colors duration-200 ${
                    viewMode === mode
                      ? 'bg-accent text-primary'
                      : 'text-tertiary hover:text-primary'
                  }`}
                >
                  {mode.charAt(0).toUpperCase() + mode.slice(1)}
                </button>
              ))}
            </div>
          </div>
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-4 gap-6">
        {/* Calendar Grid */}
        <div className="lg:col-span-3">
          <div className="bg-primary border border-primary rounded-lg overflow-hidden">
            {/* Week Days Header */}
            <div className="grid grid-cols-7 border-b border-primary">
              {weekDays.map((day) => (
                <div
                  key={day}
                  className="p-3 text-center text-tertiary font-helvetica-medium text-sm border-r border-primary last:border-r-0"
                >
                  {day}
                </div>
              ))}
            </div>

            {/* Calendar Days */}
            <div className="grid grid-cols-7">
              {days.map((day, index) => (
                <div
                  key={index}
                  onClick={() => handleDateClick(day.date)}
                  className={`
                    min-h-[120px] p-2 border-r border-b border-primary last:border-r-0 cursor-pointer transition-colors duration-200
                    ${day.isCurrentMonth ? 'bg-primary' : 'bg-primary'}
                    ${day.isToday ? 'bg-accent/10' : ''}
                    ${day.isSelected ? 'bg-accent/20' : ''}
                    hover:bg-tertiary
                  `}
                >
                  <div className="flex items-center justify-between mb-2">
                    <span
                      className={`
                        text-sm font-helvetica-medium
                        ${day.isCurrentMonth ? 'text-primary' : 'text-tertiary'}
                        ${day.isToday ? 'text-accent font-bold' : ''}
                        ${day.isSelected ? 'text-accent' : ''}
                      `}
                    >
                      {day.date.getDate()}
                    </span>
                    {day.events.length > 0 && (
                      <span className="text-xs bg-accent text-primary px-1.5 py-0.5 rounded-full">
                        {day.events.length}
                      </span>
                    )}
                  </div>

                  {/* Events */}
                  <div className="space-y-1">
                    {day.events.slice(0, 3).map((event) => (
                      <div
                        key={event.id}
                        className={`
                          text-xs px-2 py-1 rounded border flex items-center gap-1
                          ${getEventTypeColor(event.type)}
                        `}
                      >
                        {getEventIcon(event.type)}
                        <span className="truncate">{event.title}</span>
                      </div>
                    ))}
                    {day.events.length > 3 && (
                      <div className="text-xs text-tertiary px-2">
                        +{day.events.length - 3} more
                      </div>
                    )}
                  </div>
                </div>
              ))}
            </div>
          </div>
        </div>

        {/* Sidebar */}
        <div className="space-y-6">
          {/* Selected Date Info */}
          {selectedDate && (
            <div className="bg-primary border border-primary rounded-lg p-4">
              <h3 className="text-lg font-helvetica-bold text-primary mb-3">
                {formatDate(selectedDate)}
              </h3>
              <div className="space-y-2">
                {events
                  .filter(event => event.date.toDateString() === selectedDate.toDateString())
                  .map((event) => (
                    <div
                      key={event.id}
                      className="p-3 rounded-lg border border-primary bg-primary"
                    >
                      <div className="flex items-center gap-2 mb-1">
                        {getEventIcon(event.type)}
                        <span className="font-helvetica-medium text-primary text-sm">
                          {event.title}
                        </span>
                      </div>
                      {event.time && (
                        <p className="text-tertiary text-xs">{event.time}</p>
                      )}
                      {event.pnl && (
                        <p className={`text-sm font-helvetica-medium ${
                          event.pnl > 0 ? 'text-green-400' : 'text-red-400'
                        }`}>
                          P&L: ₹{event.pnl.toLocaleString()}
                        </p>
                      )}
                    </div>
                  ))}
                {events.filter(event => event.date.toDateString() === selectedDate.toDateString()).length === 0 && (
                  <p className="text-tertiary text-sm">No events for this date</p>
                )}
              </div>
            </div>
          )}

          {/* Quick Stats */}
          <div className="bg-primary border border-primary rounded-lg p-4">
            <h3 className="text-lg font-helvetica-bold text-primary mb-3">This Month</h3>
            <div className="space-y-3">
              <div className="flex items-center justify-between">
                <span className="text-tertiary text-sm">Total Trades</span>
                <span className="text-primary font-helvetica-medium">12</span>
              </div>
              <div className="flex items-center justify-between">
                <span className="text-tertiary text-sm">P&L</span>
                <span className="text-green-400 font-helvetica-medium">+₹15,420</span>
              </div>
              <div className="flex items-center justify-between">
                <span className="text-tertiary text-sm">Win Rate</span>
                <span className="text-primary font-helvetica-medium">75%</span>
              </div>
              <div className="flex items-center justify-between">
                <span className="text-tertiary text-sm">Strategies</span>
                <span className="text-primary font-helvetica-medium">3</span>
              </div>
            </div>
          </div>

          {/* Event Types Legend */}
          <div className="bg-primary border border-primary rounded-lg p-4">
            <h3 className="text-lg font-helvetica-bold text-primary mb-3">Event Types</h3>
            <div className="space-y-2">
              {[
                { type: 'trade', label: 'Trades', icon: TrendingUp },
                { type: 'strategy', label: 'Strategies', icon: Target },
                { type: 'meeting', label: 'Meetings', icon: BarChart3 },
                { type: 'reminder', label: 'Reminders', icon: CalendarIcon },
              ].map(({ type, label, icon: Icon }) => (
                <div key={type} className="flex items-center gap-2">
                  <div className={`p-1 rounded ${getEventTypeColor(type)}`}>
                    <Icon className="w-3 h-3" />
                  </div>
                  <span className="text-primary text-sm font-helvetica-light">{label}</span>
                </div>
              ))}
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
