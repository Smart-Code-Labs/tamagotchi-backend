package msg

type BathPetMsg struct {
	TargetNickname string `json:"target"`
}

type BathPetMsgReply struct {
	Hygiene  int    `json:"hygiene"`
	Activity string `json:"activity"`
	Duration int    `json:"duration"`
}
