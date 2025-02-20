package msg

type CurePetMsg struct {
	TargetNickname string `json:"target"`
	ItemId         string `json:"item_id"`
}

type CurePetMsgReply struct {
	Health int `json:"health"`
}
