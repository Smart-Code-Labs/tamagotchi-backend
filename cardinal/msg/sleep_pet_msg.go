// Package msg contains message structures for the Tamagotchi game.
package msg

/**
 * Function Flow:
 * 1. The SleepPetMsg structure is created to hold the target nickname for the sleep pet action.
 * 2. The SleepPetMsgReply structure is created to hold the reply data for the sleep pet action.
 *
 * This package provides message structures for the sleep pet action.
 */
type SleepPetMsg struct {
	/**
	 * TargetNickname is the nickname of the pet to be put to sleep.
	 */
	TargetNickname string `json:"target"`
}

/**
 * Function Flow:
 * 1. The SleepPetMsgReply structure is created to hold the reply data for the sleep pet action.
 * 2. The Energy field holds the updated energy value of the pet.
 * 3. The Activity field holds the current activity of the pet.
 * 4. The Duration field holds the duration of the sleep pet action.
 *
 * This structure provides the reply data for the sleep pet action.
 */
type SleepPetMsgReply struct {
	/**
	 * Energy is the updated energy value of the pet.
	 */
	Energy int `json:"energy"`
	/**
	 * Activity is the current activity of the pet.
	 */
	Activity string `json:"activity"`
	/**
	 * Duration is the duration of the sleep pet action.
	 */
	Duration int `json:"duration"`
}

// sleep_pet_msg.go
