package msg

type EatPetMsg struct {
	TargetNickname string `json:"target"`
}

type EatPetMsgReply struct {
	Health   int    `json:"health"`
	Activity string `json:"activity"`
	Duration int    `json:"duration"`
}
