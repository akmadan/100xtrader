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
import { ITrade, TradeDirection } from '@/types';
import { tradeApi } from '@/services/api';

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
  const [trades, setTrades] = useState<ITrade[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  
  // TODO: Get userId from authentication context
  const userId = 1;

  // Calculate P&L for a trade
  const calculatePnL = (trade: ITrade): number => {
    if (!trade.entryPrice || !trade.exitPrice || !trade.quantity) {
      return 0;
    }
    const priceDiff = trade.direction === TradeDirection.LONG
      ? trade.exitPrice - trade.entryPrice
      : trade.entryPrice - trade.exitPrice;
    return Number((priceDiff * trade.quantity).toFixed(2));
  };

  // Fetch trades on component mount and when month changes
  useEffect(() => {
    fetchTrades();
  }, []);

  const fetchTrades = async () => {
    setLoading(true);
    setError(null);
    try {
      const response = await tradeApi.getAll(userId);
      
      // Transform API response to match TypeScript interface
      const transformedTrades: ITrade[] = ((response && response.trades) || []).map((t) => ({
        id: t.id,
        userId: t.user_id.toString(),
        symbol: t.symbol,
        marketType: t.market_type as any,
        entryDate: new Date(t.entry_date),
        entryPrice: t.entry_price,
        quantity: t.quantity,
        totalAmount: t.total_amount,
        exitPrice: t.exit_price,
        direction: t.direction as any,
        stopLoss: t.stop_loss,
        target: t.target,
        strategy: t.strategy,
        outcomeSummary: t.outcome_summary as any,
        tradeAnalysis: t.trade_analysis,
        rulesFollowed: t.rules_followed,
        screenshots: t.screenshots,
        psychology: t.psychology ? {
          entryConfidence: t.psychology.entry_confidence,
          satisfactionRating: t.psychology.satisfaction_rating,
          emotionalState: t.psychology.emotional_state,
          mistakesMade: t.psychology.mistakes_made || [],
          lessonsLearned: t.psychology.lessons_learned,
        } : null,
        // Broker-specific fields
        tradingBroker: t.trading_broker as any,
        traderBrokerId: t.trader_broker_id,
        exchangeOrderId: t.exchange_order_id,
        orderId: t.order_id,
        productType: t.product_type as any,
        transactionType: t.transaction_type,
        createdAt: new Date(t.created_at),
        updatedAt: new Date(t.updated_at),
      }));
      
      setTrades(transformedTrades);
      
      // Convert trades to calendar events
      const tradeEvents: CalendarEvent[] = transformedTrades.map((trade) => {
        // Calculate P&L
        const pnl = calculatePnL(trade);
        
        return {
          id: trade.id,
          title: `${trade.symbol} ${trade.direction.toUpperCase()}`,
          type: 'trade' as const,
          date: trade.entryDate,
          time: trade.entryDate.toLocaleTimeString('en-IN', { 
            hour: '2-digit', 
            minute: '2-digit',
            hour12: false 
          }),
          description: trade.tradeAnalysis || `${trade.strategy} - ${trade.symbol}`,
          pnl: pnl,
          status: 'completed' as const
        };
      });
      
      setEvents(tradeEvents);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to fetch trades");
      console.error("Error fetching trades:", err);
    } finally {
      setLoading(false);
    }
  };

  // Calculate monthly stats
  const getMonthlyStats = () => {
    const year = currentDate.getFullYear();
    const month = currentDate.getMonth();
    
    const monthTrades = trades.filter(trade => {
      const tradeDate = trade.entryDate;
      return tradeDate.getFullYear() === year && tradeDate.getMonth() === month;
    });

    const totalTrades = monthTrades.length;
    const totalPnL = monthTrades.reduce((sum, trade) => sum + calculatePnL(trade), 0);
    const winningTrades = monthTrades.filter(trade => calculatePnL(trade) > 0).length;
    const winRate = totalTrades > 0 ? Math.round((winningTrades / totalTrades) * 100) : 0;
    
    // Get unique strategies for the month
    const uniqueStrategies = new Set(monthTrades.map(trade => trade.strategy).filter(Boolean));
    
    return {
      totalTrades,
      totalPnL,
      winRate,
      strategies: uniqueStrategies.size
    };
  };

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

  const formatCurrency = (amount: number) => {
    return new Intl.NumberFormat("en-IN", {
      style: "currency",
      currency: "INR",
    }).format(amount);
  };

  const days = getDaysInMonth(currentDate);
  const weekDays = ['Sun', 'Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat'];
  const monthlyStats = getMonthlyStats();

  return (
    <div className="p-6">
      {/* Error Message */}
      {error && (
        <div className="mb-4 p-4 bg-red-500/20 border border-red-500/30 rounded-lg text-red-400">
          <p className="font-helvetica-medium">{error}</p>
        </div>
      )}

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
            {loading ? (
              <div className="text-center py-4">
                <p className="text-tertiary text-sm">Loading...</p>
              </div>
            ) : (
              <div className="space-y-3">
                <div className="flex items-center justify-between">
                  <span className="text-tertiary text-sm">Total Trades</span>
                  <span className="text-primary font-helvetica-medium">{monthlyStats.totalTrades}</span>
                </div>
                <div className="flex items-center justify-between">
                  <span className="text-tertiary text-sm">P&L</span>
                  <span className={`font-helvetica-medium ${
                    monthlyStats.totalPnL > 0 ? 'text-green-400' : 
                    monthlyStats.totalPnL < 0 ? 'text-red-400' : 
                    'text-primary'
                  }`}>
                    {monthlyStats.totalPnL > 0 ? '+' : ''}{formatCurrency(monthlyStats.totalPnL)}
                  </span>
                </div>
                <div className="flex items-center justify-between">
                  <span className="text-tertiary text-sm">Win Rate</span>
                  <span className="text-primary font-helvetica-medium">{monthlyStats.winRate}%</span>
                </div>
                <div className="flex items-center justify-between">
                  <span className="text-tertiary text-sm">Strategies</span>
                  <span className="text-primary font-helvetica-medium">{monthlyStats.strategies}</span>
                </div>
              </div>
            )}
          </div>

          {/* Event Types Legend */}
          <div className="bg-primary border border-primary rounded-lg p-4">
            <h3 className="text-lg font-helvetica-bold text-primary mb-3">Event Types</h3>
            <div className="space-y-2">
              {[
                { type: 'trade', label: 'Trades', icon: TrendingUp },
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
