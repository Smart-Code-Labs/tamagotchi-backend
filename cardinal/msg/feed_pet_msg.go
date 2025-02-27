// Package msg contains message structures for the Tamagotchi game.
package msg

/**
 * Function Flow:
 * 1. The FeedPetMsg structure is created to hold the target nickname and item id for the feed pet action.
 * 2. The FeedPetMsgReply structure is created to hold the reply data for the feed pet action.
 *
 * This package provides message structures for the feed pet action.
 */
type FeedPetMsg struct {
	/**
	 * TargetNickname is the nickname of the pet to be fed.
	 */
	TargetNickname string `json:"target"`
	/**
	 * ItemId is the id of the item to be used for the feed pet action.
	 */
	ItemName string `json:"item_name"`
}

/**
 * Function Flow:
 * 1. The FeedPetMsgReply structure is created to hold the reply data for the feed pet action.
 * 2. The Health field holds the updated health value of the pet.
 * 3. The Activity field holds the current activity of the pet.
 * 4. The Duration field holds the duration of the feed pet action.
 *
 * This structure provides the reply data for the feed pet action.
 */
type FeedPetMsgReply struct {
	/**
	 * Health is the updated health value of the pet.
	 */
	Health int `json:"health"`
	/**
	 * Activity is the current activity of the pet.
	 */
	Activity string `json:"activity"`
	/**
	 * Duration is the duration of the feed pet action.
	 */
	Duration int `json:"duration"`
}

// feed_pet_msg.go
