'use client';

import { useState, useEffect } from 'react';
import { X, Info, Brain, Clock, Calendar, ArrowUp, ArrowDown, Upload, Plus } from 'lucide-react';
import { ITradeFormData, ITradeModalProps, MarketType, TradeDirection, TradeDuration, OutcomeSummary } from '@/types';

const marketTypes: { value: MarketType; label: string }[] = [
  { value: MarketType.INDIAN, label: 'Indian' },
  { value: MarketType.US, label: 'US' },
  { value: MarketType.CRYPTO, label: 'Crypto' },
  { value: MarketType.FOREX, label: 'Forex' },
  { value: MarketType.COMMODITIES, label: 'Commodities' },
];

const outcomeSummaries: { value: OutcomeSummary; label: string }[] = [
  { value: OutcomeSummary.PROFITABLE, label: 'Profitable' },
  { value: OutcomeSummary.LOSS, label: 'Loss' },
  { value: OutcomeSummary.BREAKEVEN, label: 'Breakeven' },
  { value: OutcomeSummary.PARTIAL_PROFIT, label: 'Partial Profit' },
  { value: OutcomeSummary.PARTIAL_LOSS, label: 'Partial Loss' },
];

const emotionalStates = [
  'Confident', 'Anxious', 'Excited', 'Nervous', 'Calm', 'Frustrated', 
  'Optimistic', 'Pessimistic', 'Focused', 'Distracted', 'Greedy', 'Fearful'
];

const commonMistakes = [
  'Overtrading', 'Risked Too Much', 'Exited Too Late', 'Ignored Signals', 
  'Ignored Stop Loss', 'Revenge Trading', 'Exited Too Early', 'FOMO Entry', 
  'No Clear Plan', 'No Mistakes'
];

