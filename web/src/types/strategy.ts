export interface IStrategy {
  id: string;
  userId: string;
  name: string;
  description: string;
  createdAt: Date;
  updatedAt: Date;
}

export interface ICreateStrategyRequest {
  userId: string;
  name: string;
  description: string;
}

export interface IStrategyFormData {
  name: string;
  description: string;
}

export interface IModalProps {
  isOpen: boolean;
  onClose: () => void;
  onSubmit: (data: IStrategyFormData) => void;
}
