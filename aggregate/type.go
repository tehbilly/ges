package aggregate

// Type is used alongside an aggregate repo
type Type struct {
	name    string
	factory func() Root
}

func NewType(name string, factory func() Root) Type {
	return Type{
		name:    name,
		factory: factory,
	}
}

func (t *Type) new() Root {
	return t.factory()
}

func (t *Type) Name() string {
	return t.name
}
