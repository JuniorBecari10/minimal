package value

type Value interface {
	value()
}

// ---

type ValueNumber struct {
	Value float64
}

type ValueString struct {
	Value string
}

type ValueBool struct {
	Value bool
}

// ---

func (x ValueNumber) value() {}
func (x ValueString) value() {}
func (x ValueBool) value()   {}
