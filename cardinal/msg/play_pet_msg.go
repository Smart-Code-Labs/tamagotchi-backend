// Package msg contains message structures for the Tamagotchi game.
package msg

/**
 * Function Flow:
 * 1. The PlayPetMsg structure is created to hold the target nickname and item name for the play pet action.
 * 2. The PlayPetMsgReply structure is created to hold the reply data for the play pet action.
 *
 * This package provides message structures for the play pet action.
 */
type PlayPetMsg struct {
	/**
	 * TargetNickname is the nickname of the pet to be played with.
	 */
	TargetNickname string `json:"target"`
	/**
	 * ItemName is the name of the item to be used for the play pet action.
	 */
	ItemName string `json:"item_name"`
}

/**
 * Function Flow:
 * 1. The PlayPetMsgReply structure is created to hold the reply data for the play pet action.
 * 2. The Energy field holds the updated energy value of the pet.
 * 3. The Hygiene field holds the updated hygiene value of the pet.
 * 4. The Wellness field holds the updated wellness value of the pet.
 * 5. The Activity field holds the current activity of the pet.
 * 6. The Duration field holds the duration of the play pet action.
 *
 * This structure provides the reply data for the play pet action.
 */
type PlayPetMsgReply struct {
	/**
	 * Energy is the updated energy value of the pet.
	 */
	Energy int `json:"energy"`
	/**
	 * Hygiene is the updated hygiene value of the pet.
	 */
	Hygiene int `json:"hygiene"`
	/**
	 * Wellness is the updated wellness value of the pet.
	 */
	Wellness int `json:"wellness"`
	/**
	 * Activity is the current activity of the pet.
	 */
	Activity string `json:"activity"`
	/**
	 * Duration is the duration of the play pet action.
	 */
	Duration int `json:"duration"`
}
