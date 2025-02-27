// Package msg contains message structures for the Tamagotchi game.
package msg

import (
	"errors"
	"regexp"
)

/**
 * Function Flow:
 * 1. The CreatePetMsg structure is created to hold the nickname for the create pet action.
 * 2. The CreatePetResult structure is created to hold the reply data for the create pet action.
 *
 * This package provides message structures for the create pet action.
 */
type CreatePetMsg struct {
	/**
	 * Nickname is the nickname of the pet to be created.
	 */
	Nickname string `json:"nickname"`
}

/**
 * Function Flow:
 * 1. The CreatePetResult structure is created to hold the reply data for the create pet action.
 * 2. The Success field holds the success status of the create pet action.
 *
 * This structure provides the reply data for the create pet action.
 */
type CreatePetReply struct {
	/**
	 * Success is the success status of the create pet action.
	 */
	Success bool `json:"success"`
}

// Validate checks if the nickname is valid.
func (m *CreatePetMsg) Validate() error {
	if m.Nickname == "" {
		return errors.New("pet nickname cannot be empty")
	}

	if len(m.Nickname) > 16 {
		return errors.New("nickname cannot be longer than 16 characters")
	}

	regex := regexp.MustCompile(`^[a-zA-Z0-9]+$`)
	if !regex.MatchString(m.Nickname) {
		return errors.New("nickname can only contain letters and numbers")
	}

	return nil
}

// TODO: create a sanity function that checks if `Nickname` is:
// - not empty
// - utf8 standard characters (only numbers and letters allow)
// - max long 16 characters
