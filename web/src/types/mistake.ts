export enum MistakeCategory {
  PSYCHOLOGICAL = "psychological",
  BEHAVIORAL = "behavioral",
}

export interface IMistake {
  id: string;
  userId: string;
  name: string;
  category: MistakeCategory;
  createdAt: Date;
  updatedAt: Date;
}

export interface IMistakeCreateRequest {
  userId: string;
  name: string;
  category: MistakeCategory;
}

export interface IMistakeUpdateRequest {
  id: string;
  name: string;
  category: MistakeCategory;
}
