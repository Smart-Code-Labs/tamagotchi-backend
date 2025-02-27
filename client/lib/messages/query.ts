import type { Item } from "../entity/item";
import type { Pet } from "../entity/pet";

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
	Pets: Pet[]
}

export interface PlayerItemsMsg {
	// The persona tag of the player to query.
	personaTag: string
}

// ItemListReply represents the response to a player items query.
export interface PlayerItemsResponse  {
	// The list of items belonging to the player.
	items: Item[]
}

export interface LeaderboardMsg {}

export interface LeaderboardReply {
	pets: Pet[]
}

export interface PlayerExistMsg  {
	personaTag: string
}

export interface PlayerExistReply  {
	exist: boolean
}
