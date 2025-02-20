export interface TxResponse {
  TxHash: string;
  Tick: number;
}

export interface Receipt {
  txHash: string;
  tick: number;
  result: any;
  errors: [];
}
export interface ReceiptsResponse {
  startTick: number;
  endTick: number;
  receipts: Receipt[];
}

export interface CreatePetMsg {
  nickname: string;
}

export interface CreatePetResult {
  success: boolean;
}

export interface BathPetMsg {
  target: string;
}

export interface BathPetMsgReply {
  hygiene: number;
  activity: number;
  duration: number;
}

export interface BreedPetMsg {
  motherName: string;
  fatherName: string;
  bornName: string;
}

export interface BreedPetMsgReply {
  success: boolean;
}

export interface FeedPetMsg {
  target: string;
}

export interface FeedPetMsgReply {
  health: number;
  activity: string;
  duration: number;
}

export interface PlayPetMsg {
  target: string;
}

export interface PlayPetMsgReply {
  energy: number;
  hygiene: number;
  wellness: number;
  activity: string;
  duration: number;
}

export interface SleepPetMsg {
  target: string;
}

export interface SleepPetMsgReply {
  energy: number;
  activity: string;
  duration: number;
}

export interface ThinkPetMsg {
  target: string;
}

export interface ThinkPetMsgReply {
  think: string;
}
