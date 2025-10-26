"use client";

import { useState } from "react";
import {
  Plus,
  TrendingUp,
  TrendingDown,
  DollarSign,
  Calendar,
  Target,
} from "lucide-react";
import { ITrade, ITradeFormData, TradeDirection } from "@/types";
import { AddTradeModal } from "@/components/Trade/AddTradeModal";

export default function TradesPage() {
  const [trades, setTrades] = useState<ITrade[]>([]);
  const [isAddModalOpen, setIsAddModalOpen] = useState(false);

  const handleAddTrade = (data: ITradeFormData) => {
    const newTrade: ITrade = {
      id: Date.now().toString(),
      userId: "current-user",
      ...data,
      entryDate: new Date(data.entryDate),
      createdAt: new Date(),
      updatedAt: new Date(),
    };

    setTrades((prev) => [newTrade, ...prev]);
    setIsAddModalOpen(false);
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

  return (
    <div className="p-6">
      {/* Header */}
      <div className="flex items-center justify-between mb-6">
        <div>
          <h1 className="text-2xl font-helvetica-bold text-primary">Trades</h1>
          <p className="text-secondary font-helvetica-light">
            Track your trading performance and psychology
          </p>
        </div>
        <button
          onClick={() => setIsAddModalOpen(true)}
          className="bg-accent hover:bg-primary text-primary font-helvetica-medium px-4 py-2 rounded-lg transition-colors duration-200 focus:outline-none focus:ring-2 focus:ring-accent flex items-center gap-2"
        >
          <Plus className="w-4 h-4" />
          Add Trade
        </button>
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
                          (sum, trade) => sum + (trade.exitPrice || 0),
                          0
                        )
                      )
                    : "text-tertiary"
                }`}
              >
                {trades.length > 0
                  ? formatCurrency(
                      trades.reduce(
                        (sum, trade) => sum + (trade.exitPrice || 0),
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
                      (trades.filter((trade) => (trade.exitPrice || 0) > 0)
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
        {trades.length === 0 ? (
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
          <div className="divide-y divide-primary">
            {trades.map((trade) => (
              <div
                key={trade.id}
                className="p-4 hover:bg-tertiary transition-colors duration-200"
              >
                <div className="flex items-center justify-between">
                  <div className="flex-1">
                    <div className="flex items-center gap-4 mb-2">
                      <h3 className="font-helvetica-bold text-primary text-lg">
                        {trade.symbol}
                      </h3>
                      <span
                        className={`px-2 py-1 rounded text-xs font-helvetica-medium ${
                          trade.direction === TradeDirection.LONG
                            ? "bg-green-500/20 text-green-400 border border-green-500/30"
                            : "bg-red-500/20 text-red-400 border border-red-500/30"
                        }`}
                      >
                        {trade.direction.toUpperCase()}
                      </span>
                      <span className="text-tertiary font-helvetica-light text-sm">
                        {formatDate(trade.entryDate)}
                      </span>
                    </div>

                    <div className="grid grid-cols-2 md:grid-cols-4 gap-4 text-sm">
                      <div>
                        <p className="text-tertiary font-helvetica-light">
                          Entry Price
                        </p>
                        <p className="text-primary font-helvetica-medium">
                          ₹{trade.entryPrice}
                        </p>
                      </div>
                      <div>
                        <p className="text-tertiary font-helvetica-light">
                          Quantity
                        </p>
                        <p className="text-primary font-helvetica-medium">
                          {trade.quantity}
                        </p>
                      </div>
                      <div>
                        <p className="text-tertiary font-helvetica-light">
                          P&L
                        </p>
                        <div
                          className={`flex items-center gap-1 ${getPnLColor(
                            trade.exitPrice || 0
                          )}`}
                        >
                          {getPnLIcon(trade.exitPrice || 0)}
                          <span className="font-helvetica-medium">
                            {trade.exitPrice
                              ? formatCurrency(trade.exitPrice)
                              : "N/A"}
                          </span>
                        </div>
                      </div>
                      <div>
                        <p className="text-tertiary font-helvetica-light">
                          Confidence
                        </p>
                        <p className="text-primary font-helvetica-medium">
                          {trade.psychology?.entryConfidence}/10
                        </p>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>

      {/* Add Trade Modal */}
      <AddTradeModal
        isOpen={isAddModalOpen}
        onClose={() => setIsAddModalOpen(false)}
        onSubmit={handleAddTrade}
      />
    </div>
  );
}