export function AddTradeModal({ isOpen, onClose, onSubmit, initialData }: ITradeModalProps) {
  const [activeTab, setActiveTab] = useState<'general' | 'psychology'>('general');
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [newRule, setNewRule] = useState('');

  // Initialize form data from initialData if provided (for editing)
  const getInitialFormData = (): ITradeFormData => {
    if (initialData) {
      return {
        symbol: initialData.symbol,
        marketType: initialData.marketType,
        entryDate: initialData.entryDate.toISOString().split('T')[0],
        entryPrice: initialData.entryPrice,
        quantity: initialData.quantity,
        totalAmount: initialData.totalAmount,
        exitPrice: initialData.exitPrice || 0,
        direction: initialData.direction,
        duration: TradeDuration.INTRADAY, // Default, not stored in ITrade
        strategy: initialData.strategy,
        outcomeSummary: initialData.outcomeSummary,
        tradeAnalysis: initialData.tradeAnalysis,
        rulesFollowed: initialData.rulesFollowed || [],
        screenshots: initialData.screenshots || [],
        psychology: initialData.psychology || {
          entryConfidence: 0,
          satisfactionRating: 0,
          emotionalState: '',
          mistakesMade: [],
          lessonsLearned: '',
        },
      };
    }
    return {
      symbol: '',
      marketType: MarketType.INDIAN,
      entryDate: '',
      entryPrice: 0,
      quantity: 0,
      totalAmount: 0,
      exitPrice: 0,
      direction: TradeDirection.LONG,
      duration: TradeDuration.INTRADAY,
      strategy: '',
      outcomeSummary: OutcomeSummary.PROFITABLE,
      tradeAnalysis: '',
      rulesFollowed: [],
      screenshots: [],
      psychology: {
        entryConfidence: 5,
        satisfactionRating: 5,
        emotionalState: '',
        mistakesMade: [],
        lessonsLearned: '',
      },
    };
  };

  const [formData, setFormData] = useState<ITradeFormData>(getInitialFormData());

  // Reset form when modal opens/closes or initialData changes
  useEffect(() => {
    if (isOpen) {
      setFormData(getInitialFormData());
      setActiveTab('general');
      setNewRule('');
    }
  }, [isOpen, initialData]);

  const handleInputChange = (field: string, value: string | number | TradeDirection | TradeDuration | OutcomeSummary | MarketType) => {
    setFormData(prev => ({
      ...prev,
      [field]: value,
    }));
  };

  // Calculate P&L based on entry price, exit price, quantity, and direction
  const calculatePnL = (): number => {
    if (!formData.entryPrice || !formData.exitPrice || !formData.quantity) {
      return 0;
    }
    const priceDiff = formData.direction === TradeDirection.LONG
      ? formData.exitPrice - formData.entryPrice
      : formData.entryPrice - formData.exitPrice;
    return Number((priceDiff * formData.quantity).toFixed(2));
  };

  // Calculate total amount based on entry price and quantity
  const calculateTotalAmount = (): number => {
    if (!formData.entryPrice || !formData.quantity) {
      return 0;
    }
    return Number((formData.entryPrice * formData.quantity).toFixed(2));
  };

  // Format number to 2 decimal places for display (except quantity which is integer)
  const formatNumber = (value: number | string | undefined, isInteger: boolean = false): string => {
    if (value === undefined || value === null || value === '') {
      return '';
    }
    const numValue = typeof value === 'string' ? parseFloat(value) : value;
    if (isNaN(numValue)) {
      return '';
    }
    return isInteger ? Math.floor(numValue).toString() : numValue.toFixed(2);
  };

  // Handle number input to prevent leading zeros
  const handleNumberInput = (field: string, value: string) => {
    // If empty, set to 0
    if (value === '' || value === null || value === undefined) {
      handleInputChange(field, 0);
      // Auto-calculate total amount if entry price or quantity changed
      if (field === 'entryPrice' || field === 'quantity') {
        const entryPrice = field === 'entryPrice' ? 0 : formData.entryPrice;
        const quantity = field === 'quantity' ? 0 : formData.quantity;
        handleInputChange('totalAmount', Number((entryPrice * quantity).toFixed(2)));
      }
      return;
    }
    
    // Remove leading zeros but preserve decimals
    // Handle cases like "0341" -> "341", "0100" -> "100", "0.5" -> "0.5", "00.5" -> "0.5"
    let cleanedValue = value.trim();
    
    // Remove leading zeros, but keep at least one digit before decimal point
    // This handles: "0341" -> "341", "0100" -> "100", "00.5" -> "0.5"
    if (cleanedValue.includes('.')) {
      // For decimal numbers, remove leading zeros before decimal point
      const parts = cleanedValue.split('.');
      parts[0] = parts[0].replace(/^0+/, '') || '0';
      cleanedValue = parts.join('.');
    } else {
      // For whole numbers, remove all leading zeros
      cleanedValue = cleanedValue.replace(/^0+/, '') || '0';
    }
    
    // Parse the number
    const numValue = cleanedValue === '' ? 0 : parseFloat(cleanedValue);
    
    // Check if field is integer type (quantity)
    if (field === 'quantity') {
      const intValue = Math.floor(Math.abs(numValue)) || 0;
      handleInputChange(field, intValue);
      // Auto-calculate total amount when quantity changes
      if (formData.entryPrice) {
        handleInputChange('totalAmount', Number((formData.entryPrice * intValue).toFixed(2)));
      }
    } else {
      const floatValue = isNaN(numValue) ? 0 : Number(numValue.toFixed(2));
      handleInputChange(field, floatValue);
      // Auto-calculate total amount when entry price changes
      if (field === 'entryPrice' && formData.quantity) {
        handleInputChange('totalAmount', Number((floatValue * formData.quantity).toFixed(2)));
      }
    }
  };

  // Handle blur event to ensure clean display
  const handleNumberBlur = (field: string, value: number) => {
    // Ensure the value is properly formatted (no leading zeros, 2 decimal places)
    if (value === 0) {
      handleInputChange(field, 0);
    } else {
      // Re-parse to ensure clean value with 2 decimal places
      const numValue = field === 'quantity' 
        ? Math.floor(value) 
        : Number(value.toFixed(2));
      handleInputChange(field, numValue);
      
      // Auto-calculate total amount if entry price or quantity changed
      if (field === 'entryPrice' && formData.quantity) {
        handleInputChange('totalAmount', Number((numValue * formData.quantity).toFixed(2)));
      } else if (field === 'quantity' && formData.entryPrice) {
        handleInputChange('totalAmount', Number((formData.entryPrice * numValue).toFixed(2)));
      }
    }
  };

  const handlePsychologyChange = (field: string, value: string | number) => {
    setFormData(prev => ({
      ...prev,
      psychology: {
        ...prev.psychology,
        [field]: value,
      },
    }));
  };

  const addRule = () => {
    if (newRule.trim() && !formData.rulesFollowed.includes(newRule.trim())) {
      setFormData(prev => ({
        ...prev,
        rulesFollowed: [...prev.rulesFollowed, newRule.trim()],
      }));
      setNewRule('');
    }
  };

  const removeRule = (rule: string) => {
    setFormData(prev => ({
      ...prev,
      rulesFollowed: prev.rulesFollowed.filter(r => r !== rule),
    }));
  };

  const toggleMistake = (mistake: string) => {
    setFormData(prev => ({
      ...prev,
      psychology: {
        ...prev.psychology,
        mistakesMade: prev.psychology.mistakesMade.includes(mistake)
          ? prev.psychology.mistakesMade.filter(m => m !== mistake)
          : [...prev.psychology.mistakesMade, mistake],
      },
    }));
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsSubmitting(true);
    
    try {
      await new Promise(resolve => setTimeout(resolve, 1000)); // Simulate API call
      onSubmit(formData);
    } catch (error) {
      console.error('Error adding trade:', error);
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleClose = () => {
    setFormData({
      symbol: '',
      marketType: MarketType.INDIAN,
      entryDate: '',
      entryPrice: 0,
      quantity: 0,
      totalAmount: 0,
      exitPrice: 0,
      direction: TradeDirection.LONG,
      duration: TradeDuration.INTRADAY,
      strategy: '',
      outcomeSummary: OutcomeSummary.PROFITABLE,
      tradeAnalysis: '',
      rulesFollowed: [],
      screenshots: [],
      psychology: {
        entryConfidence: 5,
        satisfactionRating: 5,
        emotionalState: '',
        mistakesMade: [],
        lessonsLearned: '',
      },
    });
    setActiveTab('general');
    onClose();
  };

  if (!isOpen) return null;

  return (
    <div 
      className="fixed inset-0 bg-black/30 backdrop-blur-sm flex items-center justify-center z-50"
      onClick={handleClose}
    >
      <div 
        className="bg-primary border border-primary rounded-lg shadow-theme-lg w-full max-w-4xl mx-4 modal-content max-h-[90vh] overflow-y-auto"
        onClick={(e) => e.stopPropagation()}
      >
        {/* Header */}
        <div className="flex items-center justify-between p-6 border-b border-primary">
          <h2 className="text-xl font-helvetica-bold text-primary">Add New Trade</h2>
          <button
            onClick={handleClose}
            className="text-tertiary hover:text-primary transition-colors duration-200"
          >
            <X className="w-5 h-5" />
          </button>
        </div>

        {/* Tabs */}
        <div className="flex border-b border-primary">
          <button
            onClick={() => setActiveTab('general')}
            className={`flex items-center gap-2 px-6 py-3 font-helvetica-medium transition-colors duration-200 ${
              activeTab === 'general'
                ? 'text-accent border-b-2 border-accent'
                : 'text-tertiary hover:text-primary'
            }`}
          >
            <Info className="w-4 h-4" />
            General
          </button>
          <button
            onClick={() => setActiveTab('psychology')}
            className={`flex items-center gap-2 px-6 py-3 font-helvetica-medium transition-colors duration-200 ${
              activeTab === 'psychology'
                ? 'text-accent border-b-2 border-accent'
                : 'text-tertiary hover:text-primary'
            }`}
          >
            <Brain className="w-4 h-4" />
            Psychology
          </button>
        </div>

        <form onSubmit={handleSubmit} className="p-6">
          {activeTab === 'general' && (
            <div className="space-y-6">
              {/* Trade Duration */}
              <div>
                <label className="block text-sm font-helvetica-medium text-primary mb-2">
                  Trade Duration
                </label>
                <p className="text-xs text-tertiary mb-3">Select whether this is an intraday or swing trade</p>
                <div className="flex gap-2">
                  <button
                    type="button"
                    onClick={() => handleInputChange('duration', TradeDuration.INTRADAY)}
                    className={`flex items-center gap-2 px-4 py-2 rounded-lg border transition-colors duration-200 ${
                      formData.duration === TradeDuration.INTRADAY
                        ? 'bg-secondary border-accent text-primary'
                        : 'bg-primary border-primary text-tertiary hover:text-primary'
                    }`}
                  >
                    <Clock className="w-4 h-4" />
                    Intraday
                  </button>
                  <button
                    type="button"
                    onClick={() => handleInputChange('duration', TradeDuration.SWING)}
                    className={`flex items-center gap-2 px-4 py-2 rounded-lg border transition-colors duration-200 ${
                      formData.duration === TradeDuration.SWING
                        ? 'bg-secondary border-accent text-primary'
                        : 'bg-primary border-primary text-tertiary hover:text-primary'
                    }`}
                  >
                    <Calendar className="w-4 h-4" />
                    Swing
                  </button>
                </div>
              </div>

              {/* Market Type and Symbol */}
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div>
                  <label className="block text-sm font-helvetica-medium text-primary mb-2">
                    Market type *
                  </label>
                  <select
                    value={formData.marketType}
                    onChange={(e) => handleInputChange('marketType', e.target.value)}
                    className="w-full bg-primary border border-primary text-primary placeholder-tertiary px-3 py-2 rounded-lg focus:outline-none focus:ring-2 focus:ring-accent"
                  >
                    {marketTypes.map((type) => (
                      <option key={type.value} value={type.value}>
                        {type.label}
                      </option>
                    ))}
                  </select>
                </div>
                <div>
                  <label className="block text-sm font-helvetica-medium text-primary mb-2">
                    Symbol *
                  </label>
                  <input
                    type="text"
                    value={formData.symbol}
                    onChange={(e) => handleInputChange('symbol', e.target.value)}
                    placeholder="RELIANCE, NIFTY 50, etc."
                    className="w-full bg-primary border border-primary text-primary placeholder-tertiary px-3 py-2 rounded-lg focus:outline-none focus:ring-2 focus:ring-accent"
                    required
                  />
                </div>
              </div>

              {/* Entry Date and Price */}
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div>
                  <label className="block text-sm font-helvetica-medium text-primary mb-2">
                    Entry Date *
                  </label>
                  <input
                    type="date"
                    value={formData.entryDate}
                    onChange={(e) => handleInputChange('entryDate', e.target.value)}
                    className="w-full bg-primary border border-primary text-primary placeholder-tertiary px-3 py-2 rounded-lg focus:outline-none focus:ring-2 focus:ring-accent"
                    required
                  />
                </div>
                <div>
                  <label className="block text-sm font-helvetica-medium text-primary mb-2">
                    Entry Price *
                  </label>
                  <input
                    type="number"
                    step="0.01"
                    value={formData.entryPrice || ''}
                    onChange={(e) => handleNumberInput('entryPrice', e.target.value)}
                    onBlur={(e) => handleNumberBlur('entryPrice', parseFloat(e.target.value) || 0)}
                    placeholder="Entry Price"
                    className="w-full bg-primary border border-primary text-primary placeholder-tertiary px-3 py-2 rounded-lg focus:outline-none focus:ring-2 focus:ring-accent"
                    required
                  />
                </div>
              </div>

              {/* Quantity and Total Amount */}
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div>
                  <label className="block text-sm font-helvetica-medium text-primary mb-2">
                    Quantity *
                  </label>
                  <input
                    type="number"
                    value={formData.quantity || ''}
                    onChange={(e) => handleNumberInput('quantity', e.target.value)}
                    onBlur={(e) => handleNumberBlur('quantity', parseInt(e.target.value) || 0)}
                    placeholder="Quantity"
                    className="w-full bg-primary border border-primary text-primary placeholder-tertiary px-3 py-2 rounded-lg focus:outline-none focus:ring-2 focus:ring-accent"
                    required
                    min="1"
                    step="1"
                  />
                </div>
                <div>
                  <label className="block text-sm font-helvetica-medium text-primary mb-2">
                    Total amount
                  </label>
                  <input
                    type="number"
                    step="0.01"
                    value={formatNumber(calculateTotalAmount())}
                    readOnly
                    placeholder="Auto-calculated"
                    className="w-full bg-secondary border border-primary text-primary placeholder-tertiary px-3 py-2 rounded-lg focus:outline-none focus:ring-2 focus:ring-accent cursor-not-allowed"
                  />
                </div>
              </div>

              {/* Exit Price and P&L */}
              <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                <div>
                  <label className="block text-sm font-helvetica-medium text-primary mb-2">
                    Exit Price *
                  </label>
                  <input
                    type="number"
                    step="0.01"
                    value={formData.exitPrice || ''}
                    onChange={(e) => handleNumberInput('exitPrice', e.target.value)}
                    onBlur={(e) => handleNumberBlur('exitPrice', parseFloat(e.target.value) || 0)}
                    placeholder="Exit Price"
                    className="w-full bg-primary border border-primary text-primary placeholder-tertiary px-3 py-2 rounded-lg focus:outline-none focus:ring-2 focus:ring-accent"
                    required
                  />
                </div>
                <div>
                  <label className="block text-sm font-helvetica-medium text-primary mb-2">
                    P&L Amount
                  </label>
                  <input
                    type="text"
                    value={formatNumber(calculatePnL())}
                    readOnly
                    placeholder="Auto-calculated"
                    className="w-full bg-secondary border border-primary text-primary placeholder-tertiary px-3 py-2 rounded-lg focus:outline-none focus:ring-2 focus:ring-accent cursor-not-allowed"
                  />
                </div>
              </div>

              {/* Direction */}
              <div>
                <label className="block text-sm font-helvetica-medium text-primary mb-2">
                  Direction *
                </label>
                <div className="flex gap-2">
                  <button
                    type="button"
                    onClick={() => handleInputChange('direction', TradeDirection.LONG)}
                    className={`flex items-center gap-2 px-4 py-2 rounded-lg border transition-colors duration-200 ${
                      formData.direction === TradeDirection.LONG
                        ? 'bg-green-500/20 border-green-500 text-green-400'
                        : 'bg-primary border-primary text-tertiary hover:text-primary'
                    }`}
                  >
                    <ArrowUp className="w-4 h-4" />
                    Long
                  </button>
                  <button
                    type="button"
                    onClick={() => handleInputChange('direction', TradeDirection.SHORT)}
                    className={`flex items-center gap-2 px-4 py-2 rounded-lg border transition-colors duration-200 ${
                      formData.direction === TradeDirection.SHORT
                        ? 'bg-red-500/20 border-red-500 text-red-400'
                        : 'bg-primary border-primary text-tertiary hover:text-primary'
                    }`}
                  >
                    <ArrowDown className="w-4 h-4" />
                    Short
                  </button>
                </div>
              </div>

              {/* Stop Loss and Target */}
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div>
                  <label className="block text-sm font-helvetica-medium text-primary mb-2">
                    Stop Loss
                  </label>
                  <input
                    type="number"
                    step="0.01"
                    value={formData.stopLoss || ''}
                    onChange={(e) => handleNumberInput('stopLoss', e.target.value)}
                    onBlur={(e) => handleNumberBlur('stopLoss', parseFloat(e.target.value) || 0)}
                    placeholder="Stop Loss"
                    className="w-full bg-primary border border-primary text-primary placeholder-tertiary px-3 py-2 rounded-lg focus:outline-none focus:ring-2 focus:ring-accent"
                  />
                </div>
                <div>
                  <label className="block text-sm font-helvetica-medium text-primary mb-2">
                    Target
                  </label>
                  <input
                    type="number"
                    step="0.01"
                    value={formData.target || ''}
                    onChange={(e) => handleNumberInput('target', e.target.value)}
                    onBlur={(e) => handleNumberBlur('target', parseFloat(e.target.value) || 0)}
                    placeholder="Target"
                    className="w-full bg-primary border border-primary text-primary placeholder-tertiary px-3 py-2 rounded-lg focus:outline-none focus:ring-2 focus:ring-accent"
                  />
                </div>
              </div>

              {/* Strategy and Outcome Summary */}
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div>
                  <label className="block text-sm font-helvetica-medium text-primary mb-2">
                    Strategy *
                  </label>
                  <select
                    value={formData.strategy}
                    onChange={(e) => handleInputChange('strategy', e.target.value)}
                    className="w-full bg-primary border border-primary text-primary placeholder-tertiary px-3 py-2 rounded-lg focus:outline-none focus:ring-2 focus:ring-accent"
                    required
                  >
                    <option value="">Select Strategy</option>
                    <option value="scalping">Scalping</option>
                    <option value="day-trading">Day Trading</option>
                    <option value="swing-trading">Swing Trading</option>
                    <option value="position-trading">Position Trading</option>
                    <option value="momentum">Momentum</option>
                    <option value="mean-reversion">Mean Reversion</option>
                    <option value="breakout">Breakout</option>
                    <option value="support-resistance">Support/Resistance</option>
                  </select>
                </div>
                <div>
                  <label className="block text-sm font-helvetica-medium text-primary mb-2">
                    Outcome Summary *
                  </label>
                  <select
                    value={formData.outcomeSummary}
                    onChange={(e) => handleInputChange('outcomeSummary', e.target.value)}
                    className="w-full bg-primary border border-primary text-primary placeholder-tertiary px-3 py-2 rounded-lg focus:outline-none focus:ring-2 focus:ring-accent"
                    required
                  >
                    <option value="">Select Outcome Summary</option>
                    {outcomeSummaries.map((outcome) => (
                      <option key={outcome.value} value={outcome.value}>
                        {outcome.label}
                      </option>
                    ))}
                  </select>
                </div>
              </div>

              {/* Trade Analysis */}
              <div>
                <label className="block text-sm font-helvetica-medium text-primary mb-2">
                  Trade Analysis
                </label>
                <textarea
                  value={formData.tradeAnalysis}
                  onChange={(e) => handleInputChange('tradeAnalysis', e.target.value)}
                  placeholder="Why did you take this trade? What was your analysis?"
                  rows={3}
                  className="w-full bg-primary border border-primary text-primary placeholder-tertiary px-3 py-2 rounded-lg focus:outline-none focus:ring-2 focus:ring-accent"
                />
              </div>

              {/* Rules Followed */}
              <div>
                <label className="block text-sm font-helvetica-medium text-primary mb-2">
                  Rules Followed
                </label>
                <div className="flex gap-2 mb-2">
                  <input
                    type="text"
                    value={newRule}
                    onChange={(e) => setNewRule(e.target.value)}
                    placeholder="Search or add rules..."
                    className="flex-1 bg-primary border border-primary text-primary placeholder-tertiary px-3 py-2 rounded-lg focus:outline-none focus:ring-2 focus:ring-accent"
                    onKeyPress={(e) => e.key === 'Enter' && (e.preventDefault(), addRule())}
                  />
                  <button
                    type="button"
                    onClick={addRule}
                    className="bg-accent hover:bg-accent-hover text-inverse px-3 py-2 rounded-lg transition-colors duration-200"
                  >
                    <Plus className="w-4 h-4 text-inverse" />
                  </button>
                </div>
                <div className="flex flex-wrap gap-2">
                  {formData.rulesFollowed.map((rule, index) => (
                    <span
                      key={index}
                      className="bg-primary border border-primary text-primary px-3 py-1 rounded-lg text-sm flex items-center gap-2"
                    >
                      {rule}
                      <button
                        type="button"
                        onClick={() => removeRule(rule)}
                        className="text-tertiary hover:text-primary"
                      >
                        <X className="w-3 h-3" />
                      </button>
                    </span>
                  ))}
                </div>
              </div>

              {/* Trade Screenshots */}
              <div>
                <label className="block text-sm font-helvetica-medium text-primary mb-2">
                  Trade Screenshots
                </label>
                <div className="border-2 border-dashed border-primary rounded-lg p-8 text-center hover:border-accent transition-colors duration-200">
                  <Upload className="w-8 h-8 text-tertiary mx-auto mb-2" />
                  <p className="text-primary font-helvetica-medium mb-1">Drag & drop your trade screenshots here</p>
                  <p className="text-tertiary text-sm">Supports JPG, PNG (Max 5MB each)</p>
                </div>
              </div>
            </div>
          )}

          {activeTab === 'psychology' && (
            <div className="space-y-6">
              {/* Entry Confidence Level */}
              <div>
                <label className="block text-sm font-helvetica-medium text-primary mb-2">
                  Entry Confidence Level (1-10)
                </label>
                <div className="flex items-center gap-4">
                  <input
                    type="range"
                    min="1"
                    max="10"
                    value={formData.psychology.entryConfidence}
                    onChange={(e) => handlePsychologyChange('entryConfidence', parseInt(e.target.value))}
                    className="flex-1"
                  />
                  <span className="text-primary font-helvetica-bold text-lg min-w-[2rem]">
                    {formData.psychology.entryConfidence}
                  </span>
                </div>
                <div className="flex justify-between text-xs text-tertiary mt-1">
                  <span>Low</span>
                  <span>Medium</span>
                  <span>High</span>
                </div>
              </div>

              {/* Satisfaction Rating */}
              <div>
                <label className="block text-sm font-helvetica-medium text-primary mb-2">
                  Satisfaction Rating (1-10)
                </label>
                <div className="flex items-center gap-4">
                  <input
                    type="range"
                    min="1"
                    max="10"
                    value={formData.psychology.satisfactionRating}
                    onChange={(e) => handlePsychologyChange('satisfactionRating', parseInt(e.target.value))}
                    className="flex-1"
                  />
                  <span className="text-primary font-helvetica-bold text-lg min-w-[2rem]">
                    {formData.psychology.satisfactionRating}
                  </span>
                </div>
                <div className="flex justify-between text-xs text-tertiary mt-1">
                  <span>Not Satisfied</span>
                  <span>Average</span>
                  <span>Satisfied</span>
                </div>
              </div>

              {/* Emotional State */}
              <div>
                <label className="block text-sm font-helvetica-medium text-primary mb-2">
                  Emotional State During Trade *
                </label>
                <select
                  value={formData.psychology.emotionalState}
                  onChange={(e) => handlePsychologyChange('emotionalState', e.target.value)}
                  className="w-full bg-primary border border-primary text-primary placeholder-tertiary px-3 py-2 rounded-lg focus:outline-none focus:ring-2 focus:ring-accent"
                  required
                >
                  <option value="">Select Emotional State</option>
                  {emotionalStates.map((state) => (
                    <option key={state} value={state}>
                      {state}
                    </option>
                  ))}
                </select>
              </div>

              {/* Mistakes Made */}
              <div>
                <label className="block text-sm font-helvetica-medium text-primary mb-2">
                  Mistakes Made
                </label>
                <div className="grid grid-cols-2 gap-2">
                  {commonMistakes.map((mistake) => (
                    <label key={mistake} className="flex items-center gap-2 cursor-pointer">
                      <input
                        type="checkbox"
                        checked={formData.psychology.mistakesMade.includes(mistake)}
                        onChange={() => toggleMistake(mistake)}
                        className="rounded border-primary text-accent focus:ring-accent"
                      />
                      <span className="text-primary font-helvetica-light text-sm">{mistake}</span>
                    </label>
                  ))}
                </div>
              </div>

              {/* Lessons Learned */}
              <div>
                <label className="block text-sm font-helvetica-medium text-primary mb-2">
                  Lessons Learned
                </label>
                <textarea
                  value={formData.psychology.lessonsLearned}
                  onChange={(e) => handlePsychologyChange('lessonsLearned', e.target.value)}
                  placeholder="What did you learn from this trade?"
                  rows={4}
                  className="w-full bg-primary border border-primary text-primary placeholder-tertiary px-3 py-2 rounded-lg focus:outline-none focus:ring-2 focus:ring-accent"
                />
              </div>
            </div>
          )}

          {/* Action Buttons */}
          <div className="flex justify-end gap-3 pt-6 border-t border-primary">
            <button
              type="button"
              onClick={handleClose}
              className="bg-primary hover:bg-tertiary text-primary font-helvetica-medium px-4 py-2 rounded-lg transition-colors duration-200 focus:outline-none focus:ring-2 focus:ring-primary"
            >
              Reset
            </button>
            <button
              type="submit"
              disabled={isSubmitting || !formData.symbol.trim() || !formData.strategy.trim() || !formData.outcomeSummary.trim()}
              className="bg-accent hover:bg-accent-hover disabled:bg-tertiary disabled:cursor-not-allowed text-primary font-helvetica-medium py-2 px-4 rounded-lg transition-colors duration-200 focus:outline-none focus:ring-2 focus:ring-accent"
            >
              {isSubmitting ? 'Saving Trade...' : 'Save Trade'}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
