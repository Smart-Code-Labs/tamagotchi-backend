package msg

type BreedPetMsg struct {
	MotherName string `json:"motherName"`
	FatherName string `json:"fatherName"`
	BornName   string `json:"bornName"`
}

type BreedPetMsgReply struct {
	Success bool `json:"success"`
}
