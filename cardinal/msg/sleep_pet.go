package msg

type SleepPetMsg struct {
	TargetNickname string `json:"target"`
}

type SleepPetMsgReply struct {
	Energy   int    `json:"energy"`
	Activity string `json:"activity"`
	Duration int    `json:"duration"`
}
