export interface ITrade {
  id: string;
  userId: string;
  
  // General Trade Information
  symbol: string;
  marketType: string;
  entryDate: Date;
  entryPrice: number;
  quantity: number;
  totalAmount: number;
  exitPrice?: number;
  pnlAmount?: number;
  pnlPercentage?: number;
  direction: 'long' | 'short';
  stopLoss?: number;
  target?: number;
  strategy: string;
  outcomeSummary: string;
  tradeAnalysis?: string;
  rulesFollowed?: string[];
  screenshots?: string[];
  
  // Psychology Information
  psychology: ITradePsychology;
  
  // Metadata
  createdAt: Date;
  updatedAt: Date;
}

export interface ITradePsychology {
  entryConfidence: number; // 1-10 scale
  satisfactionRating: number; // 1-10 scale
  emotionalState: string;
  mistakesMade: string[];
  lessonsLearned?: string;
}

export interface ICreateTradeRequest {
  symbol: string;
  marketType: string;
  entryDate: string;
  entryPrice: number;
  quantity: number;
  totalAmount: number;
  exitPrice?: number;
  pnlAmount?: number;
  pnlPercentage?: number;
  direction: 'long' | 'short';
  stopLoss?: number;
  target?: number;
  strategy: string;
  outcomeSummary: string;
  tradeAnalysis?: string;
  rulesFollowed?: string[];
  screenshots?: string[];
  psychology: ITradePsychology;
}

export interface ITradeFormData {
  // General Information
  symbol: string;
  marketType: string;
  entryDate: string;
  entryPrice: number;
  quantity: number;
  totalAmount: number;
  exitPrice?: number;
  pnlAmount?: number;
  pnlPercentage?: number;
  direction: 'long' | 'short';
  stopLoss?: number;
  target?: number;
  strategy: string;
  outcomeSummary: string;
  tradeAnalysis?: string;
  rulesFollowed: string[];
  screenshots: string[];
  
  // Psychology Information
  psychology: ITradePsychology;
}

export interface IModalProps {
  isOpen: boolean;
  onClose: () => void;
  onSubmit: (data: ITradeFormData) => void;
}

// Psychology-specific interfaces
export interface IEmotionalState {
  id: string;
  name: string;
  description: string;
}

export interface ITradingMistake {
  id: string;
  name: string;
  description: string;
}

// Market and strategy enums
export type MarketType = 'indian' | 'us' | 'crypto' | 'forex' | 'commodities';
export type TradeDirection = 'long' | 'short';
export type TradeDuration = 'intraday' | 'swing' | 'positional';
export type OutcomeSummary = 'profitable' | 'loss' | 'breakeven' | 'partial_profit' | 'partial_loss';
