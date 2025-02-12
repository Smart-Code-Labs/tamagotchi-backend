package msg

type PlayPetMsg struct {
	TargetNickname string `json:"target"`
}

type PlayPetMsgReply struct {
	Energy   int    `json:"energy"`
	Hygiene  int    `json:"hygiene"`
	Wellness int    `json:"wellness"`
	Activity string `json:"activity"`
	Duration int    `json:"duration"`
}
