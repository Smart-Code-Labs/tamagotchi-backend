// Package msg contains message structures for the Tamagotchi game.
package msg

/**
 * Function Flow:
 * 1. The BreedPetMsg structure is created to hold the mother name, father name, and born name for the breed pet action.
 * 2. The BreedPetMsgReply structure is created to hold the reply data for the breed pet action.
 *
 * This package provides message structures for the breed pet action.
 */
type BreedPetMsg struct {
	/**
	 * MotherName is the name of the mother pet.
	 */
	MotherName string `json:"motherName"`
	/**
	 * FatherName is the name of the father pet.
	 */
	FatherName string `json:"fatherName"`
	/**
	 * BornName is the name of the born pet.
	 */
	BornName string `json:"bornName"`
}

/**
 * Function Flow:
 * 1. The BreedPetMsgReply structure is created to hold the reply data for the breed pet action.
 * 2. The Success field holds the success status of the breed pet action.
 *
 * This structure provides the reply data for the breed pet action.
 */
type BreedPetMsgReply struct {
	/**
	 * Success is the success status of the breed pet action.
	 */
	Success bool `json:"success"`
}

// breed_pet_msg.go
