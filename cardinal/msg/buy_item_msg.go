// Package msg contains message structures for the Tamagotchi game.
package msg

/**
 * Function Flow:
 * 1. The ButItemMsg structure is created to hold the item name for the buy item action.
 * 2. The BuyItemMsgReply structure is created to hold the reply data for the buy item action.
 *
 * This package provides message structures for the buy item action.
 */
type ButItemMsg struct {
	/**
	 * Name is the name of the item to be bought.
	 */
	Name string `json:"name"`
}

/**
 * Function Flow:
 * 1. The BuyItemMsgReply structure is created to hold the reply data for the buy item action.
 * 2. The Success field holds the success status of the buy item action.
 *
 * This structure provides the reply data for the buy item action.
 */
type BuyItemMsgReply struct {
	/**
	 * Success is the success status of the buy item action.
	 */
	Success bool `json:"success"`
}

// buy_item_msg.go
