import type { Session } from "@heroiclabs/nakama-js";
import Nakama from "./lib/nakama";
import type { TxResponse } from "./lib/messages/execute";

var client = new Nakama();

const email = "email@example.com";
const password = "3bc8f72e95a9";
const personaId = "pepe5";
const petId = "Manny2";

await client.authenticate(email, password);
const id = await client.getCustomId();
await client.claimPersona(personaId);
const txResponse = await client.createPet(petId);
await client.queryTick();
await client.queryPetEnergy(petId);
await client.queryPetHealth(petId);
await client.queryPets();

if (txResponse) {
  const receipts = await client.waitForReceipt(txResponse)
  console.log("Receipts:", receipts);
}

const bathTxResponse = await client.bathPet(petId);
if (bathTxResponse) {
  const receipts = await client.waitForReceipt(bathTxResponse)
  console.log("Receipts:", receipts);
}