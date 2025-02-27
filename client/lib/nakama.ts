import {
  Client,
  Session,
  type NotificationList,
  type RpcResponse,
  type Socket,
} from "@heroiclabs/nakama-js";
import {
  type LeaderboardMsg,
  type LeaderboardReply,
  type PetEnergyRequest,
  type PetEnergyResponse,
  type PetHealthRequest,
  type PetHealthResponse,
  type PetsRequest,
  type PetsResponse,
  type PlayerExistMsg,
  type PlayerExistReply,
  type PlayerItemsMsg,
  type PlayerItemsResponse,
  type RpcCurrentTickResponse,
  type RpcFindPersonaResponse,
} from "./messages/query";
import type {
  BathPetMsg,
  BreedPetMsg,
  ButItemMsg,
  CreatePetMsg,
  CreatePlayerMsg,
  FeedPetMsg,
  PlayPetMsg,
  Receipt,
  ReceiptsResponse,
  SleepPetMsg,
  TxResponse,
} from "./messages/execute";
import type { Pet } from "./entity/pet";
import type { Item } from "./entity/item";
import { udpSocket } from "bun";

class GameState {
  public playerIndex = 0;
}

class Nakama {
  public client: Client;
  public session: Session | null = null;
  public socket: Socket | null = null;
  public matchId: string | null = null;
  public gameState: GameState = new GameState();

  constructor() {
    this.client = new Client(
      "defaultkey",
      "127.0.0.1",
      "7350",
      false,
      10000,
      true
    );
  }

  async authenticate(email: string, password: string): Promise<void> {
    this.session = await this.client.authenticateEmail(email, password);
    console.info("Successfully authenticated:", this.session);
    if (!this.session?.user_id) return;
    // Can't find variable: localStorage
    // localStorage.setItem("user_id", this.session.user_id);

    const trace = false;
    this.socket = this.client.createSocket(true, trace);
    console.log(`socket created ${JSON.stringify(this.socket)}`);
    this.socket.ondisconnect = (evt) => {
      console.info("Disconnected", evt);
    };
    this.socket.onerror = (evt) => {
      console.info("Socket error", evt);
    };

    this.socket.onchannelmessage = (message) => {
      console.info("Message received from channel", message.channel_id);
      console.info("Received message", message);
    };

    this.socket.onnotification = (notification) => {
      console.info("Notification received from channel ", notification.id);
      console.info("Received message", notification.subject);
    };

    this.socket.onstreamdata = (stream) => {
      console.info("Received stream", stream);
    };
    // const pepe = await this.socket.connect(this.session, true);
    // // Socket is open.
    // console.log(`${pepe}`);
  }

  async getCustomId(): Promise<string | undefined> {
    try {
      const account = await this.client.getAccount(this.session!);
      console.info(account.user!.id);
      console.info(account.user!.username);
      console.info(account.wallet);
      return account.custom_id!;
    } catch (error: unknown) {
      if (error instanceof Error) {
        console.error("Inner Nakama error", error.message);
      } else {
        console.error("Unknown error occurred");
      }
    }
  }

  async createMatch(): Promise<void> {
    if (!this.socket || !this.session) return;
    const match = await this.socket.createMatch();
    console.log("Match created:", match.match_id);
  }

  async findMatch(): Promise<void> {
    const rpc_name = "find_match_js";
    if (!this.session || !this.socket) {
      console.log("Session or socket not found");
      return;
    }
    const matches = await this.client.rpc(this.session, rpc_name, {});

    if (typeof matches === "object" && matches !== null) {
      const safeParsedJson = matches as {
        payload: {
          matchIds: string[];
          // height: string,
          // weight: string,
          // image: string,
        };
      };
      this.matchId = safeParsedJson.payload.matchIds[0];
      await this.socket.joinMatch(this.matchId);
      console.log("Match joined!");
    }
  }

