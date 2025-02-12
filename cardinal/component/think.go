package component

type Think struct {
	Think string
}

func (Think) Name() string {
	return "Think"
}
