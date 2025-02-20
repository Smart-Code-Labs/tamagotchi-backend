export interface RpcFindMatchRequest {
    fast: boolean;
}

export interface RpcFindMatchResponse {
    matchIds: string[];
}

export interface RpcFindPersonaResponse {
    personaTag: string;
    status: string;
    tick: number
    txHash: string;
}

export interface RpcCurrentTickResponse {
    currentTick: number
}

export interface PetEnergyRequest {
	Nickname: string
}

export interface PetEnergyResponse {
	energy: number
}

export interface PetHealthRequest  {
	Nickname: string
}

export interface PetHealthResponse {
	HP: number
}

export interface PetsRequest {}

export interface PetsResponse {
	Pets: []
}
