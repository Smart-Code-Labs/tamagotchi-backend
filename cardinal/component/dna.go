package component

type Dna struct {
	A int
	C int
	G int
	T int
}

func (Dna) Name() string {
	return "Dna"
}
