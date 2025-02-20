package msg

type ThinkPetMsg struct {
	TargetNickname string `json:"target"`
}

type ThinkPetMsgReply struct {
	Think string `json:"think"`
}
