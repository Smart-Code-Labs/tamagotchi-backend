// Package msg contains message structures for the Tamagotchi game.
package msg

/**
 * Function Flow:
 * 1. The CurePetMsg structure is created to hold the target nickname and item id for the cure pet action.
 * 2. The CurePetMsgReply structure is created to hold the reply data for the cure pet action.
 *
 * This package provides message structures for the cure pet action.
 */
type CurePetMsg struct {
	/**
	 * TargetNickname is the nickname of the pet to be cured.
	 */
	TargetNickname string `json:"target"`
	/**
	 * ItemId is the id of the item to be used for the cure pet action.
	 */
	ItemName string `json:"item_name"`
}

/**
 * Function Flow:
 * 1. The CurePetMsgReply structure is created to hold the reply data for the cure pet action.
 * 2. The Health field holds the updated health value of the pet.
 *
 * This structure provides the reply data for the cure pet action.
 */
type CurePetMsgReply struct {
	/**
	 * Health is the updated health value of the pet.
	 */
	Health int `json:"health"`
}

// cure_pet_msg.go
