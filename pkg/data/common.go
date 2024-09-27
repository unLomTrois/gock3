package data

type Common struct {
	Traits *Traits
}

func NewCommon() *Common {
	return &Common{
		Traits: NewTraits(),
	}
}

func (c *Common) Folder() string {
	return "common"
}

func (c *Common) Load() {

}
