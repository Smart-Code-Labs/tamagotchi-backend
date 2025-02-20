package component

type Activity struct {
	Activity   string
	TotalTicks int
	CountDown  int
	Percentage int
}

func (Activity) Name() string {
	return "Activity"
}