  async getPersona(): Promise<RpcFindPersonaResponse | undefined> {
    if (!this.socket || !this.session) {
      console.log("Socket or session not found");
      throw new Error("Socket or session not found");
    }
    // Check whether a session is close to expiry.
    if (this.session.isexpired(Date.now() / 1000)) {
      try {
        console.log(`Session expired ${this.session.created_at}`);
        this.session = await this.client.sessionRefresh(this.session);
      } catch (e) {
        console.info(
          "Session can no longer be refreshed. Must reauthenticate!"
        );
        throw new Error("Failed to refresh session");
      }
    }
    console.log("getPersona: Session is OK.");
    try {
      const result: RpcResponse = await this.client.rpc(
        this.session,
        "nakama/show-persona",
        {}
      );

      console.log(`getPersona: [${JSON.stringify(result)}]`);
      const persona = result.payload! as RpcFindPersonaResponse;
      console.log(persona);

      return persona;
    } catch (error: any) {
      // if (error instanceof Error) {
      //   console.error("getPersona: Inner Nakama error", error.message);
      // } else {
      //   console.error("getPersona: Unknown error occurred", error);
      // }
      return undefined
    }
  }

  async claimPersona(personaTag: string): Promise<void> {
    return new Promise<void>(async (resolve, reject) => {
      if (!this.socket || !this.session) {
        reject("Socket or session not found");
      } else {
        // Check whether a session is close to expiry.
        if (this.session.isexpired(Date.now() / 1000)) {
          try {
            console.log(`Session expired ${this.session.created_at}, refreshing`);
            this.session = await this.client.sessionRefresh(this.session);
          } catch (e) {
            reject("Session can no longer be refreshed. Must reauthenticate!");
          }
        }
        console.log("claimPersona: Session is OK.");

        try {
          const isPersona = await this.getPersona();
          console.log(`is Persona [${JSON.stringify(isPersona)}]`);
          if (isPersona == undefined) {
            const data = { personaTag: personaTag };
            console.log(`${JSON.stringify(data)}`);
            try {
              const result = await this.client.rpc(
                this.session,
                "nakama/claim-persona",
                data
              );
              console.log(`claimPersona response[${JSON.stringify(result)}]`);
              resolve()
            } catch (error: unknown) {
              if (error instanceof Error) {
                reject(`Inner Nakama error ${error.message}`);
              } else {
                reject(`Unknown error occurred ${error}`);
              }
            }
          } else {
            console.log("Persona already created.");
            resolve();
          }
        } catch (error) {
          reject(`Failed to get persona: ${error}`);
          // Handle the error or rethrow if necessary
        }
      }
    })
  }

  async createPlayer(): Promise<TxResponse | undefined> {
    if (!this.socket || !this.session) {
      console.log("Socket or session not found");
      return;
    }
    try {
      const persona = await this.getPersona();
      if (persona != undefined) {
        const playerExist = await this.queryPlayerExist(persona.personaTag);
        console.log(`queryPlayerExist Response [${JSON.stringify(playerExist)}]`)
        if (playerExist == undefined ){

          const data: CreatePlayerMsg = {};
          console.log(`${JSON.stringify(data)}`);

          const result = await this.client.rpc(
            this.session,
            "tx/game/create-player",
            data
          );
          console.log(`Create Player response [${JSON.stringify(result)}]`);
          const txResponse = result.payload! as TxResponse;
          return txResponse;
        } else {
          console.log(`Player already exist [${playerExist}]`)
        }
      } else {
        throw new Error("Persona not found");
      }
    } catch (error) {
      console.error("Unknown error occurred", error);
    }
  }

