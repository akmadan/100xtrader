"use client";

import { useState, useEffect } from "react";
import {
  Plus,
  TrendingUp,
  TrendingDown,
  DollarSign,
  Calendar,
  Target,
  RefreshCw,
} from "lucide-react";
import { ITrade, ITradeFormData, TradeDirection } from "@/types";
import { AddTradeModal } from "@/components/Trade/AddTradeModal";
import { tradeApi } from "@/services/api";

export default function TradesPage() {
  const [trades, setTrades] = useState<ITrade[]>([]);
  const [isAddModalOpen, setIsAddModalOpen] = useState(false);
  const [editingTrade, setEditingTrade] = useState<ITrade | null>(null);
  const [loading, setLoading] = useState(false);
  const [syncing, setSyncing] = useState(false);
  const [error, setError] = useState<string | null>(null);
  
  // TODO: Get userId from authentication context
  const userId = 1;

  // Fetch trades on component mount
  useEffect(() => {
    fetchTrades();
  }, []);

  const fetchTrades = async () => {
    setLoading(true);
    setError(null);
    try {
      const response = await tradeApi.getAll(userId);
      
      // Transform API response to match TypeScript interface
      // Use optional chaining and default to empty array to prevent null reference errors
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
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to fetch trades");
      console.error("Error fetching trades:", err);
    } finally {
      setLoading(false);
    }
  };

  const handleSyncDhan = async () => {
    setSyncing(true);
    setError(null);
    try {
      const result = await tradeApi.syncDhan(userId);
      // Refresh trades after sync
      await fetchTrades();
      // Show success message with date range
      const dateRange = result.date_range || `${result.from_date} to ${result.to_date}`;
      alert(`Sync completed!\nDate Range: ${dateRange}\nSaved: ${result.saved_count} new trades\nSkipped: ${result.skipped_count} existing trades\nErrors: ${result.error_count}\nTotal Fetched: ${result.total_fetched}`);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to sync trades from Dhan");
      console.error("Error syncing trades:", err);
    } finally {
      setSyncing(false);
    }
  };

  const handleEditTrade = async (data: ITradeFormData) => {
    if (!editingTrade) return;
    
    setError(null);
    try {
      const totalAmount = data.totalAmount || data.entryPrice * data.quantity;
      
      // Build psychology object only if it has valid values
      const psychology: any = {};
      if (data.psychology.entryConfidence > 0) {
        psychology.entry_confidence = data.psychology.entryConfidence;
      }
      if (data.psychology.satisfactionRating > 0) {
        psychology.satisfaction_rating = data.psychology.satisfactionRating;
      }
      if (data.psychology.emotionalState && data.psychology.emotionalState.trim() !== '') {
        psychology.emotional_state = data.psychology.emotionalState;
      }
      if (data.psychology.mistakesMade && data.psychology.mistakesMade.length > 0) {
        psychology.mistakes_made = data.psychology.mistakesMade;
      }
      if (data.psychology.lessonsLearned && data.psychology.lessonsLearned.trim() !== '') {
        psychology.lessons_learned = data.psychology.lessonsLearned;
      }
      
      const updatePayload: any = {
        symbol: data.symbol,
        market_type: data.marketType,
        entry_date: data.entryDate,
        entry_price: data.entryPrice,
        quantity: data.quantity,
        total_amount: totalAmount,
        exit_price: data.exitPrice,
        direction: data.direction,
        stop_loss: data.stopLoss,
        target: data.target,
        strategy: data.strategy,
        outcome_summary: data.outcomeSummary,
        trade_analysis: data.tradeAnalysis,
        rules_followed: data.rulesFollowed,
        screenshots: data.screenshots,
      };
      
      // Only include psychology if it has at least one valid field
      if (Object.keys(psychology).length > 0) {
        updatePayload.psychology = psychology;
      }
      
      const response = await tradeApi.update(editingTrade.id, userId, updatePayload);

      // Refresh trades after update
      await fetchTrades();
      setEditingTrade(null);
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : "Failed to update trade";
      setError(errorMessage);
      throw err;
    }
  };

  const handleAddTrade = async (data: ITradeFormData) => {
    setError(null);
    try {
      // Calculate total amount if not provided
      const totalAmount = data.totalAmount || data.entryPrice * data.quantity;
      
      const response = await tradeApi.create(userId, {
        symbol: data.symbol,
        market_type: data.marketType,
        entry_date: data.entryDate,
        entry_price: data.entryPrice,
        quantity: data.quantity,
        total_amount: totalAmount,
        exit_price: data.exitPrice,
        direction: data.direction,
        stop_loss: data.stopLoss,
        target: data.target,
        strategy: data.strategy,
        outcome_summary: data.outcomeSummary,
        trade_analysis: data.tradeAnalysis,
        rules_followed: data.rulesFollowed,
        screenshots: data.screenshots,
        psychology: {
          entry_confidence: data.psychology.entryConfidence,
          satisfaction_rating: data.psychology.satisfactionRating,
          emotional_state: data.psychology.emotionalState,
          mistakes_made: data.psychology.mistakesMade,
          lessons_learned: data.psychology.lessonsLearned,
        },
      });

      // Transform API response to match TypeScript interface
      const newTrade: ITrade = {
        id: response.id,
        userId: response.user_id.toString(),
        symbol: response.symbol,
        marketType: response.market_type as any,
        entryDate: new Date(response.entry_date),
        entryPrice: response.entry_price,
        quantity: response.quantity,
        totalAmount: response.total_amount,
        exitPrice: response.exit_price,
        direction: response.direction as any,
        stopLoss: response.stop_loss,
        target: response.target,
        strategy: response.strategy,
        outcomeSummary: response.outcome_summary as any,
        tradeAnalysis: response.trade_analysis,
        rulesFollowed: response.rules_followed,
        screenshots: response.screenshots,
        psychology: response.psychology ? {
          entryConfidence: response.psychology.entry_confidence,
          satisfactionRating: response.psychology.satisfaction_rating,
          emotionalState: response.psychology.emotional_state,
          mistakesMade: response.psychology.mistakes_made || [],
          lessonsLearned: response.psychology.lessons_learned,
        } : null,
        // Broker-specific fields
        tradingBroker: response.trading_broker as any,
        traderBrokerId: response.trader_broker_id,
        exchangeOrderId: response.exchange_order_id,
        orderId: response.order_id,
        productType: response.product_type as any,
        transactionType: response.transaction_type,
        createdAt: new Date(response.created_at),
        updatedAt: new Date(response.updated_at),
      };

      setTrades((prev) => [newTrade, ...prev]);
      setIsAddModalOpen(false);
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : "Failed to create trade";
      setError(errorMessage);
      throw err;
    }
  };

  const formatCurrency = (amount: number) => {
    return new Intl.NumberFormat("en-IN", {
      style: "currency",
      currency: "INR",
    }).format(amount);
  };

  const formatDate = (date: Date) => {
    return new Intl.DateTimeFormat("en-IN", {
      day: "2-digit",
      month: "2-digit",
      year: "numeric",
    }).format(date);
  };

  const getPnLColor = (pnl: number) => {
    if (pnl > 0) return "text-green-500";
    if (pnl < 0) return "text-red-500";
    return "text-gray-500";
  };

  const getPnLIcon = (pnl: number) => {
    if (pnl > 0) return <TrendingUp className="w-4 h-4" />;
    if (pnl < 0) return <TrendingDown className="w-4 h-4" />;
    return <Target className="w-4 h-4" />;
  };

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

  return (
    <div className="p-6">
      {/* Error Message */}
      {error && (
        <div className="mb-4 p-4 bg-red-500/20 border border-red-500/30 rounded-lg text-red-400">
          <p className="font-helvetica-medium">{error}</p>
        </div>
      )}

      {/* Header */}
      <div className="flex items-center justify-between mb-6">
        <div>
          <h1 className="text-2xl font-helvetica-bold text-primary">Trades</h1>
          <p className="text-secondary font-helvetica-light">
            Track your trading performance and psychology
          </p>
        </div>
        <div className="flex items-center gap-3">
          <button
            onClick={handleSyncDhan}
            disabled={syncing}
            className="bg-primary hover:bg-primary/80 border border-accent text-accent font-helvetica-medium px-4 py-2 rounded-lg transition-colors duration-200 focus:outline-none focus:ring-2 focus:ring-accent flex items-center gap-2 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            <RefreshCw className={`w-4 h-4 ${syncing ? 'animate-spin' : ''}`} />
            {syncing ? 'Syncing...' : 'Sync Dhan'}
          </button>
          <button
            onClick={() => setIsAddModalOpen(true)}
            className="bg-accent hover:bg-primary text-primary font-helvetica-medium px-4 py-2 rounded-lg transition-colors duration-200 focus:outline-none focus:ring-2 focus:ring-accent flex items-center gap-2"
          >
            <Plus className="w-4 h-4" />
            Add Trade
          </button>
        </div>
      </div>

      {/* Stats Cards */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-4 mb-6">
        <div className="bg-primary border border-primary rounded-lg p-4">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-tertiary font-helvetica-light text-sm">
                Total Trades
              </p>
              <p className="text-primary font-helvetica-bold text-2xl">
                {trades.length}
              </p>
            </div>
            <Calendar className="w-8 h-8 text-accent" />
          </div>
        </div>

        <div className="bg-primary border border-primary rounded-lg p-4">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-tertiary font-helvetica-light text-sm">
                Total P&L
              </p>
              <p
                className={`font-helvetica-bold text-2xl ${
                  trades.length > 0
                    ? getPnLColor(
                        trades.reduce(
                          (sum, trade) => sum + calculatePnL(trade),
                          0
                        )
                      )
                    : "text-tertiary"
                }`}
              >
                {trades.length > 0
                  ? formatCurrency(
                      trades.reduce(
                        (sum, trade) => sum + calculatePnL(trade),
                        0
                      )
                    )
                  : "₹0"}
              </p>
            </div>
            <DollarSign className="w-8 h-8 text-accent" />
          </div>
        </div>

        <div className="bg-primary border border-primary rounded-lg p-4">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-tertiary font-helvetica-light text-sm">
                Win Rate
              </p>
              <p className="text-primary font-helvetica-bold text-2xl">
                {trades.length > 0
                  ? `${Math.round(
                      (trades.filter((trade) => calculatePnL(trade) > 0)
                        .length /
                        trades.length) *
                        100
                    )}%`
                  : "0%"}
              </p>
            </div>
            <TrendingUp className="w-8 h-8 text-accent" />
          </div>
        </div>

        <div className="bg-primary border border-primary rounded-lg p-4">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-tertiary font-helvetica-light text-sm">
                Avg Confidence
              </p>
              <p className="text-primary font-helvetica-bold text-2xl">
                {trades.length > 0
                  ? `${(
                      trades.reduce(
                        (sum, trade) =>
                          sum + (trade.psychology?.entryConfidence || 0),
                        0
                      ) / trades.length
                    ).toFixed(1)}/10`
                  : "0/10"}
              </p>
            </div>
            <Target className="w-8 h-8 text-accent" />
          </div>
        </div>
      </div>

      {/* Trades List */}
      <div className="bg-primary border border-primary rounded-lg">
        {loading ? (
          <div className="p-8 text-center">
            <p className="text-tertiary font-helvetica">Loading trades...</p>
          </div>
        ) : trades.length === 0 ? (
          <div className="p-8 text-center">
            <TrendingUp className="w-12 h-12 text-tertiary mx-auto mb-4" />
            <h3 className="text-lg font-helvetica-medium text-primary mb-2">
              No trades yet
            </h3>
            <p className="text-tertiary font-helvetica-light mb-4">
              Start tracking your trades to analyze your performance
            </p>
            <button
              onClick={() => setIsAddModalOpen(true)}
              className="bg-accent hover:bg-primary text-primary font-helvetica-medium px-4 py-2 rounded-lg transition-colors duration-200 focus:outline-none focus:ring-2 focus:ring-accent"
            >
              Add Your First Trade
            </button>
          </div>
        ) : (
          <div className="overflow-x-auto">
            <table className="w-full">
              <thead>
                <tr className="border-b border-primary">
                  <th className="text-left py-3 px-4 text-sm font-helvetica-medium text-tertiary">Symbol</th>
                  <th className="text-left py-3 px-4 text-sm font-helvetica-medium text-tertiary">Direction</th>
                  <th className="text-left py-3 px-4 text-sm font-helvetica-medium text-tertiary">Date</th>
                  <th className="text-right py-3 px-4 text-sm font-helvetica-medium text-tertiary">Entry Price</th>
                  <th className="text-right py-3 px-4 text-sm font-helvetica-medium text-tertiary">Exit Price</th>
                  <th className="text-right py-3 px-4 text-sm font-helvetica-medium text-tertiary">Quantity</th>
                  <th className="text-right py-3 px-4 text-sm font-helvetica-medium text-tertiary">P&L</th>
                  <th className="text-center py-3 px-4 text-sm font-helvetica-medium text-tertiary">Confidence</th>
                  <th className="text-left py-3 px-4 text-sm font-helvetica-medium text-tertiary">Strategy</th>
                  <th className="text-left py-3 px-4 text-sm font-helvetica-medium text-tertiary">Outcome</th>
                </tr>
              </thead>
              <tbody>
                {trades.map((trade) => {
                  const pnl = calculatePnL(trade);
                  return (
                    <tr
                      key={trade.id}
                      onClick={() => setEditingTrade(trade)}
                      className="border-b border-primary hover:bg-tertiary transition-colors duration-200 cursor-pointer"
                    >
                      <td className="py-3 px-4">
                        <span className="font-helvetica-bold text-primary">
                          {trade.symbol}
                        </span>
                      </td>
                      <td className="py-3 px-4">
                        <span
                          className={`px-2 py-1 rounded text-xs font-helvetica-medium ${
                            trade.direction === TradeDirection.LONG
                              ? "bg-green-500/20 text-green-400 border border-green-500/30"
                              : "bg-red-500/20 text-red-400 border border-red-500/30"
                          }`}
                        >
                          {trade.direction.toUpperCase()}
                        </span>
                      </td>
                      <td className="py-3 px-4">
                        <span className="text-primary font-helvetica-light text-sm">
                          {formatDate(trade.entryDate)}
                        </span>
                      </td>
                      <td className="py-3 px-4 text-right">
                        <span className="text-primary font-helvetica-medium">
                          {formatCurrency(trade.entryPrice)}
                        </span>
                      </td>
                      <td className="py-3 px-4 text-right">
                        <span className="text-primary font-helvetica-medium">
                          {trade.exitPrice ? formatCurrency(trade.exitPrice) : "N/A"}
                        </span>
                      </td>
                      <td className="py-3 px-4 text-right">
                        <span className="text-primary font-helvetica-medium">
                          {trade.quantity}
                        </span>
                      </td>
                      <td className="py-3 px-4 text-right">
                        <div className={`flex items-center justify-end gap-1 ${getPnLColor(pnl)}`}>
                          {getPnLIcon(pnl)}
                          <span className="font-helvetica-medium">
                            {formatCurrency(pnl)}
                          </span>
                        </div>
                      </td>
                      <td className="py-3 px-4 text-center">
                        <span className="text-primary font-helvetica-medium">
                          {trade.psychology?.entryConfidence || 0}/10
                        </span>
                      </td>
                      <td className="py-3 px-4">
                        <span className="text-primary font-helvetica-light text-sm">
                          {trade.strategy || "N/A"}
                        </span>
                      </td>
                      <td className="py-3 px-4">
                        <span className="text-primary font-helvetica-light text-sm capitalize">
                          {trade.outcomeSummary?.replace(/_/g, " ") || "N/A"}
                        </span>
                      </td>
                    </tr>
                  );
                })}
              </tbody>
            </table>
          </div>
        )}
      </div>

      {/* Add Trade Modal */}
      <AddTradeModal
        isOpen={isAddModalOpen}
        onClose={() => setIsAddModalOpen(false)}
        onSubmit={handleAddTrade}
      />

      {/* Edit Trade Modal */}
      <AddTradeModal
        isOpen={editingTrade !== null}
        onClose={() => setEditingTrade(null)}
        onSubmit={handleEditTrade}
        initialData={editingTrade || undefined}
      />
    </div>
  );
}
