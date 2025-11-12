export interface ITrade {
  id: string;
  userId: string;

  // General Trade Information
  symbol: string;
  marketType: MarketType;
  entryDate: Date;
  entryPrice: number;
  quantity: number;
  totalAmount: number;
  exitPrice?: number;
  direction: TradeDirection;
  stopLoss?: number;
  target?: number;
  strategy: string;
  outcomeSummary: OutcomeSummary;
  tradeAnalysis?: string;
  rulesFollowed?: string[];
  screenshots?: string[];

  // Psychology Information
  psychology: ITradePsychology | null;

  // Broker-specific fields (optional, for imported trades)
  tradingBroker?: TradingBroker;
  traderBrokerId?: string;
  exchangeOrderId?: string;
  orderId?: string;
  productType?: ProductType;
  transactionType?: string; // "buy" | "sell"

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
  userId: string;
  symbol: string;
  marketType: MarketType;
  entryDate: string;
  entryPrice: number;
  quantity: number;
  totalAmount: number;
  exitPrice?: number;
  direction: TradeDirection;
  stopLoss?: number;
  target?: number;
  strategy: string;
  outcomeSummary: OutcomeSummary;
  tradeAnalysis?: string;
  rulesFollowed?: string[];
  screenshots?: string[];
  psychology: ITradePsychology;
}

export interface ITradeFormData {
  // General Information
  symbol: string;
  marketType: MarketType;
  entryDate: string;
  entryPrice: number;
  quantity: number;
  totalAmount: number;
  exitPrice?: number;
  direction: TradeDirection;
  duration: TradeDuration;
  stopLoss?: number;
  target?: number;
  strategy: string;
  outcomeSummary: OutcomeSummary;
  tradeAnalysis?: string;
  rulesFollowed: string[];
  screenshots: string[];

  // Psychology Information
  psychology: ITradePsychology;
}

export interface ITradeModalProps {
  isOpen: boolean;
  onClose: () => void;
  onSubmit: (data: ITradeFormData) => void;
  initialData?: ITrade; // For editing existing trades
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
export enum MarketType {
  INDIAN = "indian",
  US = "us",
  CRYPTO = "crypto",
  FOREX = "forex",
  COMMODITIES = "commodities",
}

export enum TradeDirection {
  LONG = "long",
  SHORT = "short",
}

export enum TradeDuration {
  INTRADAY = "intraday",
  SWING = "swing",
  POSITIONAL = "positional",
}

export enum OutcomeSummary {
  PROFITABLE = "profitable",
  LOSS = "loss",
  BREAKEVEN = "breakeven",
  PARTIAL_PROFIT = "partial_profit",
  PARTIAL_LOSS = "partial_loss",
}

export enum TradingBroker {
  ZERODHA = "zerodha",
  DHAN = "dhan",
}

export enum ProductType {
  CNC = "CNC",
  MIS = "MIS",
  NRML = "NRML",
  INTRADAY = "INTRADAY",
  OTC = "OTC",
}