  async createPet(name: string): Promise<TxResponse | undefined> {
    if (!this.socket || !this.session) {
      console.log("Socket or session not found");
      return;
    }
    try {
      const persona = await this.getPersona();
      if (persona != undefined) {
        console.log(`Creating pet [${name}] belongs to [${persona.personaTag}]`)
        const data: CreatePetMsg = { nickname: name };
        console.log(`${JSON.stringify(data)}`);

        const result = await this.client.rpc(
          this.session,
          "tx/game/create-pet",
          data
        );
        console.log(`${JSON.stringify(result)}`);
        const txResponse = result.payload! as TxResponse;
        return txResponse;
      } else {
        throw new Error("Persona not found");
      }
    } catch (error) {
      console.error("Unknown error occurred", error);
    }
  }

  async bathPet(name: string, itemName: string): Promise<TxResponse | undefined> {
    if (!this.socket || !this.session) {
      console.log("Socket or session not found");
      return;
    }
    try {
      const persona = await this.getPersona();
      if (persona != undefined) {
        const data: BathPetMsg = { target: name, item_name: itemName };
        console.log(`${JSON.stringify(data)}`);

        const result = await this.client.rpc(
          this.session,
          "tx/game/bath-pet",
          data
        );
        console.log(`${JSON.stringify(result)}`);
        const txResponse = result.payload! as TxResponse;
        return txResponse;
      } else {
        throw new Error("Persona not found");
      }
    } catch (error) {
      console.error("Unknown error occurred", error);
    }
  }

  async feedPet(name: string, itemName: string): Promise<TxResponse | undefined> {
    if (!this.socket || !this.session) {
      console.log("Socket or session not found");
      return;
    }
    try {
      const persona = await this.getPersona();
      if (persona != undefined) {
        const data: FeedPetMsg = { target: name, item_name:itemName };
        console.log(`${JSON.stringify(data)}`);

        const result = await this.client.rpc(
          this.session,
          "tx/game/feed-pet",
          data
        );
        console.log(`${JSON.stringify(result)}`);
        const txResponse = result.payload! as TxResponse;
        return txResponse;
      } else {
        throw new Error("Persona not found");
      }
    } catch (error) {
      console.error("Unknown error occurred", error);
    }
  }

  async curePet(name: string, itemName: string): Promise<TxResponse | undefined> {
    if (!this.socket || !this.session) {
      console.log("Socket or session not found");
      return;
    }
    try {
      const persona = await this.getPersona();
      if (persona != undefined) {
        const data: FeedPetMsg = { target: name, item_name:itemName };
        console.log(`${JSON.stringify(data)}`);

        const result = await this.client.rpc(
          this.session,
          "tx/game/cure-pet",
          data
        );
        console.log(`${JSON.stringify(result)}`);
        const txResponse = result.payload! as TxResponse;
        return txResponse;
      } else {
        throw new Error("Persona not found");
      }
    } catch (error) {
      console.error("Unknown error occurred", error);
    }
  }

  async sleepPet(name: string): Promise<TxResponse | undefined> {
    if (!this.socket || !this.session) {
      console.log("Socket or session not found");
      return;
    }
    try {
      const persona = await this.getPersona();
      if (persona != undefined) {
        const data: SleepPetMsg = { target: name };
        console.log(`${JSON.stringify(data)}`);

        const result = await this.client.rpc(
          this.session,
          "tx/game/sleep-pet",
          data
        );
        console.log(`${JSON.stringify(result)}`);
        const txResponse = result.payload! as TxResponse;
        return txResponse;
      } else {
        throw new Error("Persona not found");
      }
    } catch (error) {
      console.error("Unknown error occurred", error);
    }
  }

  async playPet(name: string, itemName: string): Promise<TxResponse | undefined> {
    if (!this.socket || !this.session) {
      console.log("Socket or session not found");
      return;
    }
    try {
      const persona = await this.getPersona();
      if (persona != undefined) {
        const data: PlayPetMsg = { target: name , item_name: itemName};
        console.log(`${JSON.stringify(data)}`);

        const result = await this.client.rpc(
          this.session,
          "tx/game/play-pet",
          data
        );
        console.log(`${JSON.stringify(result)}`);
        const txResponse = result.payload! as TxResponse;
        return txResponse;
      } else {
        throw new Error("Persona not found");
      }
    } catch (error) {
      console.error("Unknown error occurred", error);
    }
  }

