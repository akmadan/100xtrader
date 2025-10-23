enum RuleCategory {
  ENTRY = "entry",
  EXIT = "exit",
  STOP_LOSS = "stop_loss",
  TAKE_PROFIT = "take_profit",
  RISK_MANAGEMENT = "risk_management",
  PSYCHOLOGY = "psychology",
  OTHER = "other",
}

export interface IRule {
  id: string;
  userId: string;
  name: string;
  description: string;
  category: RuleCategory;
  createdAt: Date;
  updatedAt: Date;
}

export interface IRuleCreateRequest {
  userId: string;
  name: string;
  description: string;
  category: RuleCategory;
}

export interface IRuleUpdateRequest {
  id: string;
  name: string;
  description: string;
  category: RuleCategory;
}
