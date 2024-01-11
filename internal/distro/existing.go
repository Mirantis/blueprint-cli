package distro

type Existing struct{}

func NewExistingProvider() *Existing {
	return &Existing{}
}
