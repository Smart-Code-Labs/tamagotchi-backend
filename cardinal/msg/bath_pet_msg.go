// Package msg contains message structures for the Tamagotchi game.
package msg

/**
 * Function Flow:
 * 1. The BathPetMsg structure is created to hold the target nickname for the bath pet action.
 * 2. The BathPetMsgReply structure is created to hold the reply data for the bath pet action.
 *
 * This package provides message structures for the bath pet action.
 */
type BathPetMsg struct {
	/**
	 * TargetNickname is the nickname of the pet to be bathed.
	 */
	TargetNickname string `json:"target"`
	/**
	 * ItemName is the name of the item to be used for the play pet action.
	 */
	ItemName string `json:"item_name"`
}

/**
 * Function Flow:
 * 1. The BathPetMsgReply structure is created to hold the reply data for the bath pet action.
 * 2. The Hygiene field holds the updated hygiene value of the pet.
 * 3. The Activity field holds the current activity of the pet.
 * 4. The Duration field holds the duration of the bath pet action.
 *
 * This structure provides the reply data for the bath pet action.
 */
type BathPetMsgReply struct {
	/**
	 * Hygiene is the updated hygiene value of the pet.
	 */
	Hygiene int `json:"hygiene"`
	/**
	 * Activity is the current activity of the pet.
	 */
	Activity string `json:"activity"`
	/**
	 * Duration is the duration of the bath pet action.
	 */
	Duration int `json:"duration"`
}

// bath_pet_msg.go