  async breedPet(fatherName: string, motherName: string, bornName: string): Promise<TxResponse | undefined> {
    if (!this.socket || !this.session) {
      console.log("Socket or session not found");
      return;
    }
    try {
      const persona = await this.getPersona();
      if (persona != undefined) {
        const data: BreedPetMsg = { motherName: motherName, fatherName: fatherName, bornName: bornName};
        console.log(`${JSON.stringify(data)}`);

        const result = await this.client.rpc(
          this.session,
          "tx/game/breed-pet",
          data
        );
        console.log(`${JSON.stringify(result)}`);
        const txResponse = result.payload! as TxResponse;
        return txResponse;
      } else {
        throw new Error("Persona not found");
      }
    } catch (error) {
      console.error("Unknown error occurred", error);
    }
  }

  async buyItem(name: string): Promise<TxResponse | undefined> {
    if (!this.socket || !this.session) {
      console.log("Socket or session not found");
      return;
    }
    try {

        console.log(`Buying item [${name}]`)
        const data: ButItemMsg = { name: name };
        console.log(`${JSON.stringify(data)}`);

        const result = await this.client.rpc(
          this.session,
          "tx/game/buy-item",
          data
        );
        console.log(`${JSON.stringify(result)}`);
        const txResponse = result.payload! as TxResponse;
        return txResponse;

    } catch (error) {
      console.error("Unknown error occurred", error);
    }
  }

  async getReceipts(startTick: number): Promise<ReceiptsResponse | undefined > {
    const options = {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: `{"startTick": ${startTick}}`,
    };
    const response: Response = await fetch(
      "http://localhost:4040/query/receipts/list",
      options
    );
    const jsonResponse = await response.json();
    const receiptResponse: ReceiptsResponse = jsonResponse as ReceiptsResponse;
    console.log(`getReceipts: ${JSON.stringify(receiptResponse)}`);
    return receiptResponse;
  }

  async queryTick(): Promise<number | undefined> {
    if (!this.socket || !this.session) {
      console.log("Socket or session not found");
      return;
    }
    try {
      const result: RpcResponse = await this.client.rpc(
        this.session,
        "query/game/current-tick",
        {}
      );
      console.log(`${JSON.stringify(result)}`);
      const tickResponse = result.payload! as RpcCurrentTickResponse;
      return tickResponse.currentTick;
    } catch (error) {
      console.error("Unknown error occurred", error);
    }
  }

  async queryPetEnergy(nickname: string): Promise<number | undefined> {
    if (!this.socket || !this.session) {
      console.log("Socket or session not found");
      return;
    }
    const data: PetEnergyRequest = {
      Nickname: nickname,
    };
    try {
      const result: RpcResponse = await this.client.rpc(
        this.session,
        "query/game/pet-energy",
        data
      );
      console.log(`${JSON.stringify(result)}`);
      const energyResponse = result.payload! as PetEnergyResponse;
      return energyResponse.energy;
    } catch (error) {
      console.error("Unknown error occurred", error);
    }
  }

  async queryPetHealth(nickname: string): Promise<number | undefined> {
    if (!this.socket || !this.session) {
      console.log("Socket or session not found");
      return;
    }
    const data: PetHealthRequest = {
      Nickname: nickname,
    };
    try {
      const result: RpcResponse = await this.client.rpc(
        this.session,
        "query/game/pet-health",
        data
      );
      console.log(`${JSON.stringify(result)}`);
      const healthResponse = result.payload! as PetHealthResponse;
      return healthResponse.HP;
    } catch (error) {
      console.error("Unknown error occurred", error);
    }
  }

