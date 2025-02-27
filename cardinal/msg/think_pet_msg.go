// Package msg contains message structures for the Tamagotchi game.
package msg

/**
 * Function Flow:
 * 1. The ThinkPetMsg structure is created to hold the target nickname for the think pet action.
 * 2. The ThinkPetMsgReply structure is created to hold the reply data for the think pet action.
 *
 * This package provides message structures for the think pet action.
 */
type ThinkPetMsg struct {
	/**
	 * TargetNickname is the nickname of the pet to think.
	 */
	TargetNickname string `json:"target"`
}

/**
 * Function Flow:
 * 1. The ThinkPetMsgReply structure is created to hold the reply data for the think pet action.
 * 2. The Think field holds the thought of the pet.
 *
 * This structure provides the reply data for the think pet action.
 */
type ThinkPetMsgReply struct {
	/**
	 * Think is the thought of the pet.
	 */
	Think string `json:"think"`
}
