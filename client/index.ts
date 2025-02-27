import type { Session } from "@heroiclabs/nakama-js";
import Nakama from "./lib/nakama";
import type { TxResponse } from "./lib/messages/execute";
import type { Item } from "./lib/entity/item";

var client = new Nakama();

interface Prop {
  petName: string;
  activity: string;
  itemName: string;
}

const email = "email@example.com";
const password = "3bc8f72e95a9";
const personaTag = createRandomString(6);
const petMap: Prop[] = [
  { petName: createRandomString(5), activity: "Cure", itemName: "Vaccine" },
  { petName: createRandomString(5), activity: "Play", itemName: "Ball" },
  { petName: createRandomString(5), activity: "Bath", itemName: "Sponge" },
  { petName: createRandomString(5), activity: "Feed", itemName: "Soup" },
]

// User login in nakama and in his account
await client.authenticate(email, password);
const id = await client.getCustomId();

// Create a Persona `PersonaTag` and associate to the account 
await client.claimPersona(personaTag);

// delay 3 seconds to ensure Persona is created
await client.waitTicks(1)

// Create a Player `personaTag` associated to the Account `Persona`
const createPlayerTxResponse = await client.createPlayer();
if (createPlayerTxResponse) {
  const receipts = await client.waitForReceipt(createPlayerTxResponse)
  console.log("Receipts:", receipts);
}

// ensure Player is created
await client.waitTicks(1)


// Create Pet
await Promise.all(petMap.map(async pet => {
  const createPetTxResponse = await client.createPet(pet.petName);
  // await client.queryTick();
  // await client.queryPetEnergy(petId);
  // await client.queryPetHealth(petId);
  // await client.queryPets();

  if (createPetTxResponse) {
    const receipts = await client.waitForReceipt(createPetTxResponse)
    console.log("Receipts:", receipts);
  }

  // Buy each item to do actions
  const buyItemTxResponse = await client.buyItem(pet.itemName);

  if (buyItemTxResponse) {
    const receipts = await client.waitForReceipt(buyItemTxResponse)
    console.log("Receipts:", receipts);
  }

}));

const playerItems: Item[] | undefined = await client.queryPlayerItems(personaTag);

console.log(playerItems)

await Promise.all(petMap.map(async pet => {
  let txResponse;
  switch (pet.activity) {
    case "Play":
      txResponse = await client.playPet(pet.petName, pet.itemName);
      if (txResponse) {
        const receipts = await client.waitForReceipt(txResponse)
        console.log("Receipts:", receipts);
      }
      break
    case "Cure":
      txResponse = await client.curePet(pet.petName, pet.itemName);
      if (txResponse) {
        const receipts = await client.waitForReceipt(txResponse)
        console.log("Receipts:", receipts);
      }
      break
    case "Bath":
      txResponse = await client.bathPet(pet.petName, pet.itemName);
      if (txResponse) {
        const receipts = await client.waitForReceipt(txResponse)
        console.log("Receipts:", receipts);
      }
      break
    case "Feed":
      txResponse = await client.feedPet(pet.petName, pet.itemName);
      if (txResponse) {
        const receipts = await client.waitForReceipt(txResponse)
        console.log("Receipts:", receipts);
      }
      break
  }
}));


function createRandomString(length: number) {
  const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789";
  let result = "";
  for (let i = 0; i < length; i++) {
    result += chars.charAt(Math.floor(Math.random() * chars.length));
  }
  return result;
}
