package component

type Activity struct {
	Activity string
	Duration int
}

func (Activity) Name() string {
	return "Activity"
}
