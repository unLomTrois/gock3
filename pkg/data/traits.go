package data

type Traits struct {
	Traits []*Trait
}

func NewTraits() *Traits {
	return &Traits{}
}

func (t *Traits) Folder() string {
	return "common/traits"
}

type Trait struct {
}
