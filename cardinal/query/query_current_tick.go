// Package query contains functions for querying the state of the Tamagotchi game world.
package query

import (
	"pkg.world.dev/world-engine/cardinal"
)

// CurrentTickMsg is a request message for querying the current tick.
type CurrentTickMsg struct{}

// CurrentTickReply is a response message containing the current tick.
type CurrentTickReply struct {
	CurrentTick uint64 `json:"currentTick"`
}

/**
 * QueryCurrentTick queries the current tick of the game world.
 *
 * Flow:
 * 1. Get the current tick of the game world using the Cardinal API.
 * 2. Return a response message containing the current tick.
 */
// QueryCurrentTick queries the current tick of the game world
func QueryCurrentTick(world cardinal.WorldContext, req *CurrentTickMsg) (*CurrentTickReply, error) {
	// Step 1: Get the current tick of the game world
	//         Use the Cardinal CurrentTick API to get the current tick.
	// Step 2: Return a response message containing the current tick
	//         Return a response message containing the current tick, with no error.
	return &CurrentTickReply{
		// Get the current tick from the world context
		CurrentTick: world.CurrentTick(),
	}, nil
}