  async queryPets(): Promise<Pet[] | undefined> {
    if (!this.socket || !this.session) {
      console.log("Socket or session not found");
      return;
    }
    const data: PetsRequest = {};
    try {
      const result: RpcResponse = await this.client.rpc(
        this.session,
        "query/game/pets-list",
        data
      );
      console.log(`${JSON.stringify(result)}`);
      const petsResponse = result.payload! as PetsResponse;
      return petsResponse.Pets;
    } catch (error) {
      console.error("Unknown error occurred", error);
    }
  }

  async queryPlayerExist(personaTag: string): Promise<boolean | undefined> {
    console.log(`queryPlayerExist [${personaTag}]`)
    if (!this.socket || !this.session) {
      console.log("Socket or session not found");
      return;
    }
    const data: PlayerExistMsg = {personaTag: personaTag};
    try {
      const result: RpcResponse = await this.client.rpc(
        this.session,
        "query/game/player-exist",
        data
      );
      console.log(`player-exist [${JSON.stringify(result)}]`);
      if (result.payload == undefined) {
        return false;
      }
      const playerExistResponse = result.payload! as PlayerExistReply;
      return playerExistResponse.exist;
    } catch (error) {
      console.error("Unknown error occurred", error);
      return undefined
    }
  }

  async queryPlayerItems(personaTag: string): Promise<Item[] | undefined> {
    console.log(`queryPlayerItems [${personaTag}]`)
    if (!this.socket || !this.session) {
      console.log("Socket or session not found");
      return;
    }
    const data: PlayerItemsMsg = {personaTag: personaTag};
    try {
      const result: RpcResponse = await this.client.rpc(
        this.session,
        "query/game/personaItem-list",
        data
      );
      console.log(`queryPlayerItems [${JSON.stringify(result)}]`);
      const petsResponse = result.payload! as PlayerItemsResponse;
      return petsResponse.items;
    } catch (error) {
      console.error("Unknown error occurred", error);
    }
  }

  async queryLeaderboard(): Promise<Pet[] | undefined> {
    console.log(`queryLeaderboard.`)
    if (!this.socket || !this.session) {
      console.log("Socket or session not found");
      return;
    }
    const data: LeaderboardMsg = {};
    try {
      const result: RpcResponse = await this.client.rpc(
        this.session,
        "query/game/leaderboard",
        data
      );
      console.log(`Leaderboard [${JSON.stringify(result)}]`);
      const leaderboardResponse = result.payload! as LeaderboardReply;
      return leaderboardResponse.pets;
    } catch (error) {
      console.error("Unknown error occurred", error);
    }
  }

  async waitTicks(thicks: number): Promise<void> {
    const initialTick = await this.queryTick();
    if (initialTick) {
      let retries = 0;
      const maxRetries = 5; // Maximum number of retries (5 seconds)

      while (retries <= maxRetries) {
        await this.delay(1000); // Wait for 1 second

        const latestTick = await this.queryTick();
        if (latestTick) {
          const elapsed = latestTick - initialTick;
          if (elapsed >= thicks) {
            return;
          }
        }

        retries++;
      }

      console.log("Waited for 60 seconds without meeting the tick requirement");
    }
  }

  async delay(ms: number) {
    return new Promise((resolve) => setTimeout(resolve, ms));
  }

  async waitForReceipt(txResponse: TxResponse): Promise<Receipt[] | undefined> {
    try {
      console.log("Waiting for ticks...");
      await this.waitTicks(2);
      console.log("Fetching receipts...");
      const receiptResponse = await this.getReceipts(txResponse.Tick);
      console.log("Filtering receipts...");
      if (receiptResponse != undefined) {
        const filteredReceipts = receiptResponse.receipts.filter((receipt) => {
          return receipt.txHash === txResponse.TxHash;
        });
        console.log("Filtered receipts:", filteredReceipts);
      
        // Process filteredReceipts
        return filteredReceipts;
      }
      
    } catch (error) {
      console.error("An error occurred:", error);
    }
  }
}

export default Nakama;
